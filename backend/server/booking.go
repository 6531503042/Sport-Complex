package server

import (
	"log"
	client "main/client/payment"
	"main/modules/booking/handler"
	"main/modules/booking/repository"
	"main/modules/booking/usecase"
)

// bookingService initializes the booking module, including scheduling the midnight clearing.
func (s *server) bookingService() {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db)

	// Initialize usecases
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)

	paymentClient := client.NewPaymentClient("http://localhost:1327/payment_v1")
	// Initialize HTTP handlers
	bookingHttpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase, paymentClient)

	// Schedule midnight clearing
	// go bookingUsecase.ScheduleMidnightClearing()
	//go

	// Booking Routes
	booking := s.app.Group("/booking_v1")
	// booking.POST("/bookings", bookingHttpHandler.CreateBooking) // Create a booking
	booking.GET("/bookings/:booking_id", bookingHttpHandler.FindBooking)          // Find a specific booking
	booking.GET("/bookings/user/:user_id", bookingHttpHandler.FindOneUserBooking) // Find all bookings for a specific user

	bookingCreate := booking.Group("/:facilityName")
	bookingCreate.POST("/booking", bookingHttpHandler.CreateBooking) // Create a booking for a specific facility

	log.Println("Booking service initialized")
}
