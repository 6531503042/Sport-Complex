package server

import (
	"log"
	analyticsHandler "main/modules/analytics/handler"
	analyticsRepo "main/modules/analytics/repository"
	analyticsUsecase "main/modules/analytics/usecase"
	facilityHandler "main/modules/facility/handler"
	facilityPb "main/modules/facility/proto"
	facilityRepo "main/modules/facility/repository"
	facilityUsecase "main/modules/facility/usecase"
	"main/pkg/grpc"

	"github.com/labstack/echo/v4"
)

// Create middleware functions
func setTimeRange(timeRange string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			q := c.QueryParams()
			q.Set("time_range", timeRange)
			return next(c)
		}
	}
}

func (s *server) facilityService() {
	repo := facilityRepo.NewFacilityRepository(s.db)
	fUsecase := facilityUsecase.NewFacilityUsecase(repo)
	aRepo := analyticsRepo.NewAnalyticsRepository(s.db)
	aUsecase := analyticsUsecase.NewAnalyticsUsecase(aRepo)
	
	fHttpHandler := facilityHandler.NewFacilityHttpHandler(s.cfg, fUsecase)
	aHttpHandler := analyticsHandler.NewAnalyticsHttpHandler(s.cfg, aUsecase)
	grpcHandler := facilityHandler.NewFacilityGrpcHandler(fUsecase)

	// Start gRPC server
	go func() {
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.FacilityUrl)
		facilityPb.RegisterFacilityServiceServer(grpcServer, grpcHandler)
		log.Printf("Facility gRPC server listening on %s", s.cfg.Grpc.FacilityUrl)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// HTTP routes
	facility := s.app.Group("/facility_v1")
	facility.GET("/facility/facilities", fHttpHandler.FindManyFacility)
	facility.GET("/facility/:facility_id", fHttpHandler.FindOneFacility)
	facility.POST("/facility/facility", fHttpHandler.CreateFacility)

	// Slot Routes
	facilitySlot := facility.Group("/:facilityName/slot_v1")
	facilitySlot.POST("/slots", fHttpHandler.InsertSlot)
	facilitySlot.GET("/slots/:slot_id", fHttpHandler.FindOneSlot)
	facilitySlot.GET("/slots", fHttpHandler.FindAllSlots)

	// Badminton Routes
	badminton := facility.Group("/badminton_v1")
	badminton.POST("/court", fHttpHandler.InsertBadCourt)
	badminton.POST("/slot", fHttpHandler.InsertBadmintonSlot)
	badminton.GET("/slots", fHttpHandler.FindBadmintonSlot)
	badminton.GET("/courts", fHttpHandler.FindCourt)

	// Admin routes with analytics
	adminFacility := s.app.Group("/admin/facility_v1", s.middleware.JwtAuthorizationMiddleware(s.cfg))
	
	// Analytics routes with authentication
	analytics := s.app.Group("/analytics_v1", s.middleware.JwtAuthorizationMiddleware(s.cfg))
	
	// Dashboard overview
	analytics.GET("/dashboard/overview", aHttpHandler.GetDashboardMetrics)
	
	// Facility-specific analytics
	analytics.GET("/dashboard/facility/:facility_name", aHttpHandler.GetDashboardMetrics)
	
	// Time-based analytics
	analytics.GET("/dashboard/daily", aHttpHandler.GetDashboardMetrics, setTimeRange("daily"))
	analytics.GET("/dashboard/weekly", aHttpHandler.GetDashboardMetrics, setTimeRange("weekly"))
	analytics.GET("/dashboard/monthly", aHttpHandler.GetDashboardMetrics, setTimeRange("monthly"))
	analytics.GET("/dashboard/yearly", aHttpHandler.GetDashboardMetrics, setTimeRange("yearly"))
	
	// User analytics
	analytics.GET("/users/metrics", aHttpHandler.GetUserAnalytics)

	// Facility management routes
	adminFacility.GET("/facilities", fHttpHandler.FindManyFacility)
	adminFacility.GET("/facility/:facility_id", fHttpHandler.FindOneFacility)
	adminFacility.POST("/facility", fHttpHandler.CreateFacility)
}
