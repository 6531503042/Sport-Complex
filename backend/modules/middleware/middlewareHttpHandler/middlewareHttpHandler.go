package middlewarehttphandler

import (
	"main/config"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareHttpHandlerService interface {
		JwtAuthorizationMiddleware(cfg *config.Config) echo.MiddlewareFunc
		RbacAuthorizationMiddleware(cfg *config.Config, expected []int) echo.MiddlewareFunc
		UserIdParamValidationMiddleware() echo.MiddlewareFunc
		IsAdminRoleMiddleware(cfg *config.Config, roleCode int) echo.MiddlewareFunc
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareusecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHttpHandler(cfg *config.Config, middlewareUsecase middlewareusecase.MiddlewareUsecaseService) MiddlewareHttpHandlerService {
	return &middlewareHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}

func (m *middlewareHandler) JwtAuthorizationMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the Authorization header
			authorizationHeader := c.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
			}

			// Extract the token from the header
			parts := strings.Split(authorizationHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}
			accessToken := parts[1]

			// Call the JwtAuthorization function
			_, err := m.middlewareUsecase.JwtAuthorization(c, cfg, accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			// If everything is fine, proceed to the next handler
			return next(c)
		}
	}
}

func (m *middlewareHandler) RbacAuthorizationMiddleware(cfg *config.Config, expected []int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := m.middlewareUsecase.RbacAuthorization(c, cfg, expected)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}

			return next(c)
		}
	}
}

func (m *middlewareHandler) UserIdParamValidationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := m.middlewareUsecase.UserIdParamValidation(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return next(c)
		}
	}
}

func (m *middlewareHandler) IsAdminRoleMiddleware(cfg *config.Config, roleCode int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			admin, err := m.middlewareUsecase.IsAdminRole(c, cfg, roleCode)
			if err != nil || admin == -1 {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}

			return next(c)
		}
	}
}