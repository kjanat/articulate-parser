package exporters

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// MarkdownExporter implements the Exporter interface for Markdown format.
// It converts Articulate Rise course data into a structured Markdown document.
type MarkdownExporter struct {
	// htmlCleaner is used to convert HTML content to plain text
	htmlCleaner *services.HTMLCleaner
}

// NewMarkdownExporter creates a new MarkdownExporter instance.
// It takes an HTMLCleaner to handle HTML content conversion.
//
// Parameters:
//   - htmlCleaner: Service for cleaning HTML content in course data
//
// Returns:
//   - An implementation of the Exporter interface for Markdown format
func NewMarkdownExporter(htmlCleaner *services.HTMLCleaner) interfaces.Exporter {
	return &MarkdownExporter{
		htmlCleaner: htmlCleaner,
	}
}

// Export converts the course to Markdown format and writes it to the output path.
func (e *MarkdownExporter) Export(course *models.Course, outputPath string) error {
	var buf bytes.Buffer

	// Write course header
	buf.WriteString(fmt.Sprintf("# %s\n\n", course.Course.Title))

	if course.Course.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", e.htmlCleaner.CleanHTML(course.Course.Description)))
	}

	// Add metadata
	buf.WriteString("## Course Information\n\n")
	buf.WriteString(fmt.Sprintf("- **Course ID**: %s\n", course.Course.ID))
	buf.WriteString(fmt.Sprintf("- **Share ID**: %s\n", course.ShareID))
	buf.WriteString(fmt.Sprintf("- **Navigation Mode**: %s\n", course.Course.NavigationMode))
	if course.Course.ExportSettings != nil {
		buf.WriteString(fmt.Sprintf("- **Export Format**: %s\n", course.Course.ExportSettings.Format))
	}
	buf.WriteString("\n---\n\n")

	// Process lessons
	lessonCounter := 0
	for _, lesson := range course.Course.Lessons {
		if lesson.Type == "section" {
			buf.WriteString(fmt.Sprintf("# %s\n\n", lesson.Title))
			continue
		}

		lessonCounter++
		buf.WriteString(fmt.Sprintf("## Lesson %d: %s\n\n", lessonCounter, lesson.Title))

		if lesson.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", e.htmlCleaner.CleanHTML(lesson.Description)))
		}

		// Process lesson items
		for _, item := range lesson.Items {
			e.processItemToMarkdown(&buf, item, 3)
		}

		buf.WriteString("\n---\n\n")
	}

	// #nosec G306 - 0644 is appropriate for export files that should be readable by others
	if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}
	return nil
}

// SupportedFormat returns "markdown".
func (e *MarkdownExporter) SupportedFormat() string {
	return FormatMarkdown
}

// processItemToMarkdown converts a course item into Markdown format.
// The level parameter determines the heading level (number of # characters).
func (e *MarkdownExporter) processItemToMarkdown(buf *bytes.Buffer, item models.Item, level int) {
	headingPrefix := strings.Repeat("#", level)

	// Normalize item type to lowercase for consistent matching
	itemType := strings.ToLower(item.Type)

	switch itemType {
	case "text":
		e.processTextItem(buf, item, headingPrefix)
	case "list":
		e.processListItem(buf, item)
	case "multimedia":
		e.processMultimediaItem(buf, item, headingPrefix)
	case "image":
		e.processImageItem(buf, item, headingPrefix)
	case "knowledgecheck":
		e.processKnowledgeCheckItem(buf, item, headingPrefix)
	case "interactive":
		e.processInteractiveItem(buf, item, headingPrefix)
	case "divider":
		e.processDividerItem(buf)
	default:
		e.processUnknownItem(buf, item, headingPrefix)
	}
}

// processTextItem handles text content with headings and paragraphs.
func (e *MarkdownExporter) processTextItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	for _, subItem := range item.Items {
		if subItem.Heading != "" {
			heading := e.htmlCleaner.CleanHTML(subItem.Heading)
			if heading != "" {
				fmt.Fprintf(buf, "%s %s\n\n", headingPrefix, heading)
			}
		}
		if subItem.Paragraph != "" {
			paragraph := e.htmlCleaner.CleanHTML(subItem.Paragraph)
			if paragraph != "" {
				fmt.Fprintf(buf, "%s\n\n", paragraph)
			}
		}
	}
}

// processListItem handles list items with bullet points.
func (e *MarkdownExporter) processListItem(buf *bytes.Buffer, item models.Item) {
	for _, subItem := range item.Items {
		if subItem.Paragraph != "" {
			paragraph := e.htmlCleaner.CleanHTML(subItem.Paragraph)
			if paragraph != "" {
				fmt.Fprintf(buf, "- %s\n", paragraph)
			}
		}
	}
	buf.WriteString("\n")
}

