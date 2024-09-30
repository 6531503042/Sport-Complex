package server

import (
	"log"
	"main/modules/payment/handler"
	"main/modules/payment/repository"
	"main/modules/payment/usecase"
)



func (s *server) paymentService() {
	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(s.db)

	// Initialize usecases
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo)

	// Initialize HTTP handlers
	paymentHttpHandler := handler.NewPaymentHttpHandler(s.cfg, paymentUsecase)

	// Payment Routes
	payment := s.app.Group("/payment_v1")
	payment.POST("/payments", paymentHttpHandler.CreatePayment) // Create a payment
	payment.GET("/payments/:id", paymentHttpHandler.FindPayment) // Get payment by ID

	log.Println("Payment service initialized")
}
