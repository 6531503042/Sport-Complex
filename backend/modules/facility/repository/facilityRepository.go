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
		DeleteOneFacility(pctx context.Context, facilityId, facilityName string) error
		FindOneFacility(pctx context.Context, facilityId,facilityName string) (*facility.FacilityBson, error)
		FindManyFacility(ctx context.Context) ([]facility.FacilityBson, error)

		//Slot
		InsertSlot (pctx context.Context, facilityName string, slot facility.Slot) (*facility.Slot, error)
		FindOneSlot (ctx context.Context, facilityName,slotId string) (*facility.Slot, error)
		FindManySlot (ctx context.Context, facilityName string) ([]facility.Slot, error)
		UpdateSlot (ctx context.Context, facilityName string, req *facility.Slot) (*facility.Slot, error)
		EnableOrDisableSlot (ctx context.Context, facilityName, slotId string, status int) (*facility.Slot, error)

		ListAllFacilities(pctx context.Context) ([]string, error)
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




func (r *facilitiyReposiory) UpdateOneFacility (pctx context.Context, facilityId, facilityName string, updateFields bson.M) error {
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

	if updateResult.ModifiedCount == 0 {
		return errors.New("error: nothing to update")
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

func (r *facilitiyReposiory) UpdateSlot (ctx context.Context, facilityName string, req *facility.Slot) (*facility.Slot, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx, facilityName)
	col := db.Collection("slots")

	_, err := col.UpdateOne(ctx, bson.M{"_id": req.Id}, bson.M{"$set": req})
	if err != nil {
		log.Printf("Error: UpdateSlot: %s", err.Error())
        return nil, fmt.Errorf("error: update slot failed: %w", err)
	}
	return req, nil
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