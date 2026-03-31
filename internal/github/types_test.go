package github

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestIssue_MarshalUnmarshalJSON(t *testing.T) {
	issueTime := time.Date(2026, 3, 30, 14, 23, 45, 0, time.UTC)

	issue := Issue{
		Number:    12345,
		Title:     "Add support for custom error handlers",
		Body:      "It would be great to support custom error handlers...",
		State:     "open",
		HTMLURL:   "https://github.com/golang/go/issues/12345",
		Author:    User{Login: "johnsmith"},
		CreatedAt: issueTime,
		UpdatedAt: time.Date(2026, 3, 31, 9, 12, 33, 0, time.UTC),
		Labels: []Label{
			{Name: "enhancement", Color: "84b6eb"},
			{Name: "api", Color: "bfd4f2"},
		},
		Milestone: &Milestone{
			Title: "v1.5",
			Number: 3,
		},
		Reactions: []Reaction{
			{Content: "+1", Count: 5},
			{Content: "rocket", Count: 2},
		},
	}

	jsonData, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("Failed to marshal issue: %v", err)
	}

	var unmarshaled Issue
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal issue: %v", err)
	}

	if !reflect.DeepEqual(issue, unmarshaled) {
		t.Errorf("Issue data doesn't match after marshal/unmarshal")
	}
}

func TestComment_MarshalUnmarshalJSON(t *testing.T) {
	commentTime := time.Date(2026, 3, 30, 15, 45, 12, 0, time.UTC)

	comment := Comment{
		ID:        1234,
		Body:      "👍 This looks good. I have a few suggestions:",
		HTMLURL:   "https://github.com/golang/go/issues/12345#issuecomment-12345",
		Author:    User{Login: "janedoe"},
		CreatedAt: commentTime,
		UpdatedAt: commentTime,
		Reactions: []Reaction{
			{Content: "thumbs_up", Count: 1},
		},
		CommitSHA: "abc1234",
	}

	jsonData, err := json.Marshal(comment)
	if err != nil {
		t.Fatalf("Failed to marshal comment: %v", err)
	}

	var unmarshaled Comment
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal comment: %v", err)
	}

	if !reflect.DeepEqual(comment, unmarshaled) {
		t.Errorf("Comment data doesn't match after marshal/unmarshal")
	}
}

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    User{Login: "johnsmith", ID: 12345},
			wantErr: false,
		},
		{
			name:    "empty login",
			user:    User{Login: "", ID: 12345},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIssue_Validate(t *testing.T) {
	tests := []struct {
		name    string
		issue   Issue
		wantErr bool
	}{
		{
			name:    "valid issue",
			issue:   Issue{Number: 123, Title: "Test", State: "open", Author: User{Login: "test"}},
			wantErr: false,
		},
		{
			name:    "empty title",
			issue:   Issue{Number: 123, Title: "", State: "open", Author: User{Login: "test"}},
			wantErr: true,
		},
		{
			name:    "invalid state",
			issue:   Issue{Number: 123, Title: "Test", State: "invalid", Author: User{Login: "test"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.issue.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Issue.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}