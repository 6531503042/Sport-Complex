package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/config"
	"main/modules/booking"
	"main/modules/facility"
	"main/modules/models"
	"main/pkg/queue"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	BookingRepositoryService interface {
		// InsertBooking(ctx context.Context, booking *booking.Booking) (*booking.Booking, error)
		UpdateBooking (ctx context.Context, booking *booking.Booking) (*booking.Booking, error)
		FindBooking(ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking (ctx context.Context, userId string) ([]booking.Booking, error)
		FindOneSlotBooking(ctx context.Context, slotId string) (*facility.Slot, string, error)

		//Kafka Interface
		GetOffset(pctx context.Context) (int64, error)
		UpOffset(pctx context.Context, newOffset int64) error
		InsertBookingViaQueue(pctx context.Context, cfg *config.Config, req *booking.Booking) error
	}

	bookingRepository struct {
		db     *mongo.Client
		client *mongo.Client
	}
)

func NewBookingRepository(db *mongo.Client) BookingRepositoryService {
	return &bookingRepository{
		db:     db,
		client: db,}
}

func (r *bookingRepository) bookingDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("booking_db")
}

func (r *bookingRepository) ListAllFacilities(pctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

	dbs, err := r.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: ListAllFacilities: %s", err.Error())
		return nil, fmt.Errorf("error: list all facilities failed: %w", err)
	}

	var facilityDbs []string
    for _, dbName := range dbs {
        if strings.HasSuffix(dbName, "_facility") {
            facilityDbs = append(facilityDbs, dbName)
        }
    }

    return facilityDbs, nil
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



// func (r *bookingRepository) InsertBooking(ctx context.Context, booking *booking.Booking) (*booking.Booking, error) {
//     ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
//     defer cancel()

//     db := r.bookingDbConn(ctx)
//     col := db.Collection("bookings")

//     // Check if the slot exists using the SlotId from the booking
//     slot, facilityName, err := r.FindOneSlotBooking(ctx, booking.SlotId)
//     if err != nil {
//         return nil, fmt.Errorf("error: slot %s does not exist", booking.SlotId)
//     }

//     // Optionally, you can add additional checks here (e.g., if slot is available for booking)
//     // if slot.Status != facility.Slot {
// 	// 	return nil, fmt.Errorf("error: slot %s is not available for booking", booking.SlotId)
// 	// }

//     // Proceed to insert the booking into the bookings collection
//     result, err := col.InsertOne(ctx, booking)
//     if err != nil {
//         log.Printf("Error: InsertBooking: %s", err.Error())
//         return nil, fmt.Errorf("error: insert booking failed: %w", err)
//     }

//     // Cast the inserted ID to ObjectID
//     bookingId, ok := result.InsertedID.(primitive.ObjectID)
//     if !ok {
//         return nil, fmt.Errorf("error: failed to retrieve inserted booking ID")
//     }

//     // Assign the generated ID to the booking object
//     booking.Id = bookingId

//     // Optionally, log the facility name where the slot was found
//     log.Printf("Booking inserted for slot %s in facility: %s", booking.SlotId, facilityName)

//     return booking, nil
// }


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

func (r *bookingRepository) FindOneSlotBooking(ctx context.Context, slotId string) (*facility.Slot, string, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // List all databases
    dbs, err := r.client.ListDatabaseNames(ctx, bson.M{})
    if err != nil {
        return nil, "", fmt.Errorf("error: failed to list databases: %w", err)
    }

    // Convert string slotId to ObjectID
    objectId, err := primitive.ObjectIDFromHex(slotId)
    if err != nil {
        return nil, "", fmt.Errorf("error: invalid slot ID format: %s", slotId)
    }

    // Iterate through all facility databases
    for _, dbName := range dbs {
        // Check for facility databases that have a suffix "_db"
        if strings.HasSuffix(dbName, "_facility") {
            db := r.client.Database(dbName)
            col := db.Collection("slots")

            // Query the slots collection for the matching _id
            var slot facility.Slot
            err := col.FindOne(ctx, bson.M{"_id": objectId}).Decode(&slot)
            if err == nil {
                // Return the first matching slot and the facility name (from dbName)
                facilityName := strings.TrimSuffix(dbName, "_db")
                return &slot, facilityName, nil
            } else if err != mongo.ErrNoDocuments {
                log.Printf("Error: FindOneSlotBooking in %s: %s", dbName, err.Error())
                continue
            }
        }
    }

    return nil, "", fmt.Errorf("error: slot %s does not exist in any facility", slotId)
}

