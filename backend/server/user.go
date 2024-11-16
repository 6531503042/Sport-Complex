package server

import (
	"log"
	"main/modules/auth"
	middlewarehttphandler "main/modules/middleware/middlewareHttpHandler"
	middlewarerepository "main/modules/middleware/middlewareRepository"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"main/modules/user/handlers"
	userPb "main/modules/user/proto"
	"main/modules/user/repository"
	"main/modules/user/usecase"
	"main/pkg/grpc"
)

func (s *server) userService() {
	// Initialize repositories, usecases, and handlers
	repo := repository.NewUserRepository(s.db)
	userUsecase := usecase.NewUserUsecase(repo)
	grpcHandler := handlers.NewUserGrpcHandler(userUsecase)
	httpHandler := handlers.NewUserHttpHandler(s.cfg, userUsecase)

	// Initialize middleware
	middlewareRepo := middlewarerepository.NewMiddlewareRepository()
	middlewareUc := middlewareusecase.NewMiddlewareUsecase(middlewareRepo)
	middlewareHandler := middlewarehttphandler.NewMiddlewareHttpHandler(s.cfg, middlewareUc)

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

	// Public endpoints
	user.GET("/health", s.healthCheckService)
	user.POST("/users/register", httpHandler.CreateUser)

	// Protected endpoints (require authentication)
	protected := user.Group("", middlewareHandler.JwtAuthorizationMiddleware(s.cfg))

	// User routes - users can only read their own profile
	protected.GET("/profile/:user_id", httpHandler.FindOneUserProfile,
		middlewareHandler.RequirePermission(auth.PermissionReadUser))

	// Admin routes - requires dashboard access
	admin := protected.Group("/admin",
		middlewareHandler.RequirePermission(auth.PermissionAccessDashboard))

	// Each admin route needs both dashboard access and specific permission
	admin.GET("/users", httpHandler.FindManyUser,
		middlewareHandler.RequirePermission(auth.PermissionReadUser))

	admin.POST("/users", httpHandler.CreateUser,
		middlewareHandler.RequirePermission(auth.PermissionCreateUser))

	admin.PUT("/users/:user_id", httpHandler.UpdateUser,
		middlewareHandler.RequirePermission(auth.PermissionUpdateUser))

	admin.DELETE("/users/:user_id", httpHandler.DeleteUser,
		middlewareHandler.RequirePermission(auth.PermissionDeleteUser))
}