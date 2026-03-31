package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration
type Config struct {
	// GitHub API settings
	BaseURL  string
	Timeout  time.Duration

	// Retry settings
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration

	// Logging settings
	Verbose bool
}

// Load creates a new Config with default values and environment variables
func Load() *Config {
	config := LoadFromEnv()

	// Set defaults if not set
	if config.BaseURL == "" {
		config.BaseURL = "https://api.github.com"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.InitialBackoff == 0 {
		config.InitialBackoff = 1 * time.Second
	}
	if config.MaxBackoff == 0 {
		config.MaxBackoff = 4 * time.Second
	}

	return config
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	config := &Config{}

	// Load environment variables
	if baseURL := os.Getenv("ISSUE2MD_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if timeout := os.Getenv("ISSUE2MD_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.Timeout = d
		}
	}

	if retries := os.Getenv("ISSUE2MD_MAX_RETRIES"); retries != "" {
		if n, err := strconv.Atoi(retries); err == nil {
			config.MaxRetries = n
		}
	}

	if backoff := os.Getenv("ISSUE2MD_INITIAL_BACKOFF"); backoff != "" {
		if d, err := time.ParseDuration(backoff); err == nil {
			config.InitialBackoff = d
		}
	}

	if maxBackoff := os.Getenv("ISSUE2MD_MAX_BACKOFF"); maxBackoff != "" {
		if d, err := time.ParseDuration(maxBackoff); err == nil {
			config.MaxBackoff = d
		}
	}

	if verbose := os.Getenv("ISSUE2MD_VERBOSE"); verbose != "" {
		config.Verbose = verbose == "true" || verbose == "1"
	}

	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate base URL
	if c.BaseURL == "" {
		return fmt.Errorf("base URL cannot be empty")
	}

	// Validate timeout
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Validate max retries
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	// Validate backoff
	if c.InitialBackoff <= 0 {
		return fmt.Errorf("initial backoff must be positive")
	}

	if c.MaxBackoff <= 0 {
		return fmt.Errorf("max backoff must be positive")
	}

	if c.InitialBackoff > c.MaxBackoff {
		return fmt.Errorf("initial backoff cannot be greater than max backoff")
	}

	return nil
}