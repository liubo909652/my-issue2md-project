package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is a GitHub API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new GitHub API client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{},
		baseURL:    DefaultBaseURL,
		userAgent:  UserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetIssue fetches an issue from GitHub
func (c *Client) GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	// Build the API endpoint URL
	endpoint := fmt.Sprintf(IssuesEndpoint, owner, repo, number)
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set(HeaderAccept, "application/vnd.github.v3+json")
	req.Header.Set(HeaderUserAgent, c.userAgent)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if err := CheckResponseStatusCode(resp); err != nil {
		return nil, err
	}

	// Decode JSON response
	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}
