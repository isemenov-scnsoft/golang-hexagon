package http

import (
	"github.com/go-playground/validator/v10"
	"golang-hexagon/internal/core/domain"
)

var userRoleValidator validator.Func = func(fl validator.FieldLevel) bool {
	userRole := fl.Field().Interface().(domain.UserRole)

	switch userRole {
	case domain.Admin, domain.Basic:
		return true
	default:
		return false
	}
}
