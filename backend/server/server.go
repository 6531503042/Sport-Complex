package server

import (
	"context"
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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app *echo.Echo
		db  *mongo.Client
		cfg *config.Config
		middleware middlewarehttphandler.MiddlewareHttpHandlerService
	}
)

// newMiddleware initializes your custom middleware service
func newMiddleware(cfg *config.Config) middlewarehttphandler.MiddlewareHttpHandlerService {
	repo := middlewarerepository.NewMiddlewareRepository()
	usecase := middlewareusecase.NewMiddlewareUsecase(repo)
	return middlewarehttphandler.NewMiddlewareHttpHandler(cfg, usecase)
}

func (s *server) httpListening() {
    log.Printf("Starting HTTP server on %s", s.cfg.App.Url)
    err := s.app.Start(s.cfg.App.Url)
    if err != nil && err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }
}


func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {
	log.Printf("Start service: %s", s.cfg.App.Name)

	<-quit
	log.Printf("Shutting down service: %s", s.cfg.App.Name)

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}


func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
    s := &server{
		app:        echo.New(),
		db:         db,
		cfg:        cfg,
		middleware: newMiddleware(cfg),
	}


    jwt.SetApiKey(cfg.Jwt.ApiSecretKey)

    // Basic Middleware
	// Request Timeout
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

    // CORS
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

    // Body Limit
	s.app.Use(middleware.BodyLimit("10M"))

    switch s.cfg.App.Name {
    case "user":
        s.userService()
    case "auth":
        s.authService()
	case "booking":
		s.bookingService()
	case "facility" :
		s.facilityService()
    // Add other service cases here
    }

    // Graceful Shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    s.app.Use(middleware.Logger())

    go s.gracefulShutdown(pctx, quit)

    // Listening
    s.httpListening()
}