package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError represents a standard application error
type AppError struct {
	Code       string
	HTTPStatus int
	Message    string
	Details    map[string]interface{}
}

var (
	ErrUnauthorized  = &AppError{Code: "unauthorized", HTTPStatus: http.StatusUnauthorized, Message: "authentication required"}
	ErrForbidden     = &AppError{Code: "forbidden", HTTPStatus: http.StatusForbidden, Message: "access denied"}
	ErrNotFound      = &AppError{Code: "not_found", HTTPStatus: http.StatusNotFound, Message: "resource not found"}
	ErrBadRequest    = &AppError{Code: "bad_request", HTTPStatus: http.StatusBadRequest, Message: "invalid request"}
	ErrConflict      = &AppError{Code: "conflict", HTTPStatus: http.StatusConflict, Message: "conflict with existing resource"}
	ErrInternal      = &AppError{Code: "internal_error", HTTPStatus: http.StatusInternalServerError, Message: "internal server error"}
)

func (e *AppError) Error() string { return e.Message }

// NotFound returns a 404 error with a custom message
func NotFound(message string) *AppError {
	return &AppError{
		Code:       "not_found",
		HTTPStatus: http.StatusNotFound,
		Message:    message,
	}
}

// HandleError responds with a standardized error envelope
func HandleError(c *gin.Context, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, ErrorWithDetails(appErr.Code, appErr.Message, appErr.Details))
		return
	}

	// Default to internal error for unknown error types
	c.JSON(http.StatusInternalServerError, Error(ErrInternal.Code, ErrInternal.Message))
}

// MissingParam returns an error for missing query/form parameter
func MissingParam(param string) *AppError {
	err := &AppError{
		Code:       "missing_parameter",
		HTTPStatus: http.StatusBadRequest,
		Message:    "missing required parameter",
		Details:    make(map[string]interface{}),
	}
	err.Details["parameter"] = param
	return err
}

// InvalidParam returns an error for invalid parameter value
func InvalidParam(param, reason string) *AppError {
	err := &AppError{
		Code:       "invalid_parameter",
		HTTPStatus: http.StatusBadRequest,
		Message:    "invalid parameter value",
		Details:    make(map[string]interface{}),
	}
	err.Details["parameter"] = param
	if reason != "" {
		err.Details["reason"] = reason
	}
	return err
}
