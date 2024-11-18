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
		UpdateOneUser (pctx context.Context, userId string, updateFields map[string]interface{}) error
		DeleteOneUser (pctx context.Context, userId string) error
		FindManyUser (pctx context.Context) ([]user.User, error)

		//Kafka
		GetOffset(pctx context.Context) (int64, error)
		UpsetOffset(pctx context.Context, offset int64) error

		GetUserAnalytics(pctx context.Context, period string, startDate, endDate time.Time) (*user.UserAnalytics, error)
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

func (r *UserRepository) UpdateOneUser(ctx context.Context, userId string, updateFields map[string]interface{}) error {
	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// Convert updateFields to bson.M
	bsonUpdate := bson.M{}
	for k, v := range updateFields {
		bsonUpdate[k] = v
	}

	update := bson.M{"$set": bsonUpdate}

	result, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(userId)},
		update,
	)

	if err != nil {
		log.Printf("Error updating user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
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

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": utils.ConvertToObjectId(userId),
			},
		},
		{
			"$project": bson.M{
				"_id":        1,
				"email":      1,
				"name":       1,
				"created_at": 1,
				"updated_at": 1,
				"user_roles": 1,
				"role_code": bson.M{
					"$ifNull": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$user_roles.role_code", 0}},
						0,
					},
				},
				"role_title": bson.M{
					"$ifNull": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$user_roles.role_title", 0}},
						"user",
					},
				},
			},
		},
	}

	var result user.UserProfileBson
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error: FindOneUserProfile: %s", err.Error())
		return nil, fmt.Errorf("error: failed to fetch user profile: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Error: FindOneUserProfile decode: %s", err.Error())
			return nil, fmt.Errorf("error: failed to decode user profile: %w", err)
		}
		return &result, nil
	}

	return nil, fmt.Errorf("error: user not found")
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

func (r *UserRepository) FindManyUser(pctx context.Context) ([]user.User, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// Create a pipeline to get all fields including user_roles
	pipeline := []bson.M{
		{
			"$project": bson.M{
				"_id":        1,
				"email":      1,
				"name":       1,
				"created_at": 1,
				"updated_at": 1,
				"user_roles": 1,
			},
		},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
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

	// Debug log
	for _, u := range users {
		log.Printf("User %s has roles: %+v", u.Email, u.UserRoles)
	}

	return users, nil
}

func (r *UserRepository) GetUserAnalytics(pctx context.Context, period string, startDate, endDate time.Time) (*user.UserAnalytics, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.userDbConn(ctx)
	col := db.Collection("users")

	// Get total users
	totalUsers, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Get active users (users who have logged in within last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	activeUsers, err := col.CountDocuments(ctx, bson.M{
		"updated_at": bson.M{"$gte": thirtyDaysAgo},
	})
	if err != nil {
		return nil, err
	}

	// Get users by period
	pipeline := r.buildAnalyticsPipeline(period, startDate, endDate)
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var usersByPeriod []user.UserActivityPeriod
	if err := cursor.All(ctx, &usersByPeriod); err != nil {
		return nil, err
	}

	// Find peak activity
	peakActivity := user.PeakActivity{Count: 0}
	for _, p := range usersByPeriod {
		if p.Count > peakActivity.Count {
			peakActivity.Count = p.Count
			peakActivity.Date = p.Period
		}
	}

	return &user.UserAnalytics{
		TotalUsers:      totalUsers,
		ActiveUsers:     activeUsers,
		UsersByPeriod:   usersByPeriod,
		PeakActivityDay: peakActivity,
	}, nil
}

func (r *UserRepository) buildAnalyticsPipeline(period string, startDate, endDate time.Time) []bson.M {
	var dateFormat string
	var groupId bson.M

	switch period {
	case "daily":
		dateFormat = "%Y-%m-%d"
		groupId = bson.M{
			"$dateToString": bson.M{
				"format": dateFormat,
				"date":   "$created_at",
			},
		}
	case "weekly":
		dateFormat = "%Y-W%V"
		groupId = bson.M{
			"$dateToString": bson.M{
				"format": dateFormat,
				"date":   "$created_at",
			},
		}
	case "monthly":
		dateFormat = "%Y-%m"
		groupId = bson.M{
			"$dateToString": bson.M{
				"format": dateFormat,
				"date":   "$created_at",
			},
		}
	case "yearly":
		dateFormat = "%Y"
		groupId = bson.M{
			"$dateToString": bson.M{
				"format": dateFormat,
				"date":   "$created_at",
			},
		}
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	return []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":    groupId,
				"count":  bson.M{"$sum": 1},
				"active": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$gte": []interface{}{"$updated_at", thirtyDaysAgo}},
							1,
							0,
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":      0,
				"period":   "$_id",
				"count":    1,
				"active":   1,
				"inactive": bson.M{"$subtract": []interface{}{"$count", "$active"}},
			},
		},
		{
			"$sort": bson.M{"period": 1},
		},
	}
}
