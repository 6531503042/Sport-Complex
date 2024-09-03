package middlewarehttphandler

import (
	"main/config"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type (
	MiddlewareHttpHandlerService interface {
		JwtAuthorizationMiddleware(cfg *config.Config) fiber.Handler
		RbacAuthorizationMiddleware(cfg *config.Config, expected []int) fiber.Handler
		UserIdParamValidationMiddleware() fiber.Handler
	}

	MiddlewareHttpHandler struct {
		middlewareUsecase middlewareusecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHttpHandler(middlewareUsecase middlewareusecase.MiddlewareUsecaseService) MiddlewareHttpHandlerService {
	return &MiddlewareHttpHandler{middlewareUsecase}
}

func (m *MiddlewareHttpHandler) JwtAuthorizationMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Get("Authorization")
		if accessToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing access token"})
		}

		_, err := m.middlewareUsecase.JwtAuthorization(c, cfg, accessToken)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Next()
	}
}

func (m *MiddlewareHttpHandler) RbacAuthorizationMiddleware(cfg *config.Config, expected []int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := m.middlewareUsecase.RbacAuthorization(c, cfg, expected)
		if err != nil {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Next()
	}
}

func (m *MiddlewareHttpHandler) UserIdParamValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := m.middlewareUsecase.UserIdParamValidation(c)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Next()
	}
}
