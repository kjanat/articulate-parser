package services

import (
	"context"
	"log/slog"
	"os"

	"github.com/kjanat/articulate-parser/internal/interfaces"
)

// SlogLogger implements the Logger interface using the standard library's slog package.
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger creates a new structured logger using slog.
// The level parameter controls the minimum log level (debug, info, warn, error).
func NewSlogLogger(level slog.Level) interfaces.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return &SlogLogger{
		logger: slog.New(handler),
	}
}

// NewTextLogger creates a new structured logger with human-readable text output.
// Useful for development and debugging.
func NewTextLogger(level slog.Level) interfaces.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return &SlogLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a debug-level message with optional key-value pairs.
func (l *SlogLogger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, keysAndValues...)
}

// Info logs an info-level message with optional key-value pairs.
func (l *SlogLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

// Warn logs a warning-level message with optional key-value pairs.
func (l *SlogLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, keysAndValues...)
}

// Error logs an error-level message with optional key-value pairs.
func (l *SlogLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, keysAndValues...)
}

// With returns a new logger with the given key-value pairs added as context.
func (l *SlogLogger) With(keysAndValues ...any) interfaces.Logger {
	return &SlogLogger{
		logger: l.logger.With(keysAndValues...),
	}
}

// WithContext returns a new logger with context information.
// Currently preserves the logger as-is, but can be extended to extract
// trace IDs or other context values in the future.
func (l *SlogLogger) WithContext(ctx context.Context) interfaces.Logger {
	// Can be extended to extract trace IDs, request IDs, etc. from context
	return l
}

// NoOpLogger is a logger that discards all log messages.
// Useful for testing or when logging should be disabled.
type NoOpLogger struct{}

// NewNoOpLogger creates a logger that discards all messages.
func NewNoOpLogger() interfaces.Logger {
	return &NoOpLogger{}
}

// Debug does nothing.
func (l *NoOpLogger) Debug(msg string, keysAndValues ...any) {}

// Info does nothing.
func (l *NoOpLogger) Info(msg string, keysAndValues ...any) {}

// Warn does nothing.
func (l *NoOpLogger) Warn(msg string, keysAndValues ...any) {}

// Error does nothing.
func (l *NoOpLogger) Error(msg string, keysAndValues ...any) {}

// With returns the same no-op logger.
func (l *NoOpLogger) With(keysAndValues ...any) interfaces.Logger {
	return l
}

// WithContext returns the same no-op logger.
func (l *NoOpLogger) WithContext(ctx context.Context) interfaces.Logger {
	return l
}
