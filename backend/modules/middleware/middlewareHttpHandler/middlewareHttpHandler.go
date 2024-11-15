package middlewareHandler

import (
	"main/config"
	middlewareusecase "main/modules/middleware/middlewareUsecase"
	"main/pkg/rbac"
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
		RequirePermission(permission string) echo.MiddlewareFunc
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
			// Define public paths that don't require authentication
			publicPaths := []string{
				"/auth_v1/auth/login",
				"/auth_v1/auth/refresh-token",
				"/auth_v1/auth/health",
				"/user_v1/users/register",
				"/user_v1/health",
			}

			// Check if the current path is public
			currentPath := c.Request().URL.Path
			for _, path := range publicPaths {
				if strings.HasPrefix(currentPath, path) {
					return next(c)
				}
			}

			// For all other paths, require authentication
			authorizationHeader := c.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
			}

			parts := strings.Split(authorizationHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}
			accessToken := parts[1]

			_, err := m.middlewareUsecase.JwtAuthorization(c, cfg, accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

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

func (m *middlewareHandler) RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleCode := c.Get("role_code").(int)
			
			if !rbac.HasPermission(roleCode, permission) {
				return echo.NewHTTPError(http.StatusForbidden, "Permission denied")
			}
			
			return next(c)
		}
	}
}