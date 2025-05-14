package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFirestorePropertyRepository(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		wantErr   bool
	}{
		{
			name:      "valid project ID",
			projectID: "test-project",
			wantErr:   false,
		},
		{
			name:      "empty project ID",
			projectID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			repo, err := NewFirestorePropertyRepository(ctx, tt.projectID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, repo)

				err = repo.Close()
				require.NoError(t, err)
			}
		})
	}
}
