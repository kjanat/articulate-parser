package exporters

import (
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

	if htmlExporter.tmpl == nil {
		t.Error("template should not be nil")
	}
}

// TestHTMLExporter_SupportedFormat tests the SupportedFormat method.
func TestHTMLExporter_SupportedFormat(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner)

	expected := "html"
	result := exporter.SupportedFormat()

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
		t.Logf("Generated HTML (first 500 chars):\n%s", contentStr[:min(500, len(contentStr))])
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

// TestHTMLExporter_ComplexCourse tests export of a course with complex content.
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
									Heading:   "<h2>HTML Heading</h2>",
									Paragraph: "<p>Content with <em>emphasis</em> and <strong>strong</strong> text.</p>",
								},
							},
						},
						{
							Type: "list",
							Items: []models.SubItem{
								{Paragraph: "<p>List item with <b>bold</b> text</p>"},
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

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// HTML content in descriptions should be preserved
	if !strings.Contains(contentStr, "<b>bold</b>") {
		t.Error("Should preserve HTML formatting in descriptions")
	}

	// HTML content in headings should be preserved
	if !strings.Contains(contentStr, "<h2>HTML Heading</h2>") {
		t.Error("Should preserve HTML in headings")
	}

	// List items should have HTML tags stripped (cleaned)
	if !strings.Contains(contentStr, "List item with bold text") {
		t.Error("Should clean HTML from list items")
	}
}

// createTestCourseForHTML creates a test course for HTML export tests.
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

	tempDir := b.TempDir()

	for i := range b.N {
		outputPath := filepath.Join(tempDir, "bench-course-"+string(rune(i))+".html")
		if err := exporter.Export(course, outputPath); err != nil {
			b.Fatalf("Export failed: %v", err)
		}
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
			ID:          "lesson-" + string(rune(i)),
			Title:       "Benchmark Lesson " + string(rune(i)),
			Type:        "lesson",
			Description: "<p>Lesson description</p>",
			Items: []models.Item{
				{
					Type: "text",
					Items: []models.SubItem{
						{
							Heading:   "<h2>Heading</h2>",
							Paragraph: "<p>Paragraph with content.</p>",
						},
					},
				},
				{
					Type: "list",
					Items: []models.SubItem{
						{Paragraph: "<p>Item 1</p>"},
						{Paragraph: "<p>Item 2</p>"},
					},
				},
			},
		}
		course.Course.Lessons[i] = lesson
	}

	tempDir := b.TempDir()

	for i := range b.N {
		outputPath := filepath.Join(tempDir, "bench-complex-"+string(rune(i))+".html")
		if err := exporter.Export(course, outputPath); err != nil {
			b.Fatalf("Export failed: %v", err)
		}
	}
}
