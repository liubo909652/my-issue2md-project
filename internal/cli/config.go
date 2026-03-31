package cli

import (
	"flag"
	"fmt"
)

// minIssueURLLength is the minimum length of a valid GitHub Issue URL
const minIssueURLLength = len("https://a/b/issues/1")

// Config represents the CLI configuration
type Config struct {
	URL     string
	Verbose bool
	Version bool
	Help    bool
}

// ParseArgs parses command-line arguments and returns a Config
func ParseArgs(args []string) (*Config, error) {
	config := &Config{}

	// Parse arguments
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	fs.BoolVar(&config.Verbose, "verbose", false, "Enable verbose logging")
	fs.BoolVar(&config.Version, "version", false, "Show version information")
	fs.BoolVar(&config.Help, "help", false, "Show help information")

	// Parse positional arguments
	fs.Parse(args[1:])

	// Collect remaining arguments as URL
	if fs.NArg() > 0 {
		if fs.NArg() > 1 {
			return nil, fmt.Errorf("too many arguments")
		}
		config.URL = fs.Arg(0)
	}

	return config, nil
}

// Validate validates the CLI configuration
func (c *Config) Validate() error {
	// If --version or --help is specified, no URL is needed
	if c.Version || c.Help {
		return nil
	}

	// URL is required unless --version or --help is specified
	if c.URL == "" {
		return fmt.Errorf("GitHub Issue URL is required")
	}

	// Additional URL validation can be added here
	if len(c.URL) < minIssueURLLength {
		return fmt.Errorf("invalid GitHub Issue URL format")
	}

	return nil
}