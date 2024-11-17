package server

import (
	"log"
	client "main/client/payment"
	"main/modules/booking/handler"
	bookingPb "main/modules/booking/proto"
	"main/modules/booking/repository"
	"main/modules/booking/usecase"
	"main/pkg/grpc"
)

// bookingService initializes the booking module, including scheduling the midnight clearing.
func (s *server) bookingService() {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db)

	// Initialize usecases
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)

	// Initialize handlers
	paymentClient := client.NewPaymentClient("http://localhost:1327/payment_v1")
	bookingHttpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase, paymentClient)
	bookingGrpcHandler := handler.NewBookingGrpcHandler(bookingUsecase)

	// Initialize and start queue service
	// queueService, err := service.NewBookingQueueService(s.cfg, bookingRepo)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize queue service: %v", err)
	// }
	// go queueService.Start(context.Background())

	// Start gRPC server
	go func() {
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.BookingUrl)
		bookingPb.RegisterBookingServiceServer(grpcServer, bookingGrpcHandler)
		log.Printf("Booking gRPC server listening on %s", s.cfg.Grpc.BookingUrl)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Schedule midnight clearing
	go bookingUsecase.ScheduleMidnightClearing()

	// HTTP routes
	booking := s.app.Group("/booking_v1")
	booking.GET("/bookings/:booking_id", bookingHttpHandler.FindBooking)
	booking.GET("/bookings/user/:user_id", bookingHttpHandler.FindOneUserBooking)
	bookingCreate := booking.Group("/:facilityName")
	bookingCreate.POST("/booking", bookingHttpHandler.CreateBooking, s.middleware.JwtAuthorizationMiddleware(s.cfg))
	booking.POST("/bookings/:booking_id/pay", bookingHttpHandler.UpdateBookingStatusToPaid)

	log.Println("Booking service initialized")
}