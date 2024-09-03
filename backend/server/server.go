package server

import (
	"context"
	"fmt"
	"log"
	"main/config"
	middlewarehttphandler "main/modules/middleware/middlewareHttpHandler"
	middlewarerepository "main/modules/middleware/middlewareRepository"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"main/pkg/jwt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app        *fiber.App
		cfg        *config.Config
		db         *mongo.Client
		middleware middlewarehttphandler.MiddlewareHttpHandlerService
	}
)

func newMiddleware(cfg *config.Config) middlewarehttphandler.MiddlewareHttpHandlerService {
	repo := middlewarerepository.NewMiddlewareRepository()
	usecase := middlewareusecase.NewMiddlewareUsecase(repo)
	return middlewarehttphandler.NewMiddlewareHttpHandler(usecase)
}

func (s *server) httpListening() {
	log.Printf("Starting HTTP server on %s", s.cfg.App.Url)
	if err := s.app.Listen(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func (s *server) gracefulShutdown(quit <-chan os.Signal) {
	log.Printf("Starting graceful shutdown for service: %s", s.cfg.App.Name)

	<-quit
	log.Printf("Received shutdown signal, initiating shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Println("Shutdown timed out")
		// Handle the timeout scenario here, if needed
	}()

	if err := s.app.Shutdown(); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}

	log.Printf("Shutdown completed for service: %s", s.cfg.App.Name)
}

func Start(cfg *config.Config, db *mongo.Client) {
	s := &server{
		app:        fiber.New(),
		cfg:        cfg,
		db:         db,
		middleware: newMiddleware(cfg),
	}

	// Middleware for CORS
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
	}))

	// Middleware for request timeout using a custom handler
	s.app.Use(func(c *fiber.Ctx) error {
		timeout := 30 * time.Second
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		// Wrap the request context with the timeout context
		c.SetUserContext(ctx)

		// Continue processing the request
		return c.Next()
	})

	jwt.SetApiKey(cfg.Jwt.ApiSecretKey)

	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("Starting HTTP server on %s", address)
		if err := s.app.Listen(address); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Add service-specific routes here
	switch cfg.App.Name {
	case "user":
		s.userService()
	case "auth":
		s.authService()
	// Add other service cases here
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Middleware for logging
	s.app.Use(func(c *fiber.Ctx) error {
		log.Printf("[%s] %s %s %s", c.IP(), c.Method(), c.Path(), c.Get("User-Agent"))
		return c.Next()
	})

	go s.gracefulShutdown(quit)

	s.httpListening()
}
