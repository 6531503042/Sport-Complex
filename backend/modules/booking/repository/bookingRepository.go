package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/models"
	"main/pkg/queue"
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
		FindOneSlotBooking(ctx context.Context, slotId string) (*booking.Slot, error)
		InsertSlot(ctx context.Context, slot *booking.Slot) (*booking.Slot, error)

		//Kafka Interface
		GetOffset(pctx context.Context) (int64, error)
		UpOffset(pctx context.Context, newOffset int64) error
		InsertBookingViaQueue(pctx context.Context, cfg *config.Config, req *booking.Booking) error
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

// Kaka Repo Func
func (r *bookingRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("booking_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}


func (r *bookingRepository) UpOffset(pctx context.Context, newOffset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("booking_queue")

	filter := bson.M{} // Assuming you're updating the only document
	update := bson.M{
		"$set": bson.M{"offset": newOffset},
	}

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error: UpOffset failed: %s", err.Error())
		return errors.New("error: UpOffset failed")
	}

	return nil
}


func (r *bookingRepository) InsertBookingViaQueue(pctx context.Context, cfg *config.Config, req *booking.Booking) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: InsertBookingViaQueue failed: %s", err.Error())
		return errors.New("error: InsertBookingViaQueue failed")
	}


	bookingID := req.Id.Hex() // Convert Object to string

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"booking",
		bookingID,
		reqInBytes,
	); err != nil {
		log.Printf("Error: InsertBookingViaQueue failed: %s", err.Error())
		return errors.New("error: InsertBookingViaQueue failed")
	}

	return nil
}



func (r *bookingRepository) InsertBooking(ctx context.Context, booking *booking.Booking) (*booking.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("bookings")

	// Check if the slot exists
	_, err := r.FindOneSlotBooking(ctx, booking.SlotId)
	if err != nil {
		return nil, fmt.Errorf("error: slot %s does not exist", booking.SlotId)
	}

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

func (r *bookingRepository) FindOneSlotBooking(ctx context.Context, slotId string) (*booking.Slot, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.bookingDbConn(ctx)
    col := db.Collection("slots") // Use the correct collection

    var slot booking.Slot
    // Convert string slotId to ObjectID if needed
    objectId, err := primitive.ObjectIDFromHex(slotId)
    if err != nil {
        return nil, fmt.Errorf("error: invalid slot ID format: %s", slotId)
    }

    // Query using _id in the slots collection
    err = col.FindOne(ctx, bson.M{"_id": objectId}).Decode(&slot)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("error: slot %s does not exist", slotId)
        }
        log.Printf("Error: FindOneSlotBooking: %s", err.Error())
        return nil, fmt.Errorf("error: failed to find slot: %w", err)
    }

    return &slot, nil
}


func (r *bookingRepository) InsertSlot(ctx context.Context, slot *booking.Slot) (*booking.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("slots")

	result, err := col.InsertOne(ctx, slot)
	if err != nil {
		log.Printf("Error: InsertSlot: %s", err.Error())
		return nil, fmt.Errorf("error: insert slot failed: %w", err)
	}

	slotId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("error: insert slot failed")
	}

	slot.Id = slotId
	return slot, nil
}
