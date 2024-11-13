package handler

import (
	"encoding/json"
	"log"
	client "main/client/payment"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/usecase"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	NewBookingHttpHandlerService interface {
		// InsertBooking(c echo.Context) error

		FindBooking(c echo.Context) error
		FindOneUserBooking(c echo.Context) error
		CreateBooking(c echo.Context) error
		UpdateBookingStatusToPaid(c echo.Context) error
	}

	bookingHttpHandler struct {
		cfg            *config.Config
		bookingUsecase usecase.BookingUsecaseService
		paymentClient  *client.PaymentClient
	}
)

func NewBookingHttpHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService, paymentClient *client.PaymentClient) NewBookingHttpHandlerService {
	return &bookingHttpHandler{cfg: cfg, bookingUsecase: bookingUsecase, paymentClient: paymentClient} // ส่ง payment client
}

func (h *bookingHttpHandler) CreateBooking(c echo.Context) error {
    var createBookingReq booking.CreateBookingRequest
    if err := c.Bind(&createBookingReq); err != nil {
        log.Printf("Error binding request payload: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    if err := c.Validate(&createBookingReq); err != nil {
        log.Printf("Error validating request payload: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    facilityName := c.Param("facilityName")
    log.Printf("Received request to create booking for facility: %s", facilityName)

    bookingResponse, err := h.bookingUsecase.InsertBooking(c.Request().Context(), facilityName, &createBookingReq)
    if err != nil {
        log.Printf("Error inserting booking in database: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert booking: " + err.Error()})
    }

    const PaymentMethods = "PromptPay"
    facilityURL := "http://localhost:1335/facility_v1/facility/facilities"
    log.Printf("Attempting to fetch facility details from URL: %s", facilityURL)

    resp, err := http.Get(facilityURL)
    if err != nil {
        log.Printf("Error making GET request to facility service: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get facility: " + err.Error()})
    }
    defer resp.Body.Close()

    var facilities []struct {
        ID           string  `json:"id"`
        Name         string  `json:"name"`
        PriceInsider float64 `json:"price_insider"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&facilities); err != nil {
        log.Printf("Error decoding facility response JSON: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode facility response: " + err.Error()})
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
        log.Printf("Facility name '%s' not found in response", facilityName)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Facility not found"})
    }

    paymentRequest := client.CreatePaymentRequest{
        Amount:        priceInsider,
        UserID:        bookingResponse.UserId,
        BookingID:     bookingResponse.Id.Hex(),
        PaymentMethod: PaymentMethods,
        Currency:      "THB",
        FacilityName:  facilityName,
    }
    log.Printf("Attempting to create payment with amount: %.2f", priceInsider)

	paymentResponse, err := h.paymentClient.CreatePayment(paymentRequest)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    // Handle "PENDING" status without treating it as an error
    if paymentResponse.Status == "PENDING" {
        bookingResponse.PaymentID = paymentResponse.ID
        bookingResponse.QRCodeURL = paymentResponse.QRCodeURL

        return c.JSON(http.StatusAccepted, map[string]interface{}{
            "message":     "Payment is pending, please complete the payment using the provided QR code.",
            "booking_id":  bookingResponse.Id.Hex(),
            "payment_id":  paymentResponse.ID,
            "qr_code_url": paymentResponse.QRCodeURL,
            "status":      "PENDING",
        })
    }

    // Check if payment was successful
    if paymentResponse.Status != "PAID" {
        return c.JSON(http.StatusPaymentRequired, map[string]string{"error": "Payment is not completed"})
    }

    bookingResponse.PaymentID = paymentResponse.ID
    bookingResponse.QRCodeURL = paymentResponse.QRCodeURL

    // Final response with a successful booking
    return c.JSON(http.StatusOK, bookingResponse)
}


func (h *bookingHttpHandler) FindBooking(c echo.Context) error {
	bookingId := c.Param("booking_id") // Ensure the same parameter name as UpdateBooking
	booking, err := h.bookingUsecase.FindBooking(c.Request().Context(), bookingId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, booking)
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

	err := h.bookingUsecase.UpdateBookingStatusPaid(c.Request().Context(), bookingID)
	if err != nil {
		log.Printf("Error in UpdateBookingStatusToPaid: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update booking status to paid"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Booking status updated to paid"})
}
