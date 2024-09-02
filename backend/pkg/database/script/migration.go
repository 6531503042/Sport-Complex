package main

import (
	"context"
	"log"
	"main/config"
	"main/pkg/database/migration"
	"os"
)

func main() {
	ctx := context.Background()

	// Initialize config
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	switch cfg.App.Name {
	case "user":
		migration.UserMigrate(ctx, &cfg)
	case "auth" :
		migration.AuthMigrate(ctx, &cfg)
	//other migration db script
	}
}