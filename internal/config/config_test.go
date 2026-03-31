package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test with default values
	os.Clearenv()
	config := Load()

	if config.BaseURL != "https://api.github.com" {
		t.Errorf("Expected BaseURL to be 'https://api.github.com', got %s", config.BaseURL)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout to be 30s, got %v", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.InitialBackoff != 1*time.Second {
		t.Errorf("Expected InitialBackoff to be 1s, got %v", config.InitialBackoff)
	}

	if config.MaxBackoff != 4*time.Second {
		t.Errorf("Expected MaxBackoff to be 4s, got %v", config.MaxBackoff)
	}

	if config.Verbose {
		t.Error("Expected Verbose to be false by default")
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Clearenv()
	os.Setenv("ISSUE2MD_TIMEOUT", "60s")
	os.Setenv("ISSUE2MD_VERBOSE", "true")
	os.Setenv("ISSUE2MD_MAX_RETRIES", "5")
	os.Setenv("ISSUE2MD_INITIAL_BACKOFF", "2s")
	os.Setenv("ISSUE2MD_MAX_BACKOFF", "10s")

	config := LoadFromEnv()

	if config.Timeout != 60*time.Second {
		t.Errorf("Expected Timeout to be 60s, got %v", config.Timeout)
	}

	if !config.Verbose {
		t.Error("Expected Verbose to be true")
	}

	if config.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5, got %d", config.MaxRetries)
	}

	if config.InitialBackoff != 2*time.Second {
		t.Errorf("Expected InitialBackoff to be 2s, got %v", config.InitialBackoff)
	}

	if config.MaxBackoff != 10*time.Second {
		t.Errorf("Expected MaxBackoff to be 10s, got %v", config.MaxBackoff)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  &Config{BaseURL: "https://api.github.com", Timeout: 30 * time.Second, MaxRetries: 3, InitialBackoff: 1 * time.Second, MaxBackoff: 4 * time.Second},
			wantErr: false,
		},
		{
			name:    "invalid base URL",
			config:  &Config{BaseURL: "invalid-url", Timeout: 30 * time.Second},
			wantErr: true,
		},
		{
			name:    "negative timeout",
			config:  &Config{BaseURL: "https://api.github.com", Timeout: -1 * time.Second},
			wantErr: true,
		},
		{
			name:    "negative max retries",
			config:  &Config{BaseURL: "https://api.github.com", MaxRetries: -1},
			wantErr: true,
		},
		{
			name:    "negative backoff",
			config:  &Config{BaseURL: "https://api.github.com", InitialBackoff: -1 * time.Second},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}