package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/modules/facility"
	"main/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	FacilityRepositoryService interface {
		InsertFacility (pctx context.Context, req * facility.Facilitiy) (primitive.ObjectID, error)
		IsUniqueName (pctx context.Context, name string) bool
		UpdateOneFacility (pctx context.Context, facilityId string, updateFields bson.M) error
		DeleteOneFacility(pctx context.Context, facilityId string) error
		FindOneFacility(pctx context.Context, facilityId string) (*facility.FacilityBson, error)
		FindManyFacility(pctx context.Context) ([]facility.FacilityBson, error)
	}

	facilitiyReposiory struct {
		db *mongo.Client
	}
)

func NewFacilityRepository(db *mongo.Client) FacilityRepositoryService {
	return &facilitiyReposiory{db: db}
}

func (r *facilitiyReposiory) facilityDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("facility_db")
}

func (r *facilitiyReposiory) InsertFacility (pctx context.Context, req * facility.Facilitiy) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx)
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

func (r *facilitiyReposiory) IsUniqueName (pctx context.Context, name string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx)
	col := db.Collection("facilities")

	nameCount, err := col.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		log.Printf("Error: IsUniqueName failed: %s", err.Error())
		return false
	}

	return nameCount == 0
}



func (r *facilitiyReposiory) UpdateOneFacility (pctx context.Context, facilityId string, updateFields bson.M) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx)
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

func (r *facilitiyReposiory) DeleteOneFacility(pctx context.Context, facilityId string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx)
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


func (r *facilitiyReposiory) FindOneFacility(pctx context.Context, facilityId string) (*facility.FacilityBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx)
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

func (r *facilitiyReposiory) FindManyFacility(pctx context.Context) ([]facility.FacilityBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.facilityDbConn(ctx)
	col := db.Collection("facilities")

	cursor, err := col.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"_id": 1,
		"name": 1,
		"price_insider": 1,
		"price_outsider": 1,
		"description": 1,
		"created_at": 1,
		"updated_at": 1,
	}))
	if err != nil {
		log.Printf("Error: FindManyFacility: %s", err.Error())
		return nil, fmt.Errorf("error: failed to fetch facilities: %w", err)
	}
	defer func ()  {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error: FindManyFacility: %s", err.Error())
		}
	}()

	var facilities []facility.FacilityBson
	if err := cursor.All(ctx, &facilities); err != nil {
		log.Printf("Error: FindManyFacility: %s", err.Error())
		return nil, fmt.Errorf("error: failed to fetch facilities: %w", err)
	}

	return facilities, nil
}