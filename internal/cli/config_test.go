package cli

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Config
		wantErr bool
	}{
		{
			name:    "valid args - basic",
			args:    []string{"issue2md", "https://github.com/golang/go/issues/12345"},
			want:    &Config{URL: "https://github.com/golang/go/issues/12345", Verbose: false, Version: false, Help: false},
			wantErr: false,
		},
		{
			name:    "valid args - verbose",
			args:    []string{"issue2md", "--verbose", "https://github.com/golang/go/issues/12345"},
			want:    &Config{URL: "https://github.com/golang/go/issues/12345", Verbose: true, Version: false, Help: false},
			wantErr: false,
		},
		{
			name:    "valid args - version",
			args:    []string{"issue2md", "--version"},
			want:    &Config{URL: "", Verbose: false, Version: true, Help: false},
			wantErr: false,
		},
		{
			name:    "valid args - help",
			args:    []string{"issue2md", "--help"},
			want:    &Config{URL: "", Verbose: false, Version: false, Help: true},
			wantErr: false,
		},
		{
			name:    "no args",
			args:    []string{"issue2md"},
			want:    &Config{URL: "", Verbose: false, Version: false, Help: false},
			wantErr: false,
		},
		{
			name:    "too many args",
			args:    []string{"issue2md", "url1", "url2"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config with URL",
			config:  &Config{URL: "https://github.com/golang/go/issues/12345", Verbose: true},
			wantErr: false,
		},
		{
			name:    "valid config version flag",
			config:  &Config{Version: true},
			wantErr: false,
		},
		{
			name:    "valid config help flag",
			config:  &Config{Help: true},
			wantErr: false,
		},
		{
			name:    "invalid config - no URL, not version/help",
			config:  &Config{URL: "", Verbose: false, Version: false, Help: false},
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