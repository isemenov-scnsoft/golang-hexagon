package main

import (
	"context"
	"fmt"
	"golang-hexagon/internal/adapter/auth/paseto"
	"golang-hexagon/internal/adapter/config"
	"golang-hexagon/internal/adapter/handler/http"
	"golang-hexagon/internal/adapter/logger"
	"golang-hexagon/internal/adapter/storage/postgres"
	"golang-hexagon/internal/adapter/storage/postgres/repository"
	"golang-hexagon/internal/adapter/storage/redis"
	"golang-hexagon/internal/core/service"
	"log/slog"
	"os"
)

// @title						Go Hexagon Arch POC (Proof of concept) API
// @version					1.0
// @description				This is a simple RESTful Service API written in Go using Gin web framework, PostgreSQL database, and Redis cache.
//
// @contact.name				Ivan Semenov
// @contact.email				isemenov@scnsoft.com
//
// @host						localhost
// @BasePath					/v1
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	// Load environment variables
	conf, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger.Set(conf.App)

	slog.Info("Starting the application", "app", conf.App.Name, "env", conf.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, conf.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Successfully connected to the database", "db", conf.DB.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")

	// Init cache service
	cache, err := redis.New(ctx, conf.Redis)
	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := cache.Close(); err != nil {
			slog.Error("Error closing cache connection", err)
		}
	}()

	slog.Info("Successfully connected to the cache server")

	// Init token service
	token, err := paseto.New(conf.Token)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	// Dependency injection
	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cache)
	userHandler := http.NewUserHandler(userService)

	// Auth
	authService := service.NewAuthService(userRepo, token)
	authHandler := http.NewAuthHandler(authService)

	// Init router
	router, err := http.NewRouter(
		conf,
		token,
		*userHandler,
		*authHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", conf.HTTP.URL, conf.HTTP.Port)
	slog.Info("Starting the Config server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the Config server", "error", err)
		os.Exit(1)
	}
}
