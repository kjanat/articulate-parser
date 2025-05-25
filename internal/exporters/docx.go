// Package exporters provides implementations of the Exporter interface
// for converting Articulate Rise courses into various file formats.
package exporters

import (
	"fmt"
	"strings"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
	"github.com/unidoc/unioffice/document"
)

// DocxExporter implements the Exporter interface for DOCX format.
// It converts Articulate Rise course data into a Microsoft Word document
// using the unioffice/document package.
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
	doc := document.New()

	// Add title
	titlePara := doc.AddParagraph()
	titleRun := titlePara.AddRun()
	titleRun.AddText(course.Course.Title)
	titleRun.Properties().SetBold(true)
	titleRun.Properties().SetSize(16)

	// Add description if available
	if course.Course.Description != "" {
		descPara := doc.AddParagraph()
		descRun := descPara.AddRun()
		cleanDesc := e.htmlCleaner.CleanHTML(course.Course.Description)
		descRun.AddText(cleanDesc)
	}

	// Add each lesson
	for _, lesson := range course.Course.Lessons {
		e.exportLesson(doc, &lesson)
	}

	// Ensure output directory exists and add .docx extension
	if !strings.HasSuffix(strings.ToLower(outputPath), ".docx") {
		outputPath = outputPath + ".docx"
	}

	return doc.SaveToFile(outputPath)
}

// exportLesson adds a lesson to the document with appropriate formatting.
// It creates a lesson heading, adds the description, and processes all items in the lesson.
//
// Parameters:
//   - doc: The Word document being created
//   - lesson: The lesson data model to export
func (e *DocxExporter) exportLesson(doc *document.Document, lesson *models.Lesson) {
	// Add lesson title
	lessonPara := doc.AddParagraph()
	lessonRun := lessonPara.AddRun()
	lessonRun.AddText(fmt.Sprintf("Lesson: %s", lesson.Title))
	lessonRun.Properties().SetBold(true)
	lessonRun.Properties().SetSize(14)

	// Add lesson description if available
	if lesson.Description != "" {
		descPara := doc.AddParagraph()
		descRun := descPara.AddRun()
		cleanDesc := e.htmlCleaner.CleanHTML(lesson.Description)
		descRun.AddText(cleanDesc)
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
func (e *DocxExporter) exportItem(doc *document.Document, item *models.Item) {
	// Add item type as heading
	if item.Type != "" {
		itemPara := doc.AddParagraph()
		itemRun := itemPara.AddRun()
		itemRun.AddText(strings.Title(item.Type))
		itemRun.Properties().SetBold(true)
		itemRun.Properties().SetSize(12)
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
func (e *DocxExporter) exportSubItem(doc *document.Document, subItem *models.SubItem) {
	// Add title if available
	if subItem.Title != "" {
		subItemPara := doc.AddParagraph()
		subItemRun := subItemPara.AddRun()
		subItemRun.AddText("  " + subItem.Title) // Indented
		subItemRun.Properties().SetBold(true)
	}

	// Add heading if available
	if subItem.Heading != "" {
		headingPara := doc.AddParagraph()
		headingRun := headingPara.AddRun()
		cleanHeading := e.htmlCleaner.CleanHTML(subItem.Heading)
		headingRun.AddText("  " + cleanHeading) // Indented
		headingRun.Properties().SetBold(true)
	}

	// Add paragraph content if available
	if subItem.Paragraph != "" {
		contentPara := doc.AddParagraph()
		contentRun := contentPara.AddRun()
		cleanContent := e.htmlCleaner.CleanHTML(subItem.Paragraph)
		contentRun.AddText("  " + cleanContent) // Indented
	}

	// Add answers if this is a question
	if len(subItem.Answers) > 0 {
		answersPara := doc.AddParagraph()
		answersRun := answersPara.AddRun()
		answersRun.AddText("  Answers:")
		answersRun.Properties().SetBold(true)

		for i, answer := range subItem.Answers {
			answerPara := doc.AddParagraph()
			answerRun := answerPara.AddRun()
			prefix := fmt.Sprintf("    %d. ", i+1)
			if answer.Correct {
				prefix += "✓ "
			}
			cleanAnswer := e.htmlCleaner.CleanHTML(answer.Title)
			answerRun.AddText(prefix + cleanAnswer)
		}
	}

	// Add feedback if available
	if subItem.Feedback != "" {
		feedbackPara := doc.AddParagraph()
		feedbackRun := feedbackPara.AddRun()
		cleanFeedback := e.htmlCleaner.CleanHTML(subItem.Feedback)
		feedbackRun.AddText("  Feedback: " + cleanFeedback)
		feedbackRun.Properties().SetItalic(true)
	}
}

// GetSupportedFormat returns the format name this exporter supports.
//
// Returns:
//   - A string representing the supported format ("docx")
func (e *DocxExporter) GetSupportedFormat() string {
	return "docx"
}
