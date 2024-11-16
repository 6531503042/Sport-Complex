package server

import (
	"log"
	"main/modules/facility/handler"
	facilityPb "main/modules/facility/proto"
	"main/modules/facility/repository"
	"main/modules/facility/usecase"
	"main/pkg/grpc"
)

func (s *server) facilityService() {
	repo := repository.NewFacilityRepository(s.db)
	facilityUsecase := usecase.NewFacilityUsecase(repo)
	httpHandler := handler.NewFacilityHttpHandler(s.cfg, facilityUsecase)
	grpcHandler := handler.NewFacilityGrpcHandler(facilityUsecase)

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
	facility.GET("/facility/facilities", httpHandler.FindManyFacility)
	facility.GET("/facility/:facility_id", httpHandler.FindOneFacility)
	facility.POST("/facility/facility", httpHandler.CreateFacility)

	// Slot Routes
	facilitySlot := facility.Group("/:facilityName/slot_v1")
	facilitySlot.POST("/slots", httpHandler.InsertSlot)
	facilitySlot.GET("/slots/:slot_id", httpHandler.FindOneSlot)
	facilitySlot.GET("/slots", httpHandler.FindAllSlots)

	// Badminton Routes
	badminton := facility.Group("/badminton_v1")
	badminton.POST("/court", httpHandler.InsertBadCourt)
	badminton.POST("/slot", httpHandler.InsertBadmintonSlot)
	badminton.GET("/slots", httpHandler.FindBadmintonSlot)
	badminton.GET("/courts", httpHandler.FindCourt)
}
