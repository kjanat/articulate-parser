// Package interfaces provides the core contracts for the articulate-parser application.
// It defines interfaces for parsing and exporting Articulate Rise courses.
package interfaces

import "context"

// Logger defines the interface for structured logging.
// Implementations should provide leveled, structured logging capabilities.
type Logger interface {
	// Debug logs a debug-level message with optional key-value pairs.
	Debug(msg string, keysAndValues ...any)

	// Info logs an info-level message with optional key-value pairs.
	Info(msg string, keysAndValues ...any)

	// Warn logs a warning-level message with optional key-value pairs.
	Warn(msg string, keysAndValues ...any)

	// Error logs an error-level message with optional key-value pairs.
	Error(msg string, keysAndValues ...any)

	// With returns a new logger with the given key-value pairs added as context.
	With(keysAndValues ...any) Logger

	// WithContext returns a new logger with context information.
	WithContext(ctx context.Context) Logger
}
