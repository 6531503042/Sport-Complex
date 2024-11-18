package usecase

import (
	"context"
	"errors"
	"fmt"
	"main/modules/facility"
	"main/modules/facility/repository"
	"main/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

		UpdateSlot(ctx context.Context, facilityName string, slot *facility.Slot) (*facility.Slot, error)

		//Court - usecase
		InsertBadCourt(ctx context.Context, court *facility.BadmintonCourt) (primitive.ObjectID, error)
		FindBadCourt(ctx context.Context) ([]facility.BadmintonCourt, error)
		InsertBadmintonSlot(ctx context.Context, slot *facility.BadmintonSlot) (primitive.ObjectID, error)
		FindBadmintonSlot(ctx context.Context) ([]facility.BadmintonSlot, error)

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

func (u *facilityUsecase) InsertBadCourt(ctx context.Context, court *facility.BadmintonCourt) (primitive.ObjectID, error) {
	return u.facilityRepository.InsertBadCourt(ctx, court)
}

func (u *facilityUsecase) FindBadCourt(ctx context.Context) ([]facility.BadmintonCourt, error) {
	results, err := u.facilityRepository.FindBadmintonCourt(ctx)
	if err != nil {
		return nil, err
	}

	var badCourt []facility.BadmintonCourt
	for _, result := range results {
		badCourt = append(badCourt, facility.BadmintonCourt{
			Id:          result.Id,
			CourtNumber: result.CourtNumber,
			Status: result.Status,
		})
	}

	return badCourt, nil
}

func (u *facilityUsecase) InsertBadmintonSlot(ctx context.Context, slot *facility.BadmintonSlot) (primitive.ObjectID, error) {
	return u.facilityRepository.InsertBadmintonSlot(ctx, slot)
}

func (u *facilityUsecase) FindBadmintonSlot(ctx context.Context) ([]facility.BadmintonSlot, error) {
	slots, err := u.facilityRepository.FindBadmintonSlot(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure each slot has proper max_bookings value
	for i := range slots {
		if slots[i].MaxBookings == 0 {
			slots[i].MaxBookings = 1 // Set default max bookings to 1
		}
	}

	return slots, nil
}

// UpdateSlot updates a slot's details
func (u *facilityUsecase) UpdateSlot(ctx context.Context, facilityName string, req *facility.Slot) (*facility.Slot, error) {
	return u.facilityRepository.UpdateSlot(ctx, facilityName, req)
}

// DeleteSlot deletes a slot by its ID
func (u *facilityUsecase) DeleteSlot(ctx context.Context, facilityName, slotId string) error {
	return u.facilityRepository.DeleteSlot(ctx, facilityName, slotId)
}

// UpdateBadmintonSlot updates a badminton slot's details
func (u *facilityUsecase) UpdateBadmintonSlot(ctx context.Context, req *facility.BadmintonSlot) error {
	return u.facilityRepository.UpdateBadmintonSlot(ctx, req)
}

// DeleteBadmintonCourt deletes a badminton court by its ID
func (u *facilityUsecase) DeleteBadmintonCourt(ctx context.Context, courtId string) error {
	return u.facilityRepository.DeleteBadmintonCourt(ctx, courtId)
}

// DeleteBadmintonSlot deletes a badminton slot by its ID
func (u *facilityUsecase) DeleteBadmintonSlot(ctx context.Context, slotId string) error {
	return u.facilityRepository.DeleteBadmintonSlot(ctx, slotId)
}