package server

import (
	"log"
	"main/modules/auth/handlers"
	authPb "main/modules/auth/proto"
	"main/modules/auth/repository"
	"main/modules/auth/usecase"
	middlewarehttphandler "main/modules/middleware/middlewareHttpHandler"
	middlewarerepository "main/modules/middleware/middlewareRepository"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"main/pkg/grpc"
)

func (s *server) authService() {
	// Initialize repositories, usecases, and handlers
	repo := repository.NewAuthRepository(s.db)
	usecase := usecase.NewAuthUsecase(repo)
	httpHandler := handlers.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := handlers.NewAuthGrpcpHandler(usecase)

	// Initialize middleware
	middlewareRepo := middlewarerepository.NewMiddlewareRepository()
	middlewareUc := middlewareusecase.NewMiddlewareUsecase(middlewareRepo)
	middlewareHandler := middlewarehttphandler.NewMiddlewareHttpHandler(s.cfg, middlewareUc)

	// Start gRPC server
	go func() {
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)
		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// HTTP Endpoints
	auth := s.app.Group("/auth_v1/auth")

	// Public endpoints (no authentication required)
	auth.GET("/health", s.healthCheckService)
	auth.POST("/login", httpHandler.Login)
	auth.POST("/refresh-token", httpHandler.RefreshToken)

	// Protected endpoints (require authentication)
	protected := auth.Group("", middlewareHandler.JwtAuthorizationMiddleware(s.cfg))
	protected.POST("/logout", httpHandler.Logout)
}