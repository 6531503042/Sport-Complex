package server

import (
	"context"
	"log"
	"main/config"
	facilityRepository "main/modules/facility/repository"
	facilityUsecase "main/modules/facility/usecase"
	middlewarehttphandler "main/modules/middleware/middlewareHttpHandler"
	middlewarerepository "main/modules/middleware/middlewareRepository"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	paymentRepository "main/modules/payment/repository"
	paymentUsecase "main/modules/payment/usecase"
	"main/pkg/jwt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app             *echo.Echo
		db             *mongo.Client
		cfg            *config.Config
		middleware     middlewarehttphandler.MiddlewareHttpHandlerService
		validator      *validator.Validate
		paymentRepo    paymentRepository.PaymentRepositoryService
		facilityRepo   facilityRepository.FacilityRepositoryService
		paymentUsecase paymentUsecase.PaymentUsecaseService
		facilityUsecase facilityUsecase.FacilityUsecaseService
	}
)

// CustomValidator wraps go-playground/validator and integrates it with Echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

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
	paymentRepo := paymentRepository.NewPaymentRepository(db)
	facilityRepo := facilityRepository.NewFacilityRepository(db)

	s := &server{
		app:           echo.New(),
		db:            db,
		cfg:           cfg,
		middleware:    newMiddleware(cfg),
		validator:     validator.New(),
		paymentRepo:   paymentRepo,
		facilityRepo:  facilityRepo,
	}

	// Set API Key for JWT
	jwt.SetApiKey(cfg.Jwt.ApiSecretKey)

	// Attach the custom validator
	s.app.Validator = &CustomValidator{validator: s.validator}

	// Basic Middleware
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	// CORS setup
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend origin
		AllowMethods:     []string{echo.GET, echo.POST, echo.HEAD, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true, // Enable credentials handling (e.g., cookies)
	}))

	// Body Limit
	s.app.Use(middleware.BodyLimit("10M"))

	// Logger Middleware
	s.app.Use(middleware.Logger())

	// Route service based on App Name
	switch s.cfg.App.Name {
	case "user":
		s.userService()
	case "auth":
		s.authService()
	case "booking":
		s.bookingService()
	case "facility":
		s.facilityService()
	case "payment":
		s.paymentService()
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefulShutdown(pctx, quit)

	// Start HTTP Server
	s.httpListening()
}