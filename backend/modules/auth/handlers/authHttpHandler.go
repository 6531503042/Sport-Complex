package handlers

import (
	"main/config"
	"main/modules/auth/usecase"
)

type (

	AuthHttpHandlerService interface {
		// Blank
	}

	authHttpHandler struct {
		cfg *config.Config
		authUsecase usecase.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase usecase.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{cfg, authUsecase}
}