package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/modules/models"
	"main/modules/user"
	"main/pkg/utils"
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
		FindOneUserProfile (pctx context.Context, userId string) (*user.UserProfileBson, error)
		FindOneUserProfileRefresh (pctx context.Context, userId string) (*user.User, error)
		UpdateOneUser (pctx context.Context, userId string, updateFields bson.M) error
		DeleteOneUser (pctx context.Context, userId string) error
		FindManyUser (pctx context.Context) ([]user.UserProfileBson, error)

		//Kafka
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
	}

	UserRepository struct {
		db *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) UserRepositoryService {
	return &UserRepository{db: db}
}

func (r *UserRepository) userDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("user_db")
}

func (r *UserRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}

func (r *UserRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("user_transactions_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpserOffset failed: %s", err.Error())
		return errors.New("error: UpserOffset failed")
	}
	log.Printf("Info: UpserOffset result: %v", result)

	return nil
}

func (r *UserRepository) InsertOneUser (pctx context.Context, req * user.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneUser: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one user failed")
	}

	userId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("error: insert one user failed")
	}
	return userId, nil
}

func (r *UserRepository) UpdateOneUser (pctx context.Context, userId string, updateFields bson.M) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	updateResult, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(userId)},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		log.Printf("Error: UpdateOneUser: %s", err.Error())
		return errors.New("error: update one user failed")
	}

	if updateResult.MatchedCount == 0 {
		return errors.New("error: user not found")
	}

	if updateResult.ModifiedCount == 0 {
		return errors.New("error: nothing to update")
	}

	return nil
}

func (r *UserRepository) DeleteOneUser (pctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(userId)})

	if err != nil {
		log.Printf("Error: DeleteOneUser: %s", err.Error())
		return fmt.Errorf("error: delete one user failed: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("error: user %s not found", userId)
	}

	return nil
}

func (r * UserRepository) IsUniqueUser (pctx context.Context, email, name string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// Check for user by email
    emailCount, err := col.CountDocuments(ctx, bson.M{"email": email})
    if err != nil {
        log.Printf("Error: IsUniqueUser - Failed to count documents by email: %s", err.Error())
        return false
    }

    // Check for user by name
    nameCount, err := col.CountDocuments(ctx, bson.M{"name": name})
    if err != nil {
        log.Printf("Error: IsUniqueUser - Failed to count documents by name: %s", err.Error())
        return false
    }

    // If either count is greater than 0, the user is not unique
    return emailCount == 0 && nameCount == 0
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

func (r *UserRepository) FindOneUserProfile(pctx context.Context, userId string) (*user.UserProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.UserProfileBson)

	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(userId)},
		options.FindOne().SetProjection(
			bson.M{
				"_id":        1,
				"email":      1,
				"name":   1,
				"created_at": 1,
				"updated_at": 1,
			},
		),
	).Decode(result); err != nil {
		log.Printf("Error: FindOneUserProfile: %s", err.Error())
		return nil, errors.New("error: user profile not found")
	}

	return result, nil
}

func (r *UserRepository) FindOneUserProfileRefresh(pctx context.Context, userId string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	result := new(user.User)

	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(userId)},
	).Decode(result); err != nil {
		log.Printf("Error: FindOneUserProfile: %s", err.Error())
		return nil, errors.New("error: user profile not found")
	}

	return result, nil
}

func (r *UserRepository) FindManyUser (pctx context.Context) ([]user.UserProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	cursor, err := col.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"_id":        1,
		"email":      1,
		"name":   1,
		"created_at": 1,
		"updated_at": 1,

	}))
	if err != nil {
        log.Printf("Error: FindManyUser: %s", err.Error())
        return nil, fmt.Errorf("error: failed to fetch users: %w", err)
    }
    defer func() {
        if err := cursor.Close(ctx); err != nil {
            log.Printf("Error: FindManyUser - failed to close cursor: %s", err.Error())
        }
    }()

	var users []user.UserProfileBson
    if err = cursor.All(ctx, &users); err != nil {
        log.Printf("Error: FindManyUser: %s", err.Error())
        return nil, fmt.Errorf("error: failed to decode users: %w", err)
    }

    return users, nil
}
