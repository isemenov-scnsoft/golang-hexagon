package rmq

import (
	"os"

	"github.com/joho/godotenv"
)

const EnvProduction = "production"

// Config contains all the environment variables for the RabbitMQ
type Config struct {
	Host          string
	Port          string
	User          string
	Password      string
	Vhost         string
	InQueue       string
	ConsumerTag   string
	OutExchange   string
	OutRoutingKey string
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
		Host:          os.Getenv("RMQ_HOST"),
		Port:          os.Getenv("RMQ_PORT"),
		User:          os.Getenv("RMQ_USER"),
		Password:      os.Getenv("RMQ_PASSWORD"),
		Vhost:         os.Getenv("RMQ_VHOST"),
		InQueue:       os.Getenv("RMQ_IN_QUEUE"),
		ConsumerTag:   os.Getenv("RMQ_CONSUMER_TAG"),
		OutExchange:   os.Getenv("RMQ_OUT_EXCHANGE"),
		OutRoutingKey: os.Getenv("RMQ_OUT_QUEUE"),
	}, nil
}
