package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents the application configuration.
type (
	Config struct {
	App App
	Db   Db
	Grpc Grpc
	Kafka    Kafka
	Jwt Jwt
}

Kafka struct {
	Url    string
	ApiKey string
	Secret string
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
	FacilityUrl string
	BookingUrl string
	PaymentUrl string
}

)


// LoadConfig loads the configuration from the given .env file.
func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Config{
		App: App{
			Name: os.Getenv("APP_NAME"),
			Url:  os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
		Db: Db{
			Url: os.Getenv("DB_URL"),
		},
		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ApiSecretKey:     os.Getenv("JWT_API_SECRET_KEY"),
			AccessDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_DURATION"), 10, 64)
				if err != nil {
					log.Fatalf("Error loading access duration failed: %v", err)
				}
				return result
			}(),
			RefreshDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_DURATION"), 10, 64)
				if err != nil {
					log.Fatalf("Error loading refresh duration failed: %v", err)
				}
				return result
			}(),
		},
		Grpc: Grpc{
			AuthUrl: os.Getenv("GRPC_AUTH_URL"),
			UserUrl: os.Getenv("GRPC_USER_URL"),
			FacilityUrl: os.Getenv("GRPC_FACILITY_URL"),
			BookingUrl: os.Getenv("GRPC_BOOKING_URL"),
			PaymentUrl: os.Getenv("GRPC_PAYMENT_URL"),
		},
		Kafka: Kafka{
			Url:    os.Getenv("KAFKA_URL"),
			ApiKey: os.Getenv("KAFKA_API_KEY"),
			Secret: os.Getenv("KAFKA_SECRET"),
		},
	}
}

