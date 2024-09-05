package handlers

import (
	"log"
	"main/config"
	"main/modules/user"
	"main/modules/user/usecase"
	"main/pkg/request"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	NewUserHttpHandlerService interface{
		CreateUser(c echo.Context) error
		FindOneUserProfile(c echo.Context) error
	}

	userHttpHandler struct {
		cfg         *config.Config
		userUsecase usecase.UserUsecaseService
	}
)

// NewUserHttpHandler initializes and returns a userHttpHandler
func NewUserHttpHandler(cfg *config.Config, userUsecase usecase.UserUsecaseService) userHttpHandler {
	return &userHttpHandler{userUsecase: userUsecase}
}

func (h *userHttpHandler) CreateUser(c echo.Context) error {

	log.Println("Received request to create user")
	// Use the context from the request
	ctx := c.Request().Context()

	// Use the custom binding
	wrapper := request.ContextWrapper(c)

	req := new(user.CreateUserReq)

	// Handle binding errors
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	// Call the Usecase to create a user
	res, err := h.userUsecase.CreateUser(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *userHttpHandler) FindOneUserProfile(c echo.Context) error {
	// Use the context from the request
	ctx := c.Request().Context()

	// Extract the user ID from the URL parameter
	userId := c.Param("user_id")

	// Fetch the user profile
	res, err := h.userUsecase.FindOneUserProfile(ctx, userId)
	if err != nil {
		// Return an error response
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	// Return a successful response with the user profile
	return response.SuccessResponse(c, http.StatusOK, res)
}
