package server

import (
	"log"
	"main/modules/booking"
	"main/modules/booking/handler"
	bookingPb "main/modules/booking/proto"
	"main/modules/booking/repository"
	"main/modules/booking/usecase"
	"main/pkg/grpc"
	"net"

	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
)

// Create a wrapper struct for the HTTP handler
type bookingHandlerWrapper struct {
	handler.NewBookingHttpHandlerService
	queueHandler handler.BookingQueueHandlerService
}

// Implement CreateBooking for the wrapper
func (h *bookingHandlerWrapper) CreateBooking(c echo.Context) error {
	// Call the original handler
	if err := h.NewBookingHttpHandlerService.CreateBooking(c); err != nil {
		return err
	}

	// Get the booking from the response and publish event
	if booking := getBookingFromResponse(c); booking != nil {
		if err := h.queueHandler.PublishBookingCreated(booking); err != nil {
			log.Printf("Failed to publish booking created event: %v", err)
		}
	}
	return nil
}

func (s *server) bookingService() {
	// Check if port is already in use
	ln, err := net.Listen("tcp", s.cfg.App.Url)
	if err != nil {
		log.Fatalf("Port %s is already in use", s.cfg.App.Url)
		return
	}
	ln.Close()

	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db, s.cfg)

	// Initialize usecases
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)

	// Initialize handlers
	queueHandler := handler.NewBookingQueueHandler(s.cfg, bookingUsecase)
	httpHandler := handler.NewBookingHttpHandler(s.cfg, bookingUsecase)
	grpcHandler := handler.NewBookingGrpcHandler(bookingUsecase)

	// Start Kafka consumers in goroutines
	go queueHandler.ProcessBookingQueue()
	go queueHandler.HandleFacilityUpdates()

	// HTTP routes
	booking := s.app.Group("/booking_v1")
	booking.GET("/bookings/:booking_id", httpHandler.FindBooking)
	booking.GET("/bookings/user/:user_id", httpHandler.FindOneUserBooking)
	bookingCreate := booking.Group("/:facilityName")
	bookingCreate.POST("/booking", httpHandler.CreateBooking, s.middleware.JwtAuthorizationMiddleware(s.cfg))
	booking.POST("/bookings/:booking_id/pay", httpHandler.UpdateBookingStatusToPaid)

	// Start gRPC server
	go func() {
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.BookingUrl)
		bookingPb.RegisterBookingServiceServer(grpcServer, grpcHandler)
		log.Printf("Booking gRPC server listening on %s", s.cfg.Grpc.BookingUrl)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Schedule midnight clearing
	go bookingUsecase.ScheduleMidnightClearing()

	// Start cleanup process for expired bookings
	go bookingRepo.CleanupExpiredBookings() 
	

	// Create Kafka topics if they don't exist
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{s.cfg.Kafka.Url}, config)
	if err != nil {
		log.Printf("Warning: Failed to create Kafka admin: %v", err)
	} else {
		defer admin.Close()
		
		topics := []string{"bookings", "booking.created", "facility.updated"}
		for _, topic := range topics {
			topicDetail := &sarama.TopicDetail{
				NumPartitions:     1,
				ReplicationFactor: 1,
			}
			err = admin.CreateTopic(topic, topicDetail, false)
			if err != nil && err != sarama.ErrTopicAlreadyExists {
				log.Printf("Warning: Failed to create topic %s: %v", topic, err)
			}
		}
	}

	

	log.Println("Booking service initialized")
}

// Helper function to extract booking from response
func getBookingFromResponse(c echo.Context) *booking.Booking {
	if resp, ok := c.Get("booking").(*booking.Booking); ok {
		return resp
	}
	return nil
}