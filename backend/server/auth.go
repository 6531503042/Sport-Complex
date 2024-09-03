package server

import (
	"log"
	"main/modules/auth/handlers"
	authPb "main/modules/auth/proto"
	"main/modules/auth/repository"
	"main/modules/auth/usecase"
	"main/pkg/grpc"
)


func (s *server) authService() {
    // Initialize repository, usecase, and handlers
    repo := repository.NewAuthRepository(s.db)
    usecase := usecase.NewAuthUsecase(repo)
    grpcHandler := handlers.NewAuthGrpcpHandler(usecase)
    httpHandler := handlers.NewAuthHttpHandler(s.cfg, usecase)

    // Start the gRPC server in a new goroutine
    go func() {
        grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)
        authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
        log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("gRPC server error: %v", err)
        }
    }()

    // Fiber routes
    auth := s.app.Group("/auth_v1")
    auth.Get("/check", s.healthCheckService)
    auth.Post("/login", httpHandler.Login)
    auth.Post("/refresh-token", httpHandler.RefreshToken)
    auth.Post("/logout", httpHandler.Logout)
}
