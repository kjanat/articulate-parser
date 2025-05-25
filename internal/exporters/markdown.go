// Package exporters provides implementations of the Exporter interface
// for converting Articulate Rise courses into various file formats.
package exporters

import (
	"bytes"
	"fmt"
	"os"
	"strings"

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

// Export exports a course to Markdown format.
// It generates a structured Markdown document from the course data
// and writes it to the specified output path.
//
// Parameters:
//   - course: The course data model to export
//   - outputPath: The file path where the Markdown content will be written
//
// Returns:
//   - An error if writing to the output file fails
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

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// GetSupportedFormat returns the format name this exporter supports
// It indicates the file format that the MarkdownExporter can generate.
//
// Returns:
//   - A string representing the supported format ("markdown")
func (e *MarkdownExporter) GetSupportedFormat() string {
	return "markdown"
}

// processItemToMarkdown converts a course item into Markdown format
// and appends it to the provided buffer. It handles different item types
// with appropriate Markdown formatting.
//
// Parameters:
//   - buf: The buffer to write the Markdown content to
//   - item: The course item to process
//   - level: The heading level for the item (determines the number of # characters)
func (e *MarkdownExporter) processItemToMarkdown(buf *bytes.Buffer, item models.Item, level int) {
	headingPrefix := strings.Repeat("#", level)

	switch item.Type {
	case "text":
		e.processTextItem(buf, item, headingPrefix)
	case "list":
		e.processListItem(buf, item)
	case "multimedia":
		e.processMultimediaItem(buf, item, headingPrefix)
	case "image":
		e.processImageItem(buf, item, headingPrefix)
	case "knowledgeCheck":
		e.processKnowledgeCheckItem(buf, item, headingPrefix)
	case "interactive":
		e.processInteractiveItem(buf, item, headingPrefix)
	case "divider":
		e.processDividerItem(buf)
	default:
		e.processUnknownItem(buf, item, headingPrefix)
	}
}

// processTextItem handles text content with headings and paragraphs
func (e *MarkdownExporter) processTextItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	for _, subItem := range item.Items {
		if subItem.Heading != "" {
			heading := e.htmlCleaner.CleanHTML(subItem.Heading)
			if heading != "" {
				buf.WriteString(fmt.Sprintf("%s %s\n\n", headingPrefix, heading))
			}
		}
		if subItem.Paragraph != "" {
			paragraph := e.htmlCleaner.CleanHTML(subItem.Paragraph)
			if paragraph != "" {
				buf.WriteString(fmt.Sprintf("%s\n\n", paragraph))
			}
		}
	}
}

// processListItem handles list items with bullet points
func (e *MarkdownExporter) processListItem(buf *bytes.Buffer, item models.Item) {
	for _, subItem := range item.Items {
		if subItem.Paragraph != "" {
			paragraph := e.htmlCleaner.CleanHTML(subItem.Paragraph)
			if paragraph != "" {
				buf.WriteString(fmt.Sprintf("- %s\n", paragraph))
			}
		}
	}
	buf.WriteString("\n")
}

// processMultimediaItem handles multimedia content including videos and images
func (e *MarkdownExporter) processMultimediaItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	buf.WriteString(fmt.Sprintf("%s Media Content\n\n", headingPrefix))
	for _, subItem := range item.Items {
		e.processMediaSubItem(buf, subItem)
	}
	buf.WriteString("\n")
}

// processMediaSubItem processes individual media items (video/image)
func (e *MarkdownExporter) processMediaSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Media != nil {
		e.processVideoMedia(buf, subItem.Media)
		e.processImageMedia(buf, subItem.Media)
	}
	if subItem.Caption != "" {
		caption := e.htmlCleaner.CleanHTML(subItem.Caption)
		buf.WriteString(fmt.Sprintf("*%s*\n", caption))
	}
}

// processVideoMedia processes video media content
func (e *MarkdownExporter) processVideoMedia(buf *bytes.Buffer, media *models.Media) {
	if media.Video != nil {
		buf.WriteString(fmt.Sprintf("**Video**: %s\n", media.Video.OriginalUrl))
		if media.Video.Duration > 0 {
			buf.WriteString(fmt.Sprintf("**Duration**: %d seconds\n", media.Video.Duration))
		}
	}
}

