package github

import (
	"errors"
	"fmt"
	"time"
)

// Issue represents a GitHub issue
type Issue struct {
	Number      int
	Title       string
	Body        string
	State       string // "open" or "closed"
	HTMLURL     string
	Author      User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    *time.Time
	Labels      []Label
	Milestone   *Milestone
	Reactions   []Reaction
	CommentsURL string
}

// Comment represents a comment on a GitHub issue
type Comment struct {
	ID        int
	Body      string
	HTMLURL   string
	Author    User
	CreatedAt time.Time
	UpdatedAt time.Time
	Reactions []Reaction
	CommitSHA string // optional, for PR review comments
}

// User represents a GitHub user
type User struct {
	Login     string
	ID        int64
	AvatarURL string
	HTMLURL   string
}

// Label represents an issue label
type Label struct {
	Name        string
	Color       string
	Description string
}

// Milestone represents an issue milestone
type Milestone struct {
	Title       string
	Number      int
	State       string
	Description string
	HTMLURL     string
}

// Reaction represents a reaction on an issue or comment
type Reaction struct {
	Content string // "+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"
	Count   int
	Users   []string // list of user logins (optional, for verbose output)
}

// Validate validates the Issue struct
func (i *Issue) Validate() error {
	if i.Title == "" {
		return errors.New("issue title cannot be empty")
	}

	if i.State != "open" && i.State != "closed" {
		return fmt.Errorf("invalid issue state: %s", i.State)
	}

	if i.Author.Validate() != nil {
		return fmt.Errorf("invalid issue author: %w", i.Author.Validate())
	}

	return nil
}

// Validate validates the User struct
func (u *User) Validate() error {
	if u.Login == "" {
		return errors.New("user login cannot be empty")
	}
	return nil
}

// Validate validates the Comment struct
func (c *Comment) Validate() error {
	if c.Body == "" {
		return errors.New("comment body cannot be empty")
	}

	if c.Author.Validate() != nil {
		return fmt.Errorf("invalid comment author: %w", c.Author.Validate())
	}

	return nil
}

// Validate validates the Label struct
func (l *Label) Validate() error {
	if l.Name == "" {
		return errors.New("label name cannot be empty")
	}
	return nil
}

// Validate validates the Milestone struct
func (m *Milestone) Validate() error {
	if m.Title == "" {
		return errors.New("milestone title cannot be empty")
	}

	if m.State != "" && m.State != "open" && m.State != "closed" {
		return fmt.Errorf("invalid milestone state: %s", m.State)
	}

	return nil
}

// Validate validates the Reaction struct
func (r *Reaction) Validate() error {
	validReactions := []string{"+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"}

	for _, valid := range validReactions {
		if r.Content == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid reaction content: %s", r.Content)
}