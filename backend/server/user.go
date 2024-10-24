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
	user.GET("/check", s.healthCheckService)
	user.POST("/users/register", httpHandler.CreateUser)
	user.GET("/users/:user_id", httpHandler.FindOneUserProfile)

	//Dashboard
	user.PATCH("/users/:user_id", httpHandler.UpdateUser)
	user.DELETE("/users/:user_id", httpHandler.DeleteUser)
	user.GET("/users", httpHandler.FindManyUser)

	// Admin-only route example
	// userAdmin := user.Group("/admin")
	// userAdmin.Use(s.middleware.IsAdminRoleMiddleware(s.cfg, 1))
	// userAdmin.PATCH("/users/:user_id", httpHandler.UpdateUser)
	// userAdmin.DELETE("/users/:user_id", httpHandler.DeleteUser)
	// userAdmin.GET("/admin-only", httpHandler.AdminOnlyHandler)
}
