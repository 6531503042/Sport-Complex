package handlers

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"main/modules/auth"
	"main/modules/auth/usecase"
	"main/pkg/request"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (

	AuthHttpHandlerService interface {
		// Blank\
		Login (c echo.Context) error 
		RefreshToken (c echo.Context) error
		Logout (c echo.Context) error
	}

	authHttpHandler struct {
		cfg *config.Config
		authUsecase usecase.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase usecase.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{cfg, authUsecase}
}

func (h *authHttpHandler) Login (c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.UserLoginReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.Login(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *authHttpHandler) RefreshToken(c echo.Context) error {
    ctx := context.Background()
    wrapper := request.ContextWrapper(c)
    req := new(auth.RefreshTokenReq)

    if err := wrapper.Bind(req); err != nil {
        return response.ErrResponse(c, http.StatusBadRequest, err.Error())
    }

    log.Printf("Received refresh token request: CredentialId: %s, RefreshToken: %s", req.CredentialId, req.RefreshToken)

    res, err := h.authUsecase.RefreshToken(ctx, h.cfg, req)
    if err != nil {
        return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
    }

    return response.SuccessResponse(c, http.StatusOK, res)
}


func (h *authHttpHandler) Logout (c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.LogoutReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.Logout(ctx, req.CredentialId)
	if err != nil {
		return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, &response.MsgResponse{
		Message: fmt.Sprintf("Deleted count: %d", res),
	})
}