package repository

import (
	"context"
	"errors"
	"log"
	"main/modules/auth"
	"main/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthRepositoryService interface {
		InsertOneUserCredential (pctx context.Context, req *auth.Credential) (primitive.ObjectID, error)
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{db: db}
}

func (r * authRepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}

func (r * authRepository) InsertOneUserCredential (pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	req.CreatedAt = utils.LocalTime()
	req.UpdatedAt = utils.LocalTime()

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUserCredential: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one user credential failed")
	}

	userId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("error: insert one user credential failed")
	}
	return userId, nil
}