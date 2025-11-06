// Package config_test provides tests for the config package.
package config

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Clear environment
	os.Clearenv()

	cfg := Load()

	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("Expected BaseURL '%s', got '%s'", DefaultBaseURL, cfg.BaseURL)
	}

	if cfg.RequestTimeout != DefaultRequestTimeout {
		t.Errorf("Expected timeout %v, got %v", DefaultRequestTimeout, cfg.RequestTimeout)
	}

	if cfg.LogLevel != DefaultLogLevel {
		t.Errorf("Expected log level %v, got %v", DefaultLogLevel, cfg.LogLevel)
	}

	if cfg.LogFormat != DefaultLogFormat {
		t.Errorf("Expected log format '%s', got '%s'", DefaultLogFormat, cfg.LogFormat)
	}
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	t.Setenv("ARTICULATE_BASE_URL", "https://test.example.com")
	t.Setenv("ARTICULATE_REQUEST_TIMEOUT", "60")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "json")

	cfg := Load()

	if cfg.BaseURL != "https://test.example.com" {
		t.Errorf("Expected BaseURL 'https://test.example.com', got '%s'", cfg.BaseURL)
	}

	if cfg.RequestTimeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", cfg.RequestTimeout)
	}

	if cfg.LogLevel != slog.LevelDebug {
		t.Errorf("Expected log level Debug, got %v", cfg.LogLevel)
	}

	if cfg.LogFormat != "json" {
		t.Errorf("Expected log format 'json', got '%s'", cfg.LogFormat)
	}
}

func TestGetLogLevelEnv(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected slog.Level
	}{
		{"debug lowercase", "debug", slog.LevelDebug},
		{"debug uppercase", "DEBUG", slog.LevelDebug},
		{"info lowercase", "info", slog.LevelInfo},
		{"info uppercase", "INFO", slog.LevelInfo},
		{"warn lowercase", "warn", slog.LevelWarn},
		{"warn uppercase", "WARN", slog.LevelWarn},
		{"warning lowercase", "warning", slog.LevelWarn},
		{"error lowercase", "error", slog.LevelError},
		{"error uppercase", "ERROR", slog.LevelError},
		{"invalid value", "invalid", slog.LevelInfo},
		{"empty value", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.value != "" {
				t.Setenv("TEST_LOG_LEVEL", tt.value)
			}
			result := getLogLevelEnv("TEST_LOG_LEVEL", slog.LevelInfo)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected time.Duration
	}{
		{"valid duration", "45", 45 * time.Second},
		{"zero duration", "0", 0},
		{"invalid duration", "invalid", 30 * time.Second},
		{"empty value", "", 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.value != "" {
				t.Setenv("TEST_DURATION", tt.value)
			}
			result := getDurationEnv("TEST_DURATION", 30*time.Second)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
