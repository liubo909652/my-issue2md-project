package cli

import "fmt"

const (
	// Version is the version of issue2md
	Version = "1.0.0"

	// BuildDate is the build date
	BuildDate = "2026-03-31"
)

// Usage returns the usage information
func Usage() string {
	return `issue2md - Convert GitHub Issues to Markdown

Usage:
  issue2md [OPTIONS] <GITHUB_ISSUE_URL>

Arguments:
  GITHUB_ISSUE_URL    GitHub Issue URL to convert

Options:
  --verbose           Enable verbose logging to stderr
  --help              Show this help message
  --version           Show version information

Examples:
  issue2md https://github.com/golang/go/issues/12345
  issue2md --verbose https://github.com/golang/go/issues/12345 > issue.md
  issue2md --help

Exit Codes:
  0   Success
  1   General error (invalid URL, network failure, API error)
  2   Command-line argument error
  3   Rate limit exceeded (after retries)

Environment Variables:
  ISSUE2MD_TIMEOUT    HTTP request timeout (default: 30s)
  ISSUE2MD_VERBOSE    Enable verbose logging (default: false)

For more information, see: https://github.com/bigwhite/issue2md
`
}

// VersionInfo returns the version information
func VersionInfo() string {
	return fmt.Sprintf(`issue2md version %s (build: %s)

Copyright (c) 2026

GitHub Issue to Markdown Converter

For more information, see: https://github.com/bigwhite/issue2md
`, Version, BuildDate)
}