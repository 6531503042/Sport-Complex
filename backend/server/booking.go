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
	httpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase)  // Call the correct constructor function

	// Setup the booking routes
	booking := s.app.Group("/booking_v1")
	booking.POST("/bookings", httpHandler.InsertBooking)
	booking.GET("/bookings/:booking_id", httpHandler.FindBooking)
	booking.GET("/bookings/user/:user_id", httpHandler.FindOneUserBooking)
	booking.PUT("/bookings/:id", httpHandler.UpdateBooking)
}
