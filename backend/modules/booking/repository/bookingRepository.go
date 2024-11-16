package repository

import (
	"context"
	"encoding/json"

	"errors"
	"fmt"
	"log"

	"main/modules/booking"
	"main/modules/facility"
	facilityPb "main/modules/facility/proto"
	"main/modules/models"
	"main/pkg/grpc"
	"main/pkg/queue"
	"time"

	"main/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	BookingRepositoryService interface {

		UpdateBooking (ctx context.Context, booking *booking.Booking) (*booking.Booking, error)
		FindBooking(ctx context.Context, bookingId string) (*booking.Booking, error)
		FindOneUserBooking (ctx context.Context, userId string) ([]booking.Booking, error)

		InsertBooking(pctx context.Context, facilityName string, req *booking.Booking) (*booking.Booking, error)
		

		//Kafka Interface
		GetOffset(pctx context.Context) (int64, error)
		UpOffset(pctx context.Context, newOffset int64) error

		//Clearing system
		ClearingBookingAtMidnight(ctx context.Context) error
		MoveOldBookingTransactionToHistory(ctx context.Context) error 
        ResetFacilitySlots(ctx context.Context, facilityName string) error
        UpdateStatusPaid(ctx context.Context, bookingID string, paymentID string, qrCodeURL string) error

        //Queue with ciritcal section
        ProcessBookingQueue(pctx context.Context, cfg *config.Config, queueMsg *booking.BookingQueueMessage) error
        InsertBookingQueue(pctx context.Context, cfg *config.Config, facilityName string, req *booking.Booking) (*booking.Booking, error)
        CleanupExpiredBookings()
	}

	bookingRepository struct {
		db     *mongo.Client
		client *mongo.Client
        cfg    *config.Config
	}
)

