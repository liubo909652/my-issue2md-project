# API Sketch - Core Functionality

## Version: 1.0
## Date: 2026-03-31

本文档简要描述核心包对外暴露的主要接口，作为后续开发的参考。

---

## Package: internal/github

### Purpose
负责与 GitHub REST API 的交互，获取 Issue 和 Comment 数据。

### Core Interface

```go
// Client represents a GitHub API client
type Client struct {
    httpClient *http.Client
    baseURL    string
    verbose    bool
}

// NewClient creates a new GitHub API client
// httpClient is optional; if nil, a default client will be created
func NewClient(httpClient *http.Client, verbose bool) *Client

// FetchIssue fetches a single issue by owner, repo, and issue number
// Returns Issue data and any error encountered
func (c *Client) FetchIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)

// FetchComments fetches all comments for an issue
// Handles pagination automatically to retrieve all comments
func (c *Client) FetchComments(ctx context.Context, owner, repo string, number int) ([]*Comment, error)

// FetchAll fetches both the issue and all comments in one operation
// This is the recommended method for typical use cases
func (c *Client) FetchAll(ctx context.Context, owner, repo string, number int) (*Issue, []*Comment, error)
```

### Data Models

```go
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
    Content string // e.g., "+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"
    Count   int
    Users   []string // list of user logins (optional, for verbose output)
}
```

### Error Types

```go
// ErrorType represents the category of an API error
type ErrorType int

const (
    ErrorTypeNetwork ErrorType = iota
    ErrorTypeRateLimit
    ErrorTypeNotFound
    ErrorTypeForbidden
    ErrorTypeServerError
)

// APIError represents an error from the GitHub API
type APIError struct {
    Type    ErrorType
    Message string
    Err     error
    // Additional context for rate limiting
    RateLimit *RateLimitInfo
}

// RateLimitInfo contains rate limiting information
type RateLimitInfo struct {
    Limit     int
    Remaining int
    Reset     time.Time
    RetryAfter time.Duration
}

func (e *APIError) Error() string
func (e *APIError) Unwrap() error
func (e *APIError) IsTemporary() bool // returns true for rate limit and server errors
```

### Configuration Constants

```go
const (
    // DefaultBaseURL is the default GitHub API base URL
    DefaultBaseURL = "https://api.github.com"

    // MaxRetries is the maximum number of retry attempts for failed requests
    MaxRetries = 3

    // InitialBackoff is the initial backoff duration for retries
    InitialBackoff = 1 * time.Second

    // MaxBackoff is the maximum backoff duration for retries
    MaxBackoff = 4 * time.Second

    // DefaultTimeout is the default HTTP request timeout
    DefaultTimeout = 30 * time.Second

    // DefaultPerPage is the default number of items per page for paginated requests
    DefaultPerPage = 100
)
```

---

## Package: internal/converter

### Purpose
将 Issue 和 Comment 数据转换为 Markdown 格式输出。

### Core Interface

```go
// Converter handles conversion of GitHub data to Markdown
type Converter struct {
    // Configuration options can be added here in the future
    includeReactions bool
    includeMetadata  bool
}

// NewConverter creates a new Markdown converter
func NewConverter() *Converter

// Convert converts an issue and its comments to Markdown format
// Returns the complete Markdown document as a string
func (c *Converter) Convert(issue *github.Issue, comments []*github.Comment) string

// ConvertIssue converts only the issue (without comments) to Markdown
func (c *Converter) ConvertIssue(issue *github.Issue) string

// ConvertComment converts a single comment to Markdown
func (c *Converter) ConvertComment(comment *github.Comment) string
```

### Helper Functions

