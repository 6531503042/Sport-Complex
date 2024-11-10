package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/config"
	"main/modules/booking"
	bm "main/modules/booking"
	"main/modules/booking/repository"
	facilityUsecase "main/modules/facility/usecase"
	paymentUsecase "main/modules/payment/usecase"
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
		UpdateBookingStatusPaid(ctx context.Context, bookingID string) error
		ScheduleMidnightClearing()
	}

	bookingUsecase struct {
		cfg              *config.Config
		bookingRepository repository.BookingRepositoryService
		facilityService   facilityUsecase.FacilityUsecaseService
		paymentService    paymentUsecase.PaymentUsecaseService
	}
)

func NewBookingUsecase(
	bookingRepository repository.BookingRepositoryService,
	facilityService facilityUsecase.FacilityUsecaseService,
	paymentService paymentUsecase.PaymentUsecaseService,
) BookingUsecaseService {
	return &bookingUsecase{
		cfg: &config.Config{},
		bookingRepository: bookingRepository,
		facilityService:   facilityService,
		paymentService:    paymentService,
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

    var price float64
    dbFacilityName := facilityName + "_facility"  // Add _facility suffix for database name

    // Get facility price based on facility type
    if facilityName == "badminton" {
        // For badminton, look up in badminton_facility database, facility collection
        facility, err := u.facilityService.FindOneFacility(ctx, "", "badminton")
        if err != nil {
            log.Printf("Failed to find badminton facility: %v", err)
            return nil, fmt.Errorf("failed to get facility info: %w", err)
        }
        price = facility.PriceInsider
    } else {
        // For other facilities, look up in {facilityName}_facility database, facility collection
        facility, err := u.facilityService.FindOneFacility(ctx, "", facilityName)
        if err != nil {
            log.Printf("Failed to find facility %s in facility collection: %v", facilityName, err)
            return nil, fmt.Errorf("failed to get facility info: %w", err)
        }
        price = facility.PriceInsider
    }

    // Create payment using the payment service
    payment, err := u.paymentService.CreatePayment(
        ctx,
        req.UserId,
        facilityName,
        price,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create payment: %w", err)
    }

    // Create the booking request struct with payment ID
    bookingReq := &booking.Booking{
        UserId:          req.UserId,
        SlotId:          req.SlotId,
        BadmintonSlotId: req.BadmintonSlotId,
        SlotType:        req.SlotType,
        Status:          "pending",
        PaymentId:       payment.PaymentID,
        Facility:        facilityName,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    // Insert booking using the repository
    createdBooking, err := u.bookingRepository.InsertBooking(ctx, dbFacilityName, bookingReq)
    if err != nil {
        return nil, fmt.Errorf("failed to insert booking: %w", err)
    }

    // Map the internal booking struct to the response DTO
    bookingResponse := &bm.BookingResponse{
        Id:              createdBooking.Id,
        UserId:          createdBooking.UserId,
        SlotId:          createdBooking.SlotId,
        BadmintonSlotId: createdBooking.BadmintonSlotId,
        SlotType:        req.SlotType,
        Status:          createdBooking.Status,
        PaymentID:       payment.PaymentID,
        QRCodeURL:       payment.QRCodeURL,
        CreatedAt:       createdBooking.CreatedAt,
        UpdatedAt:       createdBooking.UpdatedAt,
    }

    return bookingResponse, nil
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

func (u *bookingUsecase) UpdateBookingStatusPaid(ctx context.Context, bookingID string) error {
	return u.bookingRepository.UpdateStatusPaid(ctx, bookingID)
}