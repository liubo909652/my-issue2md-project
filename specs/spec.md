# issue2md - GitHub Issue to Markdown Converter

## Version: 1.0
## Date: 2026-03-31

---

## 1. Overview

### 1.1 Purpose
`issue2md` is a command-line tool that converts GitHub Issues into Markdown format. It fetches issue content and all comments from public GitHub repositories and presents them in a well-structured, time-ordered Markdown document.

### 1.1.1 Primary Use Cases
- **Documentation archiving**: Export GitHub Issues for offline documentation
- **Backup and archival**: Create permanent records of important discussions
- **Migration support**: Prepare issue content for migration to other systems
- **Code review documentation**: Export issue discussions for external review processes

### 1.2 Scope - In Scope
- GitHub Issues from public repositories only
- Issue title, body, and metadata
- All comments including collapsed/hidden ones
- GitHub-specific syntax preservation
- HTML to Markdown conversion
- Time-ordered output with full metadata
- Network retry mechanism

### 1.3 Scope - Out of Scope
- Pull Requests (not supported in initial version)
- GitHub Discussions (not supported in initial version)
- Private repositories (require authentication)
- Attachment downloading (preserve URLs only)
- GitHub Enterprise instances
- Streaming output (fetch all first)
- Interactive prompts
- Output formatting options (single format only)

### 1.4 Target Users
- Open source maintainers
- Documentation specialists
- System administrators
- DevOps engineers
- Anyone who needs to preserve GitHub Issue content

---

## 2. Requirements

### 2.1 Functional Requirements

#### FR-1: Command-Line Interface
- **FR-1.1**: Accept a single GitHub Issue URL as positional argument
- **FR-1.2**: Support optional `--verbose` flag for detailed logging
- **FR-1.3**: Output converted Markdown to stdout (not files)
- **FR-1.4**: User can redirect output to file using shell redirection

#### FR-2: GitHub Issue Parsing
- **FR-2.1**: Parse GitHub Issue URLs in the format:
  ```
  https://github.com/{owner}/{repo}/issues/{number}
  ```
- **FR-2.2**: Extract owner, repository name, and issue number
- **FR-2.3**: Validate URL format and reject invalid URLs
- **FR-2.4**: Only accept `github.com` domain (not Enterprise instances)

#### FR-3: Content Fetching
- **FR-3.1**: Fetch issue metadata (title, body, status, labels, milestone)
- **FR-3.2**: Fetch ALL comments including collapsed/hidden ones
- **FR-3.3**: Use GitHub REST API for public repositories (no authentication)
- **FR-3.4**: Fetch all content before output (not streaming)
- **FR-3.5**: Implement network retry mechanism with exponential backoff

