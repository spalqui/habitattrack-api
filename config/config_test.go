package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPort(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		wantErr  error
		wantPort int
	}{
		{
			name:     "Valid port",
			port:     "8080",
			wantPort: 8080,
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

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantPort, cfg.Port)
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

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.project, cfg.ProjectID)
		})
	}
}
