package repository

import (
	"context"
	"errors"
	"log"
	"main/config"
	"main/modules/auth"
	"main/pkg/grpc"
	"main/pkg/jwt"
	"main/pkg/utils"
	"time"

	userPb "main/modules/user/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRepositoryService defines the interface for the Auth repository.
type (
	AuthRepositoryService interface {
		InsertOneUserCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error)
		CredentialSearch(pctx context.Context, grpcUrl string, req *userPb.CredentialSearchReq) (*userPb.UserProfile, error)
		FindOneUserProfileToRefresh(pctx context.Context, grpcUrl string, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error)
		FindOneUserCredential(pctx context.Context, credentialId string) (*auth.Credential, error)
		UpdateOneUserCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error
		DeleteOneUserCredential(pctx context.Context, credentialId string) (int64, error)
		FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error)
		RolesCount(pctx context.Context) (int64, error)
		AccessToken(cfg *config.Config, claims *jwt.Claims) string
		RefreshToken(cfg *config.Config, claims *jwt.Claims) string
	}

	authRepository struct {
		db *mongo.Client
	}
)

// NewAuthRepository creates a new instance of authRepository.
func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{db}
}

// authDbConn establishes a connection to the auth database.
func (r *authRepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}


// InsertOneUserCredential implements AuthRepositoryService.
func (r *authRepository) InsertOneUserCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	req.CreatedAt = utils.LocalTime()
	req.UpdatedAt = utils.LocalTime()

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error inserting user credential: %v", err) // Log the error
		return primitive.NilObjectID, errors.New("unable to create user credential") // User-friendly error message
	}

	return result.InsertedID.(primitive.ObjectID), nil
}




// InsertOneUserCredential inserts a new user credential into the database.
func (r *authRepository) CredentialSearch(pctx context.Context, grpcUrl string, req *userPb.CredentialSearchReq) (*userPb.UserProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwt.SetApiKeyInContext(&ctx)
	conn, err := grpc.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err.Error())
		return nil, errors.New("error: gRPC connection failed")
	}

	result, err := conn.User().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %s", err.Error())
		return nil, errors.New("error: email or password is incorrect")
	}

	return result, nil
}

// FindOneUserProfileToRefresh finds a user profile via gRPC using provided credentials to refresh the token.
func (r *authRepository) FindOneUserProfileToRefresh(pctx context.Context, grpcUrl string, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
    ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
    defer cancel()

    jwt.SetApiKeyInContext(&ctx)
    conn, err := grpc.NewGrpcClient(grpcUrl)
    if err != nil {
        log.Printf("Error: gRPC connection failed: %s", err.Error())
        return nil, errors.New("error: gRPC connection failed")
    }

    result, err := conn.User().FindOneUserProfileToRefresh(ctx, req)
    if err != nil {
        log.Printf("Error: FindOneUserProfileToRefresh failed: %s", err.Error())
        return nil, errors.New("error: user profile not found")
    }

    return result, nil
}

// FindOneUserCredential retrieves a user credential by ID from the database.
func (r *authRepository) FindOneUserCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result := new(auth.Credential)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)}).Decode(result); err != nil {
		log.Printf("Error: FindOneUserCredential failed: %s", err.Error())
		return nil, errors.New("error: find one user credential failed")
	}

	return result, nil
}

// UpdateOneUserCredential updates a user credential in the database.
func (r *authRepository) UpdateOneUserCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(credentialId)},
		bson.M{
			"$set": bson.M{
				"user_id":       req.UserId,
				"access_token":  req.AccessToken,
				"refresh_token": req.RefreshToken,
				"updated_at":    req.UpdatedAt,
			},
		},
	)
	if err != nil {
		log.Printf("Error: UpdateOneUserCredential: %s", err.Error())
		return errors.New("error: update one user credential failed")
	}

	return nil
}

// DeleteOneUserCredential deletes a user credential by ID from the database.
func (r *authRepository) DeleteOneUserCredential(pctx context.Context, credentialId string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)})
	if err != nil {
		log.Printf("Error: DeleteOneUserCredential: %s", err.Error())
		return 0, errors.New("error: delete one user credential failed")
	}
	log.Printf("DeleteOneUserCredential result: %v", result)

	return result.DeletedCount, nil
}

// FindOneAccessToken retrieves a user credential by access token from the database.
func (r *authRepository) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	credential := new(auth.Credential)
	if err := col.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(credential); err != nil {
		log.Printf("Error: FindOneAccessToken failed: %s", err.Error())
		return nil, errors.New("error: access token not found")
	}

	return credential, nil
}

func (r *authRepository) RolesCount(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("roles")

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: RolesCount failed: %s", err.Error())
		return -1, errors.New("error: roles count failed")
	}

	return count, nil
}

// AccessToken generates a new access token.
func (r *authRepository) AccessToken(cfg *config.Config, claims *jwt.Claims) string {
	return jwt.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwt.Claims{
		UserId:   claims.UserId,
		RoleCode: int(claims.RoleCode),
	}).SignToken()
}

// RefreshToken generates a new refresh token.
func (r * authRepository) RefreshToken(cfg *config.Config, claims *jwt.Claims) string {
	return jwt.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwt.Claims{
		UserId:   claims.UserId,
		RoleCode: int(claims.RoleCode),
	}).SignToken()
}