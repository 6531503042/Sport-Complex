package usecase

import (
	"context"
	"main/config"
	"main/modules/auth"
)

type AuthUsecaseService interface {

}

type authUsecase struct {
	authRepository authRepository.AuthRepositoryService
}

func NewAuthUsecase(authRepository authRepository.AuthRepositoryService) AuthUsecaseService {
	return &AuthUsecase{authRepository}
}

func (u * authUsecase) Login(pctx context.Context, cfg *config.Config, req*auth,UserLoginReq) (*auth.ProfileIntercepter, error) {
	profile, err := u.authRepository.Credential
}