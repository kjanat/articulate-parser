// Package services provides the core functionality for the articulate-parser application.
// It implements the interfaces defined in the interfaces package.
package services

import (
	"fmt"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
)

// App represents the main application service that coordinates the parsing
// and exporting of Articulate Rise courses. It serves as the primary entry
// point for the application's functionality.
type App struct {
	// parser is responsible for loading course data from files or URLs
	parser interfaces.CourseParser
	// exporterFactory creates the appropriate exporter for a given format
	exporterFactory interfaces.ExporterFactory
}

// NewApp creates a new application instance with dependency injection.
// It takes a CourseParser for loading courses and an ExporterFactory for
// creating the appropriate exporters.
func NewApp(parser interfaces.CourseParser, exporterFactory interfaces.ExporterFactory) *App {
	return &App{
		parser:          parser,
		exporterFactory: exporterFactory,
	}
}

// ProcessCourseFromFile loads a course from a local file and exports it to the specified format.
// It takes the path to the course file, the desired export format, and the output file path.
// Returns an error if loading or exporting fails.
func (a *App) ProcessCourseFromFile(filePath, format, outputPath string) error {
	course, err := a.parser.LoadCourseFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load course from file: %w", err)
	}

	return a.exportCourse(course, format, outputPath)
}

// ProcessCourseFromURI fetches a course from the provided URI and exports it to the specified format.
// It takes the URI to fetch the course from, the desired export format, and the output file path.
// Returns an error if fetching or exporting fails.
func (a *App) ProcessCourseFromURI(uri, format, outputPath string) error {
	course, err := a.parser.FetchCourse(uri)
	if err != nil {
		return fmt.Errorf("failed to fetch course: %w", err)
	}

	return a.exportCourse(course, format, outputPath)
}

// exportCourse exports a course to the specified format and output path.
// It's a helper method that creates the appropriate exporter and performs the export.
// Returns an error if creating the exporter or exporting the course fails.
func (a *App) exportCourse(course *models.Course, format, outputPath string) error {
	exporter, err := a.exporterFactory.CreateExporter(format)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	if err := exporter.Export(course, outputPath); err != nil {
		return fmt.Errorf("failed to export course: %w", err)
	}

	return nil
}

// GetSupportedFormats returns a list of all export formats supported by the application.
// This information is provided by the ExporterFactory.
func (a *App) GetSupportedFormats() []string {
	return a.exporterFactory.GetSupportedFormats()
}
