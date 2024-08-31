package server

import (
	"log"
	"main/modules/user/handlers"
	userPb "main/modules/user/proto"
	"main/modules/user/repository"
	"main/modules/user/usecase"
	"main/pkg/grpc"
)

func (s *server) userService() {
	// Initialize the repository, usecase, and HTTP handler
	repo := repository.NewUserRepository(s.db)
	userUsecase := usecase.NewUserUsecase(repo)
	grpcHandler := handlers.NewUserGrpcHandler(userUsecase)
	httpHandler := handlers.NewUserHttpHandler(s.cfg, userUsecase)

	// Define a route group for user-related endpoints
	user := s.app.Group("/user_v1")

	// gRPC
	go func() {
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.UserUrl)

		userPb.RegisterUserGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Player gRPC server listening on %s", s.cfg.Grpc.UserUrl)
		grpcServer.Serve(lis)
	}()


	// Health check endpoint (uncomment and implement if needed)
	user.GET("", s.healthCheckService)

	// User endpoints
	user.POST("/users/register", httpHandler.CreateUser) 
	user.GET("/users/:user_id", httpHandler.FindOneUserProfile) 
}