package handler

import (
	"main/config"
	"main/modules/booking/usecase"
)

type (
	BookingQueueHttpHandlerService interface {
		// AddBooking(newBooking *booking.Booking) error
	}

	bookingQueueHandler struct {
		cfg              *config.Config
		bookingUsecase   usecase.BookingUsecaseService
	}
)

func NewBookingQueueHandler(cfg *config.Config, bookingUsecase usecase.BookingUsecaseService) BookingQueueHttpHandlerService {
	return &bookingQueueHandler{cfg: cfg, bookingUsecase: bookingUsecase}
}
