// Package services_test provides benchmarks for the logger implementations.
package services

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

// BenchmarkSlogLogger_Info benchmarks structured JSON logging.
func BenchmarkSlogLogger_Info(b *testing.B) {
	// Create logger that writes to io.Discard to avoid benchmark noise
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := &SlogLogger{logger: slog.New(handler)}

	b.ResetTimer()
	for b.Loop() {
		logger.Info("test message", "key1", "value1", "key2", 42, "key3", true)
	}
}

// BenchmarkSlogLogger_Debug benchmarks debug level logging.
func BenchmarkSlogLogger_Debug(b *testing.B) {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := &SlogLogger{logger: slog.New(handler)}

	b.ResetTimer()
	for b.Loop() {
		logger.Debug("debug message", "operation", "test", "duration", 123)
	}
}

// BenchmarkSlogLogger_Error benchmarks error logging.
func BenchmarkSlogLogger_Error(b *testing.B) {
	opts := &slog.HandlerOptions{Level: slog.LevelError}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := &SlogLogger{logger: slog.New(handler)}

	b.ResetTimer()
	for b.Loop() {
		logger.Error("error occurred", "error", "test error", "code", 500)
	}
}

// BenchmarkTextLogger_Info benchmarks text logging.
func BenchmarkTextLogger_Info(b *testing.B) {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewTextHandler(io.Discard, opts)
	logger := &SlogLogger{logger: slog.New(handler)}

	b.ResetTimer()
	for b.Loop() {
		logger.Info("test message", "key1", "value1", "key2", 42)
	}
}

// BenchmarkNoOpLogger benchmarks the no-op logger.
func BenchmarkNoOpLogger(b *testing.B) {
	logger := NewNoOpLogger()

	b.ResetTimer()
	for b.Loop() {
		logger.Info("test message", "key1", "value1", "key2", 42)
		logger.Error("error message", "error", "test")
	}
}

// BenchmarkLogger_With benchmarks logger with context.
func BenchmarkLogger_With(b *testing.B) {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := &SlogLogger{logger: slog.New(handler)}

	b.ResetTimer()
	for b.Loop() {
		contextLogger := logger.With("request_id", "123", "user_id", "456")
		contextLogger.Info("operation completed")
	}
}

// BenchmarkLogger_WithContext benchmarks logger with Go context.
func BenchmarkLogger_WithContext(b *testing.B) {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := &SlogLogger{logger: slog.New(handler)}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		contextLogger := logger.WithContext(ctx)
		contextLogger.Info("context operation")
	}
}
