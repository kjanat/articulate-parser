// Package services_test provides tests for the services package.
package services

import (
	"errors"
	"testing"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
)

// MockCourseParser is a mock implementation of interfaces.CourseParser for testing.
type MockCourseParser struct {
	mockFetchCourse        func(uri string) (*models.Course, error)
	mockLoadCourseFromFile func(filePath string) (*models.Course, error)
}

func (m *MockCourseParser) FetchCourse(uri string) (*models.Course, error) {
	if m.mockFetchCourse != nil {
		return m.mockFetchCourse(uri)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCourseParser) LoadCourseFromFile(filePath string) (*models.Course, error) {
	if m.mockLoadCourseFromFile != nil {
		return m.mockLoadCourseFromFile(filePath)
	}
	return nil, errors.New("not implemented")
}

// MockExporter is a mock implementation of interfaces.Exporter for testing.
type MockExporter struct {
	mockExport             func(course *models.Course, outputPath string) error
	mockGetSupportedFormat func() string
}

func (m *MockExporter) Export(course *models.Course, outputPath string) error {
	if m.mockExport != nil {
		return m.mockExport(course, outputPath)
	}
	return nil
}

func (m *MockExporter) GetSupportedFormat() string {
	if m.mockGetSupportedFormat != nil {
		return m.mockGetSupportedFormat()
	}
	return "mock"
}

// MockExporterFactory is a mock implementation of interfaces.ExporterFactory for testing.
type MockExporterFactory struct {
	mockCreateExporter      func(format string) (*MockExporter, error)
	mockGetSupportedFormats func() []string
}

func (m *MockExporterFactory) CreateExporter(format string) (interfaces.Exporter, error) {
	if m.mockCreateExporter != nil {
		exporter, err := m.mockCreateExporter(format)
		return exporter, err
	}
	return &MockExporter{}, nil
}

func (m *MockExporterFactory) GetSupportedFormats() []string {
	if m.mockGetSupportedFormats != nil {
		return m.mockGetSupportedFormats()
	}
	return []string{"mock"}
}

// createTestCourse creates a sample course for testing purposes.
func createTestCourse() *models.Course {
	return &models.Course{
		ShareID: "test-share-id",
		Author:  "Test Author",
		Course: models.CourseInfo{
			ID:          "test-course-id",
			Title:       "Test Course",
			Description: "This is a test course",
			Lessons: []models.Lesson{
				{
					ID:    "lesson-1",
					Title: "Test Lesson",
					Type:  "lesson",
					Items: []models.Item{
						{
							ID:   "item-1",
							Type: "text",
							Items: []models.SubItem{
								{
									ID:        "subitem-1",
									Title:     "Test Title",
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

// TestNewApp tests the NewApp constructor.
func TestNewApp(t *testing.T) {
	parser := &MockCourseParser{}
	factory := &MockExporterFactory{}

	app := NewApp(parser, factory)

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	if app.parser != parser {
		t.Error("App parser was not set correctly")
	}

	// Test that the factory is set (we can't directly compare interface values)
	formats := app.GetSupportedFormats()
	if len(formats) == 0 {
		t.Error("App exporterFactory was not set correctly - no supported formats")
	}
}

// TestApp_ProcessCourseFromFile tests the ProcessCourseFromFile method.
func TestApp_ProcessCourseFromFile(t *testing.T) {
	testCourse := createTestCourse()

	tests := []struct {
		name          string
		filePath      string
		format        string
		outputPath    string
		setupMocks    func(*MockCourseParser, *MockExporterFactory, *MockExporter)
		expectedError string
	}{
		{
			name:       "successful processing",
			filePath:   "test.json",
			format:     "markdown",
			outputPath: "output.md",
			setupMocks: func(parser *MockCourseParser, factory *MockExporterFactory, exporter *MockExporter) {
				parser.mockLoadCourseFromFile = func(filePath string) (*models.Course, error) {
					if filePath != "test.json" {
						t.Errorf("Expected filePath 'test.json', got '%s'", filePath)
					}
					return testCourse, nil
				}

				factory.mockCreateExporter = func(format string) (*MockExporter, error) {
					if format != "markdown" {
						t.Errorf("Expected format 'markdown', got '%s'", format)
					}
					return exporter, nil
				}

				exporter.mockExport = func(course *models.Course, outputPath string) error {
					if outputPath != "output.md" {
						t.Errorf("Expected outputPath 'output.md', got '%s'", outputPath)
					}
					if course != testCourse {
						t.Error("Expected course to match testCourse")
					}
					return nil
				}
			},
		},
		{
			name:       "file loading error",
			filePath:   "nonexistent.json",
			format:     "markdown",
			outputPath: "output.md",
			setupMocks: func(parser *MockCourseParser, factory *MockExporterFactory, exporter *MockExporter) {
				parser.mockLoadCourseFromFile = func(filePath string) (*models.Course, error) {
					return nil, errors.New("file not found")
				}
			},
			expectedError: "failed to load course from file",
		},
		{
			name:       "exporter creation error",
			filePath:   "test.json",
			format:     "unsupported",
			outputPath: "output.txt",
			setupMocks: func(parser *MockCourseParser, factory *MockExporterFactory, exporter *MockExporter) {
				parser.mockLoadCourseFromFile = func(filePath string) (*models.Course, error) {
					return testCourse, nil
				}

				factory.mockCreateExporter = func(format string) (*MockExporter, error) {
					return nil, errors.New("unsupported format")
				}
			},
			expectedError: "failed to create exporter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MockCourseParser{}
			exporter := &MockExporter{}
			factory := &MockExporterFactory{}

			tt.setupMocks(parser, factory, exporter)

			app := NewApp(parser, factory)
			err := app.ProcessCourseFromFile(tt.filePath, tt.format, tt.outputPath)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", tt.expectedError)
				}
				if !contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestApp_ProcessCourseFromURI tests the ProcessCourseFromURI method.
func TestApp_ProcessCourseFromURI(t *testing.T) {
	testCourse := createTestCourse()

	tests := []struct {
		name          string
		uri           string
		format        string
		outputPath    string
		setupMocks    func(*MockCourseParser, *MockExporterFactory, *MockExporter)
		expectedError string
	}{
		{
			name:       "successful processing",
			uri:        "https://rise.articulate.com/share/test123",
			format:     "docx",
			outputPath: "output.docx",
			setupMocks: func(parser *MockCourseParser, factory *MockExporterFactory, exporter *MockExporter) {
				parser.mockFetchCourse = func(uri string) (*models.Course, error) {
					if uri != "https://rise.articulate.com/share/test123" {
						t.Errorf("Expected uri 'https://rise.articulate.com/share/test123', got '%s'", uri)
					}
					return testCourse, nil
				}

				factory.mockCreateExporter = func(format string) (*MockExporter, error) {
					if format != "docx" {
						t.Errorf("Expected format 'docx', got '%s'", format)
					}
					return exporter, nil
				}

				exporter.mockExport = func(course *models.Course, outputPath string) error {
					if outputPath != "output.docx" {
						t.Errorf("Expected outputPath 'output.docx', got '%s'", outputPath)
					}
					return nil
				}
			},
		},
		{
			name:       "fetch error",
			uri:        "invalid-uri",
			format:     "docx",
			outputPath: "output.docx",
			setupMocks: func(parser *MockCourseParser, factory *MockExporterFactory, exporter *MockExporter) {
				parser.mockFetchCourse = func(uri string) (*models.Course, error) {
					return nil, errors.New("network error")
				}
			},
			expectedError: "failed to fetch course",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MockCourseParser{}
			exporter := &MockExporter{}
			factory := &MockExporterFactory{}

			tt.setupMocks(parser, factory, exporter)

			app := NewApp(parser, factory)
			err := app.ProcessCourseFromURI(tt.uri, tt.format, tt.outputPath)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", tt.expectedError)
				}
				if !contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestApp_GetSupportedFormats tests the GetSupportedFormats method.
func TestApp_GetSupportedFormats(t *testing.T) {
	expectedFormats := []string{"markdown", "docx", "pdf"}

	parser := &MockCourseParser{}
	factory := &MockExporterFactory{
		mockGetSupportedFormats: func() []string {
			return expectedFormats
		},
	}

	app := NewApp(parser, factory)
	formats := app.GetSupportedFormats()

	if len(formats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(formats))
	}

	for i, format := range formats {
		if format != expectedFormats[i] {
			t.Errorf("Expected format '%s' at index %d, got '%s'", expectedFormats[i], i, format)
		}
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 ||
			s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

// containsSubstring checks if s contains substr as a substring.
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
