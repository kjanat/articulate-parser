// Package services provides the core functionality for the articulate-parser application.
// It implements the interfaces defined in the interfaces package.
package services

import (
	"context"
	"fmt"

	"github.com/kjanat/articulate-parser/internal/apperrors"
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
// It takes a context for cancellation, the path to the course file, the desired export format,
// and the output file path. Returns an error if loading or exporting fails, or if the context
// is canceled.
func (a *App) ProcessCourseFromFile(ctx context.Context, filePath, format, outputPath string) error {
	// Check for cancellation before loading
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrOperationCanceled, err)
	}

	course, err := a.parser.LoadCourseFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load course from file: %w", err)
	}

	return a.exportCourse(ctx, course, format, outputPath)
}

// ProcessCourseFromURI fetches a course from the provided URI and exports it to the specified format.
// It takes a context for cancellation, the URI to fetch the course from, the desired export format,
// and the output file path. Returns an error if fetching or exporting fails, or if the context
// is canceled.
func (a *App) ProcessCourseFromURI(ctx context.Context, uri, format, outputPath string) error {
	course, err := a.parser.FetchCourse(ctx, uri)
	if err != nil {
		return fmt.Errorf("failed to fetch course: %w", err)
	}

	return a.exportCourse(ctx, course, format, outputPath)
}

// exportCourse exports a course to the specified format and output path.
// It's a helper method that creates the appropriate exporter and performs the export.
// Returns an error if creating the exporter or exporting the course fails, or if the
// context is canceled.
func (a *App) exportCourse(ctx context.Context, course *models.Course, format, outputPath string) error {
	// Check for cancellation before export
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrOperationCanceled, err)
	}

	exporter, err := a.exporterFactory.CreateExporter(format)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	if err := exporter.Export(course, outputPath); err != nil {
		return fmt.Errorf("failed to export course: %w", err)
	}

	return nil
}

// SupportedFormats returns a list of all export formats supported by the application.
// This information is provided by the ExporterFactory.
func (a *App) SupportedFormats() []string {
	return a.exporterFactory.SupportedFormats()
}
