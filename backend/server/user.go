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
    // Initialize repository, usecase, and handlers
    repo := repository.NewUserRepository(s.db)
    userUsecase := usecase.NewUserUsecase(repo)
    grpcHandler := handlers.NewUserGrpcHandler(userUsecase)
    httpHandler := handlers.NewUserHttpHandler(s.cfg, userUsecase)

    // gRPC Server
    go func() {
        grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.UserUrl)

        userPb.RegisterUserGrpcServiceServer(grpcServer, grpcHandler)

        log.Printf("User gRPC server listening on %s", s.cfg.Grpc.UserUrl)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("gRPC server error: %v", err)
        }
    }()

    // HTTP Endpoints
    user := s.app.Group("/user_v1")
    user.Get("/check", s.healthCheckService)
    user.Post("/users/register", httpHandler.CreateUser) 
    user.Get("/users/:user_id", httpHandler.FindOneUserProfile)
}