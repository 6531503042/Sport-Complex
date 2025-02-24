package main

import (
	"context"
	"log"
	"main/config"
	"main/pkg/database"
	"main/pkg/database/migration"
	"main/server"
	"os"
)

func main() {
	// Initialize context
	ctx := context.Background()

	// Load configuration from .env file
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	// Connect to the database
	db := database.DbConn(ctx, &cfg)
	defer db.Disconnect(ctx)

	// Perform database migrations once
	if err := migration.SetupFacilities(ctx, &cfg, db); err != nil {
		log.Fatalf("Migration setup failed: %v", err)
	}

	// Start the server
	server.Start(ctx, &cfg, db)
}
