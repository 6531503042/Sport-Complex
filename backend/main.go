package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"main/config"
	"main/modules/user/handler"
	"main/modules/user/repository"
	"main/modules/user/usecase"
	"main/pkg/database"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig("./env/.env")

	// Initialize database connection
	client := database.DbConn(context.Background(), &cfg)

	// Ensure proper disconnection when the application exits
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error: Disconnecting database: %s", err.Error())
		}
	}()

	// Create Fiber app instance
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Initialize repository, usecase, and handler
	userRepo := repository.NewUserRepository(client)
	userUsecase := usecase.NewUserUsecase(userRepo)
	handler.NewUserHandler(app, userUsecase)

	// Start the server
	port := ":3000"
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error: Starting server: %s", err.Error())
	}
}
