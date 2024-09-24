package server

import (
	"log"
	"main/modules/booking/handler"
	"main/modules/booking/repository"
	"main/modules/booking/usecase"
)

func (s *server) bookingService() {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db)

	// Initialize usecases
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)

	// Initialize HTTP handlers
	bookingHttpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase)

	// Booking Routes
	booking := s.app.Group("/booking_v1")
	booking.POST("/bookings", bookingHttpHandler.CreateBooking) // Create a booking
	booking.GET("/bookings/:booking_id", bookingHttpHandler.FindBooking) // Find a specific booking
	booking.GET("/bookings/user/:user_id", bookingHttpHandler.FindOneUserBooking) // Find all bookings for a specific user

	log.Println("Booking service initialized")
}
