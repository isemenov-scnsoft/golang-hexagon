package rmq

import "golang-hexagon/internal/core/domain"

// asVal returns a value from pointer
func asVal[T any](val *T) T {
	var v T
	if val != nil {
		v = *val
	}
	return v
}

// asPtr returns a pointer for the value
func asPtr[T any](val T) *T {
	return &val
}

func toUser(msg *msg) *domain.User {
	return &domain.User{
		ID:       asVal(msg.UID),
		Email:    asVal(msg.Email),
		Password: string(asVal(msg.Password)),
		Role:     asVal(msg.Role),
	}
}
