package config

import (
	"errors"
	"testing"
)

func TestWithPort(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr error
	}{
		{
			name: "Valid port",
			port: "8080",
		},
		{
			name:    "Port below range",
			port:    "1023",
			wantErr: ErrPortOutOfRange,
		},
		{
			name:    "Port above range",
			port:    "65536",
			wantErr: ErrPortOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := WithPort(tt.port)(cfg)

			if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
				t.Errorf("WithPort(%s) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestWithGoogleCloudProject(t *testing.T) {
	tests := []struct {
		name    string
		project string
		wantErr error
	}{
		{
			name:    "Valid project",
			project: "my-project",
		},
		{
			name:    "Empty project",
			wantErr: ErrEmptyGoogleCloudProject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := WithGoogleCloudProject(tt.project)(cfg)

			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) || (err != nil && err.Error() != tt.wantErr.Error()) {
				t.Errorf("WithGoogleCloudProject(%q) error = %v, wantErr %v", tt.project, err, tt.wantErr)
			}

			if err == nil && cfg.GoogleCloudProject != tt.project {
				t.Errorf("WithGoogleCloudProject(%q) set GoogleCloudProject = %q, want %q", tt.project, cfg.GoogleCloudProject, tt.project)
			}
		})
	}
}
