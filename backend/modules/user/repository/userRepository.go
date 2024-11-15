package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/modules/auth"
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
		UpsetOffset(pctx context.Context, offset int64) error
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

func (r *UserRepository) UpsetOffset(pctx context.Context, offset int64) error {
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

func (r *UserRepository) InsertOneUser(pctx context.Context, req *user.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// Set role based on the request
	roleTitle := "user"
	roleCode := auth.RoleUser

	// If admin role is requested
	if req.UserRoles != nil && len(req.UserRoles) > 0 && req.UserRoles[0].RoleCode == auth.RoleAdmin {
		roleTitle = "admin"
		roleCode = auth.RoleAdmin
	}

	// Create user document with explicit role
	userDoc := &user.User{
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		UserRoles: []user.UserRole{
			{
				RoleTitle: roleTitle,
				RoleCode:  roleCode,
			},
		},
	}

	result, err := col.InsertOne(ctx, userDoc)
	if err != nil {
		log.Printf("Error: InsertOneUser: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one user failed")
	}

	userId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("error: insert one user failed")
	}

	// Verify the role was set correctly
	var insertedUser user.User
	err = col.FindOne(ctx, bson.M{"_id": userId}).Decode(&insertedUser)
	if err != nil {
		log.Printf("Error verifying user creation: %v", err)
	} else {
		log.Printf("User created with role code: %d", insertedUser.UserRoles[0].RoleCode)
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

func (r *UserRepository) FindOneUserCredential(pctx context.Context, email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// Find the user document
	var result user.User
	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(&result); err != nil {
		log.Printf("Error: FindOneUserCredential: %s", err.Error())
		return nil, errors.New("error: email is invalid")
	}

	// Find the highest role code
	highestRoleCode := 0
	for _, role := range result.UserRoles {
		if role.RoleCode > highestRoleCode {
			highestRoleCode = role.RoleCode
		}
	}

	// Update the user object with the highest role code
	result.UserRoles = []user.UserRole{{
		RoleTitle: "admin", // This will be overwritten if not admin
		RoleCode:  highestRoleCode,
	}}

	// Set the correct role title based on the role code
	if highestRoleCode == auth.RoleAdmin {
		result.UserRoles[0].RoleTitle = "admin"
	} else {
		result.UserRoles[0].RoleTitle = "user"
	}

	return &result, nil
}

func (r *UserRepository) FindOneUserProfile(pctx context.Context, userId string) (*user.UserProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// First, find the user document
	var userDoc user.User
	err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(userId)}).Decode(&userDoc)
	if err != nil {
		log.Printf("Error: FindOneUserProfile find user: %s", err.Error())
		return nil, errors.New("error: user profile not found")
	}

	// Find the highest role code
	highestRoleCode := 0
	for _, role := range userDoc.UserRoles {
		if role.RoleCode > highestRoleCode {
			highestRoleCode = role.RoleCode
		}
	}

	// Create the profile response
	result := &user.UserProfileBson{
		Id:        userDoc.Id,
		Email:     userDoc.Email,
		Name:      userDoc.Name,
		RoleCode:  highestRoleCode,  // Set the highest role code
		CreatedAt: userDoc.CreatedAt,
		UpdatedAt: userDoc.UpdatedAt,
	}

	return result, nil
}

func (r *UserRepository) FindOneUserProfileRefresh(pctx context.Context, userId string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	var result user.User
	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(userId)}).Decode(&result); err != nil {
		log.Printf("Error: FindOneUserProfileRefresh: %s", err.Error())
		return nil, errors.New("error: user profile not found")
	}

	// Find the highest role code
	highestRoleCode := 0
	for _, role := range result.UserRoles {
		if role.RoleCode > highestRoleCode {
			highestRoleCode = role.RoleCode
		}
	}

	// Create a new user object with the highest role code
	userWithHighestRole := &user.User{
		Id:        result.Id,
		Email:     result.Email,
		Password:  result.Password,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		UserRoles: []user.UserRole{{
			RoleTitle: "admin", // Will be overwritten if not admin
			RoleCode:  highestRoleCode,
		}},
	}

	// Set the correct role title based on the role code
	if highestRoleCode == auth.RoleAdmin {
		userWithHighestRole.UserRoles[0].RoleTitle = "admin"
	} else {
		userWithHighestRole.UserRoles[0].RoleTitle = "user"
	}

	return userWithHighestRole, nil
}

func (r *UserRepository) FindManyUser(pctx context.Context) ([]user.UserProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: FindManyUser: %s", err.Error())
		return nil, fmt.Errorf("error: failed to fetch users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []user.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Printf("Error: FindManyUser decode: %s", err.Error())
		return nil, fmt.Errorf("error: failed to decode users: %w", err)
	}

	// Convert User documents to UserProfileBson with correct role codes
	var profiles []user.UserProfileBson
	for _, u := range users {
		// Find highest role code
		highestRoleCode := 0
		for _, role := range u.UserRoles {
			if role.RoleCode > highestRoleCode {
				highestRoleCode = role.RoleCode
			}
		}

		profile := user.UserProfileBson{
			Id:        u.Id,
			Email:     u.Email,
			Name:      u.Name,
			RoleCode:  highestRoleCode,  // Set the highest role code
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}
