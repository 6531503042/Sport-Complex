package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/utils"
	"time"
)

type(
	BookingUsecaseService interface {
		// InsertBooking(ctx context.Context, userId, slotId string) (*booking.Booking, error)
		UpdateBooking (ctx context.Context, bookingId string, status string) (*booking.Booking, error)
		FindBooking (ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking(ctx context.Context, userId string) ([]booking.Booking, error)
		InsertBooking(ctx context.Context, facilityName string, req *booking.CreateBookingRequest) (*booking.BookingResponse, error)

		//Kafka Interface
		GetOffSet(ctx context.Context) (int64, error)
		UpOffSet(ctx context.Context, newOffset int64) error
		UpdateBookingStatusPaid(ctx context.Context, bookingID string, paymentID string, qrCodeURL string) error
		ScheduleMidnightClearing()
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

// ScheduleMidnightClearing schedules the clearing of bookings at midnight every day.
func (u *bookingUsecase) ScheduleMidnightClearing() {
    now := time.Now()
    nextMidnight := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
    duration := nextMidnight.Sub(now)

    log.Printf("Next clearing scheduled in %v", duration)

    time.AfterFunc(duration, func() {
        ctx := context.Background()

        // Execute the midnight clearing process
        if err := u.bookingRepository.ClearingBookingAtMidnight(ctx); err != nil {
            log.Printf("Error clearing bookings at midnight: %s", err.Error())
        } else {
            log.Println("Successfully cleared bookings at midnight")
        }

        // Schedule the next clearing
        u.ScheduleMidnightClearing()
    })
}


// func (u *bookingUsecase) ScheduleMidnightClearing() {
//     log.Println("Clearing process scheduled to run every 1 minute")

//     // Set up the schedule to run every 1 minute
//     time.AfterFunc(time.Minute, func() {
//         ctx := context.Background()

//         // Execute the clearing process
//         if err := u.bookingRepository.ClearingBookingAtMidnight(ctx); err != nil {
//             log.Printf("Error clearing bookings: %s", err.Error())
//         } else {
//             log.Println("Successfully cleared bookings")
//         }

//         // Schedule the next clearing after 1 minute
//         u.ScheduleMidnightClearing()
//     })
// }


//Kafka Func
func (u *bookingUsecase) GetOffSet(ctx context.Context) (int64, error) {
	return u.bookingRepository.GetOffset(ctx)
}

func (u *bookingUsecase) UpOffSet(ctx context.Context, newOffset int64) error {
    return u.bookingRepository.UpOffset(ctx, newOffset)
}

func (u *bookingUsecase) InsertBooking(ctx context.Context, facilityName string, req *booking.CreateBookingRequest) (*booking.BookingResponse, error) {
    
	// validate the request
	if req.SlotType == "normal" && req.SlotId == nil {
        return nil, errors.New("error: SlotId is required for normal bookings")
    }
    if req.SlotType == "badminton" && req.BadmintonSlotId == nil {
        return nil, errors.New("error: BadmintonSlotId is required for badminton bookings")
    }
    if req.SlotId != nil && req.BadmintonSlotId != nil {
        return nil, errors.New("error: Only one of SlotId or BadmintonSlotId should be provided")
    }
	
	// Create the booking request
    bookingReq := &booking.Booking{
        UserId:          req.UserId,
        SlotId:          req.SlotId,
        BadmintonSlotId: req.BadmintonSlotId,
        SlotType:        req.SlotType,
        Status:          "pending",
    }

    // Insert the booking using the queue
    result, err := u.bookingRepository.InsertBookingQueue(ctx, u.cfg, facilityName, bookingReq)
    if err != nil {
        return nil, err
    }

    // Convert to response
    return &booking.BookingResponse{
        Id:              result.Id,
        UserId:          result.UserId,
        SlotId:          result.SlotId,
        BadmintonSlotId: result.BadmintonSlotId,
        SlotType:        result.SlotType,
        Status:          result.Status,
        CreatedAt:       result.CreatedAt,
        UpdatedAt:       result.UpdatedAt,
    }, nil
}


func (u *bookingUsecase) UpdateBooking (ctx context.Context, bookingId string, status string) (*booking.Booking, error) {
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

func (u *bookingUsecase) UpdateBookingStatusPaid(ctx context.Context, bookingID string, paymentID string, qrCodeURL string) error {
    return u.bookingRepository.UpdateStatusPaid(ctx, bookingID, paymentID, qrCodeURL)
}