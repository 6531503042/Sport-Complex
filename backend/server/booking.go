package server

import (
	"log"
	client "main/client/payment"
	"main/modules/booking/handler"
	"main/modules/booking/repository"
	bookingUsecase "main/modules/booking/usecase"
	facilityUsecase "main/modules/facility/usecase"
	paymentUsecase "main/modules/payment/usecase"
)

// bookingService initializes the booking module, including scheduling the midnight clearing.
func (s *server) bookingService() {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db)

	// Initialize payment client
	paymentClient := client.NewPaymentClient("http://localhost:1327/payment_v1")

	// Initialize usecases using repositories from server struct
	paymentUc := paymentUsecase.NewPaymentUsecase(s.cfg, s.paymentRepo)
	facilityUc := facilityUsecase.NewFacilityUsecase(s.facilityRepo)

	// Store usecases in server struct
	s.paymentUsecase = paymentUc
	s.facilityUsecase = facilityUc

	// Initialize booking usecase with all required dependencies
	bookingUc := bookingUsecase.NewBookingUsecase(
		bookingRepo,
		facilityUc,
		paymentUc,
	)

	// Initialize HTTP handlers
	bookingHttpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUc, paymentClient)

	// Schedule midnight clearing
	go bookingUc.ScheduleMidnightClearing()

	// Booking Routes
	booking := s.app.Group("/booking_v1")
	booking.GET("/bookings/:booking_id", bookingHttpHandler.FindBooking)
	booking.GET("/bookings/user/:user_id", bookingHttpHandler.FindOneUserBooking)
	bookingCreate := booking.Group("/:facilityName")
	bookingCreate.POST("/booking", bookingHttpHandler.CreateBooking, s.middleware.JwtAuthorizationMiddleware(s.cfg)) 
	booking.POST("/bookings/:booking_id/pay", bookingHttpHandler.UpdateBookingStatusToPaid)

	log.Println("Booking service initialized")
}
