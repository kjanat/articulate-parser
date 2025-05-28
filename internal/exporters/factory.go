// Package exporters provides implementations of the Exporter interface
// for converting Articulate Rise courses into various file formats.
package exporters

import (
	"fmt"
	"strings"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/services"
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
// It returns an appropriate exporter implementation based on the format string.
// Format strings are case-insensitive.
//
// Parameters:
//   - format: The desired export format (e.g., "markdown", "docx")
//
// Returns:
//   - An implementation of the Exporter interface if the format is supported
//   - An error if the format is not supported
func (f *Factory) CreateExporter(format string) (interfaces.Exporter, error) {
	switch strings.ToLower(format) {
	case "markdown", "md":
		return NewMarkdownExporter(f.htmlCleaner), nil
	case "docx", "word":
		return NewDocxExporter(f.htmlCleaner), nil
	case "html", "htm":
		return NewHTMLExporter(f.htmlCleaner), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// GetSupportedFormats returns a list of all supported export formats.
// This includes both primary format names and their aliases.
//
// Returns:
//   - A string slice containing all supported format names
func (f *Factory) GetSupportedFormats() []string {
	return []string{"markdown", "md", "docx", "word", "html", "htm"}
}
