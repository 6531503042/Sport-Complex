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
	var req booking.CreateBookingRequest

	// Bind and validate the request
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	facilityName := c.Param("facility_name") // Assuming the facility name is part of the URL path

	// Call the usecase to insert the booking
	response, err := h.bookingUsecase.InsertBooking(c.Request().Context(), facilityName, &req)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Return the successful response
	return c.JSON(http.StatusCreated, response)
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