```go
// FormatHeader formats the issue metadata header
func (c *Converter) FormatHeader(issue *github.Issue) string

// FormatBody formats the issue body
func (c *Converter) FormatBody(issue *github.Issue) string

// FormatComment formats a comment section
func (c *Converter) FormatComment(comment *github.Comment) string

// FormatReactions formats reaction information
func (c *Converter) FormatReactions(reactions []github.Reaction) string

// FormatLabels formats label information
func (c *Converter) FormatLabels(labels []github.Label) string

// FormatMilestone formats milestone information
func (c *Converter) FormatMilestone(milestone *github.Milestone) string

// FormatTimestamp formats a time.Time to ISO 8601 format
func FormatTimestamp(t time.Time) string

// EscapeMarkdown escapes special Markdown characters in plain text
func EscapeMarkdown(text string) string
```

### HTML to Markdown Conversion

```go
// HTMLToMarkdown converts HTML content to Markdown
// This is a simplified conversion that handles common HTML elements
func HTMLToMarkdown(html string) string

// HTMLConverter handles HTML to Markdown conversion
type HTMLConverter struct {
    // Internal state for parsing
}

// NewHTMLConverter creates a new HTML converter
func NewHTMLConverter() *HTMLConverter

// Convert converts HTML to Markdown
func (h *HTMLConverter) Convert(html string) string
```

### Formatting Constants

```go
const (
    // HeaderSeparator is the separator between header and body
    HeaderSeparator = "---"

    // CommentSeparator is the separator between comments
    CommentSeparator = "---"

    // MetadataPrefix is the prefix for metadata lines
    MetadataPrefix = "**"

    // MetadataSeparator is the separator between metadata items
    MetadataSeparator = " · "
)
```

---

## Package: internal/parser

### Purpose
解析和验证 GitHub Issue URL，提取关键信息。

### Core Interface

```go
// IssueURL represents a parsed GitHub Issue URL
type IssueURL struct {
    Owner   string
    Repo    string
    Number  int
    Original string // the original URL string
}

// ParseIssueURL parses a GitHub Issue URL and extracts its components
// Returns an error if the URL is invalid
func ParseIssueURL(rawURL string) (*IssueURL, error)

// String returns the full URL string
func (u *IssueURL) String() string

// IsValid checks if the parsed URL is valid
func (u *IssueURL) IsValid() bool

// Validate validates the URL components
// Returns an error if any component is invalid
func (u *IssueURL) Validate() error
```

### Error Types

```go
// ParseError represents a URL parsing error
type ParseError struct {
    URL     string
    Message string
    Err     error
}

func (e *ParseError) Error() string
func (e *ParseError) Unwrap() error
```

### Constants

```go
const (
    // AllowedDomain is the only allowed domain for GitHub URLs
    AllowedDomain = "github.com"

    // IssuePathPattern is the regex pattern for issue URLs
    // Format: https://github.com/{owner}/{repo}/issues/{number}
    IssuePathPattern = `^https?://github\.com/([^/]+)/([^/]+)/issues/(\d+)$`
)
```

---

## Package: internal/cli

### Purpose
处理命令行参数解析和用户交互。

### Core Interface

```go
// Config holds the CLI configuration
type Config struct {
    URL      string // GitHub Issue URL
    Verbose  bool   // Enable verbose logging
    Output   string // Output file path (empty means stdout)
    Version  bool   // Show version information
    Help     bool   // Show help information
}

// ParseArgs parses command-line arguments
// Returns Config and any error encountered
func ParseArgs(args []string) (*Config, error)

// Validate validates the CLI configuration
func (c *Config) Validate() error

// Usage returns the usage message
func Usage() string

// Version returns the version string
func Version() string
```

### Error Types

```go
// CLIError represents a CLI parsing error
type CLIError struct {
    Code    int    // Exit code
    Message string // Error message
}

func (e *CLIError) Error() string
func (e *CLIError) ExitCode() int
```

### Exit Codes

```go
const (
    ExitCodeSuccess       = 0
    ExitCodeGeneralError  = 1
    ExitCodeArgumentError = 2
    ExitCodeRateLimit     = 3
)
```

---

## Package: internal/config

### Purpose
管理应用程序配置。

### Core Interface

```go
// Config holds the application configuration
type Config struct {
    // GitHub API settings
    GitHubBaseURL string
    GitHubToken   string // empty for public repos
    Timeout       time.Duration

    // Retry settings
    MaxRetries    int
    InitialBackoff time.Duration
    MaxBackoff    time.Duration

    // Output settings
    Verbose       bool
}

