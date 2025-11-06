package exporters

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// HTMLExporter implements the Exporter interface for HTML format.
// It converts Articulate Rise course data into a structured HTML document.
type HTMLExporter struct {
	// htmlCleaner is used to convert HTML content to plain text when needed
	htmlCleaner *services.HTMLCleaner
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
	return &HTMLExporter{
		htmlCleaner: htmlCleaner,
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
	var buf bytes.Buffer

	// Write HTML document structure
	buf.WriteString("<!DOCTYPE html>\n")
	buf.WriteString("<html lang=\"en\">\n")
	buf.WriteString("<head>\n")
	buf.WriteString("    <meta charset=\"UTF-8\">\n")
	buf.WriteString("    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	buf.WriteString(fmt.Sprintf("    <title>%s</title>\n", html.EscapeString(course.Course.Title)))
	buf.WriteString("    <style>\n")
	buf.WriteString(e.getDefaultCSS())
	buf.WriteString("    </style>\n")
	buf.WriteString("</head>\n")
	buf.WriteString("<body>\n")

	// Write course header
	buf.WriteString(fmt.Sprintf("    <header>\n        <h1>%s</h1>\n", html.EscapeString(course.Course.Title)))

	if course.Course.Description != "" {
		buf.WriteString(fmt.Sprintf("        <div class=\"course-description\">%s</div>\n", course.Course.Description))
	}
	buf.WriteString("    </header>\n\n")

	// Add metadata section
	buf.WriteString("    <section class=\"course-info\">\n")
	buf.WriteString("        <h2>Course Information</h2>\n")
	buf.WriteString("        <ul>\n")
	buf.WriteString(fmt.Sprintf("            <li><strong>Course ID:</strong> %s</li>\n", html.EscapeString(course.Course.ID)))
	buf.WriteString(fmt.Sprintf("            <li><strong>Share ID:</strong> %s</li>\n", html.EscapeString(course.ShareID)))
	buf.WriteString(fmt.Sprintf("            <li><strong>Navigation Mode:</strong> %s</li>\n", html.EscapeString(course.Course.NavigationMode)))
	if course.Course.ExportSettings != nil {
		buf.WriteString(fmt.Sprintf("            <li><strong>Export Format:</strong> %s</li>\n", html.EscapeString(course.Course.ExportSettings.Format)))
	}
	buf.WriteString("        </ul>\n")
	buf.WriteString("    </section>\n\n")

	// Process lessons
	lessonCounter := 0
	for _, lesson := range course.Course.Lessons {
		if lesson.Type == "section" {
			buf.WriteString(fmt.Sprintf("    <section class=\"course-section\">\n        <h2>%s</h2>\n    </section>\n\n", html.EscapeString(lesson.Title)))
			continue
		}

		lessonCounter++
		buf.WriteString(fmt.Sprintf("    <section class=\"lesson\">\n        <h3>Lesson %d: %s</h3>\n", lessonCounter, html.EscapeString(lesson.Title)))

		if lesson.Description != "" {
			buf.WriteString(fmt.Sprintf("        <div class=\"lesson-description\">%s</div>\n", lesson.Description))
		}

		// Process lesson items
		for _, item := range lesson.Items {
			e.processItemToHTML(&buf, item)
		}

		buf.WriteString("    </section>\n\n")
	}

	buf.WriteString("</body>\n")
	buf.WriteString("</html>\n")

	// #nosec G306 - 0644 is appropriate for export files that should be readable by others
	if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
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

// getDefaultCSS returns basic CSS styling for the HTML document.
func (e *HTMLExporter) getDefaultCSS() string {
	return `
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
        }
        header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem;
            border-radius: 10px;
            margin-bottom: 2rem;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        header h1 {
            margin: 0;
            font-size: 2.5rem;
            font-weight: 300;
        }
        .course-description {
            margin-top: 1rem;
            font-size: 1.1rem;
            opacity: 0.9;
        }
        .course-info {
            background: white;
            padding: 1.5rem;
            border-radius: 8px;
            margin-bottom: 2rem;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .course-info h2 {
            margin-top: 0;
            color: #4a5568;
            border-bottom: 2px solid #e2e8f0;
            padding-bottom: 0.5rem;
        }
        .course-info ul {
            list-style: none;
            padding: 0;
        }
        .course-info li {
            margin: 0.5rem 0;
            padding: 0.5rem;
            background: #f7fafc;
            border-radius: 4px;
        }
        .course-section {
            background: #4299e1;
            color: white;
            padding: 1.5rem;
            border-radius: 8px;
            margin: 2rem 0;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .course-section h2 {
            margin: 0;
            font-weight: 400;
        }
        .lesson {
            background: white;
            padding: 2rem;
            border-radius: 8px;
            margin: 2rem 0;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            border-left: 4px solid #4299e1;
        }
        .lesson h3 {
            margin-top: 0;
            color: #2d3748;
            font-size: 1.5rem;
        }
        .lesson-description {
            margin: 1rem 0;
            padding: 1rem;
            background: #f7fafc;
            border-radius: 4px;
            border-left: 3px solid #4299e1;
        }
        .item {
            margin: 1.5rem 0;
            padding: 1rem;
            border-radius: 6px;
            background: #fafafa;
            border: 1px solid #e2e8f0;
        }
        .item h4 {
            margin-top: 0;
            color: #4a5568;
            font-size: 1.2rem;
            text-transform: capitalize;
        }
        .text-item {
            background: #f0fff4;
            border-left: 3px solid #48bb78;
        }
        .list-item {
            background: #fffaf0;
            border-left: 3px solid #ed8936;
        }
        .knowledge-check {
            background: #e6fffa;
            border-left: 3px solid #38b2ac;
        }
        .multimedia-item {
            background: #faf5ff;
            border-left: 3px solid #9f7aea;
        }
        .interactive-item {
            background: #fff5f5;
            border-left: 3px solid #f56565;
        }
        .unknown-item {
            background: #f7fafc;
            border-left: 3px solid #a0aec0;
        }
        .answers {
            margin: 1rem 0;
        }
        .answers h5 {
            margin: 0.5rem 0;
            color: #4a5568;
        }
        .answers ol {
            margin: 0.5rem 0;
            padding-left: 1.5rem;
        }
        .answers li {
            margin: 0.3rem 0;
            padding: 0.3rem;
        }
        .correct-answer {
            background: #c6f6d5;
            border-radius: 3px;
            font-weight: bold;
        }
        .correct-answer::after {
            content: " ✓";
            color: #38a169;
        }
        .feedback {
            margin: 1rem 0;
            padding: 1rem;
            background: #edf2f7;
            border-radius: 4px;
            border-left: 3px solid #4299e1;
            font-style: italic;
        }
        .media-info {
            background: #edf2f7;
            padding: 1rem;
            border-radius: 4px;
            margin: 0.5rem 0;
        }
        .media-info strong {
            color: #4a5568;
        }
        hr {
            border: none;
            height: 2px;
            background: linear-gradient(to right, #667eea, #764ba2);
            margin: 2rem 0;
            border-radius: 1px;
        }
        ul {
            padding-left: 1.5rem;
        }
        li {
            margin: 0.5rem 0;
        }
    `
}

// processItemToHTML converts a course item into HTML format
// and appends it to the provided buffer. It handles different item types
// with appropriate HTML formatting.
//
// Parameters:
//   - buf: The buffer to write the HTML content to
//   - item: The course item to process
func (e *HTMLExporter) processItemToHTML(buf *bytes.Buffer, item models.Item) {
	switch strings.ToLower(item.Type) {
	case "text":
		e.processTextItem(buf, item)
	case "list":
		e.processListItem(buf, item)
	case "knowledgecheck":
		e.processKnowledgeCheckItem(buf, item)
	case "multimedia":
		e.processMultimediaItem(buf, item)
	case "image":
		e.processImageItem(buf, item)
	case "interactive":
		e.processInteractiveItem(buf, item)
	case "divider":
		e.processDividerItem(buf)
	default:
		e.processUnknownItem(buf, item)
	}
}

// processTextItem handles text content with headings and paragraphs.
func (e *HTMLExporter) processTextItem(buf *bytes.Buffer, item models.Item) {
	buf.WriteString("        <div class=\"item text-item\">\n")
	buf.WriteString("            <h4>Text Content</h4>\n")
	for _, subItem := range item.Items {
		if subItem.Heading != "" {
			fmt.Fprintf(buf, "            <h5>%s</h5>\n", subItem.Heading)
		}
		if subItem.Paragraph != "" {
			fmt.Fprintf(buf, "            <div>%s</div>\n", subItem.Paragraph)
		}
	}
	buf.WriteString("        </div>\n\n")
}

// processListItem handles list content.
func (e *HTMLExporter) processListItem(buf *bytes.Buffer, item models.Item) {
	buf.WriteString("        <div class=\"item list-item\">\n")
	buf.WriteString("            <h4>List</h4>\n")
	buf.WriteString("            <ul>\n")
	for _, subItem := range item.Items {
		if subItem.Paragraph != "" {
			cleanText := e.htmlCleaner.CleanHTML(subItem.Paragraph)
			fmt.Fprintf(buf, "                <li>%s</li>\n", html.EscapeString(cleanText))
		}
	}
	buf.WriteString("            </ul>\n")
	buf.WriteString("        </div>\n\n")
}

// processKnowledgeCheckItem handles quiz questions and answers.
func (e *HTMLExporter) processKnowledgeCheckItem(buf *bytes.Buffer, item models.Item) {
	buf.WriteString("        <div class=\"item knowledge-check\">\n")
	buf.WriteString("            <h4>Knowledge Check</h4>\n")
	for _, subItem := range item.Items {
		if subItem.Title != "" {
			fmt.Fprintf(buf, "            <p><strong>Question:</strong> %s</p>\n", subItem.Title)
		}
		if len(subItem.Answers) > 0 {
			e.processAnswers(buf, subItem.Answers)
		}
		if subItem.Feedback != "" {
			fmt.Fprintf(buf, "            <div class=\"feedback\"><strong>Feedback:</strong> %s</div>\n", subItem.Feedback)
		}
	}
	buf.WriteString("        </div>\n\n")
}

// processMultimediaItem handles multimedia content like videos.
func (e *HTMLExporter) processMultimediaItem(buf *bytes.Buffer, item models.Item) {
	buf.WriteString("        <div class=\"item multimedia-item\">\n")
	buf.WriteString("            <h4>Media Content</h4>\n")
	for _, subItem := range item.Items {
		if subItem.Title != "" {
			fmt.Fprintf(buf, "            <h5>%s</h5>\n", subItem.Title)
		}
		if subItem.Media != nil {
			if subItem.Media.Video != nil {
				buf.WriteString("            <div class=\"media-info\">\n")
				fmt.Fprintf(buf, "                <p><strong>Video:</strong> %s</p>\n", html.EscapeString(subItem.Media.Video.OriginalURL))
				if subItem.Media.Video.Duration > 0 {
					fmt.Fprintf(buf, "                <p><strong>Duration:</strong> %d seconds</p>\n", subItem.Media.Video.Duration)
				}
				buf.WriteString("            </div>\n")
			}
		}
		if subItem.Caption != "" {
			fmt.Fprintf(buf, "            <div><em>%s</em></div>\n", subItem.Caption)
		}
	}
	buf.WriteString("        </div>\n\n")
}

// processImageItem handles image content.
func (e *HTMLExporter) processImageItem(buf *bytes.Buffer, item models.Item) {
	buf.WriteString("        <div class=\"item multimedia-item\">\n")
	buf.WriteString("            <h4>Image</h4>\n")
	for _, subItem := range item.Items {
		if subItem.Media != nil && subItem.Media.Image != nil {
			buf.WriteString("            <div class=\"media-info\">\n")
			fmt.Fprintf(buf, "                <p><strong>Image:</strong> %s</p>\n", html.EscapeString(subItem.Media.Image.OriginalURL))
			buf.WriteString("            </div>\n")
		}
		if subItem.Caption != "" {
			fmt.Fprintf(buf, "            <div><em>%s</em></div>\n", subItem.Caption)
		}
	}
	buf.WriteString("        </div>\n\n")
}

// processInteractiveItem handles interactive content.
func (e *HTMLExporter) processInteractiveItem(buf *bytes.Buffer, item models.Item) {
	buf.WriteString("        <div class=\"item interactive-item\">\n")
	buf.WriteString("            <h4>Interactive Content</h4>\n")
	for _, subItem := range item.Items {
		if subItem.Title != "" {
			fmt.Fprintf(buf, "            <p><strong>%s</strong></p>\n", subItem.Title)
		}
		if subItem.Paragraph != "" {
			fmt.Fprintf(buf, "            <div>%s</div>\n", subItem.Paragraph)
		}
	}
	buf.WriteString("        </div>\n\n")
}

// processDividerItem handles divider elements.
func (e *HTMLExporter) processDividerItem(buf *bytes.Buffer) {
	buf.WriteString("        <hr>\n\n")
}

// processUnknownItem handles unknown or unsupported item types.
func (e *HTMLExporter) processUnknownItem(buf *bytes.Buffer, item models.Item) {
	if len(item.Items) > 0 {
		buf.WriteString("        <div class=\"item unknown-item\">\n")
		caser := cases.Title(language.English)
		fmt.Fprintf(buf, "            <h4>%s Content</h4>\n", caser.String(item.Type))
		for _, subItem := range item.Items {
			e.processGenericSubItem(buf, subItem)
		}
		buf.WriteString("        </div>\n\n")
	}
}

// processGenericSubItem processes sub-items for unknown types.
func (e *HTMLExporter) processGenericSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Title != "" {
		fmt.Fprintf(buf, "            <p><strong>%s</strong></p>\n", subItem.Title)
	}
	if subItem.Paragraph != "" {
		fmt.Fprintf(buf, "            <div>%s</div>\n", subItem.Paragraph)
	}
}

// processAnswers processes answer choices for quiz questions.
func (e *HTMLExporter) processAnswers(buf *bytes.Buffer, answers []models.Answer) {
	buf.WriteString("            <div class=\"answers\">\n")
	buf.WriteString("                <h5>Answers:</h5>\n")
	buf.WriteString("                <ol>\n")
	for _, answer := range answers {
		cssClass := ""
		if answer.Correct {
			cssClass = " class=\"correct-answer\""
		}
		fmt.Fprintf(buf, "                    <li%s>%s</li>\n", cssClass, html.EscapeString(answer.Title))
	}
	buf.WriteString("                </ol>\n")
	buf.WriteString("            </div>\n")
}
