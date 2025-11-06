// Package exporters_test provides benchmarks for all exporters.
package exporters

import (
	"path/filepath"
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// BenchmarkFactory_CreateExporter_Markdown benchmarks markdown exporter creation.
func BenchmarkFactory_CreateExporter_Markdown(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	b.ResetTimer()
	for b.Loop() {
		_, _ = factory.CreateExporter("markdown")
	}
}

// BenchmarkFactory_CreateExporter_All benchmarks creating all exporter types.
func BenchmarkFactory_CreateExporter_All(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)
	formats := []string{"markdown", "docx", "html"}

	b.ResetTimer()
	for b.Loop() {
		for _, format := range formats {
			_, _ = factory.CreateExporter(format)
		}
	}
}

// BenchmarkAllExporters_Export benchmarks all exporters with the same course.
func BenchmarkAllExporters_Export(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	course := createBenchmarkCourse()

	exporters := map[string]struct {
		exporter any
		ext      string
	}{
		"Markdown": {NewMarkdownExporter(htmlCleaner), ".md"},
		"Docx":     {NewDocxExporter(htmlCleaner), ".docx"},
		"HTML":     {NewHTMLExporter(htmlCleaner), ".html"},
	}

	for name, exp := range exporters {
		b.Run(name, func(b *testing.B) {
			tempDir := b.TempDir()
			exporter := exp.exporter.(interface {
				Export(*models.Course, string) error
			})

			b.ResetTimer()
			for b.Loop() {
				outputPath := filepath.Join(tempDir, "benchmark"+exp.ext)
				_ = exporter.Export(course, outputPath)
			}
		})
	}
}

// BenchmarkExporters_LargeCourse benchmarks exporters with large course data.
func BenchmarkExporters_LargeCourse(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	course := createLargeBenchmarkCourse()

	b.Run("Markdown_Large", func(b *testing.B) {
		exporter := NewMarkdownExporter(htmlCleaner)
		tempDir := b.TempDir()

		b.ResetTimer()
		for b.Loop() {
			outputPath := filepath.Join(tempDir, "large.md")
			_ = exporter.Export(course, outputPath)
		}
	})

	b.Run("Docx_Large", func(b *testing.B) {
		exporter := NewDocxExporter(htmlCleaner)
		tempDir := b.TempDir()

		b.ResetTimer()
		for b.Loop() {
			outputPath := filepath.Join(tempDir, "large.docx")
			_ = exporter.Export(course, outputPath)
		}
	})

	b.Run("HTML_Large", func(b *testing.B) {
		exporter := NewHTMLExporter(htmlCleaner)
		tempDir := b.TempDir()

		b.ResetTimer()
		for b.Loop() {
			outputPath := filepath.Join(tempDir, "large.html")
			_ = exporter.Export(course, outputPath)
		}
	})
}

// createBenchmarkCourse creates a standard-sized course for benchmarking.
func createBenchmarkCourse() *models.Course {
	return &models.Course{
		ShareID: "benchmark-id",
		Author:  "Benchmark Author",
		Course: models.CourseInfo{
			ID:             "bench-course",
			Title:          "Benchmark Course",
			Description:    "Performance testing course",
			NavigationMode: "menu",
			Lessons: []models.Lesson{
				{
					ID:    "lesson1",
					Title: "Introduction",
					Type:  "lesson",
					Items: []models.Item{
						{
							Type: "text",
							Items: []models.SubItem{
								{
									Heading:   "Welcome",
									Paragraph: "<p>This is a test paragraph with <strong>HTML</strong> content.</p>",
								},
							},
						},
						{
							Type: "list",
							Items: []models.SubItem{
								{Paragraph: "Item 1"},
								{Paragraph: "Item 2"},
								{Paragraph: "Item 3"},
							},
						},
					},
				},
			},
		},
	}
}

// createLargeBenchmarkCourse creates a large course for stress testing.
func createLargeBenchmarkCourse() *models.Course {
	lessons := make([]models.Lesson, 50)
	for i := range 50 {
		lessons[i] = models.Lesson{
			ID:          string(rune(i)),
			Title:       "Lesson " + string(rune(i)),
			Type:        "lesson",
			Description: "<p>This is lesson description with <em>formatting</em>.</p>",
			Items: []models.Item{
				{
					Type: "text",
					Items: []models.SubItem{
						{
							Heading:   "Section Heading",
							Paragraph: "<p>Content with <strong>bold</strong> and <em>italic</em> text.</p>",
						},
					},
				},
				{
					Type: "list",
					Items: []models.SubItem{
						{Paragraph: "Point 1"},
						{Paragraph: "Point 2"},
						{Paragraph: "Point 3"},
					},
				},
				{
					Type: "knowledgeCheck",
					Items: []models.SubItem{
						{
							Title: "Quiz Question",
							Answers: []models.Answer{
								{Title: "Answer A", Correct: false},
								{Title: "Answer B", Correct: true},
								{Title: "Answer C", Correct: false},
							},
							Feedback: "Good job!",
						},
					},
				},
			},
		}
	}

	return &models.Course{
		ShareID: "large-benchmark-id",
		Author:  "Benchmark Author",
		Course: models.CourseInfo{
			ID:          "large-bench-course",
			Title:       "Large Benchmark Course",
			Description: "Large performance testing course",
			Lessons:     lessons,
		},
	}
}
