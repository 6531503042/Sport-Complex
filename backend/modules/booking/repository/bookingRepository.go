package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/modules/booking"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	BookingRepositoryService interface {
		InsertBooking(ctx context.Context, booking *booking.Booking) (*booking.Booking, error)
		UpdateBooking (ctx context.Context, booking *booking.Booking) (*booking.Booking, error)
		FindBooking(ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking (ctx context.Context, userId string) ([]booking.Booking, error)
		FindOneSlotBooking(ctx context.Context, slotId string) (*booking.Booking, error)
	}

	bookingRepository struct {
		db *mongo.Client
	}
)

func NewBookingRepository(db *mongo.Client) BookingRepositoryService {
	return &bookingRepository{db}
}

func (r *bookingRepository) bookingDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("booking")
}

func (r *bookingRepository) InsertBooking(ctx context.Context, booking *booking.Booking) (*booking.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("bookings")

	result, err := col.InsertOne(ctx, booking)
	if err != nil {
		log.Printf("Error: InsertBooking: %s", err.Error())
		return nil, fmt.Errorf("error: insert booking failed: %w", err)
	}

	bookingId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("error: insert booking failed")
	}

	booking.Id = bookingId
	return booking, nil
}

func (r *bookingRepository) UpdateBooking (ctx context.Context, booking *booking.Booking) (*booking.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("bookings")
	_, err := col.UpdateOne(ctx, bson.M{"_id": booking.Id}, bson.M{"$set": booking})
    if err != nil {
        return nil, err
    }

    return booking, nil
}

// FindBooking retrieves a booking by ID from the database.
//
// Context:
//   This method is intended to be called within a context that is derived from the
//   http.Request.Context(), which is a context.Context that is associated with the request.
//
// Parameters:
//   ctx - a context that is derived from the http.Request.Context()
//   bookingId - the ID of the booking to retrieve
//
// Returns:
//   a pointer to a booking.Booking, or nil if the booking does not exist
//   an error, or nil if the booking exists and the retrieval was successful
func (r *bookingRepository) FindBooking(ctx context.Context, bookingId string) (*booking.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("bookings")
	result := new(booking.Booking)
    err := col.FindOne(ctx, bson.M{"_id": bookingId}).Decode(result)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (r*bookingRepository) FindOneUserBooking (ctx context.Context, userId string) ([]booking.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("bookings")
	cursor, err := col.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		log.Printf("Error: FindOneUserBooking: %s", err.Error())
		return nil, errors.New("error: find one user booking failed")
	}
	defer cursor.Close(ctx)

	var result []booking.Booking
	if err = cursor.All(ctx, &result); err != nil {
		log.Printf("Error: FindOneUserBooking: %s", err.Error())
		return nil, errors.New("error: find one user booking failed")
	}

	return result, nil
}

func (r *bookingRepository) FindOneSlotBooking(ctx context.Context, slotId string) (*booking.Booking, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.bookingDbConn(ctx)
    col := db.Collection("bookings")

    var booking booking.Booking
    err := col.FindOne(ctx, bson.M{"slotId": slotId}).Decode(&booking)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("error: slot %s does not exist", slotId)
        }
        log.Printf("Error: FindOneSlotBooking: %s", err.Error())
        return nil, fmt.Errorf("error: failed to find slot booking: %w", err)
    }

    return &booking, nil
}