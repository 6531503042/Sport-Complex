package usecase

import (
	"context"
	"fmt"
	"main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/utils"
)

type(
	BookingUsecaseService interface {
		InsertBooking(ctx context.Context, userId, slotId string) (*booking.Booking, error)
		UpdateBooking (ctx context.Context, bookingId string, status int) (*booking.Booking, error)
		FindBooking (ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking(ctx context.Context, userId string) ([]booking.Booking, error)
	}

	bookingUsecase struct {
		bookingRepository repository.BookingRepositoryService
	}
)

func NewBookingUsecase(bookingRepository repository.BookingRepositoryService) BookingUsecaseService {
	return &bookingUsecase{
		bookingRepository: bookingRepository,
	}
}

func (u *bookingUsecase) InsertBooking(ctx context.Context, userId, slotId string) (*booking.Booking, error) {
    // Ensure the user and slot exist, then create a booking
    newBooking := &booking.Booking{
        UserId:    userId,
        SlotId:    slotId,
        Status:    0, // Initial status
        CreatedAt: utils.LocalTime(),
        UpdatedAt: utils.LocalTime(),
    }

    // Check if the user and slot exist
    if _, err := u.bookingRepository.FindOneUserBooking(ctx, userId); err != nil {
        return nil, fmt.Errorf("error: user %s does not exist", userId)
    }

    if _, err := u.bookingRepository.FindOneSlotBooking(ctx, slotId); err != nil {
        return nil, fmt.Errorf("error: slot %s does not exist", slotId)
    }

    // Create the booking
    createdBooking, err := u.bookingRepository.InsertBooking(ctx, newBooking)
    if err != nil {
        return nil, fmt.Errorf("error: failed to create booking: %w", err)
    }

    return createdBooking, nil
}

func (u *bookingUsecase) UpdateBooking (ctx context.Context, bookingId string, status int) (*booking.Booking, error) {
	booking, err := u.bookingRepository.FindBooking(ctx, bookingId)
	if err != nil {
		return nil, fmt.Errorf("error: failed to find booking: %w", err)
	}

	booking.Status = status
	booking.UpdatedAt = utils.LocalTime()

	updatedBooking, err := u.bookingRepository.UpdateBooking(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("error: failed to update booking: %w", err)
	}

	return updatedBooking, nil
}

func (u *bookingUsecase) FindBooking (ctx context.Context, bookingId string) (*booking.Booking, error) {
	return u.bookingRepository.FindBooking(ctx, bookingId)
}

func (u * bookingUsecase) FindOneUserBooking(ctx context.Context, userId string) ([]booking.Booking, error) {
	return u.bookingRepository.FindOneUserBooking(ctx, userId)
}