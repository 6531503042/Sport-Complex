package server

import (
	"encoding/json"
	"log"
	"main/modules/auth"
	"main/modules/booking"
	middlewarehttphandler "main/modules/middleware/middlewareHttpHandler"
	middlewarerepository "main/modules/middleware/middlewareRepository"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"main/modules/user/handlers"
	userPb "main/modules/user/proto"
	"main/modules/user/repository"
	"main/modules/user/usecase"
	"main/pkg/grpc"
	"main/pkg/queue"
	"time"
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

	// Initialize Kafka server with retry logic
	kafkaServer := queue.NewKafkaServer(s.cfg)
	
	// Register handlers for specific topics with error handling
	kafkaServer.RegisterHandler("booking.created", func(data []byte) error {
		var bookingEvent booking.Booking
		if err := json.Unmarshal(data, &bookingEvent); err != nil {
			log.Printf("Error unmarshalling booking event: %v", err)
			return err
		}
		log.Printf("Processing booking created event: %+v", bookingEvent)
		return nil
	})

	// Start Kafka server with retry
	for i := 0; i < 3; i++ {
		if err := kafkaServer.Start(); err != nil {
			log.Printf("Attempt %d: Failed to start Kafka server: %v", i+1, err)
			if i == 2 {
				log.Printf("Warning: Could not start Kafka server after 3 attempts")
				break
			}
			time.Sleep(time.Second * 2)
			continue
		}
		break
	}

	// Ensure Kafka server is stopped when the service stops
	defer kafkaServer.Stop()

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

	admin.GET("/analytics", httpHandler.GetUserAnalytics,
		middlewareHandler.RequirePermission(auth.PermissionAccessDashboard))
}