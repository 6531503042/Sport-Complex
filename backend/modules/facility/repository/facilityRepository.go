package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/modules/facility"
	"main/pkg/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	FacilityRepositoryService interface {


		InsertFacility (pctx context.Context, req * facility.Facilitiy) (primitive.ObjectID, error)
		IsUniqueName(pctx context.Context, facilityName string) bool
		UpdateOneFacility (pctx context.Context, facilityId, facilityName string, updateFields bson.M) error
		FindOneFacility(pctx context.Context, facilityId,facilityName string) (*facility.FacilityBson, error)
		FindManyFacility(ctx context.Context) ([]facility.FacilityBson, error)
		DeleteOneFacility(pctx context.Context, facilityId, facilityName string) error

		//Slot
		InsertSlot (pctx context.Context, facilityName string, slot facility.Slot) (*facility.Slot, error)
		FindOneSlot (ctx context.Context, facilityName,slotId string) (*facility.Slot, error)
		FindManySlot (ctx context.Context, facilityName string) ([]facility.Slot, error)
		UpdateSlot (ctx context.Context, facilityName string, req *facility.Slot) (*facility.Slot, error)
		EnableOrDisableSlot (ctx context.Context, facilityName, slotId string, status int) (*facility.Slot, error)
		DeleteSlot(ctx context.Context, facilityName, slotId string) error

		//Badminton
		InsertBadCourt(ctx context.Context, court *facility.BadmintonCourt) (primitive.ObjectID, error)
		FindBadmintonCourt (ctx context.Context) ([]facility.BadmintonCourt, error)
		InsertBadmintonSlot(ctx context.Context, req *facility.BadmintonSlot) (primitive.ObjectID, error)
		FindBadmintonSlot (ctx context.Context) ([]facility.BadmintonSlot, error)
		UpdateBadmintonSlot(ctx context.Context, req *facility.BadmintonSlot) error
		UpdateBadCourt(ctx context.Context, courtId string, updateFields bson.M) error
		DeleteBadmintonCourt(ctx context.Context, courtId string) error
		DeleteBadmintonSlot(ctx context.Context, slotId string) error
	}

	facilitiyReposiory struct {
		db *mongo.Client
		client *mongo.Client
	}
)

func NewFacilityRepository(client *mongo.Client) *facilitiyReposiory {
	return &facilitiyReposiory{client: client}
}

func (r *facilitiyReposiory) facilityDbConn(pctx context.Context, facilityName string) *mongo.Database {
	// Use the facility name to dynamically create the database name
	databaseName := fmt.Sprintf("%s_facility", facilityName)
	return r.client.Database(databaseName) // This will create the DB if it doesn't exist
}

func (r *facilitiyReposiory) slotDbConn(pctx context.Context,facilityName string) *mongo.Database {
	// Use the existing client to connect to the facility database
	databaseName := fmt.Sprintf("%s_facility", facilityName) // Consistent naming
	return r.client.Database(databaseName) // Connect to the existing database
}

func (r *facilitiyReposiory) courtDbConn(pctx context.Context) *mongo.Database {
	// Use the existing client to connect to the facility database
	databaseName := fmt.Sprintf("badminton_facility") // Consistent naming
	return r.client.Database(databaseName) // Connect to the existing database
}

func (r *facilitiyReposiory) ListAllFacilities(pctx context.Context) ([]string, error) {
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

func (r *facilitiyReposiory) InsertFacility (pctx context.Context, req * facility.Facilitiy) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, req.Name)
	col := db.Collection("facilities")
	
	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertFacility failed: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: InsertFacility failed")
	}

	facilityId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error: InsertFacility failed: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: InsertFacility failed")
	}

	return facilityId, nil
}

func (r *facilitiyReposiory) IsUniqueName(pctx context.Context, facilityName string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, facilityName) // Pass the facility name
	col := db.Collection("facilities")

	nameCount, err := col.CountDocuments(ctx, bson.M{"name": facilityName})
	if err != nil {
		log.Printf("Error: IsUniqueName failed: %s", err.Error())
		return false
	}

	return nameCount == 0
}




