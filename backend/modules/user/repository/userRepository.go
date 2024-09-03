package repository

import (
	"context"
	"errors"
	"log"
	models "main/modules/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IUserRepository interface {
		CreateUser(ctx context.Context, user models.User) (*models.User, error)
		GetUserByID(ctx context.Context, id string) (*models.User, error)
		GetAllUsers(ctx context.Context) ([]models.User, error)
		UpdateUser(ctx context.Context, id string, user models.User) (*models.User, error)
		DeleteUser(ctx context.Context, id string) error
	}

	userRepository struct {
		db *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) userDbConnect(pctx context.Context) *mongo.Database {
	return r.db.Database("user")
}

func (r *userRepository) CreateUser (pctx context.Context, user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConnect(ctx)
	col := db.Collection("users")

	result, err := col.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error: InsertOneUser: %s", err.Error())
		return nil, errors.New("error: insert one user failed")
	}

	user.Id = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (r *userRepository) GetUserByID(pctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConnect(ctx)
	col := db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.User
	err = col.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (r *userRepository) GetAllUsers(pctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConnect(ctx)
	col := db.Collection("users")

	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates a user in the database
func (r *userRepository) UpdateUser(pctx context.Context, id string, user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConnect(ctx)
	col := db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	update := bson.M{"$set": user}
	_, err = col.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes a user by ID from the database
func (r *userRepository) DeleteUser(pctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConnect(ctx)
	col := db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	_, err = col.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}