// Load loads configuration from environment variables and defaults
func Load() *Config

// LoadFromEnv loads configuration from environment variables
// Prefix: ISSUE2MD_
// Examples:
//   ISSUE2MD_GITHUB_TOKEN=ghp_xxx
//   ISSUE2MD_TIMEOUT=30s
//   ISSUE2MD_VERBOSE=true
func LoadFromEnv() *Config

// Validate validates the configuration
func (c *Config) Validate() error
```

### Environment Variables

```go
const (
    EnvPrefix        = "ISSUE2MD_"
    EnvGitHubToken   = "ISSUE2MD_GITHUB_TOKEN"
    EnvTimeout       = "ISSUE2MD_TIMEOUT"
    EnvVerbose       = "ISSUE2MD_VERBOSE"
    EnvMaxRetries    = "ISSUE2MD_MAX_RETRIES"
    EnvBaseURL       = "ISSUE2MD_BASE_URL"
)
```

---

## Cross-Package Interactions

### Typical Usage Flow

```go
// 1. Parse CLI arguments
cliConfig := internal/cli.ParseArgs(os.Args[1:])

// 2. Parse GitHub URL
issueURL, err := internal/parser.ParseIssueURL(cliConfig.URL)
if err != nil {
    log.Fatal(err)
}

// 3. Load configuration
config := internal/config.Load()

// 4. Create GitHub client
githubClient := internal/github.NewClient(nil, cliConfig.Verbose)

// 5. Fetch issue and comments
issue, comments, err := githubClient.FetchAll(
    context.Background(),
    issueURL.Owner,
    issueURL.Repo,
    issueURL.Number,
)
if err != nil {
    log.Fatal(err)
}

// 6. Convert to Markdown
converter := internal/converter.NewConverter()
markdown := converter.Convert(issue, comments)

// 7. Output to stdout or file
if cliConfig.Output != "" {
    os.WriteFile(cliConfig.Output, []byte(markdown), 0644)
} else {
    fmt.Println(markdown)
}
```

### Package Dependencies

```
cmd/issue2md/main.go
    ├── internal/cli (argument parsing)
    ├── internal/parser (URL parsing)
    ├── internal/config (configuration)
    ├── internal/github (API client)
    │       └── uses internal/model (data structures)
    └── internal/converter (Markdown conversion)
            └── uses internal/github.Issue, github.Comment
```

---

## Testing Considerations

### Mock Interfaces

For testing, the following interfaces should be considered:

```go
// GitHubClient is an interface for GitHub API operations
// This allows mocking in tests
type GitHubClient interface {
    FetchAll(ctx context.Context, owner, repo string, number int) (*github.Issue, []*github.Comment, error)
}

// Converter is an interface for conversion operations
type Converter interface {
    Convert(issue *github.Issue, comments []*github.Comment) string
}
```

However, per the constitution, prefer integration tests with real HTTP servers (httptest.Server) over mocks.

---

## Future Extensions

### Potential Additions to `internal/converter`:

```go
// ConverterOptions allows configuration of conversion behavior
type ConverterOptions struct {
    IncludeReactions   bool
    IncludeMetadata    bool
    IncludeAttachments bool
    FormatCodeBlocks   bool
    PreserveHTML       bool
}

// NewConverterWithOptions creates a converter with custom options
func NewConverterWithOptions(opts ConverterOptions) *Converter
```

### Potential Additions to `internal/github`:

```go
// WithToken returns a new client with authentication token
func (c *Client) WithToken(token string) *Client

// WithBaseURL returns a new client with custom base URL (for Enterprise)
func (c *Client) WithBaseURL(baseURL string) *Client
```

---

**Document Status:** Draft
**Last Updated:** 2026-03-31
