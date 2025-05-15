package config

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	DefaultPort                = 8080
	DefaultFirestoreDatabaseID = "habitattrack"

	ErrPortOutOfRange          = errors.New("port must be between 1024 and 65535")
	ErrEmptyGoogleCloudProject = errors.New("google cloud project cannot be empty")
)

type Config struct {
	Port       int
	ProjectID  string
	DatabaseID string
}

type Option func(*Config) error

func New(options ...Option) (*Config, error) {
	config := &Config{
		Port:       DefaultPort, // Default port
		DatabaseID: DefaultFirestoreDatabaseID,
	}

	for _, option := range options {
		err := option(config)
		if err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}
	return config, nil
}

func WithPort(port string) Option {
	return func(c *Config) error {
		if port == "" {

			return nil
		}

		portInt, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("parsing port: %w", err)
		}

		if portInt == 0 {
			return nil
		}

		if portInt < 1024 || portInt > 65535 {
			return ErrPortOutOfRange
		}

		c.Port = portInt
		return nil
	}
}

func WithGoogleCloudProject(project string) Option {
	return func(c *Config) error {
		if project == "" {
			return ErrEmptyGoogleCloudProject
		}
		c.ProjectID = project
		return nil
	}
}

func WithFirestoreDatabase(databaseID string) Option {
	return func(c *Config) error {
		if databaseID == "" {
			return nil
		}
		c.DatabaseID = databaseID
		return nil
	}
}
