package repository

import (
	"context"
	"errors"
	"log"
	"main/modules/facility"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	FacilityRepositoryService interface {
		InsertFacility (pctx context.Context, req * facility.Facilitiy) (primitive.ObjectID, error)
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

func (r *facilitiyReposiory) UpdateOneFacility (pctx context.Context, facilityId string, )