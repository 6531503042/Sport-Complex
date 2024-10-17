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

		// Insert a new booking into the booking collection.
		// The facility name is used to determine the correct database to use.
		// The booking is inserted with the given request data.
		// The updated booking is returned with the inserted ID.
		InsertBooking(pctx context.Context, facilityName string, req *booking.Booking) (*booking.Booking, error)
		

		//Kafka Interface
		GetOffset(pctx context.Context) (int64, error)
		UpOffset(pctx context.Context, newOffset int64) error

		//Clearing system
		clearingBookingAtMidnight(ctx context.Context) error
		MoveOldBookingTransactionToHistory(ctx context.Context) error


	}

	bookingRepository struct {
		db     *mongo.Client
		client *mongo.Client
	}
)

// NewBookingRepository returns a new instance of BookingRepositoryService using the given mongo client.
// It provides access to the booking database and its collections.
func NewBookingRepository(db *mongo.Client) BookingRepositoryService {
	return &bookingRepository{
		db:     db,
		client: db,}
}

// func ScheduleMidnightClearing(bookingRepository BookingRepositoryService) {
// 	now := time.Now()
// 	nextMidnight := now.Add(time.Hour * 24).Truncate(time.Hour * 24)
// 	duration := nextMidnight.Sub(now)

// 	log.Printf("Next clearing scheduled in %v",duration)

// 	time.AfterFunc(duration, func() {
// 		ctx := context.Background()
// 		// if err := bookingRepository.Clea
// 		if err := bookingRepository.clearingBookingAtMidnight(ctx); err != nil {
// 			log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
// 		}

// 		ScheduleMidnightClearing(bookingRepository)
// 	})
// }

//
func (r *bookingRepository) bookingDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("booking_db")
}

func (r *bookingRepository) facilityDbConn(pctx context.Context, facilityName string) *mongo.Database {
	// Use the facility name to dynamically create the database name
	databaseName := fmt.Sprintf("%s_facility", facilityName)
	return r.client.Database(databaseName) // This will create the DB if it doesn't exist
}

