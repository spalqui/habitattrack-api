package config

type Config struct {
	Port               string
	GoogleCloudProject string
}

type Option func(*Config) error

func New(options ...Option) (*Config, error) {
	config := &Config{
		Port: "8080",
	}

	for _, option := range options {
		err := option(config)
		if err != nil {
			return &Config{}, err
		}
	}
	return config, nil
}

func WithPort(port string) Option {
	return func(c *Config) error {
		c.Port = port
		return nil
	}
}

func WithGoogleCloudProject(project string) Option {
	return func(c *Config) error {
		c.GoogleCloudProject = project
		return nil
	}
}
