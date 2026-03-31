package parser

import (
	"reflect"
	"testing"
)

func TestParseIssueURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *IssueURL
		wantErr bool
	}{
		{
			name: "valid URL",
			url:  "https://github.com/golang/go/issues/12345",
			want: &IssueURL{Owner: "golang", Repo: "go", Number: 12345, Original: "https://github.com/golang/go/issues/12345"},
			wantErr: false,
		},
		{
			name: "valid URL with www",
			url:  "https://www.github.com/golang/go/issues/12345",
			want: &IssueURL{Owner: "golang", Repo: "go", Number: 12345, Original: "https://www.github.com/golang/go/issues/12345"},
			wantErr: false,
		},
		{
			name:    "invalid URL - missing issue number",
			url:     "https://github.com/golang/go/issues",
			wantErr: true,
		},
		{
			name:    "invalid URL - missing components",
			url:     "https://github.com/golang/go",
			wantErr: true,
		},
		{
			name:    "invalid URL - not github.com",
			url:     "https://gitlab.com/owner/repo/issues/123",
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong path",
			url:     "https://github.com/owner/repo/pull/123",
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong path",
			url:     "https://github.com/owner/repo/discussions/123",
			wantErr: true,
		},
		{
			name: "valid URL with query params",
			url:  "https://github.com/owner/repo/issues/123?state=open",
			want: &IssueURL{Owner: "owner", Repo: "repo", Number: 123, Original: "https://github.com/owner/repo/issues/123?state=open"},
			wantErr: false,
		},
		{
			name: "valid URL with query params 2",
			url:  "https://github.com/owner/repo/issues/12345?filter=all",
			want: &IssueURL{Owner: "owner", Repo: "repo", Number: 12345, Original: "https://github.com/owner/repo/issues/12345?filter=all"},
			wantErr: false,
		},
		{
			name: "valid URL with fragment",
			url:  "https://github.com/owner/repo/issues/123#issuecomment-456",
			want: &IssueURL{Owner: "owner", Repo: "repo", Number: 123, Original: "https://github.com/owner/repo/issues/123#issuecomment-456"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIssueURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIssueURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIssueURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueURL_Validate(t *testing.T) {
	tests := []struct {
		name    string
		url     *IssueURL
		wantErr bool
	}{
		{
			name: "valid URL",
			url: &IssueURL{Owner: "golang", Repo: "go", Number: 12345},
			wantErr: false,
		},
		{
			name:    "empty owner",
			url:     &IssueURL{Owner: "", Repo: "go", Number: 12345},
			wantErr: true,
		},
		{
			name:    "empty repo",
			url:     &IssueURL{Owner: "golang", Repo: "", Number: 12345},
			wantErr: true,
		},
		{
			name:    "negative number",
			url:     &IssueURL{Owner: "golang", Repo: "go", Number: -1},
			wantErr: true,
		},
		{
			name:    "zero number",
			url:     &IssueURL{Owner: "golang", Repo: "go", Number: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.url.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueURL.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIssueURL_String(t *testing.T) {
	url := IssueURL{
		Owner:   "golang",
		Repo:    "go",
		Number:  12345,
		Original: "https://github.com/golang/go/issues/12345",
	}

	expected := "https://github.com/golang/go/issues/12345"
	result := url.String()

	if result != expected {
		t.Errorf("IssueURL.String() = %v, want %v", result, expected)
	}
}