func (r *bookingRepository) clearingBookingAtMidnight(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Step 1: Transfer old bookings to history
    if err := r.MoveOldBookingTransactionToHistory(ctx); err != nil {
        log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
        return errors.New("error: clearingBookingAtMidnight failed during transferring bookings")
    }

    // Step 2: Clear all bookings from booking_transaction
    db := r.bookingDbConn(ctx)
    col := db.Collection("booking_transaction")
    _, err := col.DeleteMany(ctx, bson.M{})
    if err != nil {
        log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
        return fmt.Errorf("error: clearingBookingAtMidnight failed during deleting bookings: %w", err)
    }

    // Step 3: Reset current_bookings in facilities_slot
    slotCol := db.Collection("facilities_slot")
    update := bson.M{"$set": bson.M{"current_bookings": 0}}
    _, err = slotCol.UpdateMany(ctx, bson.M{}, update)
    if err != nil {
        log.Printf("Error resetting current_bookings: %s", err.Error())
        return fmt.Errorf("error resetting current_bookings: %w", err)
    }

    log.Println("Successfully cleared booking transactions and reset current_bookings")
    return nil
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

// func (r *bookingRepository) clearingBookingAtMidnight(ctx context.Context) error {

// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	if err := r.MoveOldBookingTransactionToHistory(ctx); err != nil {
// 		log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
// 		return errors.New("error: clearingBookingAtMidnight failed")
// 	}

// 	db := r.bookingDbConn(ctx)
// 	col := db.Collection("booking_transaction")

// 	_, err := col.DeleteMany(ctx, bson.M{})
// 	if err != nil {
// 		log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
// 	}

// 	return nil
// }

func (r *bookingRepository) MoveOldBookingTransactionToHistory(ctx context.Context) error {
    col := r.bookingDbConn(ctx).Collection("booking_transaction")
    historyCol := r.bookingDbConn(ctx).Collection("histories_transaction")

    // Define criteria to find bookings older than today
    now := time.Now().Truncate(24 * time.Hour)
    filter := bson.M{"created_at": bson.M{"$lt": now}}

    // Find and insert old bookings into history
    cursor, err := col.Find(ctx, filter)
    if err != nil {
        log.Printf("Error retrieving bookings: %s", err.Error())
        return fmt.Errorf("failed to retrieve bookings: %w", err)
    }
    defer cursor.Close(ctx)

    var bookings []interface{}
    if err := cursor.All(ctx, &bookings); err != nil {
        log.Printf("Error reading bookings: %s", err.Error())
        return fmt.Errorf("failed to read bookings: %w", err)
    }

    if len(bookings) > 0 {
        _, err = historyCol.InsertMany(ctx, bookings)
        if err != nil {
            log.Printf("Error moving bookings to history: %s", err.Error())
            return fmt.Errorf("failed to move bookings to history: %w", err)
        }

        // Delete moved bookings from the original collection
        _, err = col.DeleteMany(ctx, filter)
        if err != nil {
            log.Printf("Error deleting old bookings: %s", err.Error())
            return fmt.Errorf("failed to delete old bookings: %w", err)
        }

        log.Printf("Moved %d bookings to history and deleted from booking_transaction", len(bookings))
    } else {
        log.Println("No old bookings to transfer")
    }

    return nil
}


func (r *bookingRepository) checkSlotAvailability(pctx context.Context, facilityName string, req *booking.Booking) (*facility.Slot, error) {
	// Connect to the facility DB (specific to each facility)
	facilityDb := r.facilityDbConn(pctx, facilityName)
	slotCol := facilityDb.Collection("slots")

	var slot facility.Slot
	var err error

	// Check for normal SlotId
	if req.SlotId != nil {
		id, _ := primitive.ObjectIDFromHex(*req.SlotId)
		err = slotCol.FindOne(pctx, bson.M{"_id": id}).Decode(&slot)
	} 
	// Check for BadmintonSlotId
	if req.BadmintonSlotId != nil {
		id, _ := primitive.ObjectIDFromHex(*req.BadmintonSlotId)
		err = slotCol.FindOne(pctx, bson.M{"_id": id}).Decode(&slot)
	}

	// If error occurred, return the appropriate message
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("slot not found or invalid slot ID")
		}
		log.Printf("Error: checkSlotAvailability: %s", err.Error())
		return nil, err
	}

	return &slot, nil
}


func (r *bookingRepository) updateSlotCurrentBooking(pctx context.Context, facilityName string, slotId primitive.ObjectID, increment int) error {
	// Connect to the facility DB
	facilityDb := r.facilityDbConn(pctx, facilityName)
	slotCol := facilityDb.Collection("slots")

	// Update the current booking count
	_, err := slotCol.UpdateOne(
		pctx,
		bson.M{"_id": slotId},
		bson.M{"$inc": bson.M{"current_bookings": increment}},
	)

	// Handle error during update
	if err != nil {
		log.Printf("Error: updateSlotCurrentBooking failed: %s", err.Error())
		return fmt.Errorf("error: updateSlotCurrentBooking failed: %w", err)
	}

	return nil
}

func (r *bookingRepository) checkUserBookingExists(pctx context.Context, userId string, facilityName string, slotId *primitive.ObjectID, badmintonSlotId *primitive.ObjectID) (bool, error) {
    log.Printf("Checking if booking exists for userId: %s, facilityName: %s, slotId: %v, badmintonSlotId: %v", userId, facilityName, slotId, badmintonSlotId)

    db := r.bookingDbConn(pctx)
    col := db.Collection("booking_transaction")

    // Initialize filter with userId
    filter := bson.M{
        "user_id": userId,
    }

    // Add either slotId or badmintonSlotId to the filter
	if slotId != nil {
		filter["slot_id"] = slotId.Hex() // Convert ObjectID to string
	}
	
	log.Printf("MongoDB query filter: %+v", filter)
    if badmintonSlotId != nil {
        filter["badminton_slot_id"] = badmintonSlotId.Hex()
    }
	log.Printf("MongoDB query filter: %+v", filter)

    // Count the number of documents matching this filter
    count, err := col.CountDocuments(pctx, filter)
    if err != nil {
        log.Printf("Error: checkUserBookingExists failed: %s", err.Error())
        return false, fmt.Errorf("error: checkUserBookingExists failed: %w", err)
    }

    log.Printf("Booking exists count for userId: %s, facilityName: %s, slotId: %v, badmintonSlotId: %v -> Count: %d", userId, facilityName, slotId, badmintonSlotId, count)
	log.Printf("Found %d documents matching the filter", count)

    return count > 0, nil
}

