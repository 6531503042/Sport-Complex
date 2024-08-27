// package server

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"main/config"
// 	"main/modules/middleware/middlewareHttpHandler"
// 	"main/modules/middleware/middlewareRepository"
// 	"main/modules/middleware/middlewareUsecase"
// 	"net/http"
// 	"time"

// 	"github.com/labstack/echo"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type (
// 	server struct {
// 		app		*echo.Echo
// 		db		*mongo.Client
// 		cfg		*config.Config
// 	}
// )

// // newMiddleware initializes your custom middleware service
// func newMiddleware(cfg *config.Config) middlewareHttpHandler.MiddlewareHttpHandlerService {
// 	repo := middlewareRepository.NewMiddlewareRepository()
// 	usecase := middlewareUsecase.NewMiddlewareUsecase(repo)
// 	return middlewareHttpHandler.NewMiddlewareHttpHandler(cfg, usecase)
// }

// func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
// 	// Initialize the server struct
// 	s := &server{
// 		app: echo.New(),
// 		db:  db,
// 		cfg: cfg,
// 	}

// 	// Set up middleware
// 	s.setupMiddleware()

// 	// Set up routes
// 	s.setupRoutes()

// 	// Start the server in a separate goroutine
// 	go func() {
// 		address := fmt.Sprintf(":%d", cfg.Server.Port) // Ensure cfg.Server.Port is defined
// 		if err := s.app.Start(address); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Shutting down the server due to error: %v", err)
// 		}
// 	}()

// 	// Wait for the context to be canceled (e.g., via OS signal)
// 	<-pctx.Done()

// 	// Attempt graceful shutdown with a timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := s.app.Shutdown(ctx); err != nil {
// 		log.Fatalf("Server Shutdown Failed:%+v", err)
// 	}
// }

// func (s *server) httpListening() {
// 	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
// 		log.Fatalf("Error: %v", err)
// 	}
// }
