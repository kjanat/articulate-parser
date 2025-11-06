// Package services_test provides examples for the services package.
package services_test

import (
	"context"
	"fmt"
	"log"

	"github.com/kjanat/articulate-parser/internal/services"
)

// ExampleNewArticulateParser demonstrates creating a new parser.
func ExampleNewArticulateParser() {
	// Create a no-op logger for this example
	logger := services.NewNoOpLogger()

	// Create parser with defaults
	parser := services.NewArticulateParser(logger, "", 0)

	fmt.Printf("Parser created: %T\n", parser)
	// Output: Parser created: *services.ArticulateParser
}

// ExampleNewArticulateParser_custom demonstrates creating a parser with custom configuration.
func ExampleNewArticulateParser_custom() {
	logger := services.NewNoOpLogger()

	// Create parser with custom base URL and timeout
	parser := services.NewArticulateParser(
		logger,
		"https://custom.articulate.com",
		60_000_000_000, // 60 seconds in nanoseconds
	)

	fmt.Printf("Parser configured: %T\n", parser)
	// Output: Parser configured: *services.ArticulateParser
}

// ExampleArticulateParser_LoadCourseFromFile demonstrates loading a course from a file.
func ExampleArticulateParser_LoadCourseFromFile() {
	logger := services.NewNoOpLogger()
	parser := services.NewArticulateParser(logger, "", 0)

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
	parser := services.NewArticulateParser(logger, "", 0)

	// Create a context with timeout
	ctx := context.Background()

	// In a real scenario, you'd use an actual share URL
	_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/YOUR_SHARE_ID")
	if err != nil {
		log.Printf("Failed to fetch course: %v", err)
	}
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