// NewBookingRepository returns a new instance of BookingRepositoryService using the given mongo client.
// It provides access to the booking database and its collections.
func NewBookingRepository(db *mongo.Client, cfg *config.Config) BookingRepositoryService {
	return &bookingRepository{
		db:     db,
		client: db,
		cfg:    cfg,
	}
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
	col := db.Collection("booking_transaction")

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
	col := db.Collection("booking_transaction")

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

func (r *bookingRepository) InsertBooking(pctx context.Context, facilityName string, req *booking.Booking) (*booking.Booking, error) {
    ctx, cancel := context.WithTimeout(pctx, 120*time.Second)
    defer cancel()

    col := r.bookingDbConn(ctx).Collection("booking_transaction")

    // Validate incoming request
    if err := validateBookingRequest(req); err != nil {
        return nil, err
    }

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

    // Check badminton booking limit
    if isBadminton {
        count, err := r.countUserBadmintonBookings(ctx, req.UserId)
        if err != nil {
            log.Printf("Error counting user's badminton bookings: %s", err)
            return nil, err
        }
        if count >= 2 {
            return nil, errors.New("error: user has reached the maximum limit of 2 badminton slots")
        }
    }

    // Check for duplicate booking
    exists, err := r.checkDuplicateBooking(ctx, req.UserId, slotIdObject, badmintonSlotIdObject)
    if err != nil {
        log.Printf("Error while checking duplicate booking: %s", err)
        return nil, err
    }
    if exists {
        log.Printf("User %s has already booked slot %v/%v", req.UserId, slotIdObject, badmintonSlotIdObject)
        return nil, errors.New("error: user has already booked this slot")
    }

    // Check slot availability
    var slot *facility.Slot
    if !isBadminton {
        slot, err = r.checkSlotAvailability(ctx, facilityName, req)
        if err != nil {
            return nil, err
        }
        if slot.CurrentBookings >= slot.MaxBookings {
            return nil, errors.New("error: Slot is full")
        }
    }

    // Create booking document with explicit fields
    bookingDoc := bson.D{
        {"user_id", req.UserId},
        {"slot_id", req.SlotId},
        {"slot_type", req.SlotType},
        {"facility_name", facilityName},  // Set facility name explicitly
        {"status", "pending"},
        {"payment_id", ""},
        {"qr_code_url", ""},
        {"created_at", time.Now().UTC()},
        {"updated_at", time.Now().UTC()},
    }

    if req.BadmintonSlotId != nil {
        bookingDoc = append(bookingDoc, bson.E{"badminton_slot_id", req.BadmintonSlotId})
    }

    // Insert booking
    result, err := col.InsertOne(ctx, bookingDoc)
    if err != nil {
        return nil, fmt.Errorf("error inserting booking: %w", err)
    }

    // Set the values back to the request object
    req.Id = result.InsertedID.(primitive.ObjectID)
    req.FacilityName = facilityName
    req.Status = "pending"
    req.PaymentID = ""
    req.QRCodeURL = ""
    req.CreatedAt = bookingDoc[7].Value.(time.Time)
    req.UpdatedAt = bookingDoc[8].Value.(time.Time)

    // Update slot booking count if needed
    if !isBadminton && slot != nil {
        if err := r.updateSlotCurrentBooking(ctx, facilityName, slot.Id, 1); err != nil {
            // Rollback if slot update fails
            if _, deleteErr := col.DeleteOne(ctx, bson.M{"_id": req.Id}); deleteErr != nil {
                log.Printf("Error rolling back booking: %v", deleteErr)
            }
            return nil, fmt.Errorf("error updating slot booking count: %w", err)
        }
    }

    // Create queue message
    queueMsg := &booking.BookingQueueMessage{
        UserId:          req.UserId,
        SlotId:          req.SlotId,
        BadmintonSlotId: req.BadmintonSlotId,
        SlotType:        req.SlotType,
        FacilityName:    facilityName,
        CreatedAt:       time.Now(),
    }

    // Send to Kafka asynchronously
    go func() {
        message, err := json.Marshal(queueMsg)
        if err != nil {
            log.Printf("Warning: Failed to marshal booking message: %v", err)
            return
        }

        for i := 0; i < 3; i++ {
            err = queue.PushMessageWithKeyToQueue(
                []string{r.cfg.Kafka.Url},
                r.cfg.Kafka.ApiKey,
                r.cfg.Kafka.Secret,
                "booking",
                req.UserId,
                message,
            )
            if err == nil {
                log.Printf("Successfully sent booking to Kafka after %d attempts", i+1)
                break
            }
            log.Printf("Warning: Failed to send booking to Kafka (attempt %d): %v", i+1, err)
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }()

    return req, nil
}

func (r *bookingRepository) checkDuplicateBooking(ctx context.Context, userId string, slotId, badmintonSlotId *primitive.ObjectID) (bool, error) {
    log.Printf("Checking duplicate booking for user %s", userId)

    // Base filter for active bookings
    filter := bson.M{
        "user_id": userId,
        "status": bson.M{"$in": []string{"pending", "paid"}},
    }

    // For normal slots
    if slotId != nil {
        filter["slot_id"] = slotId.Hex()
        
        count, err := r.bookingDbConn(ctx).Collection("booking_transaction").CountDocuments(ctx, filter)
        if err != nil {
            log.Printf("Error checking normal slot booking: %v", err)
            return false, err
        }
        
        if count > 0 {
            log.Printf("Found existing booking for slot %s", slotId.Hex())
            return true, errors.New("you already have a booking for this slot")
        }
    }

    // For badminton slots
    if badmintonSlotId != nil {
        // First check if this specific slot is already booked
        filter["badminton_slot_id"] = badmintonSlotId.Hex()
        
        count, err := r.bookingDbConn(ctx).Collection("booking_transaction").CountDocuments(ctx, filter)
        if err != nil {
            log.Printf("Error checking badminton slot booking: %v", err)
            return false, err
        }
        
        if count > 0 {
            log.Printf("Found existing booking for badminton slot %s", badmintonSlotId.Hex())
            return true, errors.New("you already have a booking for this badminton slot")
        }

        // Then check daily limit for badminton (2 slots per day)
        today := time.Now().UTC().Truncate(24 * time.Hour)
        tomorrow := today.Add(24 * time.Hour)

        dailyFilter := bson.M{
            "user_id": userId,
            "slot_type": "badminton",
            "status": bson.M{"$in": []string{"pending", "paid"}},
            "created_at": bson.M{
                "$gte": today,
                "$lt":  tomorrow,
            },
        }

        dailyCount, err := r.bookingDbConn(ctx).Collection("booking_transaction").CountDocuments(ctx, dailyFilter)
        if err != nil {
            log.Printf("Error checking daily badminton limit: %v", err)
            return false, err
        }

        if dailyCount >= 2 {
            log.Printf("User %s has reached daily badminton limit", userId)
            return true, errors.New("you have reached the maximum daily limit for badminton bookings (2 slots)")
        }
    }

    return false, nil
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

func (r *bookingRepository) UpdateStatusPaid(ctx context.Context, bookingID string, paymentID, qrCodeURL string) error {
    objID, err := primitive.ObjectIDFromHex(bookingID)
    if err != nil {
        return fmt.Errorf("invalid booking ID: %w", err)
    }

    filter := bson.M{"_id": objID}
    update := bson.M{
        "$set": bson.M{
            "status":      "PAID",
            "payment_id": paymentID,
            "qr_code_url": qrCodeURL,
            "updated_at":  time.Now().UTC(),
        },
    }

    _, err = r.bookingDbConn(ctx).Collection("booking_transaction").UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update booking status: %w", err)
    }

    return nil
}

// func (r *bookingRepository) checkDuplicateBooking(ctx context.Context, userId string, slotId, badmintonSlotId *primitive.ObjectID) (bool, error) {
//     filter := bson.M{
//         "user_id": userId,
//     }

//     if slotId != nil {
//         filter["slot_id"] = slotId
//     }
//     if badmintonSlotId != nil {
//         filter["badminton_slot_id"] = badmintonSlotId
//     }

//     count, err := r.bookingDbConn(ctx).Collection("booking_transaction").CountDocuments(ctx, filter)
//     if err != nil {
//         return false, err
//     }

//     return count > 0, nil
// }

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

func (r *bookingRepository) InsertBookingQueue(pctx context.Context, cfg *config.Config, facilityName string, req *booking.Booking) (*booking.Booking, error) {
    ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

    col := r.bookingDbConn(ctx).Collection("booking_transaction")

    queueMsg := &booking.BookingQueueMessage{
        UserId:          req.UserId,
        SlotId:          req.SlotId,
        BadmintonSlotId: req.BadmintonSlotId,
        SlotType:        req.SlotType,
        FacilityName:    facilityName,
        CreatedAt:       time.Now(),
    }

    if err := validateBookingRequest(req); err != nil {
        return nil, err
    }

    log.Printf("Attempting to insert booking for userId: %s, facilityName: %s", req.UserId, facilityName)

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

    exists, err := r.checkDuplicateBooking(ctx, req.UserId, slotIdObject, badmintonSlotIdObject)
    if err != nil {
        log.Printf("Error while checking duplicate booking: %s", err)
        return nil, err
    }
    if exists {
        log.Printf("User %s has already booked slot %v/%v", req.UserId, slotIdObject, badmintonSlotIdObject)
        return nil, errors.New("error: user has already booked this slot")
    }

    var slot *facility.Slot
    if !isBadminton {
        slot, err = r.checkSlotAvailability(ctx, facilityName, req)
        if err != nil {
            return nil, err
        }
        if slot.CurrentBookings >= slot.MaxBookings {
            log.Printf("Slot %v is full for facility %s", slot.Id, facilityName)
            return nil, errors.New("error: Slot is full")
        }
    }

    req.CreatedAt = time.Now().UTC()
    req.UpdatedAt = time.Now().UTC()

    // Send to Kafka queue with better error handling
    message, err := json.Marshal(queueMsg)
    if err != nil {
        log.Printf("Error marshalling message: %v", err)
        // Continue with booking creation even if Kafka fails
    } else {
        // Try to send to Kafka but don't fail the booking if it doesn't work
        kafkaErr := queue.PushMessageWithKeyToQueue(
            []string{cfg.Kafka.Url},
            cfg.Kafka.ApiKey,
            cfg.Kafka.Secret,
            "booking",
            req.UserId,
            message,
        )
        if kafkaErr != nil {
            log.Printf("Warning: Failed to send booking to Kafka: %v", kafkaErr)
            // Continue with booking creation
        }
    }

    // Create the booking regardless of Kafka status
    result, err := col.InsertOne(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("error inserting booking: %w", err)
    }

    req.Id = result.InsertedID.(primitive.ObjectID)
    return req, nil
}

func (r *bookingRepository) ProcessBookingQueue(pctx context.Context, cfg *config.Config, queueMsg *booking.BookingQueueMessage) error {
    ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

    // Create gRPC client for facility service
    grpcClient, err := grpc.NewGrpcClient(cfg.Grpc.FacilityUrl)
    if err != nil {
        return fmt.Errorf("error creating gRPC client: %w", err)
    }

    // Get facility client
    facilityClient := grpcClient.Facility()

    // Get the appropriate slot ID
    var slotId string
    if queueMsg.SlotId != nil {
        slotId = *queueMsg.SlotId
    } else if queueMsg.BadmintonSlotId != nil {
        slotId = *queueMsg.BadmintonSlotId
    }

    // Check slot availability using gRPC
    availabilityResp, err := facilityClient.CheckSlotAvailability(ctx, &facilityPb.CheckSlotRequest{
        SlotId:       slotId,
        FacilityName: queueMsg.FacilityName,
        SlotType:     queueMsg.SlotType,
    })
    if err != nil {
        return fmt.Errorf("error checking slot availability: %w", err)
    }

    // If slot is not available, return error
    if !availabilityResp.IsAvailable {
        return fmt.Errorf("slot %s is not available: %s", slotId, availabilityResp.ErrorMessage)
    }

    // Create booking document
    bookingDoc := bson.M{
        "user_id":     queueMsg.UserId,
        "facility":    queueMsg.FacilityName,
        "status":      "pending",
        "created_at":  queueMsg.CreatedAt,
        "updated_at":  time.Now(),
    }

    // Add the appropriate slot ID to the document
    if queueMsg.SlotId != nil {
        bookingDoc["slot_id"] = *queueMsg.SlotId
    } else if queueMsg.BadmintonSlotId != nil {
        bookingDoc["badminton_slot_id"] = *queueMsg.BadmintonSlotId
    }

    // Insert the booking
    col := r.bookingDbConn(ctx).Collection("booking_transaction")
    result, err := col.InsertOne(ctx, bookingDoc)
    if err != nil {
        return fmt.Errorf("error inserting booking: %w", err)
    }

    // Update slot booking count
    updateResp, err := facilityClient.UpdateSlotBookingCount(ctx, &facilityPb.UpdateSlotRequest{
        SlotId:       slotId,
        FacilityName: queueMsg.FacilityName,
        Increment:    1,
    })
    if err != nil {
        // If updating slot fails, try to rollback the booking
        if _, deleteErr := col.DeleteOne(ctx, bson.M{"_id": result.InsertedID}); deleteErr != nil {
            log.Printf("Failed to rollback booking after slot update failure: %v", deleteErr)
        }
        return fmt.Errorf("error updating slot booking count: %w", err)
    }

    if !updateResp.Success {
        // If update was not successful, rollback the booking
        if _, deleteErr := col.DeleteOne(ctx, bson.M{"_id": result.InsertedID}); deleteErr != nil {
            log.Printf("Failed to rollback booking after unsuccessful slot update: %v", deleteErr)
        }
        return fmt.Errorf("failed to update slot: %s", updateResp.ErrorMessage)
    }

    log.Printf("Successfully processed booking for user %s, slot %s", queueMsg.UserId, slotId)
    return nil
}

// Add cleanup function for expired pending bookings
func (r *bookingRepository) CleanupExpiredBookings() {
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for range ticker.C {
            ctx := context.Background()
            fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)

            // Find expired pending bookings
            filter := bson.M{
                "status":     "pending",
                "created_at": bson.M{"$lt": fifteenMinutesAgo},
            }

            cursor, err := r.bookingDbConn(ctx).Collection("booking_transaction").Find(ctx, filter)
            if err != nil {
                log.Printf("Error finding expired bookings: %v", err)
                continue
            }

            var expiredBookings []booking.Booking
            if err := cursor.All(ctx, &expiredBookings); err != nil {
                log.Printf("Error decoding expired bookings: %v", err)
                continue
            }

            for _, expiredBooking := range expiredBookings {
                // Move to failed bookings collection
                failedBooking := bson.M{
                    "booking_id":     expiredBooking.Id,
                    "user_id":        expiredBooking.UserId,
                    "slot_id":        expiredBooking.SlotId,
                    "facility_name":  expiredBooking.FacilityName,  // Use the stored facility name
                    "status":         "failed",
                    "failed_at":      time.Now(),
                    "reason":         "payment_timeout",
                    "created_at":     expiredBooking.CreatedAt,
                }

                _, err := r.bookingDbConn(ctx).Collection("booking_failed").InsertOne(ctx, failedBooking)
                if err != nil {
                    log.Printf("Error inserting failed booking: %v", err)
                    continue
                }

                // Delete from booking_transaction
                _, err = r.bookingDbConn(ctx).Collection("booking_transaction").DeleteOne(ctx, bson.M{"_id": expiredBooking.Id})
                if err != nil {
                    log.Printf("Error deleting expired booking: %v", err)
                    continue
                }

                // Update slot current booking count
                if expiredBooking.SlotId != nil {
                    slotId, err := primitive.ObjectIDFromHex(*expiredBooking.SlotId)
                    if err != nil {
                        log.Printf("Error converting slot ID: %v", err)
                        continue
                    }
                    if err := r.updateSlotCurrentBooking(ctx, expiredBooking.FacilityName, slotId, -1); err != nil {
                        log.Printf("Error updating slot booking count: %v", err)
                    }
                }

                // Publish event to Kafka
                message := map[string]interface{}{
                    "booking_id":    expiredBooking.Id,
                    "user_id":       expiredBooking.UserId,
                    "facility_name": expiredBooking.FacilityName,
                    "status":        "failed",
                    "reason":        "payment_timeout",
                }
                
                messageBytes, _ := json.Marshal(message)
                if err := queue.PushMessageWithKeyToQueue(
                    []string{r.cfg.Kafka.Url},
                    r.cfg.Kafka.ApiKey,
                    r.cfg.Kafka.Secret,
                    "booking.failed",
                    expiredBooking.UserId,
                    messageBytes,
                ); err != nil {
                    log.Printf("Error publishing failed booking event: %v", err)
                }
            }
        }
    }()
}
