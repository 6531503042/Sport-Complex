package usecase

import (
	"context"
	"fmt"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/utils"
)

type(
	BookingUsecaseService interface {
		// InsertBooking(ctx context.Context, userId, slotId string) (*booking.Booking, error)
		UpdateBooking (ctx context.Context, bookingId string, status int) (*booking.Booking, error)
		FindBooking (ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking(ctx context.Context, userId string) ([]booking.Booking, error)

		//Kafka Interface
		GetOffSet(ctx context.Context) (int64, error)
		UpOffSet(ctx context.Context, newOffset int64) error
	}

	bookingUsecase struct {
		cfg              *config.Config
		bookingRepository repository.BookingRepositoryService
	}
)

func NewBookingUsecase(bookingRepository repository.BookingRepositoryService) BookingUsecaseService {
	return &bookingUsecase{
		cfg: &config.Config{},
		bookingRepository: bookingRepository,
	}
}

//Kafka Func
func (u *bookingUsecase) GetOffSet(ctx context.Context) (int64, error) {
	return u.bookingRepository.GetOffset(ctx)
}

func (u *bookingUsecase) UpOffSet(ctx context.Context, newOffset int64) error {
    return u.bookingRepository.UpOffset(ctx, newOffset)
}
// func (u *bookingUsecase) InsertBooking(ctx context.Context, userId, slotId string) (*booking.Booking, error) {
//     // Ensure the user and slot exist before creating a booking
    
//     // Check if the user exists
//     _, err := u.bookingRepository.FindOneUserBooking(ctx, userId)
//     if err != nil {
//         return nil, fmt.Errorf("error: user %s does not exist", userId)
//     }

//     // Check if the slot exists
//     _, err = u.bookingRepository.FindOneSlotBooking(ctx, slotId)
//     if err != nil {
//         return nil, fmt.Errorf("error: slot %s does not exist", slotId)
//     }

//     // Create the new booking object
//     newBooking := &booking.Booking{
//         UserId:    userId,
//         SlotId:    slotId,
//         Status:    0, // Initial status
//         CreatedAt: utils.LocalTime(),
//         UpdatedAt: utils.LocalTime(),
//     }

//     // Insert the booking into MongoDB
//     createdBooking, err := u.bookingRepository.InsertBooking(ctx, newBooking)
//     if err != nil {
//         return nil, fmt.Errorf("error: failed to create booking: %w", err)
//     }

//     // Push the booking to the queue
//     err = u.bookingRepository.InsertBookingViaQueue(ctx, u.cfg, createdBooking) // Pass config here
//     if err != nil {
//         return nil, fmt.Errorf("error: failed to push booking to queue: %w", err)
//     }

//     return createdBooking, nil
// }






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
