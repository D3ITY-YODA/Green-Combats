package httpx

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	if v, exists := c.Get("request_id"); exists {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return uuid.New().String()
}

// HandleError writes a standardized error response
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = &AppError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}
	}
	requestID := GetRequestID(c)
	c.JSON(appErr.Status, Error(appErr.Code, appErr.Message, requestID))
}

// AppError represents an application error with HTTP status
type AppError struct {
	Code    string
	Message string
	Status  int
	Details map[string]interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

// Error constructors
var (
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "unauthorized", Status: http.StatusUnauthorized}
	ErrForbidden    = &AppError{Code: "FORBIDDEN", Message: "forbidden", Status: http.StatusForbidden}
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "resource not found", Status: http.StatusNotFound}
	ErrBadRequest   = &AppError{Code: "BAD_REQUEST", Message: "bad request", Status: http.StatusBadRequest}
	ErrInternal     = &AppError{Code: "INTERNAL_ERROR", Message: "internal server error", Status: http.StatusInternalServerError}
	ErrConflict     = &AppError{Code: "CONFLICT", Message: "resource conflict", Status: http.StatusConflict}
	ErrTooManyRequests = &AppError{Code: "RATE_LIMITED", Message: "too many requests", Status: http.StatusTooManyRequests}
)

func InvalidParam(param, reason string) *AppError {
	return &AppError{
		Code:    "INVALID_PARAMETER",
		Message: param + ": " + reason,
		Status:  http.StatusBadRequest,
		Details: map[string]interface{}{"parameter": param, "reason": reason},
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: message,
		Status:  http.StatusNotFound,
	}
}