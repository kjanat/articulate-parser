// Package services_test provides examples for the services package.
package services_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/kjanat/articulate-parser/internal/config"
	"github.com/kjanat/articulate-parser/internal/services"
)

// ExampleNewArticulateParser demonstrates creating a new parser.
func ExampleNewArticulateParser() {
	// Create a no-op logger for this example
	logger := services.NewNoOpLogger()

	// Create parser with defaults (nil config uses defaults)
	parser := services.NewArticulateParser(logger, nil)

	fmt.Printf("Parser created: %T\n", parser)
	// Output: Parser created: *services.ArticulateParser
}

// ExampleNewArticulateParser_custom demonstrates creating a parser with custom configuration.
func ExampleNewArticulateParser_custom() {
	logger := services.NewNoOpLogger()

	// Create parser with custom base URL and timeout
	cfg := &config.Config{
		BaseURL:        "https://custom.articulate.com",
		RequestTimeout: 60 * time.Second,
	}
	parser := services.NewArticulateParser(logger, cfg)

	fmt.Printf("Parser configured: %T\n", parser)
	// Output: Parser configured: *services.ArticulateParser
}

// ExampleArticulateParser_LoadCourseFromFile demonstrates loading a course from a file.
func ExampleArticulateParser_LoadCourseFromFile() {
	logger := services.NewNoOpLogger()
	parser := services.NewArticulateParser(logger, nil)

	// In a real scenario, you'd have an actual file
	// This example shows the API usage
	_, err := parser.LoadCourseFromFile("course.json")
	if err != nil {
		log.Printf("Failed to load course: %v", err)
	}
}

// ExampleArticulateParser_FetchCourse demonstrates fetching a course from a URI.
func ExampleArticulateParser_FetchCourse() {
	logger := services.NewNoOpLogger()
	parser := services.NewArticulateParser(logger, nil)

	// Create a context with timeout
	ctx := context.Background()

	// In a real scenario, you'd use an actual share URL
	_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/YOUR_SHARE_ID")
	if err != nil {
		log.Printf("Failed to fetch course: %v", err)
	}
}

// ExampleNewSlogLogger demonstrates creating a JSON logger.
func ExampleNewSlogLogger() {
	// Create a JSON structured logger at Info level
	// Output goes to stdout in JSON format
	logger := services.NewSlogLogger(slog.LevelInfo)

	// Logger is ready to use - calls like logger.Info("msg") output JSON
	fmt.Printf("Logger type: %T\n", logger)
	// Output: Logger type: *services.SlogLogger
}

// ExampleNewTextLogger demonstrates creating a human-readable text logger.
func ExampleNewTextLogger() {
	// Create a text logger, more readable for development
	logger := services.NewTextLogger(slog.LevelDebug)

	// Logger is ready - calls like logger.Debug("msg") output text
	fmt.Printf("Logger type: %T\n", logger)
	// Output: Logger type: *services.SlogLogger
}

// ExampleNewNoOpLogger demonstrates creating a no-op logger for testing.
func ExampleNewNoOpLogger() {
	logger := services.NewNoOpLogger()

	// All log calls are silently discarded
	logger.Debug("this won't appear")
	logger.Info("neither will this")
	logger.Error("or this")

	fmt.Println("Logging silently completed")
	// Output: Logging silently completed
}

// ExampleSlogLogger_With demonstrates creating a child logger with context.
func ExampleSlogLogger_With() {
	logger := services.NewNoOpLogger()

	// Create a child logger with additional context
	requestLogger := logger.With("request_id", "abc123", "user", "john")

	// All logs from requestLogger will include request_id and user
	requestLogger.Info("processing request")

	fmt.Printf("Child logger type: %T\n", requestLogger)
	// Output: Child logger type: *services.NoOpLogger
}

// ExampleHTMLCleaner demonstrates cleaning HTML content.
func ExampleHTMLCleaner() {
	cleaner := services.NewHTMLCleaner()

	html := "<p>This is <strong>bold</strong> text with entities.</p>"
	clean := cleaner.CleanHTML(html)

	fmt.Println(clean)
	// Output: This is bold text with entities.
}

// ExampleHTMLCleaner_CleanHTML demonstrates complex HTML cleaning.
func ExampleHTMLCleaner_CleanHTML() {
	cleaner := services.NewHTMLCleaner()

	html := `
		<div>
			<h1>Title</h1>
			<p>Paragraph with <a href="#">link</a> and &amp; entity.</p>
			<ul>
				<li>Item 1</li>
				<li>Item 2</li>
			</ul>
		</div>
	`
	clean := cleaner.CleanHTML(html)

	fmt.Println(clean)
	// Output: Title Paragraph with link and & entity. Item 1 Item 2
}