func (r *facilitiyReposiory) UpdateOneFacility(pctx context.Context, facilityId, facilityName string, updateFields bson.M) error {
    ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
    defer cancel()

    db := r.facilityDbConn(ctx, facilityName)
    col := db.Collection("facilities")

    updateResult, err := col.UpdateOne(
        ctx,
        bson.M{"_id": utils.ConvertToObjectId(facilityId)},
        bson.M{"$set": updateFields},
    )
    if err != nil {
        log.Printf("Error: UpdateOneFacility: %s", err.Error())
        return errors.New("error: update one facility failed")
    }

    if updateResult.MatchedCount == 0 {
        return errors.New("error: facility not found")
    }

    return nil
}


func (r *facilitiyReposiory) DeleteOneFacility(pctx context.Context, facilityId, facilityName string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx,facilityName)
	col := db.Collection("facilities")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(facilityId)})

	if err != nil {
		log.Printf("Error: DeleteOneFacility: %s", err.Error())
		return fmt.Errorf("error: delete one facility failed: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("error: facility %s not found", facilityId)
	}

	return nil
}


func (r *facilitiyReposiory) FindOneFacility(pctx context.Context, facilityId, facilityName string) (*facility.FacilityBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, facilityName)
	col := db.Collection("facilities")

	result := new(facility.FacilityBson)
	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(facilityId)},
		options.FindOne().SetProjection(
			bson.M{
				"_id": 1,
				"name": 1,
				"price_insider": 1,
				"price_outsider": 1,
				"description": 1,
				"created_at": 1,
				"updated_at": 1,
			},
		),
	).Decode(result); err != nil {
		log.Printf("Error: FindOneFacility: %s", err.Error())
		return nil, errors.New("error: find one facility failed")
	}

	return result, nil
}

func (r *facilitiyReposiory) FindManyFacility(ctx context.Context) ([]facility.FacilityBson, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // List all databases
    dbs, err := r.client.ListDatabaseNames(ctx, bson.M{})
    if err != nil {
        return nil, fmt.Errorf("error: failed to list databases: %w", err)
    }

    var allFacilities []facility.FacilityBson

    // Loop through all databases to find facilities
    for _, dbName := range dbs {
        if strings.HasSuffix(dbName, "_facility") {
            db := r.client.Database(dbName)
            col := db.Collection("facilities")

            var facilities []facility.FacilityBson
            cur, err := col.Find(ctx, bson.M{})
            if err != nil {
                log.Printf("Error: FindAllFacilities: %s", err.Error())
                continue
            }
            if err = cur.All(ctx, &facilities); err != nil {
                log.Printf("Error: FindAllFacilities: %s", err.Error())
                continue
            }

            allFacilities = append(allFacilities, facilities...)
        }
    }

    return allFacilities, nil
}

func (r *facilitiyReposiory) InsertSlot(pctx context.Context, facilityName string, slot facility.Slot) (*facility.Slot, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	// Use the slotDbConn to ensure we're connecting to an existing facility DB
	db := r.slotDbConn(ctx, facilityName)
	col := db.Collection("slots") // Get the "slots" collection

	slot.CreatedAt = time.Now()
	slot.UpdatedAt = time.Now()

	result, err := col.InsertOne(ctx, bson.M{
		"start_time":      slot.StartTime,
		"end_time":        slot.EndTime,
		"status":          slot.Status,
		"max_bookings":    slot.MaxBookings,
		"current_bookings": slot.CurrentBookings,
		"facility_type": slot.FacilityType,
		"created_at":      slot.CreatedAt,
		"updated_at":      slot.UpdatedAt,
	})
	if err != nil {
		log.Printf("Error: InsertSlot: %s", err.Error())
		return nil, fmt.Errorf("error: insert slot failed: %w", err)
	}

	slotId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("error: InsertSlot: %s", err.Error())
		return nil, fmt.Errorf("error: insert slot failed: %w", err)
	}

	slot.Id = slotId
	return &slot, nil
}

