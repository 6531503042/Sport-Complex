package usecase

import (
	"context"
	"errors"
	"fmt"
	"main/config"
	"main/modules/booking"
	bm "main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/utils"
	"time"
)

type(
	BookingUsecaseService interface {
		// InsertBooking(ctx context.Context, userId, slotId string) (*booking.Booking, error)
		UpdateBooking (ctx context.Context, bookingId string, status int) (*booking.Booking, error)
		FindBooking (ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking(ctx context.Context, userId string) ([]booking.Booking, error)
		InsertBooking(ctx context.Context, facilityName string, req *booking.CreateBookingRequest) (*booking.BookingResponse, error)

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

func (u *bookingUsecase) InsertBooking(ctx context.Context, facilityName string, req *booking.CreateBookingRequest) (*booking.BookingResponse, error) {
    // Validate slot type and IDs
    if req.SlotType == "normal" && req.SlotId == nil {
        return nil, errors.New("error: SlotId is required for normal bookings")
    }
    if req.SlotType == "badminton" && req.BadmintonSlotId == nil {
        return nil, errors.New("error: BadmintonSlotId is required for badminton bookings")
    }
    if req.SlotId != nil && req.BadmintonSlotId != nil {
        return nil, errors.New("error: Only one of SlotId or BadmintonSlotId should be provided")
    }

    // Create the booking request struct for repository interaction
    bookingReq := &booking.Booking{
        UserId:          req.UserId,
        SlotId:          req.SlotId,
        BadmintonSlotId: req.BadmintonSlotId,
        Status:          1, // Assuming status 1 means active or new
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    // Insert booking using the repository
        // Insert booking using the repository
		booking, err := u.bookingRepository.InsertBooking(ctx, facilityName, bookingReq)
		if err != nil {
			return nil, fmt.Errorf("failed to insert booking: %w", err)
		}
	
		// Map the internal booking struct to the response DTO
		bookingResponse := &bm.BookingResponse{
			Id:              booking.Id,
			UserId:          booking.UserId,
			SlotId:          booking.SlotId,
			BadmintonSlotId: booking.BadmintonSlotId,
			SlotType:        req.SlotType,
			Status:          booking.Status,
			CreatedAt:       booking.CreatedAt,
			UpdatedAt:       booking.UpdatedAt,
		}

    return bookingResponse, nil
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
