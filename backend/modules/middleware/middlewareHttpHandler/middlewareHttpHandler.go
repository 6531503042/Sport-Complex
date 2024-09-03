package middlewarehttphandler

import (
	"main/config"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareHttpHandlerService interface {
		JwtAuthorizationMiddleware(cfg *config.Config) echo.MiddlewareFunc
		RbacAuthorizationMiddleware(cfg *config.Config, expected []int) echo.MiddlewareFunc
		UserIdParamValidationMiddleware() echo.MiddlewareFunc
	}

	MiddlewareHttpHandler struct {
		middlewareUsecase middlewareusecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHttpHandler(middlewareUsecase middlewareusecase.MiddlewareUsecaseService) MiddlewareHttpHandlerService {
	return &MiddlewareHttpHandler{middlewareUsecase}
}

func (m *MiddlewareHttpHandler) JwtAuthorizationMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			if accessToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
			}

			_, err := m.middlewareUsecase.JwtAuthorization(c, cfg, accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			return next(c)
		}
	}
}

func (m *MiddlewareHttpHandler) RbacAuthorizationMiddleware(cfg *config.Config, expected []int) echo.MiddlewareFunc {
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

func (m *MiddlewareHttpHandler) UserIdParamValidationMiddleware() echo.MiddlewareFunc {
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
