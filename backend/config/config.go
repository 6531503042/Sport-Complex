package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (

	Config struct {
		App App
		Db   Db
		Grpc Grpc
		Server ServerConfig
		Jwt Jwt
	}

	ServerConfig struct {
		Port int
	}

	App struct {
		Name string
		Url string
		Stage string
	}

	Db struct {
		Url string
	}

	Jwt struct {
		AccessSecretKey string
		RefreshSecretKey string
		ApiSecretKey string
		AccessDuration int64
		RefreshDuration int64
		ApiDuration int64
	}

	Grpc struct {
		AuthUrl string
		UserUrl string
		GymUrl string
		BadmintonUrl string
		SwimmingUrl string
		FootballUrl string
		PaymentUrl string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config {
		App : App {
			Name : os.Getenv("APP_NAME"),
			Url : os.Getenv("APP_URL"),
			Stage : os.Getenv("APP_STAGE"),
		},
		Db : Db {
			Url : os.Getenv("DB_URL"),
		},
		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ApiSecretKey:     os.Getenv("JWT_API_SECRET_KEY"),
			AccessDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("Error loading access duration failed")
				}
				return result
			}(),
			RefreshDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("Error loading refresh duration failed")
				}
				return result
			}(),
		},
	}
}
