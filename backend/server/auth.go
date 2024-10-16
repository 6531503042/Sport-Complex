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
    repo := repository.NewAuthRepository(s.db)
    usecase := usecase.NewAuthUsecase(repo)
    httpHandler := handlers.NewAuthHttpHandler(s.cfg, usecase)
    grpcHandler := handlers.NewAuthGrpcpHandler(usecase)

    go func() {
        grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)
        authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
        log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("gRPC server error: %v", err)
        }
    }()

    auth := s.app.Group("/auth_v1")
    auth.GET("", s.healthCheckService)
    auth.POST("/auth/login", httpHandler.Login)
    auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
    auth.POST("/auth/logout", httpHandler.Logout)
}