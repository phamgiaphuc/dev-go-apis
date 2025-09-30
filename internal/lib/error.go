package lib

import (
	"dev-go-apis/internal/models"
	"net/http"
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
	MissingSignatureAndTimestampError = &models.APIError{Code: http.StatusUnauthorized, Message: "Missing signature and timestamp"}
	ExpiredTimestampError             = &models.APIError{Code: http.StatusUnauthorized, Message: "Expired timestamp"}
	InvalidSignatureError             = &models.APIError{Code: http.StatusBadRequest, Message: "Invalid signature"}
)

func NewAPIError(code int, message, stack string) *models.APIError {
	return &models.APIError{Code: code, Message: message, Stack: stack}
}
