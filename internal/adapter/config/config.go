package config

import (
	"fmt"
	"golang-hexagon/internal/adapter/config/http"
	"golang-hexagon/internal/adapter/config/rmq"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvProduction = "production"
	appTypeRMQ    = "rmq"
	appTypeHTTP   = "http"
)

type (
	// Container contains environment variables for the application, database, cache, token, and http server
	Container struct {
		App   *App
		Redis *Redis
		DB    *DB
		Token *Token
		RMQ   *rmq.Config
		HTTP  *http.Config
	}

	// App contains all the environment variables for the application
	App struct {
		Name string
		Env  string
		Type string
	}
	// Redis contains all the environment variables for the cache service
	Redis struct {
		Addr     string
		Password string
	}

	// Token contains all the environment variables for the token service
	Token struct {
		Duration string
	}

	// DB contains all the environment variables for the database
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
)

// New creates a new container instance
func New() (*Container, error) {
	if os.Getenv("APP_ENV") != EnvProduction {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
		Type: os.Getenv("APP_TYPE"),
	}

	redis := &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	token := &Token{
		Duration: os.Getenv("TOKEN_DURATION"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	container := &Container{
		App:   app,
		Redis: redis,
		DB:    db,
		Token: token,
	}

	switch app.Type {
	case appTypeRMQ:
		rmqConf, err := rmq.New()
		if err != nil {
			return nil, err
		}
		container.RMQ = rmqConf
	case appTypeHTTP:
		httpConf, err := http.New()
		if err != nil {
			return nil, err
		}
		container.HTTP = httpConf
	default:
		return nil, fmt.Errorf("invalid application type: %s", app.Type)
	}

	return container, nil
}
