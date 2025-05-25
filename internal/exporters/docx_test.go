// Package exporters_test provides tests for the docx exporter.
package exporters

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// TestNewDocxExporter tests the NewDocxExporter constructor.
func TestNewDocxExporter(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	if exporter == nil {
		t.Fatal("NewDocxExporter() returned nil")
	}

	// Type assertion to check internal structure
	docxExporter, ok := exporter.(*DocxExporter)
	if !ok {
		t.Fatal("NewDocxExporter() returned wrong type")
	}

	if docxExporter.htmlCleaner == nil {
		t.Error("htmlCleaner should not be nil")
	}
}

// TestDocxExporter_GetSupportedFormat tests the GetSupportedFormat method.
func TestDocxExporter_GetSupportedFormat(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	expected := "docx"
	result := exporter.GetSupportedFormat()

	if result != expected {
		t.Errorf("Expected format '%s', got '%s'", expected, result)
	}
}

// TestDocxExporter_Export tests the Export method.
func TestDocxExporter_Export(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	// Create test course
	testCourse := createTestCourseForDocx()

	// Create temporary directory and file
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test-course.docx")

	// Test successful export
	err := exporter.Export(testCourse, outputPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Verify file has some content (basic check)
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Output file is empty")
	}
}

// TestDocxExporter_Export_AddDocxExtension tests that the .docx extension is added automatically.
func TestDocxExporter_Export_AddDocxExtension(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	testCourse := createTestCourseForDocx()

	// Create temporary directory and file without .docx extension
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test-course")

	err := exporter.Export(testCourse, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Check that file was created with .docx extension
	expectedPath := outputPath + ".docx"
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatal("Output file with .docx extension was not created")
	}
}

// TestDocxExporter_Export_InvalidPath tests export with invalid output path.
func TestDocxExporter_Export_InvalidPath(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	testCourse := createTestCourseForDocx()

	// Try to write to invalid path
	invalidPath := "/invalid/path/that/does/not/exist/file.docx"
	err := exporter.Export(testCourse, invalidPath)

	if err == nil {
		t.Fatal("Expected error for invalid path, got nil")
	}
}

