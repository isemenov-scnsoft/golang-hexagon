package http

import (
	"os"

	"github.com/joho/godotenv"
)

const EnvProduction = "production"

// Config contains all the environment variables for the http server
type Config struct {
	Env            string
	URL            string
	Port           string
	AllowedOrigins string
}

// New creates a new container instance
func New() (*Config, error) {
	if os.Getenv("APP_ENV") != EnvProduction {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}, nil
}
