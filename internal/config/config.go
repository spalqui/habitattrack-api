package config

import (
	"os"
)

type Config struct {
	Port             string
	GoogleProject    string
	FirestoreKeyPath string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		GoogleProject:    getEnv("GOOGLE_CLOUD_PROJECT", ""),
		FirestoreKeyPath: getEnv("FIRESTORE_KEY_PATH", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
