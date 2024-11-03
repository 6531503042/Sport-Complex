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
