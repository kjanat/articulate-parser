// Package version provides version information for the Articulate Parser.
// It includes the current version, build time, and Git commit hash.
package version

// Version information.
var (
	// Version is the current version of the application.
	// Breaking changes from 0.4.x:
	// - Renamed GetSupportedFormat() -> SupportedFormat()
	// - Renamed GetSupportedFormats() -> SupportedFormats()
	// - FetchCourse now requires context.Context parameter
	// - NewArticulateParser now accepts logger, baseURL, timeout
	// New features:
	// - Structured logging with slog
	// - Configuration via environment variables
	// - Context-aware HTTP requests
	// - Comprehensive benchmarks and examples
	Version = "1.0.0"

	// BuildTime is the time the binary was built.
	BuildTime = "unknown"

	// GitCommit is the git commit hash.
	GitCommit = "unknown"
)