// TestDocxExporter_ExportLesson tests the exportLesson method indirectly through Export.
func TestDocxExporter_ExportLesson(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	// Create course with specific lesson content
	course := &models.Course{
		ShareID: "test-id",
		Course: models.CourseInfo{
			ID:    "test-course",
			Title: "Test Course",
			Lessons: []models.Lesson{
				{
					ID:          "lesson-1",
					Title:       "Test Lesson",
					Type:        "lesson",
					Description: "<p>Test lesson description with <strong>bold</strong> text.</p>",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Title:     "Test Item Title",
									Paragraph: "<p>Test paragraph content.</p>",
								},
							},
						},
					},
				},
			},
		},
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "lesson-test.docx")

	err := exporter.Export(course, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created successfully
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

// TestDocxExporter_ExportItem tests the exportItem method indirectly through Export.
func TestDocxExporter_ExportItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	// Create course with different item types
	course := &models.Course{
		ShareID: "test-id",
		Course: models.CourseInfo{
			ID:    "test-course",
			Title: "Item Test Course",
			Lessons: []models.Lesson{
				{
					ID:    "lesson-1",
					Title: "Item Types Lesson",
					Type:  "lesson",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Title:     "Text Item",
									Paragraph: "<p>Text content</p>",
								},
							},
						},
						{
							Type: "list",
							Items: []models.SubItem{
								{Paragraph: "<p>List item 1</p>"},
								{Paragraph: "<p>List item 2</p>"},
							},
						},
						{
							Type: "knowledgeCheck",
							Items: []models.SubItem{
								{
									Title: "<p>What is the answer?</p>",
									Answers: []models.Answer{
										{Title: "Option A", Correct: false},
										{Title: "Option B", Correct: true},
									},
									Feedback: "<p>Correct answer explanation</p>",
								},
							},
						},
					},
				},
			},
		},
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "items-test.docx")

	err := exporter.Export(course, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created successfully
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

// TestDocxExporter_ExportSubItem tests the exportSubItem method indirectly through Export.
func TestDocxExporter_ExportSubItem(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	// Create course with sub-item containing all possible fields
	course := &models.Course{
		ShareID: "test-id",
		Course: models.CourseInfo{
			ID:    "test-course",
			Title: "SubItem Test Course",
			Lessons: []models.Lesson{
				{
					ID:    "lesson-1",
					Title: "SubItem Test Lesson",
					Type:  "lesson",
					Items: []models.Item{
						{
							Type: "knowledgeCheck",
							Items: []models.SubItem{
								{
									Title:     "<p>Question Title</p>",
									Heading:   "<h3>Question Heading</h3>",
									Paragraph: "<p>Question description with <em>emphasis</em>.</p>",
									Answers: []models.Answer{
										{Title: "Wrong answer", Correct: false},
										{Title: "Correct answer", Correct: true},
										{Title: "Another wrong answer", Correct: false},
									},
									Feedback: "<p>Feedback with <strong>formatting</strong>.</p>",
								},
							},
						},
					},
				},
			},
		},
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "subitem-test.docx")

	err := exporter.Export(course, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created successfully
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

// TestDocxExporter_ComplexCourse tests export of a complex course structure.
func TestDocxExporter_ComplexCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	// Create complex test course
	course := &models.Course{
		ShareID: "complex-test-id",
		Course: models.CourseInfo{
			ID:          "complex-course",
			Title:       "Complex Test Course",
			Description: "<p>This is a <strong>complex</strong> course description with <em>formatting</em>.</p>",
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
					Description: "<p>Introduction to the course with <code>code</code> and <a href='#'>links</a>.</p>",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Heading:   "<h2>Welcome</h2>",
									Paragraph: "<p>Welcome to our comprehensive course!</p>",
								},
							},
						},
						{
							Type: "list",
							Items: []models.SubItem{
								{Paragraph: "<p>Learn advanced concepts</p>"},
								{Paragraph: "<p>Practice with real examples</p>"},
								{Paragraph: "<p>Apply knowledge in projects</p>"},
							},
						},
						{
							Type: "multimedia",
							Items: []models.SubItem{
								{
									Title:   "<p>Video Introduction</p>",
									Caption: "<p>Watch this introductory video</p>",
									Media: &models.Media{
										Video: &models.VideoMedia{
											OriginalUrl: "https://example.com/intro.mp4",
											Duration:    300,
										},
									},
								},
							},
						},
						{
							Type: "knowledgeCheck",
							Items: []models.SubItem{
								{
									Title: "<p>What will you learn in this course?</p>",
									Answers: []models.Answer{
										{Title: "Basic concepts only", Correct: false},
										{Title: "Advanced concepts and practical application", Correct: true},
										{Title: "Theory without practice", Correct: false},
									},
									Feedback: "<p>Excellent! This course covers both theory and practice.</p>",
								},
							},
						},
						{
							Type: "image",
							Items: []models.SubItem{
								{
									Caption: "<p>Course overview diagram</p>",
									Media: &models.Media{
										Image: &models.ImageMedia{
											OriginalUrl: "https://example.com/overview.png",
										},
									},
								},
							},
						},
						{
							Type: "interactive",
							Items: []models.SubItem{
								{
									Title: "<p>Interactive Exercise</p>",
								},
							},
						},
					},
				},
				{
					ID:    "lesson-2",
					Title: "Advanced Topics",
					Type:  "lesson",
					Items: []models.Item{
						{
							Type: "divider",
						},
						{
							Type: "unknown",
							Items: []models.SubItem{
								{
									Title:     "<p>Custom Content</p>",
									Paragraph: "<p>This is custom content type</p>",
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
	outputPath := filepath.Join(tempDir, "complex-course.docx")

	// Export course
	err := exporter.Export(course, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created and has reasonable size
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() < 1000 {
		t.Error("Output file seems too small for complex course content")
	}
}

// TestDocxExporter_EmptyCourse tests export of an empty course.
func TestDocxExporter_EmptyCourse(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

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
	outputPath := filepath.Join(tempDir, "empty-course.docx")

	err := exporter.Export(course, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

// TestDocxExporter_HTMLCleaning tests that HTML content is properly cleaned.
func TestDocxExporter_HTMLCleaning(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	// Create course with HTML content that needs cleaning
	course := &models.Course{
		ShareID: "html-test-id",
		Course: models.CourseInfo{
			ID:          "html-test-course",
			Title:       "HTML Cleaning Test",
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
	outputPath := filepath.Join(tempDir, "html-cleaning-test.docx")

	err := exporter.Export(course, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Verify file was created (basic check that HTML cleaning didn't break export)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

// TestDocxExporter_ExistingDocxExtension tests that existing .docx extension is preserved.
func TestDocxExporter_ExistingDocxExtension(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	testCourse := createTestCourseForDocx()

	// Use path that already has .docx extension
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test-course.docx")

	err := exporter.Export(testCourse, outputPath)
	if err != nil {

		t.Fatalf("Export failed: %v", err)
	}

	// Check that file was created at the exact path (no double extension)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output file was not created at expected path")
	}

	// Ensure no double extension was created
	doubleExtensionPath := outputPath + ".docx"
	if _, err := os.Stat(doubleExtensionPath); err == nil {
		t.Error("Double .docx extension file should not exist")
	}
}

// TestDocxExporter_CaseInsensitiveExtension tests that extension checking is case-insensitive.
func TestDocxExporter_CaseInsensitiveExtension(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

	testCourse := createTestCourseForDocx()

	// Test various case combinations
	testCases := []string{
		"test-course.DOCX",
		"test-course.Docx",
		"test-course.DocX",
	}

	for i, testCase := range testCases {
		tempDir := t.TempDir()
		outputPath := filepath.Join(tempDir, testCase)

		err := exporter.Export(testCourse, outputPath)
		if err != nil {

			t.Fatalf("Export failed for case %d (%s): %v", i, testCase, err)
		}

		// Check that file was created at the exact path (no additional extension)
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Fatalf("Output file was not created at expected path for case %d (%s)", i, testCase)
		}
	}
}

// createTestCourseForDocx creates a test course for DOCX export testing.
func createTestCourseForDocx() *models.Course {
	return &models.Course{
		ShareID: "test-share-id",
		Course: models.CourseInfo{
			ID:          "test-course-id",
			Title:       "Test Course",
			Description: "<p>Test course description with <strong>formatting</strong>.</p>",
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

// BenchmarkDocxExporter_Export benchmarks the Export method.
func BenchmarkDocxExporter_Export(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)
	course := createTestCourseForDocx()

	// Create temporary directory
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tempDir, "benchmark-course.docx")
		_ = exporter.Export(course, outputPath)
		// Clean up for next iteration
		os.Remove(outputPath)
	}
}

// BenchmarkDocxExporter_ComplexCourse benchmarks export of a complex course.
func BenchmarkDocxExporter_ComplexCourse(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	exporter := NewDocxExporter(htmlCleaner)

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
	for i := 0; i < 10; i++ {
		lesson := models.Lesson{
			ID:    "lesson-" + string(rune(i)),
			Title: "Lesson " + string(rune(i)),
			Type:  "lesson",
			Items: make([]models.Item, 5), // 5 items per lesson
		}

		for j := 0; j < 5; j++ {
			item := models.Item{
				Type:  "text",
				Items: make([]models.SubItem, 3), // 3 sub-items per item
			}

			for k := 0; k < 3; k++ {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tempDir, "benchmark-complex.docx")
		_ = exporter.Export(course, outputPath)
		os.Remove(outputPath)
	}
}
