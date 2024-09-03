package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App App
		Db   Db
	}

	App struct {
		Name string
		Url string
		Stage string
	}

	Db struct {
		Url string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatalf("Error: Load .env file error: %s", err.Error())
	}

	return Config {
		App: App{
			Name: os.Getenv("APP_NAME"),
			Url:  os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
		Db: Db{
			Url: os.Getenv("DB_URL"),
		},
	}
}