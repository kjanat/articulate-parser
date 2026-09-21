// Package exporters_test provides examples for the exporters package.
package exporters_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kjanat/articulate-parser/internal/exporters"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// exportToTemp runs export against a path in a fresh temp dir and removes the
// dir afterwards. Example functions can't take *testing.T, so t.TempDir is
// unavailable here.
func exportToTemp(name string, export func(path string) error) {
	dir, err := os.MkdirTemp("", "exporters-example-*")
	if err != nil {
		log.Fatal(err)
	}
	err = export(filepath.Join(dir, name))
	_ = os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleNewFactory demonstrates creating an exporter factory.
func ExampleNewFactory() {
	htmlCleaner := services.NewHTMLCleaner()
	factory := exporters.NewFactory(htmlCleaner, nil)

	// Get supported formats
	formats := factory.SupportedFormats()
	fmt.Printf("Supported formats: %d\n", len(formats))
	// Output: Supported formats: 6
}

// ExampleFactory_CreateExporter demonstrates creating exporters.
func ExampleFactory_CreateExporter() {
	htmlCleaner := services.NewHTMLCleaner()
	factory := exporters.NewFactory(htmlCleaner, nil)

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
	factory := exporters.NewFactory(htmlCleaner, nil)

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
	exporter := exporters.NewMarkdownExporter(htmlCleaner, nil)

	course := &models.Course{
		ShareID: "example-id",
		Course: models.CourseInfo{
			Title:       "Example Course",
			Description: "<p>Course description</p>",
		},
	}

	// Export to markdown file
	exportToTemp("output.md", func(path string) error {
		return exporter.Export(course, path)
	})

	fmt.Println("Export complete")
	// Output: Export complete
}

// ExampleDocxExporter_Export demonstrates exporting to DOCX.
func ExampleDocxExporter_Export() {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := exporters.NewDocxExporter(htmlCleaner, nil)

	course := &models.Course{
		ShareID: "example-id",
		Course: models.CourseInfo{
			Title: "Example Course",
		},
	}

	// Export to Word document
	exportToTemp("output.docx", func(path string) error {
		return exporter.Export(course, path)
	})

	fmt.Println("DOCX export complete")
	// Output: DOCX export complete
}

// ExampleHTMLExporter_Export demonstrates exporting to HTML.
func ExampleHTMLExporter_Export() {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := exporters.NewHTMLExporter(htmlCleaner, nil)

	course := &models.Course{
		ShareID: "example-id",
		Course: models.CourseInfo{
			Title:       "Example Course",
			Description: "<p>A sample course</p>",
			Lessons: []models.Lesson{
				{
					ID:    "lesson-1",
					Title: "Introduction",
					Type:  "lesson",
				},
			},
		},
	}

	// Export to HTML file
	exportToTemp("output.html", func(path string) error {
		return exporter.Export(course, path)
	})

	fmt.Println("HTML export complete")
	// Output: HTML export complete
}
