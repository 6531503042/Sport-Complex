package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/usecase"
	"main/pkg/response"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	NewBookingHttpHandlerService interface {
		FindBooking(c echo.Context) error
		FindOneUserBooking(c echo.Context) error
		CreateBooking(c echo.Context) error
		UpdateBookingStatusToPaid(c echo.Context) error
	}

	bookingHttpHandler struct {
		cfg            *config.Config
		bookingUsecase usecase.BookingUsecaseService
	}

	PaymentRequest struct {
		Amount        float64 `json:"amount"`
		UserID        string  `json:"user_id"`
		BookingID     string  `json:"booking_id"`
		PaymentMethod string  `json:"payment_method"`
		Currency      string  `json:"currency"`
		FacilityName  string  `json:"facility_name"`
	}

	PaymentResponse struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		QRCodeURL string `json:"qr_code_url"`
	}
)

func NewBookingHttpHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) NewBookingHttpHandlerService {
	return &bookingHttpHandler{
		cfg:            cfg,
		bookingUsecase: bookingUsecase,
	}
}

func (h *bookingHttpHandler) CreateBooking(c echo.Context) error {
    var createBookingReq booking.CreateBookingRequest
    if err := c.Bind(&createBookingReq); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    facilityName := c.Param("facilityName")
    log.Printf("Creating booking for facility: %s", facilityName)

    validFacilities := map[string]bool{
        "fitness":    true,
        "swimming":   true,
        "badminton": true,
        "football":   true,
    }

    if !validFacilities[facilityName] {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid facility name. Must be one of: fitness, swimming, badminton, football",
        })
    }

    // Make facility request with timeout
    ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
    defer cancel()

    facilityURL := "http://localhost:1335/facility_v1/facilities"
    req, _ := http.NewRequestWithContext(ctx, "GET", facilityURL, nil)
    
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get facility: " + err.Error()})
    }
    defer resp.Body.Close()

    var facilities []struct {
        Id            primitive.ObjectID `json:"_id"`
        Name          string            `json:"name"`
        PriceInsider  float64           `json:"price_insider"`
        PriceOutsider float64           `json:"price_outsider"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&facilities); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode facility response"})
    }

    var priceInsider float64
    foundFacility := false
    for _, facility := range facilities {
        if facility.Name == facilityName {
            priceInsider = facility.PriceInsider
            foundFacility = true
            break
        }
    }

    if !foundFacility {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Facility not found"})
    }

    // Create booking first
    bookingResponse, err := h.bookingUsecase.InsertBooking(ctx, facilityName, &createBookingReq)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert booking: " + err.Error()})
    }

    // Create payment request
    paymentReq := PaymentRequest{
        Amount:        priceInsider,
        UserID:        bookingResponse.UserId,
        BookingID:     bookingResponse.Id.Hex(),
        PaymentMethod: "PromptPay",
        Currency:      "THB",
        FacilityName:  facilityName,
    }

    paymentURL := "http://localhost:1327/payment_v1/payments"
    jsonData, err := json.Marshal(paymentReq)
    if err != nil {
        log.Printf("Error marshalling payment request: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create payment request"})
    }

    req, err = http.NewRequest("POST", paymentURL, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error creating payment request: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create payment request"})
    }
    req.Header.Set("Content-Type", "application/json")

    paymentClient := &http.Client{Timeout: 5 * time.Second}
    resp, err = paymentClient.Do(req)
    if err != nil {
        log.Printf("Error making payment request: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process payment"})
    }
    defer resp.Body.Close()

    var paymentResult PaymentResponse
    if err := json.NewDecoder(resp.Body).Decode(&paymentResult); err != nil {
        log.Printf("Error decoding payment response: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process payment response"})
    }

    // Update booking with payment info immediately
    if err := h.bookingUsecase.UpdateBookingStatusPaid(ctx, bookingResponse.Id.Hex(), paymentResult.ID, paymentResult.QRCodeURL); err != nil {
        log.Printf("Error updating booking with payment info: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update booking with payment info"})
    }

    // Return response with all information
    return c.JSON(http.StatusAccepted, map[string]interface{}{
        "booking_id":    bookingResponse.Id.Hex(),
        "facility_name": facilityName,
        "payment_id":    paymentResult.ID,
        "qr_code_url":   paymentResult.QRCodeURL,
        "status":        "PROCESSING",
        "message":       fmt.Sprintf("Booking created for %s facility, payment processing", facilityName),
        "amount":        priceInsider,
        "currency":      "THB",
        "created_at":    bookingResponse.CreatedAt,
    })
}

func (h *bookingHttpHandler) FindBooking(c echo.Context) error {
	bookingId := c.Param("booking_id")
	booking, err := h.bookingUsecase.FindBooking(c.Request().Context(), bookingId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"booking":       booking,
		"facility_name": booking.FacilityName,
	})
}

func (h *bookingHttpHandler) FindOneUserBooking(c echo.Context) error {
	userId := c.Param("user_id")
	bookings, err := h.bookingUsecase.FindOneUserBooking(c.Request().Context(), userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, bookings)
}

func (h *bookingHttpHandler) UpdateBookingStatusToPaid(c echo.Context) error {
	bookingID := c.Param("booking_id")
	if bookingID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "booking_id is required"})
	}

	// Pass empty strings for paymentID and qrCodeURL since they're not provided in this endpoint
	err := h.bookingUsecase.UpdateBookingStatusPaid(c.Request().Context(), bookingID, "", "")
	if err != nil {
		log.Printf("Error in UpdateBookingStatusToPaid: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update booking status to paid"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Booking status updated to paid"})
}
