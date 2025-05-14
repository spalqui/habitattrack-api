package config

import (
	"fmt"
	"strconv"
	"testing"
)

func TestWithPort(t *testing.T) {
	tests := []struct {
		name      string
		port      string
		expectErr error
	}{
		{"Valid port", "8080", nil},
		{"Zero port", "0", ErrZeroPort},
		{"Port below range", "1023", fmt.Errorf(ErrInvalidPort, 1023)},
		{"Port above range", "65536", fmt.Errorf(ErrInvalidPort, 65536)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := WithPort(tt.port)(cfg)

			if err != nil {
				if (err != nil && tt.expectErr == nil) || (err == nil && tt.expectErr != nil) || (err != nil && err.Error() != tt.expectErr.Error()) {
					t.Errorf("WithPort(%s) error = %v, wantErr %v", tt.port, err, tt.expectErr)
				}
				return
			}

			port, err := strconv.Atoi(tt.port)
			if err != nil {
				t.Fatalf("failed to convert port %s to int: %v", tt.port, err)
			}
			if err == nil && cfg.Port != port {
				t.Errorf("WithPort(%s) set Port = %d, want %d", tt.port, cfg.Port, port)
			}
		})
	}
}

func TestWithGoogleCloudProject(t *testing.T) {
	tests := []struct {
		name      string
		project   string
		expectErr error
	}{
		{"Valid project", "my-project", nil},
		{"Empty project", "", fmt.Errorf(ErrInvalidGoogleCloudProject, "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := WithGoogleCloudProject(tt.project)(cfg)

			if (err != nil && tt.expectErr == nil) || (err == nil && tt.expectErr != nil) || (err != nil && err.Error() != tt.expectErr.Error()) {
				t.Errorf("WithGoogleCloudProject(%q) error = %v, wantErr %v", tt.project, err, tt.expectErr)
			}

			if err == nil && cfg.GoogleCloudProject != tt.project {
				t.Errorf("WithGoogleCloudProject(%q) set GoogleCloudProject = %q, want %q", tt.project, cfg.GoogleCloudProject, tt.project)
			}
		})
	}
}
