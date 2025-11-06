// Package exporters_test provides tests for the exporter factory.
package exporters

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/kjanat/articulate-parser/internal/services"
)

// TestNewFactory tests the NewFactory constructor.
func TestNewFactory(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	if factory == nil {
		t.Fatal("NewFactory() returned nil")
	}

	// Type assertion to check internal structure
	factoryImpl, ok := factory.(*Factory)
	if !ok {
		t.Fatal("NewFactory() returned wrong type")
	}

	if factoryImpl.htmlCleaner == nil {
		t.Error("htmlCleaner should not be nil")
	}
}

// TestFactory_CreateExporter tests the CreateExporter method for all supported formats.
func TestFactory_CreateExporter(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	testCases := []struct {
		name           string
		format         string
		expectedType   string
		expectedFormat string
		shouldError    bool
	}{
		{
			name:           "markdown format",
			format:         "markdown",
			expectedType:   "*exporters.MarkdownExporter",
			expectedFormat: "markdown",
			shouldError:    false,
		},
		{
			name:           "md format alias",
			format:         "md",
			expectedType:   "*exporters.MarkdownExporter",
			expectedFormat: "markdown",
			shouldError:    false,
		},
		{
			name:           "docx format",
			format:         "docx",
			expectedType:   "*exporters.DocxExporter",
			expectedFormat: "docx",
			shouldError:    false,
		},
		{
			name:           "word format alias",
			format:         "word",
			expectedType:   "*exporters.DocxExporter",
			expectedFormat: "docx",
			shouldError:    false,
		},
		{
			name:           "html format",
			format:         "html",
			expectedType:   "*exporters.HTMLExporter",
			expectedFormat: "html",
			shouldError:    false,
		},
		{
			name:           "htm format alias",
			format:         "htm",
			expectedType:   "*exporters.HTMLExporter",
			expectedFormat: "html",
			shouldError:    false,
		},
		{
			name:        "unsupported format",
			format:      "pdf",
			shouldError: true,
		},
		{
			name:        "empty format",
			format:      "",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exporter, err := factory.CreateExporter(tc.format)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error for format '%s', but got nil", tc.format)
				}
				if exporter != nil {
					t.Errorf("Expected nil exporter for unsupported format '%s'", tc.format)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error creating exporter for format '%s': %v", tc.format, err)
			}

			if exporter == nil {
				t.Fatalf("CreateExporter returned nil for supported format '%s'", tc.format)
			}

			// Check type
			exporterType := reflect.TypeOf(exporter).String()
			if exporterType != tc.expectedType {
				t.Errorf("Expected exporter type '%s' for format '%s', got '%s'", tc.expectedType, tc.format, exporterType)
			}

			// Check supported format
			supportedFormat := exporter.SupportedFormat()
			if supportedFormat != tc.expectedFormat {
				t.Errorf("Expected supported format '%s' for format '%s', got '%s'", tc.expectedFormat, tc.format, supportedFormat)
			}
		})
	}
}

// TestFactory_CreateExporter_CaseInsensitive tests that format strings are case-insensitive.
func TestFactory_CreateExporter_CaseInsensitive(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	testCases := []struct {
		format         string
		expectedFormat string
	}{
		{"MARKDOWN", "markdown"},
		{"Markdown", "markdown"},
		{"MarkDown", "markdown"},
		{"MD", "markdown"},
		{"Md", "markdown"},
		{"DOCX", "docx"},
		{"Docx", "docx"},
		{"DocX", "docx"},
		{"WORD", "docx"},
		{"Word", "docx"},
		{"WoRd", "docx"},
		{"HTML", "html"},
		{"Html", "html"},
		{"HtMl", "html"},
		{"HTM", "html"},
		{"Htm", "html"},
		{"HtM", "html"},
	}

	for _, tc := range testCases {
		t.Run(tc.format, func(t *testing.T) {
			exporter, err := factory.CreateExporter(tc.format)

			if err != nil {
				t.Fatalf("Unexpected error for format '%s': %v", tc.format, err)
			}

			if exporter == nil {
				t.Fatalf("CreateExporter returned nil for format '%s'", tc.format)
			}

			supportedFormat := exporter.SupportedFormat()
			if supportedFormat != tc.expectedFormat {
				t.Errorf("Expected supported format '%s' for format '%s', got '%s'", tc.expectedFormat, tc.format, supportedFormat)
			}
		})
	}
}