// processMultimediaItem handles multimedia content including videos and images.
func (e *MarkdownExporter) processMultimediaItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	fmt.Fprintf(buf, "%s Media Content\n\n", headingPrefix)
	for _, subItem := range item.Items {
		e.processMediaSubItem(buf, subItem)
	}
	buf.WriteString("\n")
}

// processMediaSubItem processes individual media items (video/image).
func (e *MarkdownExporter) processMediaSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Media != nil {
		e.processVideoMedia(buf, subItem.Media)
		e.processImageMedia(buf, subItem.Media)
	}
	if subItem.Caption != "" {
		caption := e.htmlCleaner.CleanHTML(subItem.Caption)
		fmt.Fprintf(buf, "*%s*\n", caption)
	}
}

// processVideoMedia processes video media content.
func (e *MarkdownExporter) processVideoMedia(buf *bytes.Buffer, media *models.Media) {
	if media.Video != nil {
		fmt.Fprintf(buf, "**Video**: %s\n", media.Video.OriginalURL)
		if media.Video.Duration > 0 {
			fmt.Fprintf(buf, "**Duration**: %d seconds\n", media.Video.Duration)
		}
	}
}

// processImageMedia processes image media content.
func (e *MarkdownExporter) processImageMedia(buf *bytes.Buffer, media *models.Media) {
	if media.Image != nil {
		fmt.Fprintf(buf, "**Image**: %s\n", media.Image.OriginalURL)
	}
}

// processImageItem handles standalone image items.
func (e *MarkdownExporter) processImageItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	fmt.Fprintf(buf, "%s Image\n\n", headingPrefix)
	for _, subItem := range item.Items {
		if subItem.Media != nil && subItem.Media.Image != nil {
			fmt.Fprintf(buf, "**Image**: %s\n", subItem.Media.Image.OriginalURL)
		}
		if subItem.Caption != "" {
			caption := e.htmlCleaner.CleanHTML(subItem.Caption)
			fmt.Fprintf(buf, "*%s*\n", caption)
		}
	}
	buf.WriteString("\n")
}

// processKnowledgeCheckItem handles quiz questions and knowledge checks.
func (e *MarkdownExporter) processKnowledgeCheckItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	fmt.Fprintf(buf, "%s Knowledge Check\n\n", headingPrefix)
	for _, subItem := range item.Items {
		e.processQuestionSubItem(buf, subItem)
	}
	buf.WriteString("\n")
}

// processQuestionSubItem processes individual question items.
func (e *MarkdownExporter) processQuestionSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Title != "" {
		title := e.htmlCleaner.CleanHTML(subItem.Title)
		fmt.Fprintf(buf, "**Question**: %s\n\n", title)
	}

	e.processAnswers(buf, subItem.Answers)

	if subItem.Feedback != "" {
		feedback := e.htmlCleaner.CleanHTML(subItem.Feedback)
		fmt.Fprintf(buf, "\n**Feedback**: %s\n", feedback)
	}
}

// processAnswers processes answer choices for quiz questions.
func (e *MarkdownExporter) processAnswers(buf *bytes.Buffer, answers []models.Answer) {
	buf.WriteString("**Answers**:\n")
	for i, answer := range answers {
		correctMark := ""
		if answer.Correct {
			correctMark = " ✓"
		}
		fmt.Fprintf(buf, "%d. %s%s\n", i+1, answer.Title, correctMark)
	}
}

// processInteractiveItem handles interactive content.
func (e *MarkdownExporter) processInteractiveItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	fmt.Fprintf(buf, "%s Interactive Content\n\n", headingPrefix)
	for _, subItem := range item.Items {
		if subItem.Title != "" {
			title := e.htmlCleaner.CleanHTML(subItem.Title)
			fmt.Fprintf(buf, "**%s**\n\n", title)
		}
	}
}

// processDividerItem handles divider elements.
func (e *MarkdownExporter) processDividerItem(buf *bytes.Buffer) {
	buf.WriteString("---\n\n")
}

// processUnknownItem handles unknown or unsupported item types.
func (e *MarkdownExporter) processUnknownItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	if len(item.Items) > 0 {
		caser := cases.Title(language.English)
		fmt.Fprintf(buf, "%s %s Content\n\n", headingPrefix, caser.String(item.Type))
		for _, subItem := range item.Items {
			e.processGenericSubItem(buf, subItem)
		}
	}
}

// processGenericSubItem processes sub-items for unknown types.
func (e *MarkdownExporter) processGenericSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Title != "" {
		title := e.htmlCleaner.CleanHTML(subItem.Title)
		fmt.Fprintf(buf, "**%s**\n\n", title)
	}
	if subItem.Paragraph != "" {
		paragraph := e.htmlCleaner.CleanHTML(subItem.Paragraph)
		fmt.Fprintf(buf, "%s\n\n", paragraph)
	}
}