func (r *facilitiyReposiory) FindOneSlot (ctx context.Context, facilityName,slotId string) (*facility.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, facilityName)
	col := db.Collection("slots")

	id, err := primitive.ObjectIDFromHex(slotId)
	if err != nil {
		log.Printf("Error: FindOneSlot: %s", err.Error())
        return nil, fmt.Errorf("error: find one slot failed: %w", err)
	}

	var result bson.M
	err = col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
	    if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("error: slot %s doesn't exist", slotId)
		}
		log.Printf("Error: FindOneSlot: %s", err.Error())
        return nil, fmt.Errorf("error: find one slot failed: %w", err)
	}

	slot := &facility.Slot{
		Id: id,
		StartTime: result["start_time"].(string),
		EndTime: result["end_time"].(string),
		Status: result["status"].(int),
		MaxBookings: result["max_bookings"].(int),
		CurrentBookings: result["current_bookings"].(int),
		CreatedAt: result["created_at"].(time.Time),
		UpdatedAt: result["updated_at"].(time.Time),
	}
	return slot, nil
}

func (r *facilitiyReposiory) FindManySlot (ctx context.Context, facilityName string) ([]facility.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, facilityName)
	col := db.Collection("slots")

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: FindManySlot: %s", err.Error())
        return nil, fmt.Errorf("error: find many slot failed: %w", err)
	}
	defer cur.Close(ctx)

	var result []facility.Slot
	if err = cur.All(ctx, &result); err != nil {
		log.Printf("Error: FindManySlot: %s", err.Error())
        return nil, fmt.Errorf("error: find many slot failed: %w", err)
	}

	for _, slot := range result {
		log.Printf("Slot: start_time=%v, end_time=%v", slot.StartTime, slot.EndTime)
	}

	return result, nil
}

func (r *facilitiyReposiory) UpdateSlot(ctx context.Context, facilityName string, slot *facility.Slot) (*facility.Slot, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.slotDbConn(ctx, facilityName)
    col := db.Collection("slots")

    update := bson.M{
        "$set": bson.M{
            "current_bookings": slot.CurrentBookings,
            "updated_at":      slot.UpdatedAt,
        },
    }

    _, err := col.UpdateOne(ctx, bson.M{"_id": slot.Id}, update)
    if err != nil {
        log.Printf("Error updating slot: %v", err)
        return nil, fmt.Errorf("failed to update slot: %w", err)
    }

    return slot, nil
}


func (r *facilitiyReposiory) EnableOrDisableSlot (ctx context.Context, facilityName, slotId string, status int) (*facility.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, facilityName)
	col := db.Collection("slots")

	_, err := col.UpdateOne(ctx, bson.M{"_id": slotId}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Printf("Error: EnableOrDisableSlot: %s", err.Error())
        return nil, fmt.Errorf("error: update slot failed: %w", err)
	}
	return r.FindOneSlot(ctx, facilityName, slotId)
}

func (r *facilitiyReposiory) DeleteSlot(ctx context.Context, facilityName, slotId string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.slotDbConn(ctx, facilityName)
    col := db.Collection("slots")

    result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(slotId)})
    if err != nil {
        log.Printf("Error: DeleteSlot: %s", err.Error())
        return fmt.Errorf("error: delete slot failed: %w", err)
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("error: slot %s not found", slotId)
    }

    return nil
}


func (r *facilitiyReposiory) InsertBadCourt(ctx context.Context, court *facility.BadmintonCourt) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to the badminton facility database
	db := r.courtDbConn(ctx)
	col := db.Collection("court")

	result, err := col.InsertOne(ctx, bson.M{
		"court_number": court.CourtNumber,
		"status": court.Status,
	})
	if err != nil {
		log.Printf("Error: Insert Badminton Court: %s", err.Error())
		return primitive.NilObjectID, fmt.Errorf("error: insert badminton court failed: %w", err)
	}

	// Insert the court into the collection

	courtID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("error: insert badminton court failed")
	}

	return courtID, nil
}

func (r *facilitiyReposiory) UpdateBadCourt(ctx context.Context, courtId string, updateFields bson.M) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.courtDbConn(ctx)
    col := db.Collection("court")

    updateResult, err := col.UpdateOne(
        ctx,
        bson.M{"_id": utils.ConvertToObjectId(courtId)},
        bson.M{"$set": updateFields},
    )
    if err != nil {
        log.Printf("Error: UpdateBadCourt: %s", err.Error())
        return fmt.Errorf("error: update badminton court failed: %w", err)
    }

    if updateResult.MatchedCount == 0 {
        return errors.New("error: court not found")
    }

    return nil
}

