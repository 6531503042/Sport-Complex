package server

import (
	"main/modules/facility/handler"
	"main/modules/facility/repository"
	"main/modules/facility/usecase"
)

func (s *server) facilityService() {
	repo := repository.NewFacilityRepository(s.db)
	facilityUsecase := usecase.NewFacilityUsecase(repo)
	httpHandler := handler.NewFacilityHttpHandler(s.cfg, facilityUsecase)

	facility := s.app.Group("/facility_v1")
	facility.GET("/facility/facilities", httpHandler.FindManyFacility)
	facility.GET("/facility/:facility_id", httpHandler.FindOneFacility)
	facility.POST("/facility/facility", httpHandler.CreateFacility)

	// Slot Routes
	facilitySlot := facility.Group("/:facilityName/slot_v1") // Corrected this line
	facilitySlot.POST("/slots", httpHandler.InsertSlot)         // Create a new slot for a specific facility
	facilitySlot.GET("/slots/:slot_id", httpHandler.FindOneSlot) // GET a slot by slot_id for a specific facility
	facilitySlot.GET("/slots", httpHandler.FindAllSlots)        // GET all slots for a specific facility
}
