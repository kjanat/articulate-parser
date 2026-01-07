package exporters

import "github.com/kjanat/articulate-parser/internal/interfaces"

// baseExporter contains common fields shared by all exporter implementations.
// Embed this struct in concrete exporters to inherit HTMLCleaner and Logger.
type baseExporter struct {
	// htmlCleaner is used to convert HTML content to plain text.
	htmlCleaner interfaces.HTMLCleaner
	// logger is used for logging warnings and errors.
	logger interfaces.Logger
}

// newBaseExporter creates a new baseExporter with the given dependencies.
func newBaseExporter(htmlCleaner interfaces.HTMLCleaner, logger interfaces.Logger) baseExporter {
	return baseExporter{
		htmlCleaner: htmlCleaner,
		logger:      logger,
	}
}
