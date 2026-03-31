package github

import "time"

// API related constants
const (
	// DefaultBaseURL is the default GitHub API base URL
	DefaultBaseURL = "https://api.github.com"

	// APIVersion is the GitHub API version to use
	APIVersion = "v3"

	// UserAgent is the HTTP User-Agent header
	UserAgent = "issue2md/1.0.0"
)

// HTTP headers
const (
	HeaderUserAgent      = "User-Agent"
	HeaderRateLimit     = "X-RateLimit-Limit"
	HeaderRateRemaining = "X-RateLimit-Remaining"
	HeaderRateReset    = "X-RateLimit-Reset"
	HeaderRetryAfter    = "Retry-After"
	HeaderAccept        = "Accept"
)

// Retry configuration
const (
	// DefaultMaxRetries is the maximum number of retry attempts
	DefaultMaxRetries = 3

	// DefaultInitialBackoff is the initial backoff duration
	DefaultInitialBackoff = 1 * time.Second

	// DefaultMaxBackoff is the maximum backoff duration
	DefaultMaxBackoff = 4 * time.Second

	// MaxJitter is the maximum jitter for backoff
	MaxJitter = 0.25 // 25%
)

// API endpoints
const (
	// IssuesEndpoint is the endpoint for issues
	IssuesEndpoint = "/repos/%s/%s/issues/%d"

	// CommentsEndpoint is the endpoint for comments
	CommentsEndpoint = "/repos/%s/%s/issues/%d/comments"
)

// HTTP client configuration
const (
	// DefaultTimeout is the default HTTP request timeout
	DefaultTimeout = 30 * time.Second

	// MaxIdleConns is the maximum number of idle connections
	MaxIdleConns = 10

	// IdleConnTimeout is the timeout for idle connections
	IdleConnTimeout = 30 * time.Second
)