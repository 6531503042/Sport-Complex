package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	//TODO: Fiber Instace
	app := fiber.New()

	// TODO: Route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//TODO: Listener
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}