// Validate the booking request
func validateBookingRequest(req *booking.Booking) error {
    if req.SlotId == nil && req.BadmintonSlotId == nil {
        return errors.New("SlotId or BadmintonSlotId is required")
    }
    if req.SlotId != nil && req.BadmintonSlotId != nil {
        return errors.New("only one of SlotId or BadmintonSlotId is required")
    }
    // Add additional validations as needed
    return nil
}

func (r *bookingRepository) InsertBooking(pctx context.Context, facilityName string, req *booking.Booking) (*booking.Booking, error) {
    ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

    // Initialize database collection
    col := r.bookingDbConn(ctx).Collection("booking_transaction")

    // Validate incoming request
    if err := validateBookingRequest(req); err != nil {
        return nil, err
    }

    // Log booking attempt
    log.Printf("Attempting to insert booking for userId: %s", req.UserId)

    // Convert SlotId or BadmintonSlotId to ObjectID
    var slotIdObject *primitive.ObjectID
    var badmintonSlotIdObject *primitive.ObjectID
    isBadminton := false
    if req.SlotId != nil {
        slotId, err := primitive.ObjectIDFromHex(*req.SlotId)
        if err != nil {
            return nil, fmt.Errorf("invalid SlotId: %w", err)
        }
        slotIdObject = &slotId
    } else if req.BadmintonSlotId != nil {
        badmintonSlotId, err := primitive.ObjectIDFromHex(*req.BadmintonSlotId)
        if err != nil {
            return nil, fmt.Errorf("invalid BadmintonSlotId: %w", err)
        }
        badmintonSlotIdObject = &badmintonSlotId
        isBadminton = true
    }

    // Check if the user has already booked the same slot
    exists, err := r.checkUserBookingExists(ctx, req.UserId, facilityName, slotIdObject, badmintonSlotIdObject)
    if err != nil {
        log.Printf("Error while checking if user already booked: %s", err)
        return nil, err
    }
    if exists {
        log.Printf("User %s has already booked the slot %v/%v", req.UserId, slotIdObject, badmintonSlotIdObject)
        return nil, errors.New("error: user has already booked this slot")
    }
    log.Printf("User %s has not booked the slot, proceeding with booking", req.UserId)

    var slot *facility.Slot
    if !isBadminton {
        // Check slot availability for normal slots
        slot, err = r.checkSlotAvailability(ctx, facilityName, req)
        if err != nil {
            return nil, err
        }
        if slot.CurrentBookings >= slot.MaxBookings {
            return nil, errors.New("error: Slot is full")
        }
    }

    // Create the booking document
    bookingDoc := bson.M{
        "user_id":    req.UserId,
        "status":     "confirmed", // Assuming the initial status is 'confirmed'
        "created_at": time.Now(),
        "updated_at": time.Now(),
    }

    if req.SlotId != nil {
        bookingDoc["slot_id"] = slotIdObject
    }
    if req.BadmintonSlotId != nil {
        bookingDoc["badminton_slot_id"] = badmintonSlotIdObject
    }

    // Insert the booking document into the collection
    res, err := col.InsertOne(ctx, bookingDoc)
    if err != nil {
        log.Printf("Error inserting booking: %s", err.Error())
        return nil, fmt.Errorf("error inserting booking: %w", err)
    }

    if !isBadminton && slotIdObject != nil {
        // Update the slot's current booking count
        err = r.updateSlotCurrentBooking(ctx, facilityName, *slotIdObject, 1) // Increment by 1 for the new booking
        if err != nil {
            return nil, err
        }
    } else {
        log.Printf("Skipping slot availability and booking count update for badminton slot")
    }

    // Return the inserted booking with the new ID
    req.Id = res.InsertedID.(primitive.ObjectID) // Assign the new ID to the booking
    return req, nil
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
