package server

import (
	"context"
	"fmt"
	"log"
	"main/config"
	middlewarehttphandler "main/modules/middleware/middlewareHttpHandler"
	middlewarerepository "main/modules/middleware/middlewareRepository"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app *echo.Echo
		db  *mongo.Client
		cfg *config.Config
	}
)

// newMiddleware initializes your custom middleware service
func newMiddleware(cfg *config.Config) middlewarehttphandler.MiddlewareHttpHandlerService {
	repo := middlewarerepository.NewMiddlewareRepository()
	usecase := middlewareusecase.NewMiddlewareUsecase(repo)
	return middlewarehttphandler.NewMiddlewareUsecase(usecase)
}

func (s *server) httpListening() {
	err := s.app.Start(s.cfg.App.Url)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Shutting down the server due to error: %v", err)
	}
}

func (s *server) gracefulShutdown(pctx context.Context, quit <- chan os.Signal) {
	log.Printf("Start service: %s", s.cfg.App.Name)

	<-quit
	log.Printf("Shutting down service: %s", pctx, 10*time.Second)

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
	// Initialize the server struct
	s := &server{
		app: echo.New(),
		db:  db,
		cfg: cfg,
	}

	// Basic Middleware
	// Request Timeout
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	// CORS
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper: middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	// Start the server in a separate goroutine
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		if err := s.app.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server due to error: %v", err)
		}
	}()

	// Body Limit
	s.app.Use(middleware.BodyLimit("10M"))

	switch s.cfg.App.Name {
	case "user":
		s.userService()

	// Wait for complete service module
	// case "auth":
	// 	s.authService()
	// case "gym":
	// 	s.gymService()
	// case "swimming":
	// 	s.swimmingService()
	// case "badminton":
	// 	s.badmintonService()
	// case "football":
	// 	s.footballService()
	// case "payment":
	// 	s.paymentService()

	}

	//Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s.app.Use(middleware.Logger())

	go s.gracefulShutdown(pctx, quit)

	//Listening
	s.httpListening()
}

