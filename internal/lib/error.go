package lib

import (
	"net/http"

	"dev-go-apis/internal/models"
)

var (
	InvalidRequestError               = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid request"}
	InvalidParamRequestError          = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid param request"}
	InvalidBodyRequestError           = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid body request"}
	InternalServerError               = &models.APIError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	InvalidOAuthStateError            = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid OAuth state"}
	ResourceNotFoundError             = &models.APIError{Code: http.StatusNotFound, Message: "Resource not found"}
	ResourceForbiddenError            = &models.APIError{Code: http.StatusForbidden, Message: "Resource forbidden"}
	UnauthorizedError                 = &models.APIError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	MissingAPIKeyError                = &models.APIError{Code: http.StatusUnauthorized, Message: "Missing api key"}
	MissingSignatureAndTimestampError = &models.APIError{Code: http.StatusUnauthorized, Message: "Missing request signature and timestamp"}
	ExpiredTimestampError             = &models.APIError{Code: http.StatusUnauthorized, Message: "Expired request timestamp"}
	InvalidSignatureError             = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid request signature"}
	TooManyRequestsError              = &models.APIError{Code: http.StatusTooManyRequests, Message: "Too many requests"}
)

func NewAPIError(code int, message, stack string) *models.APIError {
	return &models.APIError{Code: code, Message: message, Stack: stack}
}
