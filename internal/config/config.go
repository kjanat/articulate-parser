// Package config provides configuration management for the articulate-parser application.
// It supports loading configuration from environment variables and command-line flags.
package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration values for the application.
type Config struct {
	// Parser configuration
	BaseURL        string
	RequestTimeout time.Duration

	// Logging configuration
	LogLevel  slog.Level
	LogFormat string // "json" or "text"
}

// Default configuration values.
const (
	DefaultBaseURL        = "https://rise.articulate.com"
	DefaultRequestTimeout = 30 * time.Second
	DefaultLogLevel       = slog.LevelInfo
	DefaultLogFormat      = "text"
)

// Load creates a new Config with values from environment variables.
// Falls back to defaults if environment variables are not set.
func Load() *Config {
	return &Config{
		BaseURL:        getEnv("ARTICULATE_BASE_URL", DefaultBaseURL),
		RequestTimeout: getDurationEnv("ARTICULATE_REQUEST_TIMEOUT", DefaultRequestTimeout),
		LogLevel:       getLogLevelEnv("LOG_LEVEL", DefaultLogLevel),
		LogFormat:      getEnv("LOG_FORMAT", DefaultLogFormat),
	}
}

// getEnv retrieves an environment variable or returns the default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv retrieves a duration from environment variable or returns default.
// The environment variable should be in seconds (e.g., "30" for 30 seconds).
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}

// getLogLevelEnv retrieves a log level from environment variable or returns default.
// Accepts: "debug", "info", "warn", "error" (case-insensitive).
func getLogLevelEnv(key string, defaultValue slog.Level) slog.Level {
	value := os.Getenv(key)
	switch value {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return defaultValue
	}
}
