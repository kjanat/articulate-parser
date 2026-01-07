package exporters

import (
	"fmt"
	"strings"

	"github.com/kjanat/articulate-parser/internal/apperrors"
	"github.com/kjanat/articulate-parser/internal/interfaces"
)

// Format constants for supported export formats.
const (
	FormatMarkdown = "markdown"
	FormatDocx     = "docx"
	FormatHTML     = "html"

	// Format aliases accepted by CreateExporter.
	formatAliasMarkdown = "md"
	formatAliasDocx     = "word"
	formatAliasHTML     = "htm"
)

// Factory implements the ExporterFactory interface.
// It creates appropriate exporter instances based on the requested format.
type Factory struct {
	// htmlCleaner is used by exporters to convert HTML content to plain text
	htmlCleaner interfaces.HTMLCleaner
	// logger is used by exporters for logging
	logger interfaces.Logger
}

// NewFactory creates a new exporter factory.
// It takes an HTMLCleaner instance and a logger that will be passed to the exporters
// created by this factory.
//
// Parameters:
//   - htmlCleaner: Service for cleaning HTML content in course data
//   - logger: Logger instance for exporter logging
//
// Returns:
//   - An implementation of the ExporterFactory interface
func NewFactory(htmlCleaner interfaces.HTMLCleaner, logger interfaces.Logger) interfaces.ExporterFactory {
	return &Factory{
		htmlCleaner: htmlCleaner,
		logger:      logger,
	}
}

// CreateExporter creates an exporter for the specified format.
// Format strings are case-insensitive (e.g., "markdown", "DOCX").
func (f *Factory) CreateExporter(format string) (interfaces.Exporter, error) {
	switch strings.ToLower(format) {
	case FormatMarkdown, formatAliasMarkdown:
		return NewMarkdownExporter(f.htmlCleaner, f.logger), nil
	case FormatDocx, formatAliasDocx:
		return NewDocxExporter(f.htmlCleaner, f.logger), nil
	case FormatHTML, formatAliasHTML:
		return NewHTMLExporter(f.htmlCleaner, f.logger), nil
	default:
		return nil, fmt.Errorf("%w: %s", apperrors.ErrUnsupportedFormat, format)
	}
}

// SupportedFormats returns a list of all supported export formats,
// including both primary format names and their aliases.
func (f *Factory) SupportedFormats() []string {
	return []string{
		FormatMarkdown, formatAliasMarkdown,
		FormatDocx, formatAliasDocx,
		FormatHTML, formatAliasHTML,
	}
}
