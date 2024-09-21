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
	// slotRepo := repository.NewSlotRepository(s.db)
	
	// Initialize usecases
	// slotUsecase := usecase.NewSlotUsecase(slotRepo)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)

	// Initialize HTTP handlers
	// slotHttpHandler := handler.NewSlotHttpHandler(s.cfg, slotUsecase)
	bookingHttpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase)
	// bookingQueueHandler := handler.NewBookingQueueHandler(s.cfg, bookingUsecase)

	// Booking Routes
	booking := s.app.Group("/booking_v1")
	booking.POST("/bookings", bookingHttpHandler.InsertBooking)       // Create a booking
	booking.GET("/bookings/:booking_id", bookingHttpHandler.FindBooking) // Find a specific booking
	booking.GET("/bookings/user/:user_id", bookingHttpHandler.FindOneUserBooking) // Find all bookings for a specific user
	booking.PUT("/bookings/:id", bookingHttpHandler.UpdateBooking)    // Update a booking

	// Start Booking Queue Consumer (Kafka)
	// go func() {
	// 	log.Println("Starting Booking Queue Consumer")
	// 	booking := &booking.Booking{} // Initialize a new Booking instance
	// 	if err := bookingQueueHandler.AddBooking(booking); err != nil {
	// 		log.Fatalf("Error running booking queue consumer: %v", err)
	// 	}
	// }()

	// Slot Routes
	// slots := s.app.Group("/slots_v1")
	// slots.POST("/slots", slotHttpHandler.InsertSlot)    // Create a slot
	// slots.PUT("/slots/:slot_id", slotHttpHandler.UpdateSlot) // Update a slot
	// slots.GET("/slots", slotHttpHandler.FindAllSlots)    // Find all slots
	// slots.GET("/slots/:slot_id", slotHttpHandler.FindSlot)  // Find a specific slot

	log.Println("Booking service initialized")
}
