package server

import (
	"log"


	client "main/client/payment"
	"main/modules/payment/handler"
	"main/modules/payment/repository"
	"main/modules/payment/usecase"
	paymentPb "main/modules/payment/proto"  // Import your generated protobuf package
	"main/pkg/grpc"

	"github.com/IBM/sarama"
)

// paymentService initializes the payment service with routes and gRPC handlers
func (s *server) paymentService() {
	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(s.db)

	// Initialize usecases with config and repository
	paymentUsecase := usecase.NewPaymentUsecase(s.cfg, paymentRepo)

	// Initialize the payment client (if needed for external API calls)
	paymentClient := client.NewPaymentClient("http://localhost:1327/payment_v1")

	// Initialize HTTP handlers
	paymentHttpHandler := handler.NewPaymentHttpHandler(s.cfg, paymentUsecase, paymentClient)

	// Payment Routes (HTTP)
	payment := s.app.Group("/payment_v1")
	payment.POST("/payments", paymentHttpHandler.CreatePayment)            // Create a payment
	payment.GET("/payments/:id", paymentHttpHandler.FindPayment)          // Get payment by ID
	payment.POST("/payments/slips", paymentHttpHandler.SaveSlip)          // Save payment slip
	payment.PUT("/payments/slips/:slipId", paymentHttpHandler.UpdateSlipStatus) // Update payment slip status
	payment.GET("/payments/slips/pending", paymentHttpHandler.GetPendingSlips)  // Get pending payment slips

	// Initialize gRPC handler
	grpcHandler := handler.NewPaymentGrpcHandler(paymentUsecase)

	// Start gRPC server in a separate goroutine
	go func() {
		// Create gRPC listener and server
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PaymentUrl)
		paymentPb.RegisterPaymentServiceServer(grpcServer, grpcHandler)
		log.Printf("Payment gRPC server listening on %s", s.cfg.Grpc.PaymentUrl)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Kafka Topic Setup (Optional)
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{s.cfg.Kafka.Url}, config)
	if err != nil {
		log.Printf("Warning: Failed to create Kafka admin: %v", err)
	} else {
		defer admin.Close()

		// Create necessary Kafka topics
		topics := []string{"payments", "payment.created", "slip.updated"}
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

	log.Println("Payment service initialized")
}
