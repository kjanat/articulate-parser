// Package exporters_test provides tests for the html exporter.
package exporters

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// TestNewHTMLExporter tests the NewHTMLExporter constructor.
func TestNewHTMLExporter(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	if exporter == nil {
		t.Fatal("NewHTMLExporter() returned nil")
	}

	// Type assertion to check internal structure
	htmlExporter, ok := exporter.(*HTMLExporter)
	if !ok {
		t.Fatal("NewHTMLExporter() returned wrong type")
	}

	if htmlExporter.htmlCleaner == nil {
		t.Error("htmlCleaner should not be nil")
	}
}

// TestHTMLExporter_GetSupportedFormat tests the GetSupportedFormat method.
func TestHTMLExporter_GetSupportedFormat(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	expected := "html"
	result := exporter.GetSupportedFormat()

	if result != expected {
		t.Errorf("Expected format '%s', got '%s'", expected, result)
	}
}

// TestHTMLExporter_Export tests the Export method.
func TestHTMLExporter_Export(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	// Create test course
	testCourse := createTestCourseForHTML()

	// Create temporary directory and file
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test-course.html")

	// Test successful export
	err := exporter.Export(testCourse, outputPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Verify HTML structure
	if !strings.Contains(contentStr, "<!DOCTYPE html>") {
		t.Error("Output should contain HTML doctype")
	}

	if !strings.Contains(contentStr, "<html lang=\"en\">") {
		t.Error("Output should contain HTML tag with lang attribute")
	}

	if !strings.Contains(contentStr, "<title>Test Course</title>") {
		t.Error("Output should contain course title in head")
	}

	// Verify main course title
	if !strings.Contains(contentStr, "<h1>Test Course</h1>") {
		t.Error("Output should contain course title as main heading")
	}

	// Verify course information section
	if !strings.Contains(contentStr, "Course Information") {
		t.Error("Output should contain course information section")
	}

	// Verify course metadata
	if !strings.Contains(contentStr, "Course ID") {
		t.Error("Output should contain course ID")
	}

	if !strings.Contains(contentStr, "Share ID") {
		t.Error("Output should contain share ID")
	}

	// Verify lesson content
	if !strings.Contains(contentStr, "Lesson 1: Test Lesson") {
		t.Error("Output should contain lesson heading")
	}

	// Verify CSS is included
	if !strings.Contains(contentStr, "<style>") {
		t.Error("Output should contain CSS styles")
	}

	if !strings.Contains(contentStr, "font-family") {
		t.Error("Output should contain CSS font-family")
	}
}

// TestHTMLExporter_Export_InvalidPath tests export with invalid output path.
func TestHTMLExporter_Export_InvalidPath(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	testCourse := createTestCourseForHTML()

	// Try to export to invalid path (non-existent directory)
	invalidPath := "/non/existent/path/test.html"
	err := exporter.Export(testCourse, invalidPath)

	if err == nil {
		t.Error("Expected error for invalid output path, but got nil")
	}
}

// TestHTMLExporter_ProcessTextItem tests the processTextItem method.
func TestHTMLExporter_ProcessTextItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "text",
		Items: []models.SubItem{
			{
				Heading:   "<h1>Test Heading</h1>",
				Paragraph: "<p>Test paragraph with <strong>bold</strong> text.</p>",
			},
			{
				Paragraph: "<p>Another paragraph.</p>",
			},
		},
	}

	exporter.processTextItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "text-item") {
		t.Error("Should contain text-item CSS class")
	}
	if !strings.Contains(result, "Text Content") {
		t.Error("Should contain text content heading")
	}
	if !strings.Contains(result, "<h1>Test Heading</h1>") {
		t.Error("Should preserve HTML heading")
	}
	if !strings.Contains(result, "<strong>bold</strong>") {
		t.Error("Should preserve HTML formatting in paragraph")
	}
}

// TestHTMLExporter_ProcessListItem tests the processListItem method.
func TestHTMLExporter_ProcessListItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "list",
		Items: []models.SubItem{
			{Paragraph: "<p>First item</p>"},
			{Paragraph: "<p>Second item with <em>emphasis</em></p>"},
			{Paragraph: "<p>Third item</p>"},
		},
	}

	exporter.processListItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "list-item") {
		t.Error("Should contain list-item CSS class")
	}
	if !strings.Contains(result, "<ul>") {
		t.Error("Should contain unordered list")
	}
	if !strings.Contains(result, "<li>First item</li>") {
		t.Error("Should contain first list item")
	}
	if !strings.Contains(result, "<li>Second item with emphasis</li>") {
		t.Error("Should contain second list item with cleaned HTML")
	}
	if !strings.Contains(result, "<li>Third item</li>") {
		t.Error("Should contain third list item")
	}
}

// TestHTMLExporter_ProcessKnowledgeCheckItem tests the processKnowledgeCheckItem method.
func TestHTMLExporter_ProcessKnowledgeCheckItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "knowledgeCheck",
		Items: []models.SubItem{
			{
				Title: "<p>What is the correct answer?</p>",
				Answers: []models.Answer{
					{Title: "Wrong answer", Correct: false},
					{Title: "Correct answer", Correct: true},
					{Title: "Another wrong answer", Correct: false},
				},
				Feedback: "<p>Great job! This is the feedback.</p>",
			},
		},
	}

	exporter.processKnowledgeCheckItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "knowledge-check") {
		t.Error("Should contain knowledge-check CSS class")
	}
	if !strings.Contains(result, "Knowledge Check") {
		t.Error("Should contain knowledge check heading")
	}
	if !strings.Contains(result, "What is the correct answer?") {
		t.Error("Should contain question text")
	}
	if !strings.Contains(result, "Wrong answer") {
		t.Error("Should contain first answer")
	}
	if !strings.Contains(result, "correct-answer") {
		t.Error("Should mark correct answer with CSS class")
	}
	if !strings.Contains(result, "Feedback") {
		t.Error("Should contain feedback section")
	}
}

// TestHTMLExporter_ProcessMultimediaItem tests the processMultimediaItem method.
func TestHTMLExporter_ProcessMultimediaItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "multimedia",
		Items: []models.SubItem{
			{
				Title: "<p>Video Title</p>",
				Media: &models.Media{
					Video: &models.VideoMedia{
						OriginalUrl: "https://example.com/video.mp4",
						Duration:    120,
					},
				},
				Caption: "<p>Video caption</p>",
			},
		},
	}

	exporter.processMultimediaItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "multimedia-item") {
		t.Error("Should contain multimedia-item CSS class")
	}
	if !strings.Contains(result, "Media Content") {
		t.Error("Should contain media content heading")
	}
	if !strings.Contains(result, "Video Title") {
		t.Error("Should contain video title")
	}
	if !strings.Contains(result, "https://example.com/video.mp4") {
		t.Error("Should contain video URL")
	}
	if !strings.Contains(result, "120 seconds") {
		t.Error("Should contain video duration")
	}
	if !strings.Contains(result, "Video caption") {
		t.Error("Should contain video caption")
	}
}

// TestHTMLExporter_ProcessImageItem tests the processImageItem method.
func TestHTMLExporter_ProcessImageItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "image",
		Items: []models.SubItem{
			{
				Media: &models.Media{
					Image: &models.ImageMedia{
						OriginalUrl: "https://example.com/image.png",
					},
				},
				Caption: "<p>Image caption</p>",
			},
		},
	}

	exporter.processImageItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "multimedia-item") {
		t.Error("Should contain multimedia-item CSS class")
	}
	if !strings.Contains(result, "Image") {
		t.Error("Should contain image heading")
	}
	if !strings.Contains(result, "https://example.com/image.png") {
		t.Error("Should contain image URL")
	}
	if !strings.Contains(result, "Image caption") {
		t.Error("Should contain image caption")
	}
}

// TestHTMLExporter_ProcessInteractiveItem tests the processInteractiveItem method.
func TestHTMLExporter_ProcessInteractiveItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "interactive",
		Items: []models.SubItem{
			{
				Title:     "<p>Interactive element title</p>",
				Paragraph: "<p>Interactive content description</p>",
			},
		},
	}

	exporter.processInteractiveItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "interactive-item") {
		t.Error("Should contain interactive-item CSS class")
	}
	if !strings.Contains(result, "Interactive Content") {
		t.Error("Should contain interactive content heading")
	}
	if !strings.Contains(result, "Interactive element title") {
		t.Error("Should contain interactive element title")
	}
	if !strings.Contains(result, "Interactive content description") {
		t.Error("Should contain interactive content description")
	}
}

// TestHTMLExporter_ProcessDividerItem tests the processDividerItem method.
func TestHTMLExporter_ProcessDividerItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	exporter.processDividerItem(&buf)

	result := buf.String()
	expected := "        <hr>\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestHTMLExporter_ProcessUnknownItem tests the processUnknownItem method.
func TestHTMLExporter_ProcessUnknownItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "unknown",
		Items: []models.SubItem{
			{
				Title:     "<p>Unknown item title</p>",
				Paragraph: "<p>Unknown item content</p>",
			},
		},
	}

	exporter.processUnknownItem(&buf, item)

	result := buf.String()

	if !strings.Contains(result, "unknown-item") {
		t.Error("Should contain unknown-item CSS class")
	}
	if !strings.Contains(result, "Unknown Content") {
		t.Error("Should contain unknown content heading")
	}
	if !strings.Contains(result, "Unknown item title") {
		t.Error("Should contain unknown item title")
	}
	if !strings.Contains(result, "Unknown item content") {
		t.Error("Should contain unknown item content")
	}
}

// TestHTMLExporter_ProcessAnswers tests the processAnswers method.
func TestHTMLExporter_ProcessAnswers(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	answers := []models.Answer{
		{Title: "Answer 1", Correct: false},
		{Title: "Answer 2", Correct: true},
		{Title: "Answer 3", Correct: false},
	}

	exporter.processAnswers(&buf, answers)

	result := buf.String()

	if !strings.Contains(result, "answers") {
		t.Error("Should contain answers CSS class")
	}
	if !strings.Contains(result, "<h5>Answers:</h5>") {
		t.Error("Should contain answers heading")
	}
	if !strings.Contains(result, "<ol>") {
		t.Error("Should contain ordered list")
	}
	if !strings.Contains(result, "<li>Answer 1</li>") {
		t.Error("Should contain first answer")
	}
	if !strings.Contains(result, "correct-answer") {
		t.Error("Should mark correct answer with CSS class")
	}
	if !strings.Contains(result, "<li class=\"correct-answer\">Answer 2</li>") {
		t.Error("Should mark correct answer properly")
	}
	if !strings.Contains(result, "<li>Answer 3</li>") {
		t.Error("Should contain third answer")
	}
}

// TestHTMLExporter_ProcessItemToHTML_AllTypes tests all item types.
func TestHTMLExporter_ProcessItemToHTML_AllTypes(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	tests := []struct {
		name         string
		itemType     string
		expectedText string
	}{
		{
			name:         "text item",
			itemType:     "text",
			expectedText: "Text Content",
		},
		{
			name:         "list item",
			itemType:     "list",
			expectedText: "List",
		},
		{
			name:         "knowledge check item",
			itemType:     "knowledgeCheck",
			expectedText: "Knowledge Check",
		},
		{
			name:         "multimedia item",
			itemType:     "multimedia",
			expectedText: "Media Content",
		},
		{
			name:         "image item",
			itemType:     "image",
			expectedText: "Image",
		},
		{
			name:         "interactive item",
			itemType:     "interactive",
			expectedText: "Interactive Content",
		},
		{
			name:         "divider item",
			itemType:     "divider",
			expectedText: "<hr>",
		},
		{
			name:         "unknown item",
			itemType:     "unknown",
			expectedText: "Unknown Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			item := models.Item{
				Type: tt.itemType,
				Items: []models.SubItem{
					{Title: "Test title", Paragraph: "Test content"},
				},
			}

			// Handle empty unknown items
			if tt.itemType == "unknown" && tt.expectedText == "" {
				item.Items = []models.SubItem{}
			}

			exporter.processItemToHTML(&buf, item)

			result := buf.String()
			if tt.expectedText != "" && !strings.Contains(result, tt.expectedText) {
				t.Errorf("Expected content to contain: %q\nGot: %q", tt.expectedText, result)
			}
		})
	}
}

// TestHTMLExporter_ComplexCourse tests export of a complex course structure.
func TestHTMLExporter_ComplexCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	// Create complex test course
	course := &models.Course{
		ShareID: "complex-test-id",
		Author:  "Test Author",
		Course: models.CourseInfo{
			ID:             "complex-course",
			Title:          "Complex Test Course",
			Description:    "<p>This is a <strong>complex</strong> course description.</p>",
			NavigationMode: "menu",
			ExportSettings: &models.ExportSettings{
				Format: "scorm",
			},
			Lessons: []models.Lesson{
				{
					ID:    "section-1",
					Title: "Course Section",
					Type:  "section",
				},
				{
					ID:          "lesson-1",
					Title:       "Introduction Lesson",
					Type:        "lesson",
					Description: "<p>Introduction to the course</p>",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Heading:   "<h2>Welcome</h2>",
									Paragraph: "<p>Welcome to our course!</p>",
								},
							},
						},
						{
							Type: "list",
							Items: []models.SubItem{
								{Paragraph: "<p>First objective</p>"},
								{Paragraph: "<p>Second objective</p>"},
							},
						},
						{
							Type: "knowledgeCheck",
							Items: []models.SubItem{
								{
									Title: "<p>What will you learn?</p>",
									Answers: []models.Answer{
										{Title: "Nothing", Correct: false},
										{Title: "Everything", Correct: true},
									},
									Feedback: "<p>Great choice!</p>",
								},
							},
						},
					},
				},
			},
		},
	}

	// Create temporary output file
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "complex-course.html")

	// Export course
	err := exporter.Export(course, outputPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Verify various elements are present
	checks := []string{
		"<title>Complex Test Course</title>",
		"<h1>Complex Test Course</h1>",
		"This is a <strong>complex</strong> course description.",
		"Course Information",
		"complex-course",
		"complex-test-id",
		"menu",
		"scorm",
		"Course Section",
		"Lesson 1: Introduction Lesson",
		"Introduction to the course",
		"<h2>Welcome</h2>",
		"Welcome to our course!",
		"First objective",
		"Second objective",
		"Knowledge Check",
		"What will you learn?",
		"Nothing",
		"Everything",
		"correct-answer",
		"Great choice!",
	}

	for _, check := range checks {
		if !strings.Contains(contentStr, check) {
			t.Errorf("Output should contain: %q", check)
		}
	}

	// Verify HTML structure
	structureChecks := []string{
		"<!DOCTYPE html>",
		"<html lang=\"en\">",
		"<head>",
		"<body>",
		"</html>",
		"<style>",
		"font-family",
	}

	for _, check := range structureChecks {
		if !strings.Contains(contentStr, check) {
			t.Errorf("Output should contain HTML structure element: %q", check)
		}
	}
}

// TestHTMLExporter_EmptyCourse tests export of an empty course.
func TestHTMLExporter_EmptyCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	// Create minimal course
	course := &models.Course{
		ShareID: "empty-id",
		Course: models.CourseInfo{
			ID:      "empty-course",
			Title:   "Empty Course",
			Lessons: []models.Lesson{},
		},
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "empty-course.html")

	err := exporter.Export(course, outputPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read and verify basic structure
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Verify basic HTML structure even for empty course
	if !strings.Contains(contentStr, "<!DOCTYPE html>") {
		t.Error("Output should contain HTML doctype")
	}
	if !strings.Contains(contentStr, "<title>Empty Course</title>") {
		t.Error("Output should contain course title")
	}
	if !strings.Contains(contentStr, "<h1>Empty Course</h1>") {
		t.Error("Output should contain course heading")
	}
}

// TestHTMLExporter_HTMLCleaning tests that HTML content is properly handled.
func TestHTMLExporter_HTMLCleaning(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	// Create course with HTML content that needs cleaning in some places
	course := &models.Course{
		ShareID: "html-test-id",
		Course: models.CourseInfo{
			ID:          "html-test-course",
			Title:       "HTML Test Course",
			Description: "<p>Description with <script>alert('xss')</script> and <b>bold</b> text.</p>",
			Lessons: []models.Lesson{
				{
					ID:          "lesson-1",
					Title:       "Test Lesson",
					Type:        "lesson",
					Description: "<div>Lesson description with <span style='color:red'>styled</span> content.</div>",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Heading:   "<h1>Heading with <em>emphasis</em> and &amp; entities</h1>",
									Paragraph: "<p>Paragraph with &lt;code&gt; entities and <strong>formatting</strong>.</p>",
								},
							},
						},
					},
				},
			},
		},
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "html-test.html")

	err := exporter.Export(course, outputPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created (basic check that HTML handling didn't break export)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read content and verify some HTML is preserved (descriptions, headings, paragraphs)
	// while list items are cleaned for safety
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// HTML should be preserved in some places
	if !strings.Contains(contentStr, "<b>bold</b>") {
		t.Error("Should preserve HTML formatting in descriptions")
	}
	if !strings.Contains(contentStr, "<h1>Heading with <em>emphasis</em>") {
		t.Error("Should preserve HTML in headings")
	}
	if !strings.Contains(contentStr, "<strong>formatting</strong>") {
		t.Error("Should preserve HTML in paragraphs")
	}
}

// createTestCourseForHTML creates a test course for HTML export testing.
func createTestCourseForHTML() *models.Course {
	return &models.Course{
		ShareID: "test-share-id",
		Course: models.CourseInfo{
			ID:             "test-course-id",
			Title:          "Test Course",
			Description:    "<p>Test course description with <strong>formatting</strong>.</p>",
			NavigationMode: "free",
			Lessons: []models.Lesson{
				{
					ID:    "section-1",
					Title: "Test Section",
					Type:  "section",
				},
				{
					ID:          "lesson-1",
					Title:       "Test Lesson",
					Type:        "lesson",
					Description: "<p>Test lesson description</p>",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Heading:   "<h2>Test Heading</h2>",
									Paragraph: "<p>Test paragraph content.</p>",
								},
							},
						},
						{
							Type: "list",
							Items: []models.SubItem{
								{Paragraph: "<p>First list item</p>"},
								{Paragraph: "<p>Second list item</p>"},
							},
						},
					},
				},
			},
		},
	}
}

// BenchmarkHTMLExporter_Export benchmarks the Export method.
func BenchmarkHTMLExporter_Export(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)
	course := createTestCourseForHTML()

	// Create temporary directory
	tempDir := b.TempDir()

	for b.Loop() {
		outputPath := filepath.Join(tempDir, "benchmark-course.html")
		_ = exporter.Export(course, outputPath)
		// Clean up for next iteration
		os.Remove(outputPath)
	}
}

// BenchmarkHTMLExporter_ProcessTextItem benchmarks text item processing.
func BenchmarkHTMLExporter_ProcessTextItem(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &HTMLExporter{htmlCleaner: htmlCleaner}

	item := models.Item{
		Type: "text",
		Items: []models.SubItem{
			{
				Heading:   "<h1>Benchmark Heading</h1>",
				Paragraph: "<p>Benchmark paragraph with <strong>formatting</strong>.</p>",
			},
		},
	}

	for b.Loop() {
		var buf bytes.Buffer
		exporter.processTextItem(&buf, item)
	}
}

// BenchmarkHTMLExporter_ComplexCourse benchmarks export of a complex course.
func BenchmarkHTMLExporter_ComplexCourse(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	// Create complex course for benchmarking
	course := &models.Course{
		ShareID: "benchmark-id",
		Course: models.CourseInfo{
			ID:          "benchmark-course",
			Title:       "Benchmark Course",
			Description: "<p>Complex course for performance testing</p>",
			Lessons:     make([]models.Lesson, 10), // 10 lessons
		},
	}

	// Fill with test data
	for i := range 10 {
		lesson := models.Lesson{
			ID:    "lesson-" + string(rune(i)),
			Title: "Lesson " + string(rune(i)),
			Type:  "lesson",
			Items: make([]models.Item, 5), // 5 items per lesson
		}

		for j := range 5 {
			item := models.Item{
				Type:  "text",
				Items: make([]models.SubItem, 3), // 3 sub-items per item
			}

			for k := range 3 {
				item.Items[k] = models.SubItem{
					Heading:   "<h3>Heading " + string(rune(k)) + "</h3>",
					Paragraph: "<p>Paragraph content with <strong>formatting</strong> for performance testing.</p>",
				}
			}

			lesson.Items[j] = item
		}

		course.Course.Lessons[i] = lesson
	}

	tempDir := b.TempDir()

	for b.Loop() {
		outputPath := filepath.Join(tempDir, "benchmark-complex.html")
		_ = exporter.Export(course, outputPath)
		os.Remove(outputPath)
	}
}
