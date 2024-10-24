package handler

import (
	"encoding/json"
	client "main/client/payment"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/usecase"
	"main/pkg/response"
	"net/http"

	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4"
)

type (
	NewBookingHttpHandlerService interface {
		// InsertBooking(c echo.Context) error

		FindBooking(c echo.Context) error
		FindOneUserBooking(c echo.Context) error
		CreateBooking(c echo.Context) error
	}

	NewBookingHttpHandlerService interface {
		// InsertBooking(c echo.Context) error

		FindBooking(c echo.Context) error
		FindOneUserBooking(c echo.Context) error
		CreateBooking(c echo.Context) error
	}

	bookingHttpHandler struct {
		cfg            *config.Config
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := c.Validate(&createBookingReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	facilityName := c.Param("facilityName")
	bookingResponse, err := h.bookingUsecase.InsertBooking(c.Request().Context(), facilityName, &createBookingReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert booking: " + err.Error()})
	}

	// เรียก API เพื่อรับข้อมูลของสนาม
	facilityURL := "http://localhost:1335/facility_v1/facility/facilities"
	resp, err := http.Get(facilityURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get facility: " + err.Error()})
	}
	defer resp.Body.Close()

	var facilities []struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		PriceInsider float64 `json:"price_insider"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&facilities); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode facility response: " + err.Error()})
	}

	var priceInsider float64
	for _, facility := range facilities {
		if facility.Name == facilityName {
			priceInsider = facility.PriceInsider
			break
		}
	}

	// ใช้ bookingResponse.Id เพื่อดึง ID ที่ถูกต้อง
	paymentRequest := client.CreatePaymentRequest{
		Amount:        priceInsider,
		UserID:        bookingResponse.UserId,
		BookingID:     bookingResponse.Id.Hex(),
		PaymentMethod: "PromptPay",
		Currency:      "THB",
	}

	_, err = h.paymentClient.CreatePayment(paymentRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create payment: " + err.Error()})
	}

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
