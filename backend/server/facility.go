package server

import (
	"main/modules/facility/handler"
	"main/modules/facility/repository"
	"main/modules/facility/usecase"
)

func (s *server) facilityService() {

	// TODO
	repo := repository.NewFacilityRepository(s.db)
	facilityUsecase := usecase.NewFacilityUsecase(repo)
	httpHander := handler.NewFacilityHttpHandler(s.cfg, facilityUsecase)

	facility := s.app.Group("/facility_v1")
	facility.GET("/facility/facilities", httpHander.FindManyFacility)
	facility.GET("/facility/:facility_id", httpHander.FindOneFacility)
	facility.POST("/facility/facility", httpHander.CreateFacility)

}