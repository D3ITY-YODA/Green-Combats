package httpx

import (
	"time"
)

type Envelope struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Meta   *Meta       `json:"meta,omitempty"`
	Error  *ErrorDetail `json:"error,omitempty"`
}

type Meta struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
}

type PaginationMeta struct {
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	HasNext bool `json:"has_next"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Success returns a successful response envelope
func Success(data interface{}) *Envelope {
	return &Envelope{
		Status: "success",
		Data:   data,
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
	}
}

// SuccessWithPagination returns a paginated response envelope
func SuccessWithPagination(data interface{}, pagination *PaginationMeta) *Envelope {
	return &Envelope{
		Status: "success",
		Data:   data,
		Meta: &Meta{
			Pagination: pagination,
			Timestamp:  time.Now().UTC(),
		},
	}
}

// Error returns an error response envelope
func Error(code, message string) *Envelope {
	return &Envelope{
		Status: "error",
		Data:   nil,
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// ErrorWithDetails returns an error response with additional context
func ErrorWithDetails(code, message string, details map[string]interface{}) *Envelope {
	return &Envelope{
		Status: "error",
		Data:   nil,
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