// TestFactory_CreateExporter_ErrorMessages tests error messages for unsupported formats.
func TestFactory_CreateExporter_ErrorMessages(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	testCases := []string{
		"pdf",
		"txt",
		"json",
		"xml",
		"unknown",
		"123",
		"markdown-invalid",
	}

	for _, format := range testCases {
		t.Run(format, func(t *testing.T) {
			exporter, err := factory.CreateExporter(format)

			if err == nil {
				t.Errorf("Expected error for unsupported format '%s', got nil", format)
			}

			if exporter != nil {
				t.Errorf("Expected nil exporter for unsupported format '%s', got %v", format, exporter)
			}

			// Check error message contains the format
			if err != nil && !strings.Contains(err.Error(), format) {
				t.Errorf("Error message should contain the unsupported format '%s', got: %s", format, err.Error())
			}

			// Check error message has expected prefix
			if err != nil && !strings.Contains(err.Error(), "unsupported export format") {
				t.Errorf("Error message should contain 'unsupported export format', got: %s", err.Error())
			}
		})
	}
}

// TestFactory_SupportedFormats tests the SupportedFormats method.
func TestFactory_SupportedFormats(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	formats := factory.SupportedFormats()

	if formats == nil {
		t.Fatal("SupportedFormats() returned nil")
	}

	expected := []string{"markdown", "md", "docx", "word", "html", "htm"}

	// Sort both slices for comparison
	sort.Strings(formats)
	sort.Strings(expected)

	if !reflect.DeepEqual(formats, expected) {
		t.Errorf("Expected formats %v, got %v", expected, formats)
	}

	// Verify all returned formats can create exporters
	for _, format := range formats {
		exporter, err := factory.CreateExporter(format)
		if err != nil {
			t.Errorf("Format '%s' from SupportedFormats() should be creatable, got error: %v", format, err)
		}
		if exporter == nil {
			t.Errorf("Format '%s' from SupportedFormats() should create non-nil exporter", format)
		}
	}
}

// TestFactory_SupportedFormats_Immutable tests that the returned slice is safe to modify.
func TestFactory_SupportedFormats_Immutable(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	// Get formats twice
	formats1 := factory.SupportedFormats()
	formats2 := factory.SupportedFormats()

	// Modify first slice
	if len(formats1) > 0 {
		formats1[0] = "modified"
	}

	// Check that second call returns unmodified data
	if len(formats2) > 0 && formats2[0] == "modified" {
		t.Error("SupportedFormats() should return independent slices")
	}

	// Verify original functionality still works
	formats3 := factory.SupportedFormats()
	if len(formats3) == 0 {
		t.Error("SupportedFormats() should still return formats after modification")
	}
}

// TestFactory_ExporterTypes tests that created exporters are of correct types.
func TestFactory_ExporterTypes(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	// Test markdown exporter
	markdownExporter, err := factory.CreateExporter("markdown")
	if err != nil {
		t.Fatalf("Failed to create markdown exporter: %v", err)
	}

	if _, ok := markdownExporter.(*MarkdownExporter); !ok {
		t.Error("Markdown exporter should be of type *MarkdownExporter")
	}

	// Test docx exporter
	docxExporter, err := factory.CreateExporter("docx")
	if err != nil {
		t.Fatalf("Failed to create docx exporter: %v", err)
	}

	if _, ok := docxExporter.(*DocxExporter); !ok {
		t.Error("DOCX exporter should be of type *DocxExporter")
	}
}

