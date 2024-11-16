package handlers

import (
	"log"
	"main/config"
	"main/modules/user"
	"main/modules/user/usecase"
	"main/pkg/request"
	"main/pkg/response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	NewUserHttpHandlerService interface{
		CreateUser(c echo.Context) error
		FindOneUserProfile(c echo.Context) error
		FindManyUser (c echo.Context) error
		DeleteUser(c echo.Context) error
		UpdateUser(c echo.Context) error
		GetUserAnalytics(c echo.Context) error
	}

	userHttpHandler struct {
		cfg         *config.Config
		userUsecase usecase.UserUsecaseService
	}
)

// NewUserHttpHandler initializes and returns a userHttpHandler

func NewUserHttpHandler(cfg *config.Config, userUsecase usecase.UserUsecaseService) NewUserHttpHandlerService {
	return &userHttpHandler{cfg: cfg, userUsecase: userUsecase}
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
	userId := strings.TrimPrefix(c.Param("user_id"), "user:")

	// Fetch the user profile
	res, err := h.userUsecase.FindOneUserProfile(ctx, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	// Return a successful response with the user profile
	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *userHttpHandler) FindManyUser (c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.userUsecase.FindManyUser(ctx)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, users)
}

func (h *userHttpHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Param("user_id")

	var updateReq map[string]interface{}
	if err := c.Bind(&updateReq); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	err := h.userUsecase.UpdateOneUser(ctx, userId, updateReq)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "User updated successfully")
}

func (h *userHttpHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Param("user_id")

	err := h.userUsecase.DeleteUser(ctx, userId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "User deleted successfully")
}

func (h *userHttpHandler) GetUserAnalytics(c echo.Context) error {
	ctx := c.Request().Context()
	
	query := new(user.AnalyticsQuery)
	if err := c.Bind(query); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, "Invalid query parameters")
	}

	if err := c.Validate(query); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	analytics, err := h.userUsecase.GetUserAnalytics(ctx, query)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, analytics)
}