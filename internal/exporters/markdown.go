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
)

// MarkdownExporter implements the Exporter interface for Markdown format.
// It converts Articulate Rise course data into a structured Markdown document.
type MarkdownExporter struct {
	// htmlCleaner is used to convert HTML content to plain text
	htmlCleaner interfaces.HTMLCleaner
	// logger is used for logging warnings and errors
	logger interfaces.Logger
}

// NewMarkdownExporter creates a new MarkdownExporter instance.
// It takes an HTMLCleaner to handle HTML content conversion and a logger for logging.
//
// Parameters:
//   - htmlCleaner: Service for cleaning HTML content in course data
//   - logger: Logger instance for warnings and errors
//
// Returns:
//   - An implementation of the Exporter interface for Markdown format
func NewMarkdownExporter(htmlCleaner interfaces.HTMLCleaner, logger interfaces.Logger) interfaces.Exporter {
	return &MarkdownExporter{
		htmlCleaner: htmlCleaner,
		logger:      logger,
	}
}

// Export converts the course to Markdown format and writes it to the output path.
func (e *MarkdownExporter) Export(course *models.Course, outputPath string) error {
	var buf bytes.Buffer

	// Write course header
	fmt.Fprintf(&buf, "# %s\n\n", strings.TrimSpace(course.Course.Title))

	if course.Course.Description != "" {
		if desc := e.htmlCleaner.CleanAndTrim(course.Course.Description); desc != "" {
			fmt.Fprintf(&buf, "%s\n\n", desc)
		}
	}

	// Add metadata
	buf.WriteString("## Course Information\n\n")
	fmt.Fprintf(&buf, "- **Course ID**: %s\n", strings.TrimSpace(course.Course.ID))
	fmt.Fprintf(&buf, "- **Share ID**: %s\n", strings.TrimSpace(course.ShareID))
	fmt.Fprintf(&buf, "- **Navigation Mode**: %s\n", strings.TrimSpace(course.Course.NavigationMode))
	if course.Course.ExportSettings != nil {
		fmt.Fprintf(&buf, "- **Export Format**: %s\n", strings.TrimSpace(course.Course.ExportSettings.Format))
	}
	buf.WriteString("\n---\n\n")

	// Process lessons
	lessonCounter := 0
	for _, lesson := range course.Course.Lessons {
		if lesson.Type == models.LessonTypeSection {
			fmt.Fprintf(&buf, "# %s\n\n", strings.TrimSpace(lesson.Title))
			continue
		}

		lessonCounter++
		fmt.Fprintf(&buf, "## Lesson %d: %s\n\n", lessonCounter, strings.TrimSpace(lesson.Title))

		if lesson.Description != "" {
			if desc := e.htmlCleaner.CleanAndTrim(lesson.Description); desc != "" {
				fmt.Fprintf(&buf, "%s\n\n", desc)
			}
		}

		// Process lesson items
		for _, item := range lesson.Items {
			e.processItemToMarkdown(&buf, item, 3)
		}

		buf.WriteString("---\n\n")
	}

	// Clean up multiple consecutive blank lines before writing
	output := cleanMarkdownWhitespace(buf.Bytes())

	// #nosec G306 - 0644 is appropriate for export files that should be readable by others
	if err := os.WriteFile(outputPath, output, 0o644); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}
	return nil
}

// cleanMarkdownWhitespace normalizes whitespace in markdown output.
// It removes trailing whitespace from lines and collapses multiple blank lines into one.
func cleanMarkdownWhitespace(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	result := make([]string, 0, len(lines))
	prevBlank := false

	for _, line := range lines {
		// Trim trailing whitespace from each line
		trimmed := strings.TrimRight(line, " \t")

		// Collapse multiple consecutive blank lines
		isBlank := trimmed == ""
		if isBlank && prevBlank {
			continue
		}

		result = append(result, trimmed)
		prevBlank = isBlank
	}

	// Remove trailing blank lines
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}

	// Ensure file ends with newline
	return []byte(strings.Join(result, "\n") + "\n")
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
	case models.ItemTypeText:
		e.processTextItem(buf, item, headingPrefix)
	case models.ItemTypeList:
		e.processListItem(buf, item)
	case models.ItemTypeMultimedia:
		e.processMultimediaItem(buf, item, headingPrefix)
	case models.ItemTypeImage:
		e.processImageItem(buf, item, headingPrefix)
	case models.ItemTypeKnowledgeCheck:
		e.processKnowledgeCheckItem(buf, item, headingPrefix)
	case models.ItemTypeInteractive:
		e.processInteractiveItem(buf, item, headingPrefix)
	case models.ItemTypeDivider:
		e.processDividerItem(buf)
	default:
		e.processUnknownItem(buf, item, headingPrefix)
	}
}

// processTextItem handles text content with headings and paragraphs.
func (e *MarkdownExporter) processTextItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	for _, subItem := range item.Items {
		if subItem.Heading != "" {
			if heading := e.htmlCleaner.CleanAndTrim(subItem.Heading); heading != "" {
				fmt.Fprintf(buf, "%s %s\n\n", headingPrefix, heading)
			}
		}
		if subItem.Paragraph != "" {
			if paragraph := e.htmlCleaner.CleanAndTrim(subItem.Paragraph); paragraph != "" {
				fmt.Fprintf(buf, "%s\n\n", paragraph)
			}
		}
	}
}

// processListItem handles list items with bullet points.
func (e *MarkdownExporter) processListItem(buf *bytes.Buffer, item models.Item) {
	for _, subItem := range item.Items {
		if subItem.Paragraph != "" {
			if paragraph := e.htmlCleaner.CleanAndTrim(subItem.Paragraph); paragraph != "" {
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
		if caption := e.htmlCleaner.CleanAndTrim(subItem.Caption); caption != "" {
			fmt.Fprintf(buf, "*%s*\n", caption)
		}
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
			if caption := e.htmlCleaner.CleanAndTrim(subItem.Caption); caption != "" {
				fmt.Fprintf(buf, "*%s*\n", caption)
			}
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
		if title := e.htmlCleaner.CleanAndTrim(subItem.Title); title != "" {
			fmt.Fprintf(buf, "**Question**: %s\n\n", title)
		}
	}

	e.processAnswers(buf, subItem.Answers)

	if subItem.Feedback != "" {
		if feedback := e.htmlCleaner.CleanAndTrim(subItem.Feedback); feedback != "" {
			fmt.Fprintf(buf, "\n**Feedback**: %s\n", feedback)
		}
	}
}

// processAnswers processes answer choices for quiz questions.
func (e *MarkdownExporter) processAnswers(buf *bytes.Buffer, answers []models.Answer) {
	if len(answers) == 0 {
		return
	}
	buf.WriteString("**Answers**:\n")
	for i, answer := range answers {
		correctMark := ""
		if answer.Correct {
			correctMark = " [correct]"
		}
		title := strings.TrimSpace(answer.Title)
		fmt.Fprintf(buf, "%d. %s%s\n", i+1, title, correctMark)
	}
}

// processInteractiveItem handles interactive content.
func (e *MarkdownExporter) processInteractiveItem(buf *bytes.Buffer, item models.Item, headingPrefix string) {
	fmt.Fprintf(buf, "%s Interactive Content\n\n", headingPrefix)
	for _, subItem := range item.Items {
		if subItem.Title != "" {
			if title := e.htmlCleaner.CleanAndTrim(subItem.Title); title != "" {
				fmt.Fprintf(buf, "**%s**\n\n", title)
			}
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
		if title := e.htmlCleaner.CleanAndTrim(subItem.Title); title != "" {
			fmt.Fprintf(buf, "**%s**\n\n", title)
		}
	}
	if subItem.Paragraph != "" {
		if paragraph := e.htmlCleaner.CleanAndTrim(subItem.Paragraph); paragraph != "" {
			fmt.Fprintf(buf, "%s\n\n", paragraph)
		}
	}
}
