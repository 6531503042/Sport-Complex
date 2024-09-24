package usecase

import (
	"context"
	"fmt"
	"main/config"
	"main/modules/booking"
	"main/modules/booking/repository"
	"main/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// Convert UserId from string to ObjectID
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var slotId, badmintonSlotId *primitive.ObjectID
	if req.SlotId != nil {
		id, err := primitive.ObjectIDFromHex(*req.SlotId)
		if err != nil {
			return nil, fmt.Errorf("invalid slot ID: %w", err)
		}
		slotId = &id
	} else if req.BadmintonSlotId != nil {
		id, err := primitive.ObjectIDFromHex(*req.BadmintonSlotId)
		if err != nil {
			return nil, fmt.Errorf("invalid badminton slot ID: %w", err)
		}
		badmintonSlotId = &id
	}

	// Create the booking object to pass to the repository
	bookingDoc := &booking.Booking{
		UserId:         userId,
		SlotId:         slotId,
		BadmintonSlotId: badmintonSlotId,
		Status:         1, // Set the initial status
		CreatedAt:      utils.LocalTime(),
		UpdatedAt:      utils.LocalTime(),
	}

	// Insert the booking via repository
	bookingResult, err := u.bookingRepository.InsertBooking(ctx, facilityName, bookingDoc)
	if err != nil {
		return nil, fmt.Errorf("error inserting booking: %w", err)
	}

	// Prepare and return response
	// Prepare and return response
	response := &booking.BookingResponse{ // Change here to use a pointer
		Id:              bookingResult.Id,
		UserId:          bookingResult.UserId.Hex(),
		SlotId:          nil,
		BadmintonSlotId: nil, // Initialize as nil
		Status:          bookingResult.Status,
		CreatedAt:       bookingResult.CreatedAt,
		UpdatedAt:       bookingResult.UpdatedAt,
	}

	if bookingResult.SlotId != nil {
		slotIdStr := bookingResult.SlotId.Hex() // Convert ObjectID to string
		response.SlotId = &slotIdStr // Set the SlotId field
	}

	if bookingResult.BadmintonSlotId != nil {
		badmintonSlotIdStr := bookingResult.BadmintonSlotId.Hex() // Convert ObjectID to string
		response.BadmintonSlotId = &badmintonSlotIdStr // Set the BadmintonSlotId field
	}

	return response, nil // Return a pointer
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
