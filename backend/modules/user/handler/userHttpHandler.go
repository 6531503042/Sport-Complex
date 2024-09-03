package handler

import (
	"net/http"

	models "main/modules/user/model"
	"main/modules/user/usecase"

	fiber "github.com/gofiber/fiber/v2"
)

// UserHandler is a struct that contains the UserUsecase
type UserHandler struct {
	UserUsecase usecase.IUserUsecase
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userProfile, err := h.UserUsecase.CreateUser(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(userProfile)
}

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	userProfile, err := h.UserUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(userProfile)
}

// GetAllUsers retrieves all users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.UserUsecase.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(users)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userProfile, err := h.UserUsecase.UpdateUser(c.Context(), id, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(userProfile)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.UserUsecase.DeleteUser(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusNoContent)
}

