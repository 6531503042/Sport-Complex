package middlewarerepository

import (
	"context"
	"errors"
	"log"
	authPb "main/modules/auth/proto"
	"main/pkg/grpc"
	"main/pkg/jwt"
	"time"
)

type (
	MiddlewareRepositoryService interface {
		AccessTokenSearch (pctx context.Context, grpcUrl, accessToken string) error
		RolesCount(pctx context.Context, grpcUrl string) (int64, error)
	}

	middlewareRepository struct {
		//Empty
	}
)

func NewMiddlewareRepository() MiddlewareRepositoryService {
	return &middlewareRepository{}
}

func (r *middlewareRepository) AccessTokenSearch (pctx context.Context, grpcUrl, accessToken string) error {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpc.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err.Error())
		return errors.New("error: gRPC connection failed")
	}

	jwt.SetApiKeyInContext(&ctx)
	result, err := conn.Auth().AccessTokenSearch(ctx, &authPb.AccessTokenSearchReq{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("Error: AccessTokenSearch failed: %s", err.Error())
		return errors.New("error: AccessTokenSearch failed")
	}

	if result == nil {
		log.Printf("Error: AccessTokenSearch result is nil")
		return errors.New("error: AccessTokenSearch result is nil")
	}

	if (!result.IsValid) {
		log.Printf("Error: AccessTokenSearch result is invalid")
		return errors.New("error: AccessTokenSearch result is invalid")
	}

	return nil
}

func (r *middlewareRepository) RolesCount(pctx context.Context, grpcUrl string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpc.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err.Error())
		return 0, errors.New("error: gRPC connection failed")
	}
	
	jwt.SetApiKeyInContext(&ctx)
	result, err := conn.Auth().RolesCount(ctx, &authPb.RolesCountReq{})
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %s", err.Error())
		return -1, errors.New("error: email or password is incorrect")
	}

	if result == nil {
		log.Printf("Error: roles count failed")
		return -1, errors.New("error: roles count failed")
	}

	return result.Count, nil
}