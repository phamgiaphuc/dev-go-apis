package lib

import (
	"dev-go-apis/internal/models"
	"net/http"
)

var (
	InvalidBodyRequestError = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid body request"}
	InternalServerError     = &models.APIError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ResourceNotFoundError   = &models.APIError{Code: http.StatusNotFound, Message: "Resource not found"}
)

func NewAPIError(code int, message, stack string) *models.APIError {
	return &models.APIError{Code: code, Message: message, Stack: stack}
}
