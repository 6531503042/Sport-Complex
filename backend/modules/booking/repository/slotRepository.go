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
	SlotRepositoryService interface {
		InsertSlot(ctx context.Context, slot *booking.Slot) (*booking.Slot, error)
		FindOneSlot(ctx context.Context, slotId string) (*booking.Slot, error)
		FindAllSlots(ctx context.Context) ([]booking.Slot, error)
		EnableOrDisableSlot(ctx context.Context, slotId string, status int) (*booking.Slot, error)
		UpdateSlot(ctx context.Context, slot *booking.Slot) (*booking.Slot, error)
	}

	slotRepository struct {
		db *mongo.Client
	}
)

func NewSlotRepository(db *mongo.Client) SlotRepositoryService {
	return &slotRepository{db}
}


func (r *slotRepository) slotDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("booking")
}

func (r *slotRepository) UpdateSlot(ctx context.Context, slot *booking.Slot) (*booking.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	col := r.slotDbConn(ctx).Collection("slots")
	slot.UpdatedAt = time.Now()

	_, err := col.UpdateOne(ctx, bson.M{"_id": slot.Id}, bson.M{"$set": slot})
	if err != nil {
		log.Printf("Error: UpdateSlot: %s", err.Error())
		return nil, fmt.Errorf("error: update slot failed: %w", err)
	}
	return slot, nil
}


func (r *slotRepository) InsertSlot(ctx context.Context, slot *booking.Slot) (*booking.Slot, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.slotDbConn(ctx)
    col := db.Collection("slots")

    // Set the creation and update time
    slot.CreatedAt = time.Now()
    slot.UpdatedAt = time.Now()

    // Insert the slot into the database
    result, err := col.InsertOne(ctx, bson.M{
        "start_time": slot.StartTime, // Store as string "HH:mm"
        "end_time":   slot.EndTime,   // Store as string "HH:mm"
        "status":     slot.Status,
        "created_at": slot.CreatedAt,
        "updated_at": slot.UpdatedAt,
    })
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




func (r *slotRepository) FindOneSlot(ctx context.Context, slotId string) (*booking.Slot, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.slotDbConn(ctx)
    col := db.Collection("slots")

    id, err := primitive.ObjectIDFromHex(slotId)
    if err != nil {
        return nil, fmt.Errorf("error: failed to parse slot ID: %w", err)
    }

    var result bson.M
    err = col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("error: slot %s does not exist", slotId)
        }
        log.Printf("Error: FindOneSlot: %s", err.Error())
        return nil, fmt.Errorf("error: failed to find slot: %w", err)
    }

    slot := &booking.Slot{
        Id:        id,
        StartTime: result["start_time"].(string),  // Use string for "HH:mm"
        EndTime:   result["end_time"].(string),    // Use string for "HH:mm"
        Status:    result["status"].(int),
        CreatedAt: result["created_at"].(time.Time),
        UpdatedAt: result["updated_at"].(time.Time),
    }

    return slot, nil
}



func (r *slotRepository) FindAllSlots(ctx context.Context) ([]booking.Slot, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.slotDbConn(ctx)
    col := db.Collection("slots")

    cursor, err := col.Find(ctx, bson.M{})
    if err != nil {
        log.Printf("Error: FindAllSlots: %s", err.Error())
        return nil, errors.New("error: find all slots failed")
    }
    defer cursor.Close(ctx)

    var result []booking.Slot
    if err = cursor.All(ctx, &result); err != nil {
        log.Printf("Error: FindAllSlots: %s", err.Error())
        return nil, errors.New("error: find all slots failed")
    }

    // Ensure all time values are properly converted
    for _, slot := range result {
		log.Printf("Slot: start_time=%v, end_time=%v", slot.StartTime, slot.EndTime)
	}

    return result, nil
}



func (r *slotRepository) EnableOrDisableSlot(ctx context.Context, slotId string, status int) (*booking.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	db := r.slotDbConn(ctx)
	col := db.Collection("slots")
	_, err := col.UpdateOne(ctx, bson.M{"_id": slotId}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Printf("Error: EnableOrDisableSlot: %s", err.Error())
		return nil, fmt.Errorf("error: enable/disable slot failed: %w", err)
	}
	return r.FindOneSlot(ctx, slotId) // Return the updated slot
}