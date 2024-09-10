package handler

import (
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/usecase"
	"main/pkg/request"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	NewBookingHttpHandlerService interface {
		InsertBooking(c echo.Context) error
		UpdateBooking(c echo.Context) error
		FindBooking(c echo.Context) error
		FindOneUserBooking(c echo.Context) error
	}

	bookingHttpHandler struct {
		cfg            *config.Config
		bookingUsecase usecase.BookingUsecaseService
	}
)

func NewBookingHttpHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) NewBookingHttpHandlerService {
	return &bookingHttpHandler{cfg: cfg, bookingUsecase: bookingUsecase}
}

func (h *bookingHttpHandler) InsertBooking(c echo.Context) error {
	log.Println("Received request to create booking")

	ctx := c.Request().Context()
	wrapper := request.ContextWrapper(c)

	req := new(booking.CreateBookingReq)

	// Bind and validate the incoming request payload
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// Use the usecase to insert a new booking
	res, err := h.bookingUsecase.InsertBooking(ctx, req.UserId, req.SlotId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *bookingHttpHandler) UpdateBooking(c echo.Context) error {
	log.Println("Received request to update booking")

	ctx := c.Request().Context()
	bookingId := c.Param("booking_id")

	// Bind the request to the UpdateBookingReq struct
	req := new(booking.BookingUpdateReq)

	if err := c.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// Call the usecase with the required parameters
	updatedBooking, err := h.bookingUsecase.UpdateBooking(ctx, bookingId, int(req.SlotId))
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, updatedBooking)
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
