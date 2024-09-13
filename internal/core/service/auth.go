package service

import (
	"context"
	"errors"
	"golang-hexagon/internal/core/domain"
	"golang-hexagon/internal/core/port"
	"golang-hexagon/internal/core/util"
)

// AuthService implements port.AuthService interface
// and provides access to the user repository
// and token service
type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

// NewAuthService creates a new auth service instance
func NewAuthService(repo port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		repo,
		ts,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, email, password string) ([]byte, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, domain.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return nil, domain.ErrTokenCreation
	}

	return []byte(accessToken), nil
}
