package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParsedURLType represents the type of a parsed GitHub URL
type ParsedURLType int

const (
	ParsedURLTypeUnknown ParsedURLType = iota
	ParsedURLTypeIssue
	ParsedURLTypePullRequest
	ParsedURLTypeDiscussion
)

// ParsedURL represents a parsed GitHub URL with its components
type ParsedURL struct {
	Owner    string
	Repo     string
	Number   int
	Type     ParsedURLType
	Original string
}

// Validate validates the parsed URL components
func (u *ParsedURL) Validate() error {
	if u.Owner == "" {
		return &ParseError{Reason: "owner cannot be empty"}
	}

	if u.Repo == "" {
		return &ParseError{Reason: "repo cannot be empty"}
	}

	if u.Number <= 0 {
		return &ParseError{Reason: "number must be positive"}
	}

	if u.Type == ParsedURLTypeUnknown {
		return &ParseError{Reason: "URL type must be specified"}
	}

	return nil
}

// ParseError represents a URL parsing error
type ParseError struct {
	Reason string
}

func (e *ParseError) Error() string {
	return "parse error: " + e.Reason
}

// Parse parses a GitHub URL and extracts its components
// Supports Issue URLs (/issues/{number}), Pull Request URLs (/pull/{number}),
// and Discussion URLs (/discussions/{number})
func Parse(rawURL string) (*ParsedURL, error) {
	if rawURL == "" {
		return nil, &ParseError{Reason: "empty URL"}
	}

	// Parse the URL using net/url
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, &ParseError{Reason: fmt.Sprintf("invalid URL format: %v", err)}
	}

	// Validate scheme
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, &ParseError{Reason: fmt.Sprintf("URL scheme must be http or https, got: %s", parsed.Scheme)}
	}

	// Validate host (must be github.com or www.github.com)
	host := strings.ToLower(parsed.Host)
	if host != "github.com" && host != "www.github.com" {
		return nil, &ParseError{Reason: fmt.Sprintf("URL must be from github.com, got: %s", parsed.Host)}
	}

	// Split path into segments
	path := parsed.Path
	if path == "" {
		return nil, &ParseError{Reason: "URL path is empty"}
	}

	// Remove leading slash and split
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(segments) < 3 {
		return nil, &ParseError{Reason: fmt.Sprintf("URL must contain at least owner/repo/type/number, got: %s", rawURL)}
	}

	owner := segments[0]
	repo := segments[1]
	urlType := segments[2]

	if owner == "" || repo == "" {
		return nil, &ParseError{Reason: fmt.Sprintf("owner and repo cannot be empty, got: %s/%s", owner, repo)}
	}

	// Determine URL type and extract number
	var numType ParsedURLType
	var numberStr string

	switch urlType {
	case "issues":
		numType = ParsedURLTypeIssue
	case "pull":
		numType = ParsedURLTypePullRequest
	case "discussions":
		numType = ParsedURLTypeDiscussion
	default:
		return nil, &ParseError{Reason: fmt.Sprintf("unsupported URL type: %s (must be issues, pull, or discussions)", urlType)}
	}

	if len(segments) < 4 {
		return nil, &ParseError{Reason: fmt.Sprintf("URL must contain a number: %s", rawURL)}
	}
	numberStr = segments[3]

	// Parse number
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, &ParseError{Reason: fmt.Sprintf("invalid number: %s", numberStr)}
	}

	if number <= 0 {
		return nil, &ParseError{Reason: fmt.Sprintf("number must be positive, got: %d", number)}
	}

	return &ParsedURL{
		Owner:    owner,
		Repo:     repo,
		Number:   number,
		Type:     numType,
		Original: rawURL,
	}, nil
}
