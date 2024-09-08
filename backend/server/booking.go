package server

import (
	"main/modules/booking/handler"
	"main/modules/booking/repository"
	"main/modules/booking/usecase"
)

func (s *server) bookingService() {
	// Initialize repository, usecase, and handlers
	repo := repository.NewBookingRepository(s.db)
	bookingUsecase := usecase.NewBookingUsecase(repo)
	httpHander := handler.NewBookingHttpHandlerService(s.cfg, bookingUsecase)

	booking := s.app.Group("/booking_v1")
	booking.POST("/bookings/create", httpHander.InsertBooking)
	booking.GET("/bookings/:booking_id", httpHander.FindBooking)
	booking.GET("/bookings/user/:user_id", httpHander.FindOneUserBooking)
	booking.PUT("/bookings/:id", httpHander.UpdateBooking)
}