// processImageMedia processes image media content
func (e *MarkdownExporter) processImageMedia(buf *bytes.Buffer, media *models.Media) {
	if media.Image != nil {
		buf.WriteString(fmt.Sprintf("**Image**: %s\n", media.Image.OriginalUrl))
	}
}

// processImageItem handles standalone image items
func (e *MarkdownExporter) processImageItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	buf.WriteString(fmt.Sprintf("%s Image\n\n", headingPrefix))
	for _, subItem := range item.Items {
		if subItem.Media != nil && subItem.Media.Image != nil {
			buf.WriteString(fmt.Sprintf("**Image**: %s\n", subItem.Media.Image.OriginalUrl))
		}
		if subItem.Caption != "" {
			caption := e.htmlCleaner.CleanHTML(subItem.Caption)
			buf.WriteString(fmt.Sprintf("*%s*\n", caption))
		}
	}
	buf.WriteString("\n")
}

// processKnowledgeCheckItem handles quiz questions and knowledge checks
func (e *MarkdownExporter) processKnowledgeCheckItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	buf.WriteString(fmt.Sprintf("%s Knowledge Check\n\n", headingPrefix))
	for _, subItem := range item.Items {
		e.processQuestionSubItem(buf, subItem)
	}
	buf.WriteString("\n")
}

// processQuestionSubItem processes individual question items
func (e *MarkdownExporter) processQuestionSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Title != "" {
		title := e.htmlCleaner.CleanHTML(subItem.Title)
		buf.WriteString(fmt.Sprintf("**Question**: %s\n\n", title))
	}

	e.processAnswers(buf, subItem.Answers)

	if subItem.Feedback != "" {
		feedback := e.htmlCleaner.CleanHTML(subItem.Feedback)
		buf.WriteString(fmt.Sprintf("\n**Feedback**: %s\n", feedback))
	}
}

// processAnswers processes answer choices for quiz questions
func (e *MarkdownExporter) processAnswers(buf *bytes.Buffer, answers []models.Answer) {
	buf.WriteString("**Answers**:\n")
	for i, answer := range answers {
		correctMark := ""
		if answer.Correct {
			correctMark = " ✓"
		}
		buf.WriteString(fmt.Sprintf("%d. %s%s\n", i+1, answer.Title, correctMark))
	}
}

// processInteractiveItem handles interactive content
func (e *MarkdownExporter) processInteractiveItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	buf.WriteString(fmt.Sprintf("%s Interactive Content\n\n", headingPrefix))
	for _, subItem := range item.Items {
		if subItem.Title != "" {
			title := e.htmlCleaner.CleanHTML(subItem.Title)
			buf.WriteString(fmt.Sprintf("**%s**\n\n", title))
		}
	}
}

// processDividerItem handles divider elements
func (e *MarkdownExporter) processDividerItem(buf *bytes.Buffer) {
	buf.WriteString("---\n\n")
}

// processUnknownItem handles unknown or unsupported item types
func (e *MarkdownExporter) processUnknownItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	if len(item.Items) > 0 {
		buf.WriteString(fmt.Sprintf("%s %s Content\n\n", headingPrefix, strings.Title(item.Type)))
		for _, subItem := range item.Items {
			e.processGenericSubItem(buf, subItem)
		}
	}
}

// processGenericSubItem processes sub-items for unknown types
func (e *MarkdownExporter) processGenericSubItem(buf *bytes.Buffer, subItem models.SubItem) {
	if subItem.Title != "" {
		title := e.htmlCleaner.CleanHTML(subItem.Title)
		buf.WriteString(fmt.Sprintf("**%s**\n\n", title))
	}
	if subItem.Paragraph != "" {
		paragraph := e.htmlCleaner.CleanHTML(subItem.Paragraph)
		buf.WriteString(fmt.Sprintf("%s\n\n", paragraph))
	}
}
