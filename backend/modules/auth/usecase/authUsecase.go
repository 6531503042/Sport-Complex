package usecase

import (
	"main/modules/auth/repository"
)

type AuthUsecaseService interface {

}

type authUsecase struct {
	authRepository repository.AuthRepositoryService
}

func NewAuthUsecase(authRepository repository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{authRepository}
}
