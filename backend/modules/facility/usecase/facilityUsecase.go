package usecase

import (
	"context"
	"errors"
	"fmt"
	"main/modules/facility"
	"main/modules/facility/repository"
	"main/pkg/utils"
)

type (
	FacilityUsecaseService interface {

	}

	facilityUsecase struct {
		facilityRepository repository.FacilityRepositoryService
	}
)

func NewFacilityUsecase(facilityRepository repository.FacilityRepositoryService) FacilityUsecaseService {
	return &facilityUsecase{
		facilityRepository: facilityRepository,
	}
}

func (u *facilityUsecase) CreateFacility (pctx context.Context, req *facility.CreateFaciliityRequest) (facility.FacilityBson, error) {
	if !u.facilityRepository.IsUniqueName(pctx, req.Name) {
		return facility.FacilityBson{}, errors.New("error: name already existing")
	}

	//Insert Facility
	facilityId, err := u.facilityRepository.InsertFacility(pctx, &facility.Facilitiy{
		Name: req.Name,
		PriceInsider: req.PriceInsider,
		PriceOutsider: req.PriceOutsider,
		Description: req.Description,
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	})
	if err != nil {
		return facility.FacilityBson{}, fmt.Errorf("error: insert facility failed: %w", err)
	}

	return *u.facilityRepository.FindOneFacility(pctx, facilityId.Hex())
}

