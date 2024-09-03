package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type healthCheck struct {
	App    string `json:"app"`
	Status string `json:"status"`
}

func (s *server) healthCheckService(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	return c.JSON(&healthCheck{
		App:    s.cfg.App.Name,
		Status: "OK",
	})
}
