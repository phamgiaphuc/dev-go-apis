package lib

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"dev-go-apis/internal/models"

	"github.com/go-playground/validator/v10"
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

func GetValidationErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "email":
		return "Invalid email format"
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "lte":
		return fmt.Sprintf("%s should be less than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s should be greater than %s", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

func ParseValidationErrors(err error) []models.ValidationError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}

	errorsList := make([]models.ValidationError, len(ve))
	for i, fe := range ve {
		errorsList[i] = models.ValidationError{
			Field:   strings.ToLower(fe.Field()),
			Message: GetValidationErrorMsg(fe),
		}
	}

	return errorsList
}
