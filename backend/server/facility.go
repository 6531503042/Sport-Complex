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
	facilitySlot := facility.Group("/:facilityName/slot_v1") 
	facilitySlot.POST("/slots", httpHandler.InsertSlot)         
	facilitySlot.GET("/slots/:slot_id", httpHandler.FindOneSlot) 
	facilitySlot.GET("/slots", httpHandler.FindAllSlots)      
	
	//Badminton
	badminton := facility.Group("/badminton_v1")
	badminton.POST("/court", httpHandler.InsertBadCourt)
	badminton.POST("/slot", httpHandler.InsertBadmintonSlot)
	badminton.GET("/slots", httpHandler.FindBadmintonSlot)
	badminton.GET("/courts", httpHandler.FindCourt)
}