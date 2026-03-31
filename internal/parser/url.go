package parser

import (
	"errors"
	"fmt"
)

// IssueURL represents a parsed GitHub Issue URL
type IssueURL struct {
	Owner   string
	Repo    string
	Number  int
	Original string
}

// Validate validates the URL components
func (u *IssueURL) Validate() error {
	if u.Owner == "" {
		return errors.New("owner cannot be empty")
	}

	if u.Repo == "" {
		return errors.New("repo cannot be empty")
	}

	if u.Number <= 0 {
		return errors.New("issue number must be positive")
	}

	return nil
}

// String returns the full URL string
func (u *IssueURL) String() string {
	return fmt.Sprintf("https://github.com/%s/%s/issues/%d", u.Owner, u.Repo, u.Number)
}

// ParseIssueURL parses a GitHub Issue URL and extracts its components
func ParseIssueURL(rawURL string) (*IssueURL, error) {
	// Reuse the generic Parse function
	parsed, err := Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Verify it's an issue URL
	if parsed.Type != ParsedURLTypeIssue {
		return nil, &ParseError{Reason: fmt.Sprintf("URL is not an issue URL, got type: %v", parsed.Type)}
	}

	return &IssueURL{
		Owner:    parsed.Owner,
		Repo:     parsed.Repo,
		Number:   parsed.Number,
		Original: parsed.Original,
	}, nil
}