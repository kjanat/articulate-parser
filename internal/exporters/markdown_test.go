// Package exporters_test provides tests for the markdown exporter.
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

// TestNewMarkdownExporter tests the NewMarkdownExporter constructor.
func TestNewMarkdownExporter(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewMarkdownExporter(htmlCleaner)

	if exporter == nil {
		t.Fatal("NewMarkdownExporter() returned nil")
	}

	// Type assertion to check internal structure
	markdownExporter, ok := exporter.(*MarkdownExporter)
	if !ok {
		t.Fatal("NewMarkdownExporter() returned wrong type")
	}

	if markdownExporter.htmlCleaner == nil {
		t.Error("htmlCleaner should not be nil")
	}
}

// TestMarkdownExporter_GetSupportedFormat tests the GetSupportedFormat method.
func TestMarkdownExporter_GetSupportedFormat(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewMarkdownExporter(htmlCleaner)

	expected := "markdown"
	result := exporter.GetSupportedFormat()

	if result != expected {
		t.Errorf("Expected format '%s', got '%s'", expected, result)
	}
}

// TestMarkdownExporter_Export tests the Export method.
func TestMarkdownExporter_Export(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewMarkdownExporter(htmlCleaner)

	// Create test course
	testCourse := createTestCourseForMarkdown()

	// Create temporary directory and file
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test-course.md")

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

	// Verify main course title
	if !strings.Contains(contentStr, "# Test Course") {
		t.Error("Output should contain course title as main heading")
	}

	// Verify course information section
	if !strings.Contains(contentStr, "## Course Information") {
		t.Error("Output should contain course information section")
	}

	// Verify course metadata
	if !strings.Contains(contentStr, "- **Course ID**: test-course-id") {
		t.Error("Output should contain course ID")
	}

	if !strings.Contains(contentStr, "- **Share ID**: test-share-id") {
		t.Error("Output should contain share ID")
	}

	// Verify lesson content
	if !strings.Contains(contentStr, "## Lesson 1: Test Lesson") {
		t.Error("Output should contain lesson heading")
	}

	// Verify section handling
	if !strings.Contains(contentStr, "# Test Section") {
		t.Error("Output should contain section as main heading")
	}
}

// TestMarkdownExporter_Export_InvalidPath tests export with invalid output path.
func TestMarkdownExporter_Export_InvalidPath(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewMarkdownExporter(htmlCleaner)

	testCourse := createTestCourseForMarkdown()

	// Try to write to invalid path
	invalidPath := "/invalid/path/that/does/not/exist/file.md"
	err := exporter.Export(testCourse, invalidPath)

	if err == nil {
		t.Fatal("Expected error for invalid path, got nil")
	}
}

// TestMarkdownExporter_ProcessTextItem tests the processTextItem method.
func TestMarkdownExporter_ProcessTextItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

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

	exporter.processTextItem(&buf, item, "###")

	result := buf.String()
	expected := "### Test Heading\n\nTest paragraph with bold text.\n\nAnother paragraph.\n\n"

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

// TestMarkdownExporter_ProcessListItem tests the processListItem method.
func TestMarkdownExporter_ProcessListItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

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
	expected := "- First item\n- Second item with emphasis\n- Third item\n\n"

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

// TestMarkdownExporter_ProcessMultimediaItem tests the processMultimediaItem method.
func TestMarkdownExporter_ProcessMultimediaItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "multimedia",
		Items: []models.SubItem{
			{
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

	exporter.processMultimediaItem(&buf, item, "###")

	result := buf.String()

	if !strings.Contains(result, "### Media Content") {
		t.Error("Should contain media content heading")
	}
	if !strings.Contains(result, "**Video**: https://example.com/video.mp4") {
		t.Error("Should contain video URL")
	}
	if !strings.Contains(result, "**Duration**: 120 seconds") {
		t.Error("Should contain video duration")
	}
	if !strings.Contains(result, "*Video caption*") {
		t.Error("Should contain video caption")
	}
}

// TestMarkdownExporter_ProcessImageItem tests the processImageItem method.
func TestMarkdownExporter_ProcessImageItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "image",
		Items: []models.SubItem{
			{
				Media: &models.Media{
					Image: &models.ImageMedia{
						OriginalUrl: "https://example.com/image.jpg",
					},
				},
				Caption: "<p>Image caption</p>",
			},
		},
	}

	exporter.processImageItem(&buf, item, "###")

	result := buf.String()

	if !strings.Contains(result, "### Image") {
		t.Error("Should contain image heading")
	}
	if !strings.Contains(result, "**Image**: https://example.com/image.jpg") {
		t.Error("Should contain image URL")
	}
	if !strings.Contains(result, "*Image caption*") {
		t.Error("Should contain image caption")
	}
}

// TestMarkdownExporter_ProcessKnowledgeCheckItem tests the processKnowledgeCheckItem method.
func TestMarkdownExporter_ProcessKnowledgeCheckItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "knowledgeCheck",
		Items: []models.SubItem{
			{
				Title: "<p>What is the capital of France?</p>",
				Answers: []models.Answer{
					{Title: "London", Correct: false},
					{Title: "Paris", Correct: true},
					{Title: "Berlin", Correct: false},
				},
				Feedback: "<p>Paris is the capital of France.</p>",
			},
		},
	}

	exporter.processKnowledgeCheckItem(&buf, item, "###")

	result := buf.String()

	if !strings.Contains(result, "### Knowledge Check") {
		t.Error("Should contain knowledge check heading")
	}
	if !strings.Contains(result, "**Question**: What is the capital of France?") {
		t.Error("Should contain question")
	}
	if !strings.Contains(result, "**Answers**:") {
		t.Error("Should contain answers heading")
	}
	if !strings.Contains(result, "2. Paris ✓") {
		t.Error("Should mark correct answer")
	}
	if !strings.Contains(result, "**Feedback**: Paris is the capital of France.") {
		t.Error("Should contain feedback")
	}
}

// TestMarkdownExporter_ProcessInteractiveItem tests the processInteractiveItem method.
func TestMarkdownExporter_ProcessInteractiveItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	item := models.Item{
		Type: "interactive",
		Items: []models.SubItem{
			{Title: "<p>Interactive element title</p>"},
		},
	}

	exporter.processInteractiveItem(&buf, item, "###")

	result := buf.String()

	if !strings.Contains(result, "### Interactive Content") {
		t.Error("Should contain interactive content heading")
	}
	if !strings.Contains(result, "**Interactive element title**") {
		t.Error("Should contain interactive element title")
	}
}

// TestMarkdownExporter_ProcessDividerItem tests the processDividerItem method.
func TestMarkdownExporter_ProcessDividerItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	exporter.processDividerItem(&buf)

	result := buf.String()
	expected := "---\n\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestMarkdownExporter_ProcessUnknownItem tests the processUnknownItem method.
func TestMarkdownExporter_ProcessUnknownItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

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

	exporter.processUnknownItem(&buf, item, "###")

	result := buf.String()

	if !strings.Contains(result, "### Unknown Content") {
		t.Error("Should contain unknown content heading")
	}
	if !strings.Contains(result, "**Unknown item title**") {
		t.Error("Should contain unknown item title")
	}
	if !strings.Contains(result, "Unknown item content") {
		t.Error("Should contain unknown item content")
	}
}

// TestMarkdownExporter_ProcessVideoMedia tests the processVideoMedia method.
func TestMarkdownExporter_ProcessVideoMedia(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	media := &models.Media{
		Video: &models.VideoMedia{
			OriginalUrl: "https://example.com/video.mp4",
			Duration:    300,
		},
	}

	exporter.processVideoMedia(&buf, media)

	result := buf.String()

	if !strings.Contains(result, "**Video**: https://example.com/video.mp4") {
		t.Error("Should contain video URL")
	}
	if !strings.Contains(result, "**Duration**: 300 seconds") {
		t.Error("Should contain video duration")
	}
}

// TestMarkdownExporter_ProcessImageMedia tests the processImageMedia method.
func TestMarkdownExporter_ProcessImageMedia(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	media := &models.Media{
		Image: &models.ImageMedia{
			OriginalUrl: "https://example.com/image.jpg",
		},
	}

	exporter.processImageMedia(&buf, media)

	result := buf.String()
	expected := "**Image**: https://example.com/image.jpg\n"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestMarkdownExporter_ProcessAnswers tests the processAnswers method.
func TestMarkdownExporter_ProcessAnswers(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	var buf bytes.Buffer
	answers := []models.Answer{
		{Title: "Answer 1", Correct: false},
		{Title: "Answer 2", Correct: true},
		{Title: "Answer 3", Correct: false},
	}

	exporter.processAnswers(&buf, answers)

	result := buf.String()

	if !strings.Contains(result, "**Answers**:") {
		t.Error("Should contain answers heading")
	}
	if !strings.Contains(result, "1. Answer 1") {
		t.Error("Should contain first answer")
	}
	if !strings.Contains(result, "2. Answer 2 ✓") {
		t.Error("Should mark correct answer")
	}
	if !strings.Contains(result, "3. Answer 3") {
		t.Error("Should contain third answer")
	}
}

// TestMarkdownExporter_ProcessItemToMarkdown_AllTypes tests all item types.
func TestMarkdownExporter_ProcessItemToMarkdown_AllTypes(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	tests := []struct {
		name         string
		itemType     string
		expectedText string
	}{
		{
			name:         "text item",
			itemType:     "text",
			expectedText: "", // processTextItem handles empty items
		},
		{
			name:         "list item",
			itemType:     "list",
			expectedText: "\n", // Empty list adds newline
		},
		{
			name:         "multimedia item",
			itemType:     "multimedia",
			expectedText: "### Media Content",
		},
		{
			name:         "image item",
			itemType:     "image",
			expectedText: "### Image",
		},
		{
			name:         "knowledgeCheck item",
			itemType:     "knowledgeCheck",
			expectedText: "### Knowledge Check",
		},
		{
			name:         "interactive item",
			itemType:     "interactive",
			expectedText: "### Interactive Content",
		},
		{
			name:         "divider item",
			itemType:     "divider",
			expectedText: "---",
		},
		{
			name:         "unknown item",
			itemType:     "unknown",
			expectedText: "", // Empty unknown items don't add content
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			item := models.Item{Type: tt.itemType}

			exporter.processItemToMarkdown(&buf, item, 3)

			result := buf.String()
			if tt.expectedText != "" && !strings.Contains(result, tt.expectedText) {
				t.Errorf("Expected result to contain %q, got %q", tt.expectedText, result)
			}
		})
	}
}

// TestMarkdownExporter_ComplexCourse tests export of a complex course structure.
func TestMarkdownExporter_ComplexCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewMarkdownExporter(htmlCleaner)

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
	outputPath := filepath.Join(tempDir, "complex-course.md")

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
		"# Complex Test Course",
		"This is a complex course description.",
		"- **Export Format**: scorm",
		"# Course Section",
		"## Lesson 1: Introduction Lesson",
		"Introduction to the course",
		"### Welcome",
		"Welcome to our course!",
		"- First objective",
		"- Second objective",
		"### Knowledge Check",
		"**Question**: What will you learn?",
		"2. Everything ✓",
		"**Feedback**: Great choice!",
	}

	for _, check := range checks {
		if !strings.Contains(contentStr, check) {
			t.Errorf("Output should contain: %q", check)
		}
	}
}

// createTestCourseForMarkdown creates a test course for markdown export testing.
func createTestCourseForMarkdown() *models.Course {
	return &models.Course{
		ShareID: "test-share-id",
		Author:  "Test Author",
		Course: models.CourseInfo{
			ID:             "test-course-id",
			Title:          "Test Course",
			Description:    "Test course description",
			NavigationMode: "menu",
			Lessons: []models.Lesson{
				{
					ID:    "section-1",
					Title: "Test Section",
					Type:  "section",
				},
				{
					ID:    "lesson-1",
					Title: "Test Lesson",
					Type:  "lesson",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Heading:   "Test Heading",
									Paragraph: "Test paragraph content",
								},
							},
						},
					},
				},
			},
		},
	}
}

// BenchmarkMarkdownExporter_Export benchmarks the Export method.
func BenchmarkMarkdownExporter_Export(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewMarkdownExporter(htmlCleaner)
	course := createTestCourseForMarkdown()

	// Create temporary directory
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tempDir, "benchmark-course.md")
		_ = exporter.Export(course, outputPath)
		// Clean up for next iteration
		os.Remove(outputPath)
	}
}

// BenchmarkMarkdownExporter_ProcessTextItem benchmarks the processTextItem method.
func BenchmarkMarkdownExporter_ProcessTextItem(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := &MarkdownExporter{htmlCleaner: htmlCleaner}

	item := models.Item{
		Type: "text",
		Items: []models.SubItem{
			{
				Heading:   "<h1>Benchmark Heading</h1>",
				Paragraph: "<p>Benchmark paragraph with <strong>bold</strong> text.</p>",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		exporter.processTextItem(&buf, item, "###")
	}
}
