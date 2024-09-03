package response

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type (
	MsgResponse struct {
		Message string `json:"message"`
	}
)

// Custom logging function, or replace with your preferred logger
func logError(message string) {
	logger := log.New(log.Writer(), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
}

func ErrResponse(c *fiber.Ctx, statusCode int, message string) error {
	if statusCode == fiber.StatusInternalServerError {
		// log error
		logError(message)
	}

	return c.Status(statusCode).JSON(&MsgResponse{Message: message})
}

func SuccessResponse(c *fiber.Ctx, statusCode int, data any) error {
	return c.Status(statusCode).JSON(data)
}
