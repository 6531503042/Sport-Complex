package handlers

import (
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