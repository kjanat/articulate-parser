package apperrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		sentinel error
		wrapped  error
	}{
		{"ErrUnsupportedFormat", ErrUnsupportedFormat, fmt.Errorf("%w: xyz", ErrUnsupportedFormat)},
		{"ErrInvalidURI", ErrInvalidURI, fmt.Errorf("%w: bad", ErrInvalidURI)},
		{"ErrInvalidDomain", ErrInvalidDomain, fmt.Errorf("%w: evil.com", ErrInvalidDomain)},
		{"ErrNoShareID", ErrNoShareID, fmt.Errorf("%w: no match", ErrNoShareID)},
		{"ErrOperationCanceled", ErrOperationCanceled, fmt.Errorf("%w: ctx done", ErrOperationCanceled)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.wrapped, tt.sentinel) {
				t.Errorf("errors.Is(%v, %v) = false, want true", tt.wrapped, tt.sentinel)
			}
		})
	}
}

func TestSentinelErrors_Message(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{ErrUnsupportedFormat, "unsupported export format"},
		{ErrInvalidURI, "invalid URI"},
		{ErrInvalidDomain, "invalid domain for Articulate Rise URI"},
		{ErrNoShareID, "could not extract share ID from URI"},
		{ErrOperationCanceled, "operation canceled"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("got %q, want %q", tt.err.Error(), tt.expected)
			}
		})
	}
}
