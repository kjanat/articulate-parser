// Package interfaces provides the core contracts for the articulate-parser application.
// It defines interfaces for parsing and exporting Articulate Rise courses.
package interfaces

import "github.com/kjanat/articulate-parser/internal/models"

// Exporter defines the interface for exporting courses to different formats.
// Implementations of this interface handle the conversion of course data to
// specific output formats like Markdown or DOCX.
type Exporter interface {
	// Export converts a course to the supported format and writes it to the
	// specified output path. It returns an error if the export operation fails.
	Export(course *models.Course, outputPath string) error

	// SupportedFormat returns the name of the format this exporter supports.
	// This is used to identify which exporter to use for a given format.
	SupportedFormat() string
}

// ExporterFactory creates exporters for different formats.
// It acts as a factory for creating appropriate Exporter implementations
// based on the requested format.
type ExporterFactory interface {
	// CreateExporter instantiates an exporter for the specified format.
	// It returns the appropriate exporter or an error if the format is not supported.
	CreateExporter(format string) (Exporter, error)

	// SupportedFormats returns a list of all export formats supported by this factory.
	// This is used to inform users of available export options.
	SupportedFormats() []string
}
