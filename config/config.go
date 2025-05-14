package config

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrZeroPort                  = errors.New("port cannot be zero")
	ErrInvalidPort               = "invalid port: %d it must be between 1024 and 65535"
	ErrInvalidGoogleCloudProject = "invalid google cloud project: %s"
)

type Config struct {
	Port               int
	GoogleCloudProject string
}

type Option func(*Config) error

func New(options ...Option) (*Config, error) {
	config := &Config{
		Port: 8080,
	}

	for _, option := range options {
		err := option(config)
		if err != nil {
			return nil, err
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
			return fmt.Errorf(ErrInvalidPort, portInt)
		}

		if portInt == 0 {
			return nil
		}

		if portInt < 1024 || portInt > 65535 {
			return fmt.Errorf(ErrInvalidPort, portInt)
		}

		c.Port = portInt
		return nil
	}
}

func WithGoogleCloudProject(project string) Option {
	return func(c *Config) error {
		if project == "" {
			return fmt.Errorf(ErrInvalidGoogleCloudProject, project)
		}
		c.GoogleCloudProject = project
		return nil
	}
}
