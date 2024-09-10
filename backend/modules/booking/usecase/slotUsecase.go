package usecase

import (
	"context"
	"fmt"
	"main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/utils"
	"time"
)

type (
	SlotUsecaseService interface {
		InsertSlot(ctx context.Context, startTime, endTime string) (*booking.Slot, error)
		UpdateSlot(ctx context.Context, slotId string, startTime, endTime string) (*booking.Slot, error)
		FindOneSlot(ctx context.Context, slotId string) (*booking.Slot, error)
		FindAllSlots(ctx context.Context) ([]booking.Slot, error)
		EnableOrDisableSlot(ctx context.Context, slotId string, status int) (*booking.Slot, error)
	}

	slotUsecase struct {
		slotRepository repository.SlotRepositoryService
	}
)

func NewSlotUsecase(slotRepo repository.SlotRepositoryService) SlotUsecaseService {
	return &slotUsecase{slotRepository: slotRepo}
}



func (u *slotUsecase) InsertSlot(ctx context.Context, startTime, endTime string) (*booking.Slot, error) {
	slot := &booking.Slot{
		StartTime: utils.ParseTimeOnly(startTime),
		EndTime:   utils.ParseTimeOnly(endTime),
		Status:    1, // Enabled by default
	}

	return u.slotRepository.InsertSlot(ctx, slot)

}

func (u *slotUsecase) UpdateSlot(ctx context.Context, slotId string, startTime, endTime string) (*booking.Slot, error) {
	slot, err := u.slotRepository.FindOneSlot(ctx, slotId)
	if err != nil {
		return nil, fmt.Errorf("error: failed to find slot: %w", err)
	}

	slot.StartTime = utils.ParseTimeOnly(startTime)
	slot.EndTime = utils.ParseTimeOnly(endTime)
	slot.UpdatedAt = time.Now()

	return u.slotRepository.UpdateSlot(ctx, slot)
}

func (u *slotUsecase) FindOneSlot(ctx context.Context, slotId string) (*booking.Slot, error) {
	return u.slotRepository.FindOneSlot(ctx, slotId)
}

func (u *slotUsecase) FindAllSlots(ctx context.Context) ([]booking.Slot, error) {
	return u.slotRepository.FindAllSlots(ctx)
}

func (u *slotUsecase) EnableOrDisableSlot(ctx context.Context, slotId string, status int) (*booking.Slot, error) {
	return u.slotRepository.EnableOrDisableSlot(ctx, slotId, status)
}