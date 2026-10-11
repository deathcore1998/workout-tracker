package config

import "os"

type Config struct {
	DatabaseURL string
	ServerPort  string
}

// Load reads configuration from environment variables.
// If a variable is not set, it uses the default value.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://workout:secret@localhost:5432/workout?sslmode=disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

// getEnv returns the value of an environment variable, or a default if empty.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
