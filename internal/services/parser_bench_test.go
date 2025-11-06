package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
)

// BenchmarkArticulateParser_FetchCourse benchmarks the FetchCourse method.
func BenchmarkArticulateParser_FetchCourse(b *testing.B) {
	testCourse := &models.Course{
		ShareID: "benchmark-id",
		Author:  "Benchmark Author",
		Course: models.CourseInfo{
			ID:          "bench-course",
			Title:       "Benchmark Course",
			Description: "Testing performance",
			Lessons: []models.Lesson{
				{
					ID:    "lesson1",
					Title: "Lesson 1",
					Type:  "lesson",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Encode errors are ignored in benchmarks; the test server's ResponseWriter
		// writes are reliable and any encoding error would be a test setup issue
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client:  &http.Client{},
		Logger:  NewNoOpLogger(),
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := parser.FetchCourse(context.Background(), "https://rise.articulate.com/share/benchmark-id")
		if err != nil {
			b.Fatalf("FetchCourse failed: %v", err)
		}
	}
}

// BenchmarkArticulateParser_FetchCourse_LargeCourse benchmarks with a large course.
func BenchmarkArticulateParser_FetchCourse_LargeCourse(b *testing.B) {
	// Create a large course with many lessons
	lessons := make([]models.Lesson, 100)
	for i := range 100 {
		lessons[i] = models.Lesson{
			ID:          string(rune(i)),
			Title:       "Lesson " + string(rune(i)),
			Type:        "lesson",
			Description: "This is a test lesson with some description",
			Items: []models.Item{
				{
					Type: "text",
					Items: []models.SubItem{
						{
							Heading:   "Test Heading",
							Paragraph: "Test paragraph content with some text",
						},
					},
				},
			},
		}
	}

	testCourse := &models.Course{
		ShareID: "large-course-id",
		Author:  "Benchmark Author",
		Course: models.CourseInfo{
			ID:          "large-course",
			Title:       "Large Benchmark Course",
			Description: "Testing performance with large course",
			Lessons:     lessons,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Encode errors are ignored in benchmarks; the test server's ResponseWriter
		// writes are reliable and any encoding error would be a test setup issue
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client:  &http.Client{},
		Logger:  NewNoOpLogger(),
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := parser.FetchCourse(context.Background(), "https://rise.articulate.com/share/large-course-id")
		if err != nil {
			b.Fatalf("FetchCourse failed: %v", err)
		}
	}
}

// BenchmarkArticulateParser_LoadCourseFromFile benchmarks loading from file.
func BenchmarkArticulateParser_LoadCourseFromFile(b *testing.B) {
	testCourse := &models.Course{
		ShareID: "file-test-id",
		Course: models.CourseInfo{
			Title: "File Test Course",
		},
	}

	tempDir := b.TempDir()
	tempFile := filepath.Join(tempDir, "benchmark.json")

	data, err := json.Marshal(testCourse)
	if err != nil {
		b.Fatalf("Failed to marshal: %v", err)
	}

	if err := os.WriteFile(tempFile, data, 0o644); err != nil {
		b.Fatalf("Failed to write file: %v", err)
	}

	parser := NewArticulateParser(nil, "", 0)

	b.ResetTimer()
	for b.Loop() {
		_, err := parser.LoadCourseFromFile(tempFile)
		if err != nil {
			b.Fatalf("LoadCourseFromFile failed: %v", err)
		}
	}
}

// BenchmarkArticulateParser_LoadCourseFromFile_Large benchmarks with large file.
func BenchmarkArticulateParser_LoadCourseFromFile_Large(b *testing.B) {
	// Create a large course
	lessons := make([]models.Lesson, 200)
	for i := range 200 {
		lessons[i] = models.Lesson{
			ID:    string(rune(i)),
			Title: "Lesson " + string(rune(i)),
			Type:  "lesson",
			Items: []models.Item{
				{Type: "text", Items: []models.SubItem{{Heading: "H", Paragraph: "P"}}},
				{Type: "list", Items: []models.SubItem{{Paragraph: "Item 1"}, {Paragraph: "Item 2"}}},
			},
		}
	}

	testCourse := &models.Course{
		ShareID: "large-file-id",
		Course: models.CourseInfo{
			Title:   "Large File Course",
			Lessons: lessons,
		},
	}

	tempDir := b.TempDir()
	tempFile := filepath.Join(tempDir, "large-benchmark.json")

	data, err := json.Marshal(testCourse)
	if err != nil {
		b.Fatalf("Failed to marshal: %v", err)
	}

	if err := os.WriteFile(tempFile, data, 0o644); err != nil {
		b.Fatalf("Failed to write file: %v", err)
	}

	parser := NewArticulateParser(nil, "", 0)

	b.ResetTimer()
	for b.Loop() {
		_, err := parser.LoadCourseFromFile(tempFile)
		if err != nil {
			b.Fatalf("LoadCourseFromFile failed: %v", err)
		}
	}
}

// BenchmarkArticulateParser_ExtractShareID benchmarks share ID extraction.
func BenchmarkArticulateParser_ExtractShareID(b *testing.B) {
	parser := &ArticulateParser{}
	uri := "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/"

	b.ResetTimer()
	for b.Loop() {
		_, err := parser.extractShareID(uri)
		if err != nil {
			b.Fatalf("extractShareID failed: %v", err)
		}
	}
}

// BenchmarkArticulateParser_BuildAPIURL benchmarks API URL building.
func BenchmarkArticulateParser_BuildAPIURL(b *testing.B) {
	parser := &ArticulateParser{
		BaseURL: "https://rise.articulate.com",
	}
	shareID := "test-share-id-12345"

	b.ResetTimer()
	for b.Loop() {
		_ = parser.buildAPIURL(shareID)
	}
}
