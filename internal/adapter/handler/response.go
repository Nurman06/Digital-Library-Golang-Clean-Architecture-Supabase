package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response represents a standard API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorData  `json:"error,omitempty"`
	Meta      *MetaData   `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorData represents error information in the response
type ErrorData struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// MetaData represents pagination metadata
type MetaData struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SuccessResponse writes a successful JSON response
func SuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response := Response{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
	WriteJSON(w, statusCode, response)
}

// SuccessResponseWithMeta writes a successful JSON response with pagination metadata
func SuccessResponseWithMeta(w http.ResponseWriter, statusCode int, data interface{}, meta *MetaData) {
	response := Response{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}
	WriteJSON(w, statusCode, response)
}

// ErrorResponse writes an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	response := Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	}
	WriteJSON(w, statusCode, response)
}

// ErrorResponseWithDetails writes an error JSON response with details
func ErrorResponseWithDetails(w http.ResponseWriter, statusCode int, code, message string, details map[string]interface{}) {
	response := Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
	WriteJSON(w, statusCode, response)
}

// Error codes
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeBadRequest    = "BAD_REQUEST"
	ErrCodeUnavailable   = "UNAVAILABLE"
)

// CalculateTotalPages calculates total pages based on total items and page size
func CalculateTotalPages(total int64, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return totalPages
}