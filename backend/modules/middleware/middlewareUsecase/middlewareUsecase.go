package middlewareusecase

import (
	"errors"
	"log"
	"main/config"
	middlewareRepository "main/modules/middleware/middlewareRepository"
	"main/pkg/jwt"
	"main/pkg/rbac"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
		RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error)
		IsAdminRole(c echo.Context, cfg *config.Config, roleCode int) (int64, error)
		UserIdParamValidation(c echo.Context) (echo.Context, error)
	}

	middlewareUsecase struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewareRepository}
}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {
	ctx := c.Request().Context()

	claims, err := jwt.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	if err := u.middlewareRepository.AccessTokenSearch(ctx, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}

	c.Set("user_id", claims.UserId)
	c.Set("role_code", claims.RoleCode)

	return c, nil
}

//Set authorization for rbac.
func (u *middlewareUsecase) RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error) {
	ctx := c.Request().Context()

	userRoleCode := c.Get("role_code").(int)

	rolesCount, err := u.middlewareRepository.RolesCount(ctx, cfg.Grpc.AuthUrl)
	if err != nil {
		return nil, err
	}

	userRoleBinary := rbac.IntToBinary(userRoleCode, int(rolesCount))

	for i := 0; i < int(rolesCount); i++ {
		if userRoleBinary[i]&expected[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("error: permission denied")
}

func (u *middlewareUsecase) UserIdParamValidation(c echo.Context) (echo.Context, error) {
	userIdReq := c.Param("user_id")
	userIdToken := c.Get("user_id").(string)

	if userIdToken == "" {
		log.Printf("Error: user_id not found")
		return nil, errors.New("error: user_id is required")
	}

	if userIdToken != userIdReq {
		log.Printf("Error: user_id not match, user_id_req: %s, user_id_token: %s", userIdReq, userIdToken)
		return nil, errors.New("error: user_id not match")
	}

	return c, nil
}

func (u *middlewareUsecase) IsAdminRole(c echo.Context, cfg *config.Config, roleCode int) (int64, error) {
	ctx := c.Request().Context()
	return u.middlewareRepository.IsAdminRole(ctx, cfg.Grpc.AuthUrl, roleCode)
}
