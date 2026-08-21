package httpx

import (
	"time"
	"github.com/google/uuid"
)

type Envelope struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta,omitempty"`
	Error *ErrorDetail `json:"error,omitempty"`
}

type Meta struct {
	RequestID  string           `json:"request_id"`
	Pagination *PaginationMeta  `json:"pagination,omitempty"`
	Timestamp  time.Time        `json:"generated_at"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type ErrorDetail struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// Success returns a successful response envelope
func Success(data interface{}, requestID string) *Envelope {
	return &Envelope{
		Data: data,
		Meta: &Meta{
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
	}
}

// SuccessWithPagination returns a paginated response envelope
func SuccessWithPagination(data interface{}, pagination *PaginationMeta, requestID string) *Envelope {
	return &Envelope{
		Data: data,
		Meta: &Meta{
			RequestID:  requestID,
			Pagination: pagination,
			Timestamp:  time.Now().UTC(),
		},
	}
}

// Error returns an error response envelope
func Error(code, message, requestID string) *Envelope {
	return &Envelope{
		Data: nil,
		Meta: &Meta{
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
		Error: &ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}
}

// ErrorWithDetails returns an error response with additional context
func ErrorWithDetails(code, message, requestID string, details map[string]interface{}) *Envelope {
	return &Envelope{
		Data: nil,
		Meta: &Meta{
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
		Error: &ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
			Details:   details,
		},
	}
}
