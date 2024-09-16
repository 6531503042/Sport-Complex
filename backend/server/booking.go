package server

import (
	"main/modules/booking/handler"
	"main/modules/booking/repository"
	"main/modules/booking/usecase"
)

func (s *server) bookingService() {

	// Initialize repository, usecase, and handlers
	bookingRepo := repository.NewBookingRepository(s.db)
	slotRepo := repository.NewSlotRepository(s.db)
	slotUsecase := usecase.NewSlotUsecase(slotRepo)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)

	// Initialize the booking HTTP handler
	slotHttpHandler := handler.NewSlotHttpHandler(s.cfg, slotUsecase)
	bookingHttpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase)
	bookingQueueHandler := handler.NewBookingQueueHandler(s.cfg, bookingUsecase)

	//Booking Route
	booking := s.app.Group("/booking_v1")
	booking.POST("/bookings", bookingHttpHandler.InsertBooking)
	booking.GET("/bookings/:booking_id", bookingHttpHandler.FindBooking)
	booking.GET("/bookings/user/:user_id", bookingHttpHandler.FindOneUserBooking)
	booking.PUT("/bookings/:id", bookingHttpHandler.UpdateBooking)
	
	//Message queue technical
	go bookingQueueHandler.AddBooking()
	
	//Slot Route
	slots := s.app.Group("/slots_v1")
	slots.POST("/slots/slots", slotHttpHandler.InsertSlot)
	slots.PUT("/slots/:slot_id", slotHttpHandler.UpdateSlot)
	slots.GET("/slots", slotHttpHandler.FindAllSlots)
	slots.GET("/slots/:slot_id", slotHttpHandler.FindSlot)
}
