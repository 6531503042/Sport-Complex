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
)

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

	// Regular facility routes
	facility.GET("/facilities", fHttpHandler.FindManyFacility)
	facility.GET("/facility/:facility_id", fHttpHandler.FindOneFacility)
	facility.POST("/facility", fHttpHandler.CreateFacility)

	// Slot Routes for each facility type
	facilitySlot := facility.Group("/:facilityName/slot_v1")
	facilitySlot.POST("/slots", fHttpHandler.InsertSlot)
	facilitySlot.GET("/slots", fHttpHandler.FindAllSlots)
	facilitySlot.GET("/slots/:slot_id", fHttpHandler.FindOneSlot)

	// Badminton specific routes
	badminton := facility.Group("/badminton_v1")
	badminton.GET("/courts", fHttpHandler.FindCourt)
	badminton.POST("/court", fHttpHandler.InsertBadCourt)
	badminton.GET("/slots", fHttpHandler.FindBadmintonSlot)
	badminton.POST("/slot", fHttpHandler.InsertBadmintonSlot)

	// Admin routes with analytics
	adminFacility := s.app.Group("/admin/facility_v1", s.middleware.JwtAuthorizationMiddleware(s.cfg))
	
	// Analytics routes
	adminFacility.GET("/analytics/dashboard", aHttpHandler.GetDashboardMetrics)
	adminFacility.GET("/analytics/users", aHttpHandler.GetUserAnalytics)
	
	// Facility management routes
	adminFacility.GET("/facilities", fHttpHandler.FindManyFacility)
	adminFacility.GET("/facility/:facility_id", fHttpHandler.FindOneFacility)
	adminFacility.POST("/facility", fHttpHandler.CreateFacility)
}
