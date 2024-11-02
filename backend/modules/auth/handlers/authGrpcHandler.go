package handlers

import (
	"context"
	"main/modules/auth/usecase"

	authPb "main/modules/auth/proto"
)

type (
	authGrpcHandler struct {
		authPb.UnimplementedAuthGrpcServiceServer
		authUsecase usecase.AuthUsecaseService
	}
)

func NewAuthGrpcpHandler (usecase usecase.AuthUsecaseService) *authGrpcHandler {
	return &authGrpcHandler{authUsecase: usecase}
}

func (g *authGrpcHandler) AccessTokenSearch(ctx context.Context, req *authPb.AccessTokenSearchReq) (*authPb.AccessTokenSearchRes, error) {
	return g.authUsecase.AccessTokenSearch(ctx, req.AccessToken)
}

func (g *authGrpcHandler) RolesCount(ctx context.Context, req *authPb.RolesCountReq) (*authPb.RolesCountRes, error) {
	return g.authUsecase.RolesCount(ctx)
}

// func (g *authGrpcHandler) FindOneUserProfileToRefresh(ctx context.Context, req *authPb.FindOneUserProfileToRefreshReq) (*authPb.UserProfile, error) {
//     // Use your usecase to fetch the user profile based on the userId from the request
//     userProfile, err := g.authUsecase.GetUserProfileByID(req.UserId)
//     if err != nil {
//         return nil, status.Errorf(codes.NotFound, "user profile not found: %v", err)
//     }
//     return &authPb.UserProfile{
//         UserId:   userProfile.ID,
//         Username: userProfile.Username, // Adjust according to your User model
//     }, nil
// }