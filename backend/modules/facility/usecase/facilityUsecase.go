package usecase

import (
	"context"
	"errors"
	"fmt"
	"main/modules/facility"
	"main/modules/facility/repository"
	"main/pkg/utils"
	"time"
)

type (
	FacilityUsecaseService interface {
		CreateFacility (pctx context.Context, req *facility.CreateFaciliityRequest) (facility.FacilityBson, error)
		FindOneFacility (pctx  context.Context, facilityId string) (*facility.FacilityBson, error)
		FindManyFacility (pctx context.Context) ([]facility.FacilityBson, error)
		UpdateOneFacility (pctx context.Context, facilityId string, updateFields map[string]interface{}) error
		DeleteOneFacility (pctx context.Context, facilityId string) error
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
	// Find the newly inserted facility
    facilityBson, err := u.facilityRepository.FindOneFacility(pctx, facilityId.Hex())
    if err != nil {
        return facility.FacilityBson{}, fmt.Errorf("error: find facility failed: %w", err)
    }
    if facilityBson == nil {
        return facility.FacilityBson{}, errors.New("error: facility not found")
    }

    return *facilityBson, nil
}


func (u *facilityUsecase) FindOneFacility (pctx  context.Context, facilityId string) (*facility.FacilityBson, error) {
	result, err := u.facilityRepository.FindOneFacility(pctx, facilityId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	if err != nil {
        return nil, fmt.Errorf("error: unable to load time location: %w", err)
    }

	return &facility.FacilityBson{
        Id:            result.Id, // No need to convert ObjectID to Hex manually
        Name:          result.Name,
        PriceInsider:  result.PriceInsider,
        PriceOutsider: result.PriceOutsider,
        Description:   result.Description,
        CreatedAt:     result.CreatedAt.In(loc),
        UpdatedAt:     result.UpdatedAt.In(loc),
    }, nil
}

func (u *facilityUsecase) FindManyFacility (pctx context.Context) ([]facility.FacilityBson, error) {
	results, err := u.facilityRepository.FindManyFacility(pctx)
	if err != nil {
		return nil, err
	}

	var facilityProfile []facility.FacilityBson
	for _, result := range results {
		facilityProfile = append(facilityProfile, facility.FacilityBson{
			Id:            result.Id, // No need to convert ObjectID to Hex manually
			Name:          result.Name,
			PriceInsider:  result.PriceInsider,
			PriceOutsider: result.PriceOutsider,
			Description:   result.Description,
			CreatedAt:     result.CreatedAt,
			UpdatedAt:     result.UpdatedAt,
		})
	}

	return facilityProfile, nil
}

func (u *facilityUsecase) UpdateOneFacility (pctx context.Context, facilityId string, updateFields map[string]interface{}) error {
	if _, err := u.facilityRepository.FindOneFacility(pctx, facilityId); err != nil {
		return err
	}

	updateFields["updated_at"] = utils.LocalTime().Format(time.RFC3339)
	return u.facilityRepository.UpdateOneFacility(pctx, facilityId, updateFields)
}

func (u *facilityUsecase) DeleteOneFacility (pctx context.Context, facilityId string) error {
	_, err := u.facilityRepository.FindOneFacility(pctx, facilityId)
	if err != nil {
		return err
	}

	return u.facilityRepository.DeleteOneFacility(pctx,facilityId)
}