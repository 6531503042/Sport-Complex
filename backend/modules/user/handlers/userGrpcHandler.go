package handlers

import (
	"context"
	userPb "main/modules/user/proto"
	"main/modules/user/usecase"
)

type (
    userGrpcHandler struct {
        userUsecase usecase.UserUsecaseService
        userPb.UnimplementedUserGrpcServiceServer
    }
)

func NewUserGrpcHandler(userUsecase usecase.UserUsecaseService) *userGrpcHandler {
    return &userGrpcHandler{userUsecase: userUsecase}
}

func (g *userGrpcHandler) CredentialSearch(ctx context.Context, req *userPb.CredentialSearchReq) (*userPb.UserProfile, error) {
    return g.userUsecase.FindOneUserCredential(ctx, req.Password, req.Email)
}

func (g *userGrpcHandler) FindOneUserProfileRefresh(ctx context.Context, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
    return g.userUsecase.FindOneUserProfileToRefresh(ctx, req.UserId)
}
