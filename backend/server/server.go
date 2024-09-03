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
	return middlewarehttphandler.NewMiddlewareHttpHandler(usecase)
}

func (s *server) httpListening() {
    log.Printf("Starting HTTP server on %s", s.cfg.App.Url)
    err := s.app.Start(s.cfg.App.Url)
    if err != nil && err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }
}


func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {
    log.Printf("Starting graceful shutdown for service: %s", s.cfg.App.Name)

    <-quit
    log.Printf("Received shutdown signal, initiating shutdown...")

    ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

    if err := s.app.Shutdown(ctx); err != nil {
        log.Fatalf("Error during shutdown: %v", err)
    }

    log.Printf("Shutdown completed for service: %s", s.cfg.App.Name)
}


func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
    s := &server{
        app: echo.New(),
        db:  db,
        cfg: cfg,
        middleware: newMiddleware(cfg),
    }

    jwt.SetApiKey(cfg.Jwt.ApiSecretKey)

    s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
        Skipper:      middleware.DefaultSkipper,
        ErrorMessage: "Error: Request Timeout",
        Timeout:      30 * time.Second,
    }))

    s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        Skipper: middleware.DefaultSkipper,
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
    }))

    go func() {
        address := fmt.Sprintf(":%d", cfg.Server.Port)
        log.Printf("Starting HTTP server on %s", address)
        if err := s.app.Start(address); err != nil && err != http.ErrServerClosed {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()

    switch s.cfg.App.Name {
    case "user":
        s.userService()
    case "auth":
        s.authService()
    // Add other service cases here
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    s.app.Use(middleware.Logger())

    go s.gracefulShutdown(pctx, quit)

    s.httpListening()
}
