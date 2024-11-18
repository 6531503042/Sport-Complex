package repository

import (
	"context"

	"errors"
	"fmt"
	"log"

	"main/modules/booking"
	"main/modules/facility"
	"main/modules/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	BookingRepositoryService interface {

		UpdateBooking (ctx context.Context, booking *booking.Booking) (*booking.Booking, error)
		FindBooking(ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking (ctx context.Context, userId string) ([]booking.Booking, error)
        updateBadmintonSlotBookings(ctx context.Context, slotId primitive.ObjectID, increment int) error
		InsertBooking(pctx context.Context, facilityName string, req *booking.Booking) (*booking.Booking, error)
		

		//Kafka Interface
		GetOffset(pctx context.Context) (int64, error)
		UpOffset(pctx context.Context, newOffset int64) error

		//Clearing system
		ClearingBookingAtMidnight(ctx context.Context) error
		MoveOldBookingTransactionToHistory(ctx context.Context) error 
        ResetFacilitySlots(ctx context.Context, facilityName string) error
        UpdateStatusPaid(ctx context.Context, bookingID string) error

        //Queue with ciritcal section
        // ProcessBookingQueue(pctx context.Context, cfg *config.Config, queueMsg *booking.BookingQueueMessage) error
        // InsertBookingQueue(pctx context.Context, cfg *config.Config, facilityName string, req *booking.Booking) (*booking.Booking, error)
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

func (r *bookingRepository) bookingDbConn(ctx context.Context) *mongo.Database {
	return r.db.Database("booking_db")
}

func (r *bookingRepository) facilityDbConn(ctx context.Context, facilityName string) *mongo.Database {
	// Use the facility name to dynamically create the database name
	databaseName := fmt.Sprintf("%s_facility", facilityName)
	return r.client.Database(databaseName) // This will create the DB if it doesn't exist
}


func (r *bookingRepository) ClearingBookingAtMidnight(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Step 1: Transfer old bookings to history
    if err := r.MoveOldBookingTransactionToHistory(ctx); err != nil {
        log.Printf("Error clearing bookings at midnight: %s", err.Error())
        return err
    }

    // Step 2: Clear all bookings from booking_transaction
    db := r.bookingDbConn(ctx)
    col := db.Collection("booking_transaction")
    _, err := col.DeleteMany(ctx, bson.M{})
    if err != nil {
        log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
        return fmt.Errorf("error: clearingBookingAtMidnight failed during deleting bookings: %w", err)
    }

    //Reset facilitty 
     // Step 3: Reset facility slots for each facility type
     facilities := []string{"fitness", "swimming", "badminton", "football"}

     for _, facilityName := range facilities {
         if err := r.ResetFacilitySlots(ctx, facilityName); err != nil {
             log.Printf("Error resetting slots for facility %s: %s", facilityName, err.Error())
             return fmt.Errorf("error resetting slots for facility %s: %w", facilityName, err)
         }
     }
     log.Println("Successfully cleared booking transactions and reset facility slots")
     return nil
}


func (r *bookingRepository) MoveOldBookingTransactionToHistory(ctx context.Context) error {
    col := r.bookingDbConn(ctx).Collection("booking_transaction")
    historyCol := r.bookingDbConn(ctx).Collection("histories_transaction")

    // Define criteria to find bookings older than today
    now := time.Now().Truncate(24 * time.Hour)
    filter := bson.M{"created_at": bson.M{"$lt": now}}

    // Find old bookings
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
        // Transfer bookings to the history collection
        _, err = historyCol.InsertMany(ctx, bookings)
        if err != nil {
            log.Printf("Error moving bookings to history: %s", err.Error())
            return fmt.Errorf("failed to move bookings to history: %w", err)
        }

        // Only delete if the transfer was successful
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


func (r *bookingRepository) ResetFacilitySlots(ctx context.Context, facilityName string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Connect to the correct facility database
    db := r.facilityDbConn(ctx, facilityName)
    col := db.Collection("slots") // Use the "slots" collection for the facility

    var update bson.M
    var filter bson.M

    if facilityName == "badminton" {
        // For badminton, reset the status field to "available" only if max bookings have been reached
        filter = bson.M{"current_bookings": bson.M{"$gte": 10}} // Example max bookings threshold
        update = bson.M{"$set": bson.M{"status": "available", "current_bookings": 0}} // Reset to available, clear bookings
    } else {
        // For other facilities, reset current_bookings to 0
        filter = bson.M{}
        update = bson.M{"$set": bson.M{"current_bookings": 0}}
    }

    // Update the slots collection for the specific facility
    _, err := col.UpdateMany(ctx, filter, update)
    if err != nil {
        log.Printf("Error resetting slots for facility %s: %s", facilityName, err.Error())
        return fmt.Errorf("error resetting slots for facility %s: %w", facilityName, err)
    }

    log.Printf("Successfully reset slots for facility %s", facilityName)
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

	if err := r.MoveOldBookingTransactionToHistory(ctx); err != nil {
		log.Printf("Error: clearingBookingAtMidnight: %s", err.Error())
		return errors.New("error: clearingBookingAtMidnight failed")
	}

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


func (r *bookingRepository) updateSlotCurrentBooking(ctx context.Context, facilityName string, slotId primitive.ObjectID, increment int) error {
    log.Printf("Updating %s slot %s bookings by %d", facilityName, slotId.Hex(), increment)
    
    db := r.facilityDbConn(ctx, facilityName)
    col := db.Collection("slots")

    // First get current slot state
    var currentSlot facility.Slot
    err := col.FindOne(ctx, bson.M{"_id": slotId}).Decode(&currentSlot)
    if err != nil {
        return fmt.Errorf("failed to get current slot state: %w", err)
    }

    log.Printf("Current bookings before update: %d", currentSlot.CurrentBookings)

    // Update the slot
    update := bson.M{
        "$set": bson.M{
            "current_bookings": currentSlot.CurrentBookings + increment,
            "updated_at": time.Now(),
        },
    }

    result, err := col.UpdateOne(ctx, bson.M{"_id": slotId}, update)
    if err != nil {
        return fmt.Errorf("failed to update slot bookings: %w", err)
    }

    if result.MatchedCount == 0 {
        return errors.New("slot not found")
    }

    // Verify the update
    var updatedSlot facility.Slot
    err = col.FindOne(ctx, bson.M{"_id": slotId}).Decode(&updatedSlot)
    if err != nil {
        return fmt.Errorf("failed to verify update: %w", err)
    }

    log.Printf("Current bookings after update: %d", updatedSlot.CurrentBookings)
    return nil
}

func (r *bookingRepository) InsertBooking(pctx context.Context, facilityName string, req *booking.Booking) (*booking.Booking, error) {
    ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

    col := r.bookingDbConn(ctx).Collection("booking_transaction")

    if err := validateBookingRequest(req); err != nil {
        return nil, err
    }

    log.Printf("Attempting to insert booking for userId: %s", req.UserId)

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

    // Check for duplicate bookings
    exists, err := r.checkDuplicateBooking(ctx, req.UserId, slotIdObject, badmintonSlotIdObject)
    if err != nil {
        log.Printf("Error while checking duplicate booking: %s", err)
        return nil, err
    }
    if exists {
        log.Printf("User %s has already booked slot %v/%v", req.UserId, slotIdObject, badmintonSlotIdObject)
        return nil, errors.New("error: user has already booked this slot")
    }

    var updatedSlot *facility.Slot
    // Check slot availability and update bookings
    if isBadminton {
        // Check badminton slot availability
        count, err := r.countUserBadmintonBookings(ctx, req.UserId)
        if err != nil {
            return nil, err
        }
        if count >= 2 {
            return nil, errors.New("error: user has reached the maximum limit of 2 badminton slots")
        }

        // Get current slot state before update
        updatedSlot, err = r.getBadmintonSlot(ctx, *badmintonSlotIdObject)
        if err != nil {
            return nil, err
        }

        // Update badminton slot booking count
        err = r.updateBadmintonSlotBookings(ctx, *badmintonSlotIdObject, 1)
        if err != nil {
            return nil, fmt.Errorf("failed to update badminton slot booking count: %w", err)
        }

        // Get updated slot state
        updatedSlot, err = r.getBadmintonSlot(ctx, *badmintonSlotIdObject)
        if err != nil {
            return nil, err
        }
    } else {
        // Handle normal slot booking
        slot, err := r.checkSlotAvailability(ctx, facilityName, req)
        if err != nil {
            return nil, err
        }
        if slot.CurrentBookings >= slot.MaxBookings {
            return nil, errors.New("error: Slot is full")
        }
        
        // Update normal slot booking count
        err = r.updateSlotCurrentBooking(ctx, facilityName, *slotIdObject, 1)
        if err != nil {
            return nil, err
        }

        // Get updated slot state
        updatedSlot, err = r.getSlot(ctx, facilityName, *slotIdObject)
        if err != nil {
            return nil, err
        }
    }

    // Create booking document with all necessary fields
    bookingDoc := bson.M{
        "user_id":          req.UserId,
        "facility":         facilityName,
        "status":          "pending",
        "created_at":      time.Now(),
        "updated_at":      time.Now(),
        "slot_type":       isBadminton,
        "current_bookings": updatedSlot.CurrentBookings,
        "max_bookings":    updatedSlot.MaxBookings,
    }

    if req.SlotId != nil {
        bookingDoc["slot_id"] = slotIdObject
    }
    if req.BadmintonSlotId != nil {
        bookingDoc["badminton_slot_id"] = badmintonSlotIdObject
    }

    // Insert booking
    res, err := col.InsertOne(ctx, bookingDoc)
    if err != nil {
        log.Printf("Error inserting booking: %s", err.Error())
        return nil, fmt.Errorf("error inserting booking: %w", err)
    }

    // Create complete response
    req.Id = res.InsertedID.(primitive.ObjectID)
    req.Status = "pending"
    req.CreatedAt = bookingDoc["created_at"].(time.Time)
    req.UpdatedAt = bookingDoc["updated_at"].(time.Time)

    return req, nil
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

func (r *bookingRepository) UpdateStatusPaid(ctx context.Context, bookingID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("booking_transaction")

	// Convert bookingID to ObjectID
	objID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		log.Printf("Error: Invalid booking ID format: %s", err.Error())
		return errors.New("invalid booking ID format")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": "PAID"}}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error: UpdateStatusPaid: %s", err.Error())
		return errors.New("error: update status to paid failed")
	}

	return nil
}

func (r *bookingRepository) checkDuplicateBooking(ctx context.Context, userId string, slotId, badmintonSlotId *primitive.ObjectID) (bool, error) {
    filter := bson.M{
        "user_id": userId,
    }

    // Add slot filter depending on slot type
    if slotId != nil {
        filter["slot_id"] = slotId
    } else if badmintonSlotId != nil {
        filter["badminton_slot_id"] = badmintonSlotId
    }

    // Count documents matching the filter
    count, err := r.bookingDbConn(ctx).Collection("booking_transaction").CountDocuments(ctx, filter)
    if err != nil {
        return false, err
    }

    return count > 0, nil // True if a duplicate booking exists
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



func (r *bookingRepository) countUserBadmintonBookings(ctx context.Context, userId string) (int, error) {
    col := r.bookingDbConn(ctx).Collection("booking_transaction")

    filter := bson.M{
        "user_id":          userId,
        "badminton_slot_id": bson.M{"$exists": true}, // Count only badminton slots
    }

    count, err := col.CountDocuments(ctx, filter)
    if err != nil {
        return 0, fmt.Errorf("error counting badminton bookings: %w", err)
    }

    return int(count), nil
}

// func (r *bookingRepository) InsertBookingQueue(pctx context.Context, cfg *config.Config, facilityName string, req *booking.Booking) (*booking.Booking, error) {
//     ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
//     defer cancel()

//     col := r.bookingDbConn(ctx).Collection("booking_transaction")

//     queueMsg := &booking.BookingQueueMessage{
//         UserId:          req.UserId,
//         SlotId:          req.SlotId,
//         BadmintonSlotId: req.BadmintonSlotId,
//         SlotType:        req.SlotType,
//         FacilityName:    facilityName,
//         CreatedAt:       time.Now(),
//     }

//     if err := validateBookingRequest(req); err != nil {
//         return nil, err
//     }

//     log.Printf("Attempting to insert booking for userId: %s, facilityName: %s", req.UserId, facilityName)

//     var slotIdObject *primitive.ObjectID
//     var badmintonSlotIdObject *primitive.ObjectID
//     isBadminton := false
//     if req.SlotId != nil {
//         slotId, err := primitive.ObjectIDFromHex(*req.SlotId)
//         if err != nil {
//             return nil, fmt.Errorf("invalid SlotId: %w", err)
//         }
//         slotIdObject = &slotId
//     } else if req.BadmintonSlotId != nil {
//         badmintonSlotId, err := primitive.ObjectIDFromHex(*req.BadmintonSlotId)
//         if err != nil {
//             return nil, fmt.Errorf("invalid BadmintonSlotId: %w", err)
//         }
//         badmintonSlotIdObject = &badmintonSlotId
//         isBadminton = true
//     }

//     exists, err := r.checkDuplicateBooking(ctx, req.UserId, slotIdObject, badmintonSlotIdObject)
//     if err != nil {
//         log.Printf("Error while checking duplicate booking: %s", err)
//         return nil, err
//     }
//     if exists {
//         log.Printf("User %s has already booked slot %v/%v", req.UserId, slotIdObject, badmintonSlotIdObject)
//         return nil, errors.New("error: user has already booked this slot")
//     }

//     var slot *facility.Slot
//     if !isBadminton {
//         slot, err = r.checkSlotAvailability(ctx, facilityName, req)
//         if err != nil {
//             return nil, err
//         }
//         if slot.CurrentBookings >= slot.MaxBookings {
//             log.Printf("Slot %v is full for facility %s", slot.Id, facilityName)
//             return nil, errors.New("error: Slot is full")
//         }
//     }

//     req.CreatedAt = time.Now().UTC()
//     req.UpdatedAt = time.Now().UTC()

//     res, err := col.InsertOne(ctx, req)
//     if err != nil {
//         return nil, fmt.Errorf("error inserting booking: %w", err)
//     }

//     req.Id = res.InsertedID.(primitive.ObjectID)

//     message, err := json.Marshal(req)
//     if err != nil {
//         return nil, fmt.Errorf("error marshalling booking for Kafka: %w", err)
//     }

//     //Convert
//     message, err = json.Marshal(queueMsg)
//     if err != nil {
//         return nil, fmt.Errorf("error marshalling booking for Kafka: %w", err)
//     }

//     retryCount := 3
//     for i := 0; i < retryCount; i++ {
//         kafkaErr := queue.PushMessageWithKeyToQueue(
//             []string{cfg.Kafka.Url},
//             cfg.Kafka.ApiKey,
//             cfg.Kafka.Secret,
//             "booking",
//             req.UserId,
//             message,
//         )
//         if kafkaErr == nil {
//             log.Printf("Booking sent to Kafka: UserId=%s, BookingId=%s", req.UserId, req.Id.Hex())
//             break
//         }
//         log.Printf("Failed to send booking to Kafka (attempt %d): %s", i+1, kafkaErr.Error())
//         if i == retryCount-1 {
//             return nil, fmt.Errorf("error sending booking to Kafka after %d attempts: %w", retryCount, kafkaErr)
//         }
//         time.Sleep(time.Second)
//     }

//     if !isBadminton && slotIdObject != nil {
//         err = r.updateSlotCurrentBooking(ctx, facilityName, *slotIdObject, 1)
//         if err != nil {
//             return nil, fmt.Errorf("failed to update slot booking count: %w", err)
//         }
//     } else {
//         log.Printf("Skipping slot availability and booking count update for badminton slot")
//     }

//     return req, nil
// }

// func (r *bookingRepository) ProcessBookingQueue(pctx context.Context, cfg *config.Config, queueMsg *booking.BookingQueueMessage) error {
//     ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
//     defer cancel()

//     conn, err := grpc.NewGrpcClient(cfg.Grpc.FacilityUrl)
//     if err != nil {
//         return fmt.Errorf("error creating gRPC client: %w", err)
//     }

//     var slotId string
//     if queueMsg.SlotId != nil {
//         slotId = *queueMsg.SlotId
//     } else if queueMsg.BadmintonSlotId != nil {
//         slotId = *queueMsg.BadmintonSlotId
//     }

//     bookingDoc := bson.M{
//         "user_id":           queueMsg.UserId,
//         "facility":          queueMsg.FacilityName,
//         "status":           "pending",
//         "created_at":       queueMsg.CreatedAt,
//         "updated_at":       time.Now(),
//     }

//     if queueMsg.SlotId != nil {
//         bookingDoc["slot_id"] = *queueMsg.SlotId
//     } else if queueMsg.BadmintonSlotId != nil {
//         bookingDoc["badminton_slot_id"] = *queueMsg.BadmintonSlotId
//     }

//     col := r.bookingDbConn(ctx).Collection("booking_transaction")
//     _, err = col.InsertOne(ctx, bookingDoc)
//     if err != nil {
//         return fmt.Errorf("error inserting booking: %w", err)
//     }

//     return nil
// }

func (r *bookingRepository) updateBadmintonSlotBookings(ctx context.Context, slotId primitive.ObjectID, increment int) error {
    log.Printf("Updating badminton slot %s bookings by %d", slotId.Hex(), increment)
    
    db := r.facilityDbConn(ctx, "badminton")
    col := db.Collection("slots")

    // First get current slot state
    var currentSlot facility.Slot
    err := col.FindOne(ctx, bson.M{"_id": slotId}).Decode(&currentSlot)
    if err != nil {
        return fmt.Errorf("failed to get current slot state: %w", err)
    }

    log.Printf("Current bookings before update: %d", currentSlot.CurrentBookings)

    // Update the slot
    update := bson.M{
        "$set": bson.M{
            "current_bookings": currentSlot.CurrentBookings + increment,
            "updated_at": time.Now(),
        },
    }

    result, err := col.UpdateOne(ctx, bson.M{"_id": slotId}, update)
    if err != nil {
        return fmt.Errorf("failed to update badminton slot bookings: %w", err)
    }

    if result.MatchedCount == 0 {
        return errors.New("badminton slot not found")
    }

    // Verify the update
    var updatedSlot facility.Slot
    err = col.FindOne(ctx, bson.M{"_id": slotId}).Decode(&updatedSlot)
    if err != nil {
        return fmt.Errorf("failed to verify update: %w", err)
    }

    log.Printf("Current bookings after update: %d", updatedSlot.CurrentBookings)
    return nil
}

// Add these helper functions to get slot information
func (r *bookingRepository) getBadmintonSlot(ctx context.Context, slotId primitive.ObjectID) (*facility.Slot, error) {
    db := r.facilityDbConn(ctx, "badminton")
    col := db.Collection("slots")

    var slot facility.Slot
    err := col.FindOne(ctx, bson.M{"_id": slotId}).Decode(&slot)
    if err != nil {
        return nil, fmt.Errorf("failed to get badminton slot: %w", err)
    }

    return &slot, nil
}

func (r *bookingRepository) getSlot(ctx context.Context, facilityName string, slotId primitive.ObjectID) (*facility.Slot, error) {
    db := r.facilityDbConn(ctx, facilityName)
    col := db.Collection("slots")

    var slot facility.Slot
    err := col.FindOne(ctx, bson.M{"_id": slotId}).Decode(&slot)
    if err != nil {
        return nil, fmt.Errorf("failed to get slot: %w", err)
    }

    return &slot, nil
}
