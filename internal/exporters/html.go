package exporters

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

//go:embed html_styles.css
var defaultCSS string

//go:embed html_template.html
var htmlTemplate string

// HTMLExporter implements the Exporter interface for HTML format.
// It converts Articulate Rise course data into a structured HTML document using templates.
type HTMLExporter struct {
	// htmlCleaner is used to convert HTML content to plain text when needed
	htmlCleaner *services.HTMLCleaner
	// tmpl holds the parsed HTML template
	tmpl *template.Template
}

// NewHTMLExporter creates a new HTMLExporter instance.
// It takes an HTMLCleaner to handle HTML content conversion when plain text is needed.
//
// Parameters:
//   - htmlCleaner: Service for cleaning HTML content in course data
//
// Returns:
//   - An implementation of the Exporter interface for HTML format
func NewHTMLExporter(htmlCleaner *services.HTMLCleaner) interfaces.Exporter {
	// Parse the template with custom functions
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s) // #nosec G203 - HTML content is from trusted course data
		},
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s) // #nosec G203 - CSS content is from trusted embedded file
		},
	}

	tmpl := template.Must(template.New("html").Funcs(funcMap).Parse(htmlTemplate))

	return &HTMLExporter{
		htmlCleaner: htmlCleaner,
		tmpl:        tmpl,
	}
}

// Export exports a course to HTML format.
// It generates a structured HTML document from the course data
// and writes it to the specified output path.
//
// Parameters:
//   - course: The course data model to export
//   - outputPath: The file path where the HTML content will be written
//
// Returns:
//   - An error if writing to the output file fails
func (e *HTMLExporter) Export(course *models.Course, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	return e.WriteHTML(f, course)
}

// WriteHTML writes the HTML content to an io.Writer.
// This allows for better testability and flexibility in output destinations.
//
// Parameters:
//   - w: The writer to output HTML content to
//   - course: The course data model to export
//
// Returns:
//   - An error if writing fails
func (e *HTMLExporter) WriteHTML(w io.Writer, course *models.Course) error {
	// Prepare template data
	data := prepareTemplateData(course, e.htmlCleaner)

	// Execute template
	if err := e.tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// SupportedFormat returns the format name this exporter supports
// It indicates the file format that the HTMLExporter can generate.
//
// Returns:
//   - A string representing the supported format ("html")
func (e *HTMLExporter) SupportedFormat() string {
	return FormatHTML
}
