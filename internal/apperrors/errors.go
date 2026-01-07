// Package apperrors provides sentinel errors for the articulate-parser application.
package apperrors

import "errors"

// Export errors.
var (
	// ErrUnsupportedFormat is returned when an unknown export format is requested.
	ErrUnsupportedFormat = errors.New("unsupported export format")
)

// URI parsing errors.
var (
	// ErrInvalidURI is returned when a URI cannot be parsed.
	ErrInvalidURI = errors.New("invalid URI")

	// ErrInvalidDomain is returned when the URI domain is not an Articulate Rise domain.
	ErrInvalidDomain = errors.New("invalid domain for Articulate Rise URI")

	// ErrNoShareID is returned when a share ID cannot be extracted from the URI.
	ErrNoShareID = errors.New("could not extract share ID from URI")
)

// Operation errors.
var (
	// ErrOperationCanceled is returned when an operation is canceled via context.
	ErrOperationCanceled = errors.New("operation canceled")
)
