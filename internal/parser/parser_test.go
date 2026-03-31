package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *ParsedURL
		wantErr bool
	}{
		{
			name: "valid Issue URL",
			url:  "https://github.com/golang/go/issues/12345",
			want: &ParsedURL{
				Owner:    "golang",
				Repo:     "go",
				Number:   12345,
				Type:     ParsedURLTypeIssue,
				Original: "https://github.com/golang/go/issues/12345",
			},
			wantErr: false,
		},
		{
			name: "valid Pull Request URL",
			url:  "https://github.com/golang/go/pull/67890",
			want: &ParsedURL{
				Owner:    "golang",
				Repo:     "go",
				Number:   67890,
				Type:     ParsedURLTypePullRequest,
				Original: "https://github.com/golang/go/pull/67890",
			},
			wantErr: false,
		},
		{
			name: "valid Discussion URL",
			url:  "https://github.com/golang/go/discussions/999",
			want: &ParsedURL{
				Owner:    "golang",
				Repo:     "go",
				Number:   999,
				Type:     ParsedURLTypeDiscussion,
				Original: "https://github.com/golang/go/discussions/999",
			},
			wantErr: false,
		},
		{
			name:    "invalid URL - empty string",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL - malformed",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "invalid URL - missing components",
			url:     "https://github.com/golang/go",
			wantErr: true,
		},
		{
			name:    "invalid URL - not github.com domain",
			url:     "https://gitlab.com/owner/repo/issues/123",
			wantErr: true,
		},
		{
			name:    "unsupported URL - repository home",
			url:     "https://github.com/golang/go",
			wantErr: true,
		},
		{
			name:    "unsupported URL - releases page",
			url:     "https://github.com/golang/go/releases",
			wantErr: true,
		},
		{
			name:    "unsupported URL - actions page",
			url:     "https://github.com/golang/go/actions",
			wantErr: true,
		},
		{
			name: "valid Issue URL with query params",
			url:  "https://github.com/owner/repo/issues/123?state=open",
			want: &ParsedURL{
				Owner:    "owner",
				Repo:     "repo",
				Number:   123,
				Type:     ParsedURLTypeIssue,
				Original: "https://github.com/owner/repo/issues/123?state=open",
			},
			wantErr: false,
		},
		{
			name: "valid Issue URL with fragment",
			url:  "https://github.com/owner/repo/issues/123#issuecomment-456",
			want: &ParsedURL{
				Owner:    "owner",
				Repo:     "repo",
				Number:   123,
				Type:     ParsedURLTypeIssue,
				Original: "https://github.com/owner/repo/issues/123#issuecomment-456",
			},
			wantErr: false,
		},
		{
			name: "valid URL with www subdomain",
			url:  "https://www.github.com/golang/go/issues/12345",
			want: &ParsedURL{
				Owner:    "golang",
				Repo:     "go",
				Number:   12345,
				Type:     ParsedURLTypeIssue,
				Original: "https://www.github.com/golang/go/issues/12345",
			},
			wantErr: false,
		},
		{
			name: "valid Issue URL with HTTP",
			url:  "http://github.com/golang/go/issues/12345",
			want: &ParsedURL{
				Owner:    "golang",
				Repo:     "go",
				Number:   12345,
				Type:     ParsedURLTypeIssue,
				Original: "http://github.com/golang/go/issues/12345",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil {
				t.Errorf("Parse() returned nil, want %v", tt.want)
				return
			}

			if got.Owner != tt.want.Owner {
				t.Errorf("Parse() Owner = %v, want %v", got.Owner, tt.want.Owner)
			}
			if got.Repo != tt.want.Repo {
				t.Errorf("Parse() Repo = %v, want %v", got.Repo, tt.want.Repo)
			}
			if got.Number != tt.want.Number {
				t.Errorf("Parse() Number = %v, want %v", got.Number, tt.want.Number)
			}
			if got.Type != tt.want.Type {
				t.Errorf("Parse() Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Original != tt.want.Original {
				t.Errorf("Parse() Original = %v, want %v", got.Original, tt.want.Original)
			}
		})
	}
}
