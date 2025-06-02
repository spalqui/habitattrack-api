package constants

import "time"

const (
	DefaultTimeout = 30 * time.Second
	// DefaultPageSize and MaxPageSize are now defined locally within feature handlers/services
	// to allow for feature-specific pagination settings if needed.
)
