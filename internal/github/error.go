package github

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ErrorType represents the category of an API error
type ErrorType int

const (
	// ErrorTypeNetwork represents network-related errors
	ErrorTypeNetwork ErrorType = iota

	// ErrorTypeRateLimit represents rate limit errors
	ErrorTypeRateLimit

	// ErrorTypeNotFound represents 404 errors
	ErrorTypeNotFound

	// ErrorTypeForbidden represents 403 errors
	ErrorTypeForbidden

	// ErrorTypeServerError represents 5xx server errors
	ErrorTypeServerError

	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation

	// ErrorTypeUnknown represents unknown errors
	ErrorTypeUnknown
)

// APIError represents an error from the GitHub API
type APIError struct {
	Type    ErrorType
	Message string
	Err     error

	// Additional context for rate limiting
	StatusCode int
	RateLimit  *RateLimitInfo
}

// RateLimitInfo contains rate limiting information
type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     time.Time
	RetryAfter time.Duration
}

// Error returns the error message
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *APIError) Unwrap() error {
	return e.Err
}

// IsTemporary returns true for temporary errors that can be retried
func (e *APIError) IsTemporary() bool {
	switch e.Type {
	case ErrorTypeNetwork, ErrorTypeRateLimit, ErrorTypeServerError:
		return true
	default:
		return false
	}
}

// IsPermanent returns true for errors that should not be retried
func (e *APIError) IsPermanent() bool {
	switch e.Type {
	case ErrorTypeNotFound, ErrorTypeForbidden:
		return true
	default:
		return false
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(err error) *APIError {
	return &APIError{
		Type:    ErrorTypeNetwork,
		Message: "network error occurred",
		Err:     err,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(remaining int, reset time.Time, retryAfter time.Duration) *APIError {
	return &APIError{
		Type:    ErrorTypeRateLimit,
		Message: "GitHub API rate limit exceeded",
		RateLimit: &RateLimitInfo{
			Remaining: remaining,
			Reset:     reset,
			RetryAfter: retryAfter,
		},
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(owner, repo string, number int) *APIError {
	return &APIError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("issue not found: %s/%s#%d", owner, repo, number),
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(owner, repo string, number int) *APIError {
	return &APIError{
		Type:    ErrorTypeForbidden,
		Message: fmt.Sprintf("repository is private or access denied: %s/%s#%d", owner, repo, number),
	}
}

// NewServerError creates a new server error
func NewServerError(statusCode int) *APIError {
	return &APIError{
		Type:       ErrorTypeServerError,
		Message:    fmt.Sprintf("GitHub server returned %d", statusCode),
		StatusCode: statusCode,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *APIError {
	return &APIError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// CheckResponseStatusCode checks the HTTP response status code and returns appropriate errors
func CheckResponseStatusCode(resp *http.Response) error {
	if resp == nil {
		return NewValidationError("nil response")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return NewNotFoundError("", "", 0) // Owner/repo/number should be added by caller
	case http.StatusForbidden:
		return NewForbiddenError("", "", 0) // Owner/repo/number should be added by caller
	case http.StatusTooManyRequests:
		// Parse rate limit headers if available
		remaining := parseHeaderInt(resp.Header, HeaderRateRemaining)
		reset := parseHeaderTime(resp.Header, HeaderRateReset)
		retryAfter := parseHeaderDuration(resp.Header, HeaderRetryAfter)
		return NewRateLimitError(remaining, reset, retryAfter)
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return NewServerError(resp.StatusCode)
	default:
		return NewServerError(resp.StatusCode)
	}
}

// Helper functions to parse headers
func parseHeaderInt(headers http.Header, key string) int {
	if values := headers[key]; len(values) > 0 {
		if i, err := strconv.Atoi(values[0]); err == nil {
			return i
		}
	}
	return 0
}

func parseHeaderTime(headers http.Header, key string) time.Time {
	if values := headers[key]; len(values) > 0 {
		if i, err := strconv.ParseInt(values[0], 10, 64); err == nil {
			return time.Unix(i, 0)
		}
	}
	return time.Time{}
}

func parseHeaderDuration(headers http.Header, key string) time.Duration {
	if values := headers[key]; len(values) > 0 {
		if d, err := time.ParseDuration(values[0]); err == nil {
			return d
		}
	}
	return 0
}