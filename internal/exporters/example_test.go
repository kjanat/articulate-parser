// Package exporters_test provides examples for the exporters package.
package exporters_test

import (
	"fmt"
	"log"

	"github.com/kjanat/articulate-parser/internal/exporters"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// ExampleNewFactory demonstrates creating an exporter factory.
func ExampleNewFactory() {
	htmlCleaner := services.NewHTMLCleaner()
	factory := exporters.NewFactory(htmlCleaner)

	// Get supported formats
	formats := factory.SupportedFormats()
	fmt.Printf("Supported formats: %d\n", len(formats))
	// Output: Supported formats: 6
}

// ExampleFactory_CreateExporter demonstrates creating exporters.
func ExampleFactory_CreateExporter() {
	htmlCleaner := services.NewHTMLCleaner()
	factory := exporters.NewFactory(htmlCleaner)

	// Create a markdown exporter
	exporter, err := factory.CreateExporter("markdown")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created: %s exporter\n", exporter.SupportedFormat())
	// Output: Created: markdown exporter
}

// ExampleFactory_CreateExporter_caseInsensitive demonstrates case-insensitive format names.
func ExampleFactory_CreateExporter_caseInsensitive() {
	htmlCleaner := services.NewHTMLCleaner()
	factory := exporters.NewFactory(htmlCleaner)

	// All these work (case-insensitive)
	formats := []string{"MARKDOWN", "Markdown", "markdown", "MD"}

	for _, format := range formats {
		exporter, _ := factory.CreateExporter(format)
		fmt.Printf("%s -> %s\n", format, exporter.SupportedFormat())
	}
	// Output:
	// MARKDOWN -> markdown
	// Markdown -> markdown
	// markdown -> markdown
	// MD -> markdown
}

// ExampleMarkdownExporter_Export demonstrates exporting to Markdown.
func ExampleMarkdownExporter_Export() {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := exporters.NewMarkdownExporter(htmlCleaner)

	course := &models.Course{
		ShareID: "example-id",
		Course: models.CourseInfo{
			Title:       "Example Course",
			Description: "<p>Course description</p>",
		},
	}

	// Export to markdown file
	err := exporter.Export(course, "output.md")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Export complete")
	// Output: Export complete
}

// ExampleDocxExporter_Export demonstrates exporting to DOCX.
func ExampleDocxExporter_Export() {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := exporters.NewDocxExporter(htmlCleaner)

	course := &models.Course{
		ShareID: "example-id",
		Course: models.CourseInfo{
			Title: "Example Course",
		},
	}

	// Export to Word document
	err := exporter.Export(course, "output.docx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DOCX export complete")
	// Output: DOCX export complete
}
