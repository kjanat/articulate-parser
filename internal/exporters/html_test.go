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
	exporter := NewHTMLExporter(htmlCleaner, nil)

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
	exporter := NewHTMLExporter(htmlCleaner, nil)

	expected := "html"
	result := exporter.SupportedFormat()

	if result != expected {
		t.Errorf("Expected format '%s', got '%s'", expected, result)
	}
}

// TestHTMLExporter_Export tests the Export method.
func TestHTMLExporter_Export(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner, nil)

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
	exporter := NewHTMLExporter(htmlCleaner, nil)

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
	exporter := NewHTMLExporter(htmlCleaner, nil)

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
	exporter := NewHTMLExporter(htmlCleaner, nil)

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
	exporter := NewHTMLExporter(htmlCleaner, nil)

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

// TestHTMLExporter_Export_ErrorPaths tests various error conditions for the HTML exporter.
func TestHTMLExporter_Export_ErrorPaths(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner, nil)

	tests := []struct {
		name        string
		course      *models.Course
		outputPath  string
		wantErr     bool
		errContains string
	}{
		{
			name:        "non-existent directory",
			course:      createTestCourseForHTML(),
			outputPath:  "/nonexistent/deeply/nested/path/output.html",
			wantErr:     true,
			errContains: "failed to create file",
		},
		{
			name:        "empty output path",
			course:      createTestCourseForHTML(),
			outputPath:  "",
			wantErr:     true,
			errContains: "failed to create file",
		},
		{
			name:        "output path is directory",
			course:      createTestCourseForHTML(),
			outputPath:  t.TempDir(), // directory, not a file
			wantErr:     true,
			errContains: "failed to create file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exporter.Export(tt.course, tt.outputPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Export() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Export() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("Export() unexpected error = %v", err)
			}
		})
	}
}

// TestHTMLExporter_WriteHTML_ErrorPaths tests WriteHTML error conditions.
func TestHTMLExporter_WriteHTML_ErrorPaths(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner, nil)
	htmlExporter := exporter.(*HTMLExporter)

	tests := []struct {
		name        string
		course      *models.Course
		writer      *failingWriter
		wantErr     bool
		errContains string
	}{
		{
			name:        "writer fails immediately",
			course:      createTestCourseForHTML(),
			writer:      &failingWriter{failAfter: 0},
			wantErr:     true,
			errContains: "failed to execute template",
		},
		{
			name:        "writer fails mid-write",
			course:      createTestCourseForHTML(),
			writer:      &failingWriter{failAfter: 100},
			wantErr:     true,
			errContains: "failed to execute template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := htmlExporter.WriteHTML(tt.writer, tt.course)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WriteHTML() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("WriteHTML() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("WriteHTML() unexpected error = %v", err)
			}
		})
	}
}

// TestHTMLExporter_Export_NilCourse verifies that exporting a nil course panics.
// This documents the current behavior - nil course is not a supported input.
func TestHTMLExporter_Export_NilCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner, nil)

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "nil-course.html")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Export(nil) expected panic, got none")
		}
	}()

	// Export nil course - should panic
	_ = exporter.Export(nil, outputPath)
}

// TestHTMLExporter_WriteHTML_NilCourse verifies that WriteHTML with nil course panics.
func TestHTMLExporter_WriteHTML_NilCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner, nil)
	htmlExporter := exporter.(*HTMLExporter)

	defer func() {
		if r := recover(); r == nil {
			t.Error("WriteHTML(nil) expected panic, got none")
		}
	}()

	var buf strings.Builder
	_ = htmlExporter.WriteHTML(&buf, nil)
}

// TestHTMLExporter_WriteHTML_SuccessWithBuffer tests WriteHTML with a buffer.
func TestHTMLExporter_WriteHTML_SuccessWithBuffer(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewHTMLExporter(htmlCleaner, nil)
	htmlExporter := exporter.(*HTMLExporter)

	tests := []struct {
		name     string
		course   *models.Course
		contains []string
	}{
		{
			name:   "basic course",
			course: createTestCourseForHTML(),
			contains: []string{
				"<!DOCTYPE html>",
				"<html lang=\"en\">",
				"Test Course",
			},
		},
		{
			name: "course with special characters",
			course: &models.Course{
				ShareID: "special-chars-id",
				Course: models.CourseInfo{
					ID:          "special-course",
					Title:       "Course with <special> & \"characters\"",
					Description: "<p>Description with &amp; entities</p>",
					Lessons:     []models.Lesson{},
				},
			},
			contains: []string{
				"<!DOCTYPE html>",
				"Course with",
				"special",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := htmlExporter.WriteHTML(&buf, tt.course)
			if err != nil {
				t.Fatalf("WriteHTML() unexpected error: %v", err)
			}

			output := buf.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("WriteHTML() output missing %q", want)
				}
			}
		})
	}
}

// failingWriter is a test helper that fails after writing a certain number of bytes.
type failingWriter struct {
	written   int
	failAfter int
}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	if w.written >= w.failAfter {
		return 0, errWriteFailed
	}
	w.written += len(p)
	if w.written > w.failAfter {
		return 0, errWriteFailed
	}
	return len(p), nil
}

var errWriteFailed = &writeError{msg: "simulated write failure"}

type writeError struct {
	msg string
}

func (e *writeError) Error() string {
	return e.msg
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
	exporter := NewHTMLExporter(htmlCleaner, nil)
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
	exporter := NewHTMLExporter(htmlCleaner, nil)

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
