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
	// facility.PUT("/facility/:facility_id",httpHandler.UpdateOneFacility)
	// facility.DELETE("/facility/:facility_id",httpHandler.DeleteOneFacility)

	// Slot Routes
	facilitySlot := facility.Group("/:facilityName/slot_v1") 
	facilitySlot.POST("/slots", httpHandler.InsertSlot)         
	facilitySlot.GET("/slots/:slot_id", httpHandler.FindOneSlot) 
	facilitySlot.GET("/slots", httpHandler.FindAllSlots)      
	// facilitySlot.PUT("/slots/:slot_id", httpHandler.UpdateOneSlot)
	// facilitySlot.DELETE("/slots/:slot_id", httpHandler.DeleteOneSlot)
	
	//Badminton
	badminton := facility.Group("/badminton_v1")
	badminton.POST("/court", httpHandler.InsertBadCourt)
	badminton.POST("/slot", httpHandler.InsertBadmintonSlot)
	badminton.GET("/slots", httpHandler.FindBadmintonSlot)
	badminton.GET("/courts", httpHandler.FindCourt)
	// badminton.GET("/court/:court_id", httpHandler.FindOneCourt)
	// badminton.PUT("/court/:court_id", httpHandler.UpdateOneCourt)
	// badminton.DELETE("/court/:court_id", httpHandler.DeleteOneCourt)
	

}
