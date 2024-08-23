package repository

import (
	"context"
	"errors"
	"log"
	"main/pkg/utils"
	"os/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (

	UserRepositoryService interface {
		InsertOneUser (pctx context.Context, req * user.User) (primitive.ObjectID, error)
		IsUniqueUser (pctx context.Context, email, name string) bool
		FindOneUserCredential (pctx context.Context, email string) (*user.User, error)
		
	}

	UserRepository struct {
		db *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) UserRepositoryService {
	return &UserRepository{db: db}
}

func (r *UserRepository) userDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("user")
}

func (r *UserRepository) InsertOneUser (pctx context.Context, req * user.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	userId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUser: %s", err.Error())
		return primitive.NilObjectID, errors.New("Error: Insert One User Failed")
	}

	return userId.InsertedID.(primitive.ObjectID), nil
}

func (r * UserRepository) IsUniqueUser (pctx context.Context, email, name string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	user := new(user.User)
	if err := col.FindOne(
		ctx,
		bson.M{"email": email, "name": name},
	).Decode(user); err != nil {
		log.Printf("Error: IsUniqueUser: %s", err.Error())
		return false
	}

	return false
}

func (r *UserRepository) FindOneUserCredential (pctx context.Context, email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.User)

	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(result); err != nil {
		log.Printf("Error: FindOneUserCredential: %s", err.Error())
		return nil, errors.New("error: email is invalid")
	}
	return result, nil
}

func (r *UserRepository) FindOneUserProfile(pctx context.Context, userId string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.User)

	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(userId)},
		options.FindOne().SetProjection(
			bson.M{
				"_id":		1,
				"email":	1,
				"name":		1,
				"created_at": 	1,
				"updated_at": 	1,
			},
		),
	).Decode(result); err != nil {
		log.Printf("Error: FindOneUserProfile: %s", err.Error())
		return nil, errors.New("error: user not found")
	}
	return result, nil
}