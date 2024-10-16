package handler

import (
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
		CreateBooking (c echo.Context) error
	}

	bookingHttpHandler struct {
		cfg            *config.Config
		bookingUsecase usecase.BookingUsecaseService
	}
)

func NewBookingHttpHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) NewBookingHttpHandlerService {
	return &bookingHttpHandler{cfg: cfg, bookingUsecase: bookingUsecase}
}

func (h *bookingHttpHandler) CreateBooking (c echo.Context) error {
	// Bind the request body to the CreateBookingRequest DTO
    var createBookingReq booking.CreateBookingRequest
    if err := c.Bind(&createBookingReq); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    // Validate the request
    if err := c.Validate(&createBookingReq); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    // Extract the facility name (from URL or context)
    facilityName := c.Param("facilityName")

    // Call usecase to insert the booking
    bookingResponse, err := h.bookingUsecase.InsertBooking(c.Request().Context(), facilityName, &createBookingReq)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    // Return the booking response
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
