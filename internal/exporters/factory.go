package exporters

import (
	"fmt"
	"strings"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/services"
)

// Format constants for supported export formats.
const (
	FormatMarkdown = "markdown"
	FormatDocx     = "docx"
	FormatHTML     = "html"
)

// Factory implements the ExporterFactory interface.
// It creates appropriate exporter instances based on the requested format.
type Factory struct {
	// htmlCleaner is used by exporters to convert HTML content to plain text
	htmlCleaner *services.HTMLCleaner
}

// NewFactory creates a new exporter factory.
// It takes an HTMLCleaner instance that will be passed to the exporters
// created by this factory.
//
// Parameters:
//   - htmlCleaner: Service for cleaning HTML content in course data
//
// Returns:
//   - An implementation of the ExporterFactory interface
func NewFactory(htmlCleaner *services.HTMLCleaner) interfaces.ExporterFactory {
	return &Factory{
		htmlCleaner: htmlCleaner,
	}
}

// CreateExporter creates an exporter for the specified format.
// Format strings are case-insensitive (e.g., "markdown", "DOCX").
func (f *Factory) CreateExporter(format string) (interfaces.Exporter, error) {
	switch strings.ToLower(format) {
	case FormatMarkdown, "md":
		return NewMarkdownExporter(f.htmlCleaner), nil
	case FormatDocx, "word":
		return NewDocxExporter(f.htmlCleaner), nil
	case FormatHTML, "htm":
		return NewHTMLExporter(f.htmlCleaner), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// SupportedFormats returns a list of all supported export formats,
// including both primary format names and their aliases.
func (f *Factory) SupportedFormats() []string {
	return []string{FormatMarkdown, "md", FormatDocx, "word", FormatHTML, "htm"}
}