#### FR-4: Content Conversion
- **FR-4.1**: Preserve GitHub-specific syntax:
  - `@mentions` (e.g., @username)
  - Emoji (e.g., :thumbs_up:)
  - Issue references (e.g., #123)
  - Commit references (e.g., abc1234)
- **FR-4.2**: Convert HTML content to Markdown syntax
- **FR-4.3**: Preserve attachment URLs as-is (no downloading)
- **FR-4.4**: Maintain code blocks and inline code formatting
- **FR-4.5**: Preserve links and references

#### FR-5: Output Formatting
- **FR-5.1**: Generate time-ordered linear display of all content
- **FR-5.2**: Include comprehensive metadata header:
  - Issue title and number
  - Issue status (open/closed)
  - Labels (if any)
  - Milestone (if any)
  - Creation timestamp
  - Last updated timestamp
  - Author information
- **FR-5.3**: For each comment, include:
  - Author username
  - Timestamp
  - Associated commit (if referenced)
  - Reaction counts (if any)
- **FR-5.4**: Follow Markdown CommonMark specification

#### FR-6: Error Handling
- **FR-6.1**: Exit with non-zero status on errors
- **FR-6.2**: Provide clear, actionable error messages
- **FR-6.3**: Handle network failures gracefully with retry
- **FR-6.4**: Validate URLs before making API calls
- **FR-6.5**: Handle API rate limiting appropriately

### 2.2 Non-Functional Requirements

#### NFR-1: Performance
- **NFR-1.1**: Fetch complete issue within 30 seconds (typical issue with <100 comments)
- **NFR-1.2**: Support issues with up to 1000 comments
- **NFR-1.3**: Memory usage should not grow linearly with comment count (streaming processing)

#### NFR-2: Reliability
- **NFR-2.1**: Implement retry mechanism with up to 3 attempts
- **NFR-2.2**: Use exponential backoff for retries (1s, 2s, 4s)
- **NFR-2.3**: Handle network timeouts (default 30s)
- **NFR-2.4**: Gracefully handle GitHub API rate limits (429 responses)

#### NFR-3: Usability
- **NFR-3.1**: Clear, concise help message on `--help` flag
- **NFR-3.2**: Verbose logging should not interfere with Markdown output
- **NFR-3.3**: Error messages should guide users to resolution
- **NFR-3.4**: Output should be human-readable and machine-parseable

#### NFR-4: Maintainability
- **NFR-4.1**: Code must follow Go best practices and idioms
- **NFR-4.2**: Comprehensive test coverage (>80%)
- **NFR-4.3**: Clear separation of concerns (CLI, API client, converter)
- **NFR-4.4**: Minimal external dependencies (prefer standard library)

#### NFR-5: Compatibility
- **NFR-5.1**: Go version >= 1.24
- **NFR-5.2**: Support macOS, Linux, and Windows
- **NFR-5.3**: Output compatible with CommonMark-compliant Markdown parsers
- **NFR-5.4**: Handle Unicode characters correctly

---

## 3. CLI Interface Specification

### 3.1 Command Syntax

```
issue2md [OPTIONS] <GITHUB_ISSUE_URL>
```

### 3.2 Arguments

| Position | Name | Type | Required | Description |
|----------|------|------|----------|-------------|
| 1 | GITHUB_ISSUE_URL | string | Yes | Full GitHub Issue URL (https://github.com/{owner}/{repo}/issues/{number}) |

### 3.3 Options

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--verbose` | boolean | false | Enable detailed logging to stderr |
| `--help` | boolean | false | Display usage information |
| `--version` | boolean | false | Display version information |

### 3.4 Usage Examples

**Basic usage:**
```bash
issue2md https://github.com/golang/go/issues/12345
```

**Save to file:**
```bash
issue2md https://github.com/golang/go/issues/12345 > issue.md
```

**Verbose mode:**
```bash
issue2md --verbose https://github.com/golang/go/issues/12345 > issue.md 2> log.txt
```

**Display help:**
```bash
issue2md --help
```

### 3.5 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (invalid URL, network failure, API error) |
| 2 | Command-line argument error |
| 3 | Rate limit exceeded (after retries) |

### 3.6 Error Messages

| Scenario | Error Message |
|----------|---------------|
| Invalid URL format | `error: invalid GitHub Issue URL: {url}. Expected format: https://github.com/{owner}/{repo}/issues/{number}` |
| Issue not found | `error: issue not found or is private: {owner}/{repo}#{number}` |
| Network failure (after retries) | `error: failed to fetch issue after 3 attempts: {details}` |
| Rate limit exceeded | `error: GitHub API rate limit exceeded. Please wait before retrying.` |
| Invalid repository | `error: repository {owner}/{repo} not found or is private` |

---

## 4. Output Format Specification

### 4.1 Overall Structure

```
# {Issue Title}

[Issue Metadata Header]

## Issue Body

{Issue body content in Markdown}

---

{Comments section (if any)}
```

### 4.2 Issue Metadata Header

```
**Issue:** {owner}/{repo}#{number} · **Status:** {Open/Closed} · **Author:** @{author} · **Created:** {ISO 8601 timestamp}
**Updated:** {ISO 8601 timestamp}
{Labels line (if any)}
{Milestone line (if any)}
```

### 4.3 Issue Body Section

```
## Issue Body

{Converted issue body}
```

### 4.4 Comment Sections

Each comment follows this format:

```
### Comment by @{username} · {ISO 8601 timestamp}

{Reaction counts (if any)}

{Converted comment body}

{Associated commit (if referenced)}
```

### 4.5 Complete Example Output

```markdown
# Add support for custom error handlers

**Issue:** golang/go/issues/12345 · **Status:** Open · **Author:** @johnsmith · **Created:** 2026-03-30T14:23:45Z
**Updated:** 2026-03-31T09:12:33Z
**Labels:** enhancement, api · **Milestone:** v1.5

## Issue Body

It would be great to support custom error handlers in the middleware. Currently, errors are handled using a default handler, but users may want to implement custom logic.

### Proposed API

```go
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func WithErrorHandler(h ErrorHandler) Option {
    // implementation
}
```

This would allow users to:
1. Format error responses differently
2. Add custom logging
3. Integrate with monitoring systems

---

### Comment by @janedoe · 2026-03-30T15:45:12Z

👍 This looks good. I have a few suggestions:

1. Consider adding a context to the error handler
2. Maybe support chaining multiple handlers?

### Comment by @johnsmith · 2026-03-30T16:30:45Z

@janedoe Good points! Updated the proposal:

```go
type ErrorHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)
```

For chaining, we could add a `ChainErrorHandler` helper. Thoughts?

*Associated commit: abc1234*

---

### Comment by @devops-bot · 2026-03-31T09:12:33Z

🎉 Automated test completed successfully.
```

### 4.6 Markdown Formatting Rules

1. **Headers**: Issue title uses H1, sections use H2/H3
2. **Code blocks**: Preserve triple-backtick blocks with language specifiers
3. **Inline code**: Preserve single-backtick formatting
4. **Bold/Italic**: Preserve **bold** and *italic* formatting
5. **Links**: Convert HTML `<a>` tags to Markdown `[text](url)` format
6. **Lists**: Preserve ordered and unordered lists
7. **Blockquotes**: Preserve `>` blockquote formatting
8. **Horizontal rules**: Use `---` for separators
9. **Tables**: Preserve Markdown table syntax
10. **Task lists**: Preserve `- [ ]` and `- [x]` format

### 4.7 GitHub-Specific Elements

| Element | Format | Example |
|---------|--------|---------|
| Mentions | `@{username}` | @johnsmith |
| Emoji | `:{emoji_name}:` | :thumbs_up: |
| Issue references | `#{number}` | #12345 |
| Commit references | `{short_hash}` | abc1234 |
| PR references | `#{number}` (context-dependent) | #54321 |

---

## 5. Technical Constraints and Considerations

### 5.1 Technology Stack

#### 5.1.1 Go Version and Requirements
- **Minimum Go version**: 1.24
- **Module path**: `github.com/bigwhite/issue2md`
- **Build system**: Standard Go toolchain (`go build`, `go test`)
- **Package structure**: Standard Go layout

#### 5.1.2 Dependencies Philosophy (Constitution-driven)
- **Standard Library First**: Prefer `net/http`, `encoding/json`, `time`, `strings`
- **No External Dependencies**: Avoid third-party packages unless absolutely necessary
- **Allowed External Dependencies**:
  - None for initial version
  - If needed, must be justified in design review

#### 5.1.3 Core Packages
```
cmd/issue2md/     - CLI entry point
internal/         - Internal packages
  ├── api/        - GitHub API client
  ├── converter/  - HTML to Markdown conversion
  ├── parser/     - URL parsing and validation
  └── model/      - Data structures
```

### 5.2 GitHub API Integration

#### 5.2.1 API Endpoints Used
- **Issue details**: `GET /repos/{owner}/{repo}/issues/{number}`
- **Comments**: `GET /repos/{owner}/{repo}/issues/{number}/comments`

#### 5.2.2 API Response Handling
- **Rate limiting**: Respect `X-RateLimit-Remaining` headers
- **Pagination**: Handle Link header for comment pagination
- **Error codes**:
  - 404: Issue or repository not found
  - 403: Private repository (unsupported)
  - 429: Rate limit exceeded
  - 5xx: Server error (retry)

#### 5.2.3 Authentication
- **Initial version**: No authentication (public repositories only)
- **Future consideration**: Support for personal access tokens for private repos

### 5.3 Network and Retry Strategy

#### 5.3.1 Retry Configuration
```go
const (
    MaxRetries = 3
    InitialBackoff = 1 * time.Second
    MaxBackoff = 4 * time.Second
    Timeout = 30 * time.Second
)
```

#### 5.3.2 Retry Logic
1. **Retry on**: Network errors, 5xx responses, rate limiting (429)
2. **Exponential backoff**: 1s, 2s, 4s
3. **Jitter**: Add random +/- 25% to backoff to avoid thundering herd
4. **No retry on**: 4xx client errors (except 429), 404, 403

#### 5.3.3 HTTP Client Configuration
```go
http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 10,
        IdleConnTimeout: 30 * time.Second,
    },
}
```

### 5.4 HTML to Markdown Conversion

#### 5.4.1 Conversion Strategy
- **Parser**: Use `golang.org/x/net/html` for HTML parsing
- **Converter**: Custom implementation following CommonMark specification
- **Preservation**: Maintain original structure as much as possible

#### 5.4.2 Supported HTML Elements
| HTML Tag | Markdown Equivalent |
|----------|---------------------|
| `<h1>`-`<h6>` | `#` to `######` |
| `<p>` | Paragraph (blank line) |
| `<code>` | `` `inline code` `` |
| `<pre><code>` | ```language |
| `<strong>` | `**bold**` |
| `<em>` | `*italic*` |
| `<a href="url">` | `[text](url)` |
| `<ul>`/`<ol>` | `- item` / `1. item` |
| `<li>` | List item |
| `<blockquote>` | `> quote` |
| `<hr>` | `---` |
| `<img>` | `
![alt](url)
` |

#### 5.4.3 Special Handling
- **GitHub-flavored markdown**: Preserve GFM extensions (tables, task lists)
- **Escaping**: Escape special characters where needed
- **Whitespace**: Preserve significant whitespace in code blocks
- **Entities**: Convert HTML entities to Unicode characters

### 5.5 Data Structures

#### 5.5.1 Issue Model
```go
type Issue struct {
    Number      int
    Title       string
    Body        string
    State       string // "open" or "closed"
    Author      User
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Labels      []Label
    Milestone   *Milestone
    CommentsURL string
}

type User struct {
    Login string
    URL   string
}

type Label struct {
    Name  string
    Color string
}

type Milestone struct {
    Title string
    URL   string
}
```

#### 5.5.2 Comment Model
```go
type Comment struct {
    ID        int
    Body      string
    Author    User
    CreatedAt time.Time
    UpdatedAt time.Time
    Reactions []Reaction
    CommitSHA string // optional
}

type Reaction struct {
    Content string
    Count   int
}
```

#### 5.5.3 Conversion Model
```go
type Converter interface {
    ConvertIssue(issue *Issue, comments []*Comment) string
}
```

### 5.6 Performance Considerations

#### 5.6.1 Memory Management
- **Streaming**: Process HTML elements as they're parsed (don't build full DOM)
- **Buffer reuse**: Use string builders efficiently
- **Comment batch processing**: Process comments in batches to reduce memory pressure

#### 5.6.2 Concurrency
- **Sequential fetching**: Initial version uses sequential API calls (simple)
- **Future optimization**: Consider parallel fetching of issue and comments

#### 5.6.3 Caching
- **No caching**: Initial version does not implement caching
- **Future consideration**: In-memory cache for repeated runs

### 5.7 Error Handling Strategy

#### 5.7.1 Error Types
```go
type ErrorType int

const (
    ErrorTypeInvalidURL ErrorType = iota
    ErrorTypeNetwork
    ErrorTypeAPI
    ErrorTypeRateLimit
    ErrorTypeNotFound
)

type AppError struct {
    Type    ErrorType
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
    return e.Err
}
```

#### 5.7.2 Error Propagation
- **Wrap all errors**: Use `fmt.Errorf("context: %w", err)`
- **Explicit handling**: Never ignore errors
- **Context preservation**: Maintain error context through the call stack

### 5.8 Logging Strategy

#### 5.8.1 Log Levels
- **Info**: High-level progress (issue fetched, converting, done)
- **Debug**: Detailed information (API calls, retry attempts)
- **Error**: Error conditions

#### 5.8.2 Log Output
- **Destination**: stderr (to avoid interfering with stdout)
- **Format**: Human-readable text
- **Conditional**: Only when `--verbose` flag is set

#### 5.8.3 Example Log Messages
```
[INFO] Fetching issue: golang/go#12345
[DEBUG] GET https://api.github.com/repos/golang/go/issues/12345
[INFO] Found 45 comments
[DEBUG] GET https://api.github.com/repos/golang/go/issues/12345/comments?page=1
[INFO] Converting to Markdown...
[INFO] Done
```

---

## 6. Testing Requirements

### 6.1 Testing Philosophy (Constitution-driven)

#### 6.1.1 Test-First Development (TDD)
- **Mandatory**: All new features must start with failing tests
- **TDD Cycle**: Red (failing test) → Green (passing test) → Refactor
- **No exception**: Even for "simple" features, write tests first

#### 6.1.2 Test Organization
```
internal/
  ├── api/
  │   ├── client.go
  │   └── client_test.go
  ├── converter/
  │   ├── markdown.go
  │   └── markdown_test.go
  ├── parser/
  │   ├── url.go
  │   └── url_test.go
  └── model/
      ├── issue.go
      └── issue_test.go
```

### 6.2 Unit Tests

#### 6.2.1 Test Coverage Goals
- **Overall coverage**: >80%
- **Critical paths**: 100% (URL parsing, API client, converter)
- **Error handling**: 100% coverage

#### 6.2.2 Table-Driven Tests (Preferred)
Example structure for URL parsing tests:
```go
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
            want: &IssueURL{Owner: "golang", Repo: "go", Number: 12345},
            wantErr: false,
        },
        {
            name:    "invalid URL - missing issue number",
            url:     "https://github.com/golang/go/issues",
            wantErr: true,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseIssueURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseIssueURL() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseIssueURL() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### 6.2.3 Test Categories

**Parser Tests (`internal/parser/`):**
- Valid GitHub Issue URLs
- Invalid URL formats
- Missing components
- Wrong domain (not github.com)
- PR URLs (should be rejected)
- Discussion URLs (should be rejected)

**API Client Tests (`internal/api/`):**
- HTTP client configuration
- Request building
- Response parsing
- Error handling (404, 403, 429, 5xx)
- Retry logic
- Rate limit handling
- Pagination handling

**Converter Tests (`internal/converter/`):**
- Basic HTML to Markdown conversion
- Code blocks preservation
- Link conversion
- List conversion
- Table conversion
- GitHub-specific elements
- Special character escaping
- Complex nested HTML

**Model Tests (`internal/model/`):**
- JSON unmarshaling
- Time parsing
- Default values
- Validation

### 6.3 Integration Tests

#### 6.3.1 Real GitHub API Tests
- **Test repository**: Use a dedicated test repository with known issues
- **Test fixtures**: Create specific issues with various content types
- **Rate limiting**: Space out requests to avoid hitting limits
- **Conditional execution**: Skip if `GITHUB_TOKEN` is not set or network is unavailable

#### 6.3.2 End-to-End Tests
- CLI argument parsing
- Complete workflow (URL → API → Conversion → Output)
- Error scenarios
- Verbose logging

### 6.4 Mocking Strategy (Constitution-driven)

#### 6.4.1 Minimize Mocks
- **Preference**: Use integration tests with real GitHub API
- **Mock only when**: Testing error scenarios that are hard to reproduce
- **Mock approach**: Use `httptest.Server` for HTTP mocking

#### 6.4.2 HTTP Server Mocking
```go
func TestFetchIssue_RateLimit(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-RateLimit-Remaining", "0")
        w.WriteHeader(http.StatusTooManyRequests)
    }))
    defer server.Close()

    // Test client with mock server URL
    // ...
}
```

### 6.5 Test Data

#### 6.5.1 Test Fixtures
```
testdata/
  ├── issues/
  │   ├── simple_issue.json
  │   ├── complex_issue.json
  │   └── with_comments.json
  ├── html/
  │   ├── simple.html
  │   ├── code_blocks.html
  │   └── tables.html
  └── expected/
      ├── simple.md
      ├── complex.md
      └── with_comments.md
```

#### 6.5.2 Golden File Testing
- Use golden files for converter output comparison
- Update golden files with `UPDATE_GOLDEN=1` flag
- Golden files live in `testdata/expected/`

### 6.6 Performance Tests

#### 6.6.1 Benchmark Tests
```go
func BenchmarkParseIssueURL(b *testing.B) {
    url := "https://github.com/golang/go/issues/12345"
    for i := 0; i < b.N; i++ {
        ParseIssueURL(url)
    }
}

func BenchmarkConvertHTMLToMarkdown(b *testing.B) {
    html := loadTestData("complex.html")
    for i := 0; i < b.N; i++ {
        ConvertHTMLToMarkdown(html)
    }
}
```

#### 6.6.2 Performance Goals
- **URL parsing**: <1µs per URL
- **HTML conversion**: <10ms for typical issue body
- **Memory usage**: <10MB for typical issue

### 6.7 Test Execution

#### 6.7.1 Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...

# Run verbose tests
go test -v ./...

# Run specific package
go test ./internal/api/

# Run specific test
go test -run TestParseIssueURL ./internal/parser/
```

#### 6.7.2 CI/CD Integration
- All tests must pass before merging
- Coverage report generated on each run
- Benchmark tests run on main branch

---

## 7. Implementation Phases

### 7.1 Phase 1: Foundation (MVP)
**Duration**: 1-2 weeks

**Tasks**:
1. Project structure setup
2. URL parsing and validation
3. Basic GitHub API client (issue fetch only)
4. Simple HTML to Markdown converter
5. Basic CLI with single URL argument
6. Error handling foundation

**Deliverable**: Working prototype that converts issue title and body to Markdown

### 7.2 Phase 2: Content Enhancement
**Duration**: 1-2 weeks

**Tasks**:
1. Comment fetching with pagination
2. Enhanced HTML converter (code blocks, lists, links)
3. Metadata extraction (labels, milestone, reactions)
4. GitHub-specific syntax preservation
5. Complete output formatting

**Deliverable**: Full-featured issue and comment conversion

### 7.3 Phase 3: Robustness
**Duration**: 1 week

**Tasks**:
1. Network retry mechanism
2. Rate limit handling
3. Comprehensive error handling
4. Verbose logging implementation
5. Edge case handling

**Deliverable**: Production-ready tool with robust error handling

### 7.4 Phase 4: Polish
**Duration**: 1 week

**Tasks**:
1. Documentation (README, examples)
2. Comprehensive test coverage
3. Performance optimization
4. Code review and refactoring
5. Release preparation

**Deliverable**: Version 1.0 release

---

## 8. Documentation Requirements

### 8.1 README.md Structure
1. Project description and purpose
2. Installation instructions
3. Usage examples
4. Command-line options
5. Output format description
6. Limitations and future work
7. Contributing guidelines
8. License

### 8.2 Code Documentation
- Public APIs: Complete Go doc comments
- Internal functions: Comments explaining non-obvious logic
- Examples: Where helpful, add example code in doc comments

### 8.3 Changelog
- Maintain CHANGELOG.md following Keep a Changelog format
- Document all user-facing changes
- Version numbers follow Semantic Versioning

---

## 9. Success Criteria

### 9.1 Technical Success
- [ ] All functional requirements implemented
- [ ] >80% test coverage achieved
- [ ] All tests pass consistently
- [ ] No external dependencies (Go standard library only)
- [ ] Code follows Go best practices and project constitution

### 9.2 User Success
- [ ] Can convert typical GitHub Issues successfully
- [ ] Output is readable and properly formatted
- [ ] Error messages are clear and actionable
- [ ] Performance meets requirements (30s for typical issue)

### 9.3 Project Success
- [ ] Specification fully implemented
- [ ] Documentation complete
- [ ] Ready for v1.0 release
- [ ] Foundation laid for future enhancements (private repos, PRs)

---

## 10. Future Enhancements (Out of Scope for v1.0)

### 10.1 Planned Features
- Pull Request support
- GitHub Discussions support
- Private repository support (with authentication)
- Output format options (HTML, JSON)
- Attachment downloading
- Batch processing multiple issues
- GitHub Enterprise support

### 10.2 Technical Improvements
- Parallel API calls for performance
- Caching layer
- Configurable retry strategy
- Custom output templates
- Interactive mode

---

## Appendix A: GitHub API Reference

### A.1 Issue Details Endpoint
```
GET /repos/{owner}/{repo}/issues/{number}
Response: JSON object with issue details
```

### A.2 Comments Endpoint
```
GET /repos/{owner}/{repo}/issues/{number}/comments
Response: JSON array of comments (supports pagination)
```

### A.3 Rate Limits
- Unauthenticated: 60 requests per hour
- Authenticated: 5,000 requests per hour
- Headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## Appendix B: Markdown CommonMark Reference

This tool should output Markdown compliant with the CommonMark specification. For details, see: https://spec.commonmark.org/

---

## Appendix C: Glossary

| Term | Definition |
|------|------------|
| Issue | GitHub Issue (not Pull Request or Discussion) |
| Comment | User comment on an issue |
| Collapsed comment | Comments that are hidden by default in GitHub UI |
| Attachment | File attached to an issue or comment |
| Mention | Reference to a GitHub user (@username) |
| Emoji | GitHub emoji shortcodes (:emoji_name:) |

---

**Document Status:** Draft
**Last Updated:** 2026-03-31
**Next Review:** After Phase 1 completion

---

## Critical Files for Implementation

Based on this specification, here are the 5 most critical files for implementing the issue2md tool:

- `cmd/issue2md/main.go` - CLI entry point and command-line argument handling
- `internal/api/client.go` - GitHub API client for fetching issues and comments
- `internal/parser/url.go` - URL parsing and validation logic
- `internal/converter/markdown.go` - HTML to Markdown conversion engine
- `internal/model/issue.go` - Data structures for issues, comments, and related entities
