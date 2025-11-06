// Package exporters provides implementations of the Exporter interface
// for converting Articulate Rise courses into various file formats.
package exporters

import (
	"fmt"
	"os"
	"strings"

	"github.com/fumiama/go-docx"
	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DocxExporter implements the Exporter interface for DOCX format.
// It converts Articulate Rise course data into a Microsoft Word document
// using the go-docx package.
type DocxExporter struct {
	// htmlCleaner is used to convert HTML content to plain text
	htmlCleaner *services.HTMLCleaner
}

// NewDocxExporter creates a new DocxExporter instance.
// It takes an HTMLCleaner to handle HTML content conversion.
//
// Parameters:
//   - htmlCleaner: Service for cleaning HTML content in course data
//
// Returns:
//   - An implementation of the Exporter interface for DOCX format
func NewDocxExporter(htmlCleaner *services.HTMLCleaner) interfaces.Exporter {
	return &DocxExporter{
		htmlCleaner: htmlCleaner,
	}
}

// Export exports the course to a DOCX file.
// It creates a Word document with formatted content based on the course data
// and saves it to the specified output path.
//
// Parameters:
//   - course: The course data model to export
//   - outputPath: The file path where the DOCX content will be written
//
// Returns:
//   - An error if creating or saving the document fails
func (e *DocxExporter) Export(course *models.Course, outputPath string) error {
	doc := docx.New()

	// Add title
	titlePara := doc.AddParagraph()
	titlePara.AddText(course.Course.Title).Size("32").Bold()

	// Add description if available
	if course.Course.Description != "" {
		descPara := doc.AddParagraph()
		cleanDesc := e.htmlCleaner.CleanHTML(course.Course.Description)
		descPara.AddText(cleanDesc)
	}

	// Add each lesson
	for _, lesson := range course.Course.Lessons {
		e.exportLesson(doc, &lesson)
	}

	// Ensure output directory exists and add .docx extension
	if !strings.HasSuffix(strings.ToLower(outputPath), ".docx") {
		outputPath = outputPath + ".docx"
	}

	// Create the file
	// #nosec G304 - Output path is provided by user via CLI argument, which is expected behavior
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	// Ensure file is closed even if WriteTo fails. Close errors are logged but not
	// fatal since the document content has already been written to disk. A close
	// error typically indicates a filesystem synchronization issue that doesn't
	// affect the validity of the exported file.
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close output file: %v\n", err)
		}
	}()

	// Save the document
	_, err = doc.WriteTo(file)
	if err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	return nil
}

// exportLesson adds a lesson to the document with appropriate formatting.
// It creates a lesson heading, adds the description, and processes all items in the lesson.
//
// Parameters:
//   - doc: The Word document being created
//   - lesson: The lesson data model to export
func (e *DocxExporter) exportLesson(doc *docx.Docx, lesson *models.Lesson) {
	// Add lesson title
	lessonPara := doc.AddParagraph()
	lessonPara.AddText(fmt.Sprintf("Lesson: %s", lesson.Title)).Size("28").Bold()

	// Add lesson description if available
	if lesson.Description != "" {
		descPara := doc.AddParagraph()
		cleanDesc := e.htmlCleaner.CleanHTML(lesson.Description)
		descPara.AddText(cleanDesc)
	}

	// Add each item in the lesson
	for _, item := range lesson.Items {
		e.exportItem(doc, &item)
	}
}

// exportItem adds an item to the document.
// It creates an item heading and processes all sub-items within the item.
//
// Parameters:
//   - doc: The Word document being created
//   - item: The item data model to export
func (e *DocxExporter) exportItem(doc *docx.Docx, item *models.Item) {
	// Add item type as heading
	if item.Type != "" {
		itemPara := doc.AddParagraph()
		caser := cases.Title(language.English)
		itemPara.AddText(caser.String(item.Type)).Size("24").Bold()
	}

	// Add sub-items
	for _, subItem := range item.Items {
		e.exportSubItem(doc, &subItem)
	}
}

// exportSubItem adds a sub-item to the document.
// It handles different components of a sub-item like title, heading,
// paragraph content, answers, and feedback.
//
// Parameters:
//   - doc: The Word document being created
//   - subItem: The sub-item data model to export
func (e *DocxExporter) exportSubItem(doc *docx.Docx, subItem *models.SubItem) {
	// Add title if available
	if subItem.Title != "" {
		subItemPara := doc.AddParagraph()
		subItemPara.AddText("  " + subItem.Title).Bold() // Indented
	}

	// Add heading if available
	if subItem.Heading != "" {
		headingPara := doc.AddParagraph()
		cleanHeading := e.htmlCleaner.CleanHTML(subItem.Heading)
		headingPara.AddText("  " + cleanHeading).Bold() // Indented
	}

	// Add paragraph content if available
	if subItem.Paragraph != "" {
		contentPara := doc.AddParagraph()
		cleanContent := e.htmlCleaner.CleanHTML(subItem.Paragraph)
		contentPara.AddText("  " + cleanContent) // Indented
	}

	// Add answers if this is a question
	if len(subItem.Answers) > 0 {
		answersPara := doc.AddParagraph()
		answersPara.AddText("  Answers:").Bold()

		for i, answer := range subItem.Answers {
			answerPara := doc.AddParagraph()
			prefix := fmt.Sprintf("    %d. ", i+1)
			if answer.Correct {
				prefix += "✓ "
			}
			cleanAnswer := e.htmlCleaner.CleanHTML(answer.Title)
			answerPara.AddText(prefix + cleanAnswer)
		}
	}

	// Add feedback if available
	if subItem.Feedback != "" {
		feedbackPara := doc.AddParagraph()
		cleanFeedback := e.htmlCleaner.CleanHTML(subItem.Feedback)
		feedbackPara.AddText("  Feedback: " + cleanFeedback).Italic()
	}
}

// SupportedFormat returns the format name this exporter supports.
//
// Returns:
//   - A string representing the supported format ("docx")
func (e *DocxExporter) SupportedFormat() string {
	return "docx"
}
