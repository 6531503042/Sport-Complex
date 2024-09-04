package server

import (
	"log"
	"main/modules/auth/handlers"
	authPb "main/modules/auth/proto"
	"main/modules/auth/repository"
	"main/modules/auth/usecase"
	"main/pkg/grpc"

	"github.com/gofiber/fiber/v2"
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

    // Route to handle GET /auth_v1/ and return a JSON response
    auth.Get("/", func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "message": "Auth Service v1 - Available endpoints: /check, /auth/login, /auth/refresh-token, /auth/logout",
        })
    })

    // Other routes
    auth.Get("/check", s.healthCheckService)  // Health check route for /auth_v1/check
    auth.Post("/auth/login", httpHandler.Login)
    auth.Post("/auth/refresh-token", httpHandler.RefreshToken)
    auth.Post("/auth/logout", httpHandler.Logout)
}
