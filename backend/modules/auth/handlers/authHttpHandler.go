package handlers

import (
	"context"
	"fmt"
	"main/config"
	"main/modules/auth"
	"main/modules/auth/usecase"
	"main/pkg/request"
	"main/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type (
	AuthHttpHandlerService interface {
		Login(c *fiber.Ctx) error
		RefreshToken(c *fiber.Ctx) error
		Logout(c *fiber.Ctx) error
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase usecase.AuthUsecaseService
	}
)

// NewAuthHttpHandler initializes and returns an AuthHttpHandlerService
func NewAuthHttpHandler(cfg *config.Config, authUsecase usecase.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUsecase: authUsecase,
	}
}

func (h *authHttpHandler) Login(c *fiber.Ctx) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.UserLoginReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, fiber.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.Login(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, fiber.StatusOK, res)
}

func (h *authHttpHandler) RefreshToken(c *fiber.Ctx) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.RefreshTokenReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, fiber.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.RefreshToken(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, fiber.StatusOK, res)
}

func (h *authHttpHandler) Logout(c *fiber.Ctx) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.LogoutReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, fiber.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.Logout(ctx, req.CredentialId)
	if err != nil {
		return response.ErrResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, fiber.StatusOK, &response.MsgResponse{
		Message: fmt.Sprintf("Deleted count: %d", res),
	})
}
