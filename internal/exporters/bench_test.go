package exporters

import (
	"path/filepath"
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
	"github.com/kjanat/articulate-parser/internal/testdata"
)

// BenchmarkFactory_CreateExporter_Markdown benchmarks markdown exporter creation.
func BenchmarkFactory_CreateExporter_Markdown(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner, nil)

	b.ResetTimer()
	for b.Loop() {
		_, _ = factory.CreateExporter("markdown")
	}
}

// BenchmarkFactory_CreateExporter_All benchmarks creating all exporter types.
func BenchmarkFactory_CreateExporter_All(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner, nil)
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
	course := testdata.BasicCourse()

	exporters := map[string]struct {
		exporter any
		ext      string
	}{
		"Markdown": {NewMarkdownExporter(htmlCleaner, nil), ".md"},
		"Docx":     {NewDocxExporter(htmlCleaner, nil), ".docx"},
		"HTML":     {NewHTMLExporter(htmlCleaner, nil), ".html"},
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
	course := testdata.LargeCourse(50)

	b.Run("Markdown_Large", func(b *testing.B) {
		exporter := NewMarkdownExporter(htmlCleaner, nil)
		tempDir := b.TempDir()

		b.ResetTimer()
		for b.Loop() {
			outputPath := filepath.Join(tempDir, "large.md")
			_ = exporter.Export(course, outputPath)
		}
	})

	b.Run("Docx_Large", func(b *testing.B) {
		exporter := NewDocxExporter(htmlCleaner, nil)
		tempDir := b.TempDir()

		b.ResetTimer()
		for b.Loop() {
			outputPath := filepath.Join(tempDir, "large.docx")
			_ = exporter.Export(course, outputPath)
		}
	})

	b.Run("HTML_Large", func(b *testing.B) {
		exporter := NewHTMLExporter(htmlCleaner, nil)
		tempDir := b.TempDir()

		b.ResetTimer()
		for b.Loop() {
			outputPath := filepath.Join(tempDir, "large.html")
			_ = exporter.Export(course, outputPath)
		}
	})
}