func (r *facilitiyReposiory) DeleteBadmintonCourt(ctx context.Context, courtId string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.courtDbConn(ctx)
    col := db.Collection("court")

    result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(courtId)})
    if err != nil {
        log.Printf("Error: DeleteBadmintonCourt: %s", err.Error())
        return fmt.Errorf("error: delete badminton court failed: %w", err)
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("error: badminton court %s not found", courtId)
    }

    return nil
}



func (r *facilitiyReposiory) FindBadmintonCourt (ctx context.Context) ([]facility.BadmintonCourt, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.courtDbConn(ctx)
	col := db.Collection("court")

	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: Find Many Badminton Court: %s", err.Error())
        return nil, fmt.Errorf("error: find many court failed: %w", err)
	}
	defer cur.Close(ctx)

	var result []facility.BadmintonCourt
	if err = cur.All(ctx, &result); err != nil {
		log.Printf("Error: FindManySlot: %s", err.Error())
        return nil, fmt.Errorf("error: find many slot failed: %w", err)
	}

	return result, nil
}

func (r *facilitiyReposiory) InsertBadmintonSlot(ctx context.Context, req *facility.BadmintonSlot) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.courtDbConn(ctx)
	col := db.Collection("slots")

	// Set default values for new fields
	if req.MaxBookings == 0 {
		req.MaxBookings = 1 // Default max bookings for badminton is 1
	}
	req.CurrentBookings = 0 // Initialize current bookings to 0

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertBadmintonSlot: %s", err.Error())
		return primitive.NilObjectID, fmt.Errorf("error: insert badminton slot failed: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *facilitiyReposiory) UpdateBadmintonSlot(ctx context.Context, req *facility.BadmintonSlot) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.courtDbConn(ctx)
    col := db.Collection("slots")

    _, err := col.UpdateOne(ctx, bson.M{"_id": req.Id}, bson.M{"$set": req})
    if err != nil {
        log.Printf("Error: UpdateBadmintonSlot: %s", err.Error())
        return fmt.Errorf("error: update badminton slot failed: %w", err)
    }
    return nil
}


func (r *facilitiyReposiory) FindBadmintonSlot(ctx context.Context) ([]facility.BadmintonSlot, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.courtDbConn(ctx)
    col := db.Collection("slots")

    // Find all slots
    cursor, err := col.Find(ctx, bson.M{})
    if err != nil {
        log.Printf("Error: Find Badminton Slot: %s", err.Error())
        return nil, fmt.Errorf("error: find badminton slot failed: %w", err)
    }
    defer cursor.Close(ctx)

    var slots []facility.BadmintonSlot
    if err = cursor.All(ctx, &slots); err != nil {
        log.Printf("Error decoding slots: %s", err.Error())
        return nil, fmt.Errorf("error decoding slots: %w", err)
    }

    // Set default max_bookings to 1 if it's 0
    for i := range slots {
        if slots[i].MaxBookings == 0 {
            slots[i].MaxBookings = 1 // Set default max bookings to 1
        }
    }

    // Log for debugging
    for _, slot := range slots {
        log.Printf("Slot %s: max_bookings=%d, current_bookings=%d", 
            slot.Id.Hex(), 
            slot.MaxBookings,
            slot.CurrentBookings)
    }

    return slots, nil
}

func (r *facilitiyReposiory) DeleteBadmintonSlot(ctx context.Context, slotId string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.courtDbConn(ctx)
    col := db.Collection("slots")

    result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(slotId)})
    if err != nil {
        log.Printf("Error: DeleteBadmintonSlot: %s", err.Error())
        return fmt.Errorf("error: delete badminton slot failed: %w", err)
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("error: badminton slot %s not found", slotId)
    }

    return nil
}

// Add a new method to update badminton slot bookings
func (r *facilitiyReposiory) UpdateBadmintonSlotBookings(ctx context.Context, slotId primitive.ObjectID, increment int) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.courtDbConn(ctx)
    col := db.Collection("slots")

    update := bson.M{
        "$inc": bson.M{"current_bookings": increment},
    }

    result, err := col.UpdateOne(ctx, bson.M{"_id": slotId}, update)
    if err != nil {
        return fmt.Errorf("failed to update badminton slot bookings: %w", err)
    }

    if result.MatchedCount == 0 {
        return errors.New("badminton slot not found")
    }

    return nil
}