// TestFactory_HTMLCleanerPropagation tests that HTMLCleaner is properly passed to exporters.
func TestFactory_HTMLCleanerPropagation(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	// Test with markdown exporter
	markdownExporter, err := factory.CreateExporter("markdown")
	if err != nil {
		t.Fatalf("Failed to create markdown exporter: %v", err)
	}

	markdownImpl, ok := markdownExporter.(*MarkdownExporter)
	if !ok {
		t.Fatal("Failed to cast to MarkdownExporter")
	}

	if markdownImpl.htmlCleaner == nil {
		t.Error("HTMLCleaner should be propagated to MarkdownExporter")
	}

	// Test with docx exporter
	docxExporter, err := factory.CreateExporter("docx")
	if err != nil {
		t.Fatalf("Failed to create docx exporter: %v", err)
	}

	docxImpl, ok := docxExporter.(*DocxExporter)
	if !ok {
		t.Fatal("Failed to cast to DocxExporter")
	}

	if docxImpl.htmlCleaner == nil {
		t.Error("HTMLCleaner should be propagated to DocxExporter")
	}

	// Test with html exporter
	htmlExporter, err := factory.CreateExporter("html")
	if err != nil {
		t.Fatalf("Failed to create html exporter: %v", err)
	}

	htmlImpl, ok := htmlExporter.(*HTMLExporter)
	if !ok {
		t.Fatal("Failed to cast to HTMLExporter")
	}

	if htmlImpl.htmlCleaner == nil {
		t.Error("HTMLCleaner should be propagated to HTMLExporter")
	}
}

// TestFactory_MultipleExporterCreation tests creating multiple exporters of same type.
func TestFactory_MultipleExporterCreation(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	// Create multiple markdown exporters
	exporter1, err := factory.CreateExporter("markdown")
	if err != nil {
		t.Fatalf("Failed to create first markdown exporter: %v", err)
	}

	exporter2, err := factory.CreateExporter("md")
	if err != nil {
		t.Fatalf("Failed to create second markdown exporter: %v", err)
	}

	// They should be different instances
	if exporter1 == exporter2 {
		t.Error("Factory should create independent exporter instances")
	}

	// But both should be MarkdownExporter type
	if _, ok := exporter1.(*MarkdownExporter); !ok {
		t.Error("First exporter should be MarkdownExporter")
	}

	if _, ok := exporter2.(*MarkdownExporter); !ok {
		t.Error("Second exporter should be MarkdownExporter")
	}
}

// TestFactory_WithNilHTMLCleaner tests factory behavior with nil HTMLCleaner.
func TestFactory_WithNilHTMLCleaner(t *testing.T) {
	// This tests edge case - should not panic but behavior may vary
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Factory should handle nil HTMLCleaner gracefully, but panicked: %v", r)
		}
	}()

	factory := NewFactory(nil)

	if factory == nil {
		t.Fatal("NewFactory(nil) returned nil")
	}

	// Try to create an exporter - this might fail or succeed depending on implementation
	_, err := factory.CreateExporter("markdown")

	// We don't assert on the error since nil HTMLCleaner handling is implementation-dependent
	// The important thing is that it doesn't panic
	_ = err
}

// TestFactory_FormatNormalization tests that format strings are properly normalized.
func TestFactory_FormatNormalization(t *testing.T) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	// Test formats with extra whitespace
	testCases := []struct {
		input    string
		expected string
	}{
		{"markdown", "markdown"},
		{"MARKDOWN", "markdown"},
		{"Markdown", "markdown"},
		{"docx", "docx"},
		{"DOCX", "docx"},
		{"Docx", "docx"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			exporter, err := factory.CreateExporter(tc.input)
			if err != nil {
				t.Fatalf("Failed to create exporter for '%s': %v", tc.input, err)
			}

			format := exporter.SupportedFormat()
			if format != tc.expected {
				t.Errorf("Expected format '%s' for input '%s', got '%s'", tc.expected, tc.input, format)
			}
		})
	}
}

// BenchmarkFactory_CreateExporter benchmarks the CreateExporter method.
func BenchmarkFactory_CreateExporter(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	for b.Loop() {
		_, _ = factory.CreateExporter("markdown")
	}
}

// BenchmarkFactory_CreateExporter_Docx benchmarks creating DOCX exporters.
func BenchmarkFactory_CreateExporter_Docx(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	for b.Loop() {
		_, _ = factory.CreateExporter("docx")
	}
}

// BenchmarkFactory_SupportedFormats benchmarks the SupportedFormats method.
func BenchmarkFactory_SupportedFormats(b *testing.B) {
	htmlCleaner := services.NewHTMLCleaner()
	factory := NewFactory(htmlCleaner)

	for b.Loop() {
		_ = factory.SupportedFormats()
	}
}
