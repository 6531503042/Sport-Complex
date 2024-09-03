package middlewareusecase

import (
	"errors"
	"log"
	"main/config"
	middlewareRepository "main/modules/middleware/middlewareRepository"
	"main/pkg/jwt"
	"main/pkg/rbac"

	"github.com/gofiber/fiber/v2"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c *fiber.Ctx, cfg *config.Config, accessToken string) (*fiber.Ctx, error)
		RbacAuthorization(c *fiber.Ctx, cfg *config.Config, expected []int) (*fiber.Ctx, error)
		UserIdParamValidation(c *fiber.Ctx) (*fiber.Ctx, error)
	}

	middlewareUsecase struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewareRepository}
}

func (u *middlewareUsecase) JwtAuthorization(c *fiber.Ctx, cfg *config.Config, accessToken string) (*fiber.Ctx, error) {
	claims, err := jwt.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	ctx := c.Context() // This is *fasthttp.RequestCtx
	if err := u.middlewareRepository.AccessTokenSearch(ctx, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}

	c.Locals("user_id", claims.UserId)
	c.Locals("role_code", claims.RoleCode)

	return c, nil
}

// Set authorization for RBAC
func (u *middlewareUsecase) RbacAuthorization(c *fiber.Ctx, cfg *config.Config, expected []int) (*fiber.Ctx, error) {
	ctx := c.Context() // This is *fasthttp.RequestCtx

	userRoleCode := c.Locals("role_code").(int)

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

func (u *middlewareUsecase) UserIdParamValidation(c *fiber.Ctx) (*fiber.Ctx, error) {
	userIdReq := c.Params("user_id")
	userIdToken := c.Locals("user_id")

	if userIdToken == nil {
		log.Printf("Error: user_id not found")
		return nil, errors.New("error: user_id is required")
	}

	userIdTokenStr := userIdToken.(string)

	if userIdTokenStr != userIdReq {
		log.Printf("Error: user_id not match, user_id_req: %s, user_id_token: %s", userIdReq, userIdTokenStr)
		return nil, errors.New("error: user_id not match")
	}

	return c, nil
}
