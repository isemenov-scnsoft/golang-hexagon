package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang-hexagon/internal/adapter/config"
	"golang-hexagon/internal/core/domain"
	"golang-hexagon/internal/core/port"
	"time"
)

// PasetoToken implements port.TokenService interface
// and provides an access to the paseto library
type PasetoToken struct {
	token    *paseto.Token
	key      *paseto.V4SymmetricKey
	parser   *paseto.Parser
	duration time.Duration
}

// New creates a new paseto instance
func New(config *config.Token) (port.TokenService, error) {
	durationStr := config.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, domain.ErrTokenDuration
	}

	token := paseto.NewToken()
	key := paseto.NewV4SymmetricKey()
	parser := paseto.NewParser()

	return &PasetoToken{
		&token,
		&key,
		&parser,
		duration,
	}, nil
}

// CreateToken creates a new paseto token
func (pt *PasetoToken) CreateToken(user *domain.User) ([]byte, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, domain.ErrTokenCreation
	}

	payload := &domain.TokenPayload{
		ID:     id,
		UserID: user.ID,
		Role:   user.Role,
	}

	err = pt.token.Set("payload", payload)
	if err != nil {
		return nil, domain.ErrTokenCreation
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(pt.duration)

	pt.token.SetIssuedAt(issuedAt)
	pt.token.SetNotBefore(issuedAt)
	pt.token.SetExpiration(expiredAt)

	token := []byte(pt.token.V4Encrypt(*pt.key, nil))

	return token, nil
}

// VerifyToken verifies the paseto token
func (pt *PasetoToken) VerifyToken(token []byte) (*domain.TokenPayload, error) {
	var payload *domain.TokenPayload

	parsedToken, err := pt.parser.ParseV4Local(*pt.key, string(token), nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	err = parsedToken.Get("payload", &payload)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	return payload, nil
}
