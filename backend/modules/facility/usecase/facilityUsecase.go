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
		CreateFacility(pctx context.Context, req *facility.CreateFaciliityRequest) (facility.FacilityBson, error)
		FindOneFacility(pctx context.Context, facilityId, facilityName string) (*facility.FacilityBson, error)
		FindManyFacility(pctx context.Context) ([]facility.FacilityBson, error)
		UpdateOneFacility(pctx context.Context, facilityId, facilityName string, updateFields map[string]interface{}) error
		DeleteOneFacility(pctx context.Context, facilityId, facilityName string) error

		//Slot - usecase
		InsertSlot(ctx context.Context, startTime, endTime, facilityName string, maxBookings, currentBookings int, facilityType string) (*facility.Slot, error)
		FindOneSlot(ctx context.Context, facilityName, slotId string) (*facility.Slot, error)
		FindManySlot(ctx context.Context,facilityName string) ([]facility.Slot, error)
		EnableOrDisableSlot(ctx context.Context, facilityName, slotId string, status int) (*facility.Slot, error)

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

func (u *facilityUsecase) CreateFacility(pctx context.Context, req *facility.CreateFaciliityRequest) (facility.FacilityBson, error) {
	// Check if the facility name is unique
	if !u.facilityRepository.IsUniqueName(pctx, req.Name) {
		return facility.FacilityBson{}, errors.New("error: name already exists")
	}

	// Insert Facility
	facilityId, err := u.facilityRepository.InsertFacility(pctx, &facility.Facilitiy{
		Name:          req.Name,
		PriceInsider:  req.PriceInsider,
		PriceOutsider: req.PriceOutsider,
		Description:   req.Description,
		CreatedAt:     utils.LocalTime(),
		UpdatedAt:     utils.LocalTime(),
	})
	if err != nil {
		return facility.FacilityBson{}, fmt.Errorf("error: failed to create facility: %w", err)
	}

	// Find the newly inserted facility
	facilityBson, err := u.facilityRepository.FindOneFacility(pctx, facilityId.Hex(), req.Name)
	if err != nil {
		return facility.FacilityBson{}, fmt.Errorf("error: find facility failed: %w", err)
	}

	return *facilityBson, nil
}

func (u *facilityUsecase) FindOneFacility(pctx context.Context, facilityId, facilityName string) (*facility.FacilityBson, error) {
	result, err := u.facilityRepository.FindOneFacility(pctx, facilityId, facilityName)
	if err != nil {
		return nil, err
	}

	// Set the location to Asia/Bangkok
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, fmt.Errorf("error: unable to load time location: %w", err)
	}

	// Return facility details with localized time
	return &facility.FacilityBson{
		Id:            result.Id,
		Name:          result.Name,
		PriceInsider:  result.PriceInsider,
		PriceOutsider: result.PriceOutsider,
		Description:   result.Description,
		CreatedAt:     result.CreatedAt.In(loc),
		UpdatedAt:     result.UpdatedAt.In(loc),
	}, nil
}

func (u *facilityUsecase) FindManyFacility(pctx context.Context) ([]facility.FacilityBson, error) {
	results, err := u.facilityRepository.FindManyFacility(pctx)
	if err != nil {
		return nil, err
	}

	var facilityProfile []facility.FacilityBson
	for _, result := range results {
		facilityProfile = append(facilityProfile, facility.FacilityBson{
			Id:            result.Id,
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

func (u *facilityUsecase) UpdateOneFacility(pctx context.Context, facilityId, facilityName string, updateFields map[string]interface{}) error {
	// Check if the facility exists
	if _, err := u.facilityRepository.FindOneFacility(pctx, facilityId, facilityName); err != nil {
		return err
	}

	// Update the updated_at field
	updateFields["updated_at"] = utils.LocalTime().Format(time.RFC3339)

	// Update the facility
	return u.facilityRepository.UpdateOneFacility(pctx, facilityId, facilityName, updateFields)
}

func (u *facilityUsecase) DeleteOneFacility(pctx context.Context, facilityId, facilityName string) error {
	// Ensure the facility exists before deleting
	_, err := u.facilityRepository.FindOneFacility(pctx, facilityId, facilityName)
	if err != nil {
		return err
	}

	// Delete the facility
	return u.facilityRepository.DeleteOneFacility(pctx, facilityId, facilityName)
}

func (u *facilityUsecase) InsertSlot(ctx context.Context, startTime, endTime, facilityName string, maxBookings, currentBookings int, facilityType string) (*facility.Slot, error) {
	slot := facility.Slot{
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          1,
		MaxBookings:     maxBookings,
		CurrentBookings: currentBookings,
		FacilityType:    facilityType,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Call the repository to insert the slot
	return u.facilityRepository.InsertSlot(ctx, facilityName, slot)
}


func (u *facilityUsecase) FindOneSlot(ctx context.Context, facilityName, slotId string) (*facility.Slot, error) {
	return u.facilityRepository.FindOneSlot(ctx, facilityName, slotId)
}

func (u *facilityUsecase) FindManySlot(ctx context.Context,facilityName string) ([]facility.Slot, error) {
	return u.facilityRepository.FindManySlot(ctx, facilityName)
}

func (u *facilityUsecase) EnableOrDisableSlot(ctx context.Context, facilityName, slotId string, status int) (*facility.Slot, error) {
	return u.facilityRepository.EnableOrDisableSlot(ctx, facilityName, slotId, status)
}