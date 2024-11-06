package server

import (
	"log"
	client "main/client/payment"
	"main/modules/payment/handler"
	"main/modules/payment/repository"
	"main/modules/payment/usecase"
)

// paymentService initializes the payment service with routes and handlers
func (s *server) paymentService() {
	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(s.db)

	// Initialize usecases with config and repository
	paymentUsecase := usecase.NewPaymentUsecase(s.cfg, paymentRepo)

	// Initialize the payment client
	paymentClient := client.NewPaymentClient("http://localhost:1327/payment_v1")

	// Initialize HTTP handlers
	paymentHttpHandler := handler.NewPaymentHttpHandler(s.cfg, paymentUsecase, paymentClient)


	// Payment Routes
	payment := s.app.Group("/payment_v1")
	payment.POST("/payments", paymentHttpHandler.CreatePayment) // Create a payment
	payment.GET("/payments/:id", paymentHttpHandler.FindPayment) // Get payment by ID

	// Payment Service
	payment.POST("/payments/slips", paymentHttpHandler.SaveSlip) // Save payment slip
	payment.PUT("/payments/slips/:slipId", paymentHttpHandler.UpdateSlipStatus) // Update payment slip status
	payment.GET("/payments/slips/pending", paymentHttpHandler.GetPendingSlips) // Get pending payment slips



	log.Println("Payment service initialized")
}
