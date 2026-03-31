package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetIssue(t *testing.T) {
	// Create a mock GitHub API server
	testIssue := Issue{
		Number:    12345,
		Title:     "Add support for custom error handlers",
		Body:      "It would be great to support custom error handlers in the middleware.",
		State:     "open",
		HTMLURL:   "https://github.com/golang/go/issues/12345",
		Author: User{
			Login:     "johnsmith",
			ID:        123456,
			AvatarURL: "https://avatars.githubusercontent.com/u/123456",
			HTMLURL:   "https://github.com/johnsmith",
		},
		CreatedAt: time.Date(2026, 3, 30, 14, 23, 45, 0, time.UTC),
		UpdatedAt: time.Date(2026, 3, 31, 9, 12, 33, 0, time.UTC),
		Labels: []Label{
			{Name: "enhancement", Color: "a2eeef", Description: "New feature or request"},
			{Name: "api", Color: "007acc", Description: "API related"},
		},
		Milestone: &Milestone{
			Title:     "v1.5",
			Number:    10,
			State:     "open",
			HTMLURL:   "https://github.com/golang/go/milestone/10",
		},
		CommentsURL: "https://api.github.com/repos/golang/go/issues/12345/comments",
	}

	tests := []struct {
		name           string
		owner          string
		repo           string
		number         int
		responseStatus int
		responseBody   interface{}
		expectedError  bool
		verifyIssue    func(*testing.T, *Issue)
	}{
		{
			name:           "successful fetch",
			owner:          "golang",
			repo:           "go",
			number:         12345,
			responseStatus: http.StatusOK,
			responseBody:   testIssue,
			expectedError:  false,
			verifyIssue: func(t *testing.T, issue *Issue) {
				if issue.Number != 12345 {
					t.Errorf("Number = %d, want %d", issue.Number, 12345)
				}
				if issue.Title != "Add support for custom error handlers" {
					t.Errorf("Title = %s, want %s", issue.Title, "Add support for custom error handlers")
				}
				if issue.Author.Login != "johnsmith" {
					t.Errorf("Author.Login = %s, want %s", issue.Author.Login, "johnsmith")
				}
				if len(issue.Labels) != 2 {
					t.Errorf("Labels count = %d, want %d", len(issue.Labels), 2)
				}
				if issue.Milestone == nil {
					t.Error("Milestone = nil, want non-nil")
				}
				if issue.Milestone.Title != "v1.5" {
					t.Errorf("Milestone.Title = %s, want %s", issue.Milestone.Title, "v1.5")
				}
			},
		},
		{
			name:           "not found",
			owner:          "golang",
			repo:           "go",
			number:         99999,
			responseStatus: http.StatusNotFound,
			responseBody:   map[string]string{"message": "Not Found"},
			expectedError:  true,
		},
		{
			name:           "forbidden",
			owner:          "private",
			repo:           "private",
			number:         1,
			responseStatus: http.StatusForbidden,
			responseBody:   map[string]string{"message": "Forbidden"},
			expectedError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expectedPath := createIssuesEndpoint(tc.owner, tc.repo, tc.number)

			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request path
				if r.URL.Path != expectedPath {
					t.Errorf("unexpected request path: got %s, want %s", r.URL.Path, expectedPath)
				}

				// Verify User-Agent header
				if ua := r.Header.Get("User-Agent"); ua != UserAgent {
					t.Errorf("unexpected User-Agent: got %s, want %s", ua, UserAgent)
				}

				// Verify Accept header
				if accept := r.Header.Get("Accept"); accept != "application/vnd.github.v3+json" {
					t.Errorf("unexpected Accept header: got %s, want %s", accept, "application/vnd.github.v3+json")
				}

				// Return mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				if tc.responseBody != nil {
					json.NewEncoder(w).Encode(tc.responseBody)
				}
			}))
			defer server.Close()

			// Create client with mock server
			client := NewClient(
				WithBaseURL(server.URL),
			)

			// Call GetIssue
			ctx := context.Background()
			issue, err := client.GetIssue(ctx, tc.owner, tc.repo, tc.number)

			if tc.expectedError {
				if err == nil {
					t.Errorf("GetIssue() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetIssue() unexpected error: %v", err)
				return
			}

			// Verify issue
			if tc.verifyIssue != nil {
				tc.verifyIssue(t, issue)
			}
		})
	}
}

// Helper function to create endpoint path
func createIssuesEndpoint(owner, repo string, number int) string {
	return "/repos/" + owner + "/" + repo + "/issues/" + itoa(number)
}

// Simple integer to string conversion
func itoa(n int) string {
	digits := make([]byte, 0, 10)
	if n == 0 {
		return "0"
	}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
