package port

import (
	"context"
	"golang-hexagon/internal/core/domain"
)

//go:generate mockgen -source=user.go -destination=mock/user.go -package=mock

type (
	// UserRepository is an interface for interacting with user-related data
	UserRepository interface {
		// CreateUser inserts a new user into the database
		CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
		// GetUserByID selects a user by id
		GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
		// GetUserByEmail selects a user by email
		GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
		// ListUsers selects a list of users with pagination
		ListUsers(ctx context.Context, skip, limit uint64) ([]*domain.User, error)
		// UpdateUser updates a user
		UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
		// DeleteUser deletes a user
		DeleteUser(ctx context.Context, id uint64) error
	}

	UserService interface {
		// Register registers a new user
		Register(ctx context.Context, user *domain.User) (*domain.User, error)
		// GetUser returns a user by id
		GetUser(ctx context.Context, id uint64) (*domain.User, error)
		// ListUsers returns a list of users with pagination
		ListUsers(ctx context.Context, offset, limit uint64) ([]*domain.User, error)
		// UpdateUser updates a user
		UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
		// DeleteUser deletes a user
		DeleteUser(ctx context.Context, id uint64) error
	}
)
