package lib

import (
	"dev-go-apis/internal/models"
	"net/http"
)

var (
	InvalidParamRequestError = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid param request"}
	InvalidBodyRequestError  = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid body request"}
	InternalServerError      = &models.APIError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ResourceNotFoundError    = &models.APIError{Code: http.StatusNotFound, Message: "Resource not found"}
	ResourceForbiddenError   = &models.APIError{Code: http.StatusForbidden, Message: "Resource forbidden"}
	UnauthorizedError        = &models.APIError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
)

func NewAPIError(code int, message, stack string) *models.APIError {
	return &models.APIError{Code: code, Message: message, Stack: stack}
}
