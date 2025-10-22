package domain

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	Internal   string
	HTTPStatus int
	Metadata   map[string]interface{}
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Internal, e.Err)
	}
	return e.Internal
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithContext(key string, value interface{}) *AppError {
	newErr := *e
	if newErr.Metadata == nil {
		newErr.Metadata = make(map[string]interface{})
	}
	newErr.Metadata[key] = value
	return &newErr
}

func (e *AppError) Wrap(err error) *AppError {
	newErr := *e
	newErr.Err = err
	return &newErr
}

func (e *AppError) WithMessage(msg string) *AppError {
	newErr := *e
	newErr.Message = msg
	return &newErr
}

var (
	ErrURLNotFound = &AppError{
		Code:       "URL_NOT_FOUND",
		Message:    "The short URL you're looking for does not exist",
		Internal:   "url not found in database",
		HTTPStatus: http.StatusNotFound,
	}

	ErrURLExpired = &AppError{
		Code:       "URL_EXPIRED",
		Message:    "This short URL has expired",
		Internal:   "url expiration date has passed",
		HTTPStatus: http.StatusGone,
	}

	ErrURLInactive = &AppError{
		Code:       "URL_INACTIVE",
		Message:    "This short URL is currently inactive",
		Internal:   "url is marked as inactive",
		HTTPStatus: http.StatusForbidden,
	}
)

var (
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required",
		Internal:   "missing or invalid authentication token",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenExpired = &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    "Your session has expired. Please log in again",
		Internal:   "jwt token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidToken = &AppError{
		Code:       "INVALID_TOKEN",
		Message:    "Invalid authentication token",
		Internal:   "jwt token validation failed",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidSigningKey = &AppError{
		Code:       "INVALID_SIGNING_KEY",
		Message:    "Authentication system error",
		Internal:   "jwt signing key validation failed",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "You don't have permission to access this resource",
		Internal:   "insufficient permissions",
		HTTPStatus: http.StatusForbidden,
	}
)

var (
	ErrInvalidInput = &AppError{
		Code:       "INVALID_INPUT",
		Message:    "The provided input is invalid",
		Internal:   "input validation failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrMissingRequired = &AppError{
		Code:       "MISSING_REQUIRED_FIELD",
		Message:    "Required field is missing",
		Internal:   "required field validation failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidFormat = &AppError{
		Code:       "INVALID_FORMAT",
		Message:    "The provided data format is invalid",
		Internal:   "data format validation failed",
		HTTPStatus: http.StatusBadRequest,
	}
)

var (
	ErrDatabaseError = &AppError{
		Code:       "DATABASE_ERROR",
		Message:    "A database error occurred. Please try again later",
		Internal:   "database operation failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrCacheError = &AppError{
		Code:       "CACHE_ERROR",
		Message:    "A caching error occurred",
		Internal:   "redis operation failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrQueueError = &AppError{
		Code:       "QUEUE_ERROR",
		Message:    "A messaging queue error occurred",
		Internal:   "rabbitmq operation failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrExternalService = &AppError{
		Code:       "EXTERNAL_SERVICE_ERROR",
		Message:    "An external service error occurred",
		Internal:   "external API call failed",
		HTTPStatus: http.StatusServiceUnavailable,
	}

	ErrInternalServer = &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An unexpected error occurred. Please try again later",
		Internal:   "internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrServiceUnavailable = &AppError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Service is temporarily unavailable. Please try again later",
		Internal:   "service unavailable",
		HTTPStatus: http.StatusServiceUnavailable,
	}

	ErrTimeout = &AppError{
		Code:       "REQUEST_TIMEOUT",
		Message:    "Request timed out. Please try again",
		Internal:   "operation timeout",
		HTTPStatus: http.StatusRequestTimeout,
	}
)

var (
	ErrRateLimitExceeded = &AppError{
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    "Too many requests. Please try again later",
		Internal:   "rate limit exceeded",
		HTTPStatus: http.StatusTooManyRequests,
	}
)

var (
	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		Internal:   "requested resource not found",
		HTTPStatus: http.StatusNotFound,
	}
)
