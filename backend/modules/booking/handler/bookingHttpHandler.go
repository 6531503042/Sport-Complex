package handler

import (
	"main/config"
	"main/modules/booking/usecase"
)

type (
	bookingHttpHandler struct {
		cfg		*config.Config
		bookingUsecase usecase.BookingUsecaseService
	}
)

func NewBookingHttpHandler (cfg config.Config, bookingUsecase usecase.BookingUsecaseService) *bookingHttpHandler {

	return &bookingHttpHandler{
		cfg:         &cfg,
		bookingUsecase: bookingUsecase,
	}
}