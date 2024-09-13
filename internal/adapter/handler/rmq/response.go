package rmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"golang-hexagon/internal/core/domain"
	"net/http"
)

const contentType = "application/json"

type (
	ResponseMessage struct {
		Success bool   `json:"success"`
		Status  int    `json:"statusCode"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
)

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrConflictingData:            http.StatusConflict,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
}

// newResponseMessage creates a new response message for RMQ sending
func newResponseMessage(message string, err error) *amqp.Publishing {
	var (
		statusCode int
		ok         bool
		errMsg     string
	)

	if err != nil {
		statusCode, ok = errorStatusMap[err]
		if !ok {
			statusCode = http.StatusInternalServerError
		}
		errMsg = err.Error()
	} else {
		statusCode = http.StatusOK
	}
	rsp := ResponseMessage{
		Success: statusCode == http.StatusOK,
		Status:  statusCode,
		Message: message,
		Error:   errMsg,
	}
	// error is muted here because we know that there will be no encoding errors
	body, _ := json.Marshal(rsp)

	return &amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	}
}
