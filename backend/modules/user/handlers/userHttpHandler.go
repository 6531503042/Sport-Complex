package handlers

import (
	"log"
	"main/config"
	"main/modules/user"
	"main/modules/user/usecase"
	"main/pkg/request"
	"main/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type (
    userHttpHandler struct {
        cfg         *config.Config
        userUsecase usecase.UserUsecaseService
    }
)

// NewUserHttpHandler initializes and returns a userHttpHandler
func NewUserHttpHandler(cfg *config.Config, userUsecase usecase.UserUsecaseService) *userHttpHandler {
    return &userHttpHandler{
        cfg:         cfg,
        userUsecase: userUsecase,
    }
}

func (h *userHttpHandler) CreateUser(c *fiber.Ctx) error {
    log.Println("Received request to create user")

    // Use the custom binding
    wrapper := request.ContextWrapper(c)
    req := new(user.CreateUserReq)

    // Handle binding errors
    if err := wrapper.Bind(req); err != nil {
        return response.ErrResponse(c, fiber.StatusBadRequest, "Invalid request payload")
    }

    // Call the Usecase to create a user
    ctx := c.Context() // Fiber context
    res, err := h.userUsecase.CreateUser(ctx, req)
    if err != nil {
        return response.ErrResponse(c, fiber.StatusBadRequest, err.Error())
    }

    return response.SuccessResponse(c, fiber.StatusCreated, res)
}

func (h *userHttpHandler) FindOneUserProfile(c *fiber.Ctx) error {
    // Use the custom context
    ctx := c.Context()

    // Extract the user ID from the URL parameter
    userId := c.Params("user_id")

    // Fetch the user profile
    res, err := h.userUsecase.FindOneUserProfile(ctx, userId)
    if err != nil {
        // Return an error response
        return response.ErrResponse(c, fiber.StatusBadRequest, err.Error())
    }

    // Return a successful response with the user profile
    return response.SuccessResponse(c, fiber.StatusOK, res)
}
