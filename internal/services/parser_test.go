package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kjanat/articulate-parser/internal/models"
)

// TestNewArticulateParser tests the NewArticulateParser constructor.
func TestNewArticulateParser(t *testing.T) {
	parser := NewArticulateParser(nil, "", 0)

	if parser == nil {
		t.Fatal("NewArticulateParser() returned nil")
	}

	// Type assertion to check internal structure
	articulateParser, ok := parser.(*ArticulateParser)
	if !ok {
		t.Fatal("NewArticulateParser() returned wrong type")
	}

	expectedBaseURL := "https://rise.articulate.com"
	if articulateParser.BaseURL != expectedBaseURL {
		t.Errorf("Expected BaseURL '%s', got '%s'", expectedBaseURL, articulateParser.BaseURL)
	}

	if articulateParser.Client == nil {
		t.Error("Client should not be nil")
	}

	expectedTimeout := 30 * time.Second
	if articulateParser.Client.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, articulateParser.Client.Timeout)
	}
}

// TestArticulateParser_FetchCourse tests the FetchCourse method.
func TestArticulateParser_FetchCourse(t *testing.T) {
	// Create a test course object
	testCourse := &models.Course{
		ShareID: "test-share-id",
		Author:  "Test Author",
		Course: models.CourseInfo{
			ID:          "test-course-id",
			Title:       "Test Course",
			Description: "Test Description",
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		expectedPath := "/api/rise-runtime/boot/share/test-share-id"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Check request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(testCourse); err != nil {
			t.Fatalf("Failed to encode test course: %v", err)
		}
	}))
	defer server.Close()

	// Create parser with test server URL
	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	tests := []struct {
		name          string
		uri           string
		expectedError string
	}{
		{
			name: "valid articulate rise URI",
			uri:  "https://rise.articulate.com/share/test-share-id#/",
		},
		{
			name: "valid articulate rise URI without fragment",
			uri:  "https://rise.articulate.com/share/test-share-id",
		},
		{
			name:          "invalid URI format",
			uri:           "invalid-uri",
			expectedError: "invalid domain for Articulate Rise URI:",
		},
		{
			name:          "empty URI",
			uri:           "",
			expectedError: "invalid domain for Articulate Rise URI:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			course, err := parser.FetchCourse(context.Background(), tt.uri)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}
				if course == nil {
					t.Fatal("Expected course, got nil")
				}
				if course.ShareID != testCourse.ShareID {
					t.Errorf("Expected ShareID '%s', got '%s'", testCourse.ShareID, course.ShareID)
				}
			}
		})
	}
}

// TestArticulateParser_FetchCourse_NetworkError tests network error handling.
func TestArticulateParser_FetchCourse_NetworkError(t *testing.T) {
	// Create parser with invalid URL to simulate network error
	parser := &ArticulateParser{
		BaseURL: "http://localhost:99999", // Invalid port
		Client: &http.Client{
			Timeout: 1 * time.Millisecond, // Very short timeout
		},
	}

	_, err := parser.FetchCourse(context.Background(), "https://rise.articulate.com/share/test-share-id")
	if err == nil {
		t.Fatal("Expected network error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to fetch course data") {
		t.Errorf("Expected error to contain 'failed to fetch course data', got '%s'", err.Error())
	}
}

// TestArticulateParser_FetchCourse_InvalidJSON tests invalid JSON response handling.
func TestArticulateParser_FetchCourse_InvalidJSON(t *testing.T) {
	// Create test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Write is used for its side effect; the test verifies error handling on
		// the client side, not whether the write succeeds. Ignore the error since
		// httptest.ResponseWriter writes are rarely problematic in test contexts.
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	_, err := parser.FetchCourse(context.Background(), "https://rise.articulate.com/share/test-share-id")
	if err == nil {
		t.Fatal("Expected JSON parsing error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal JSON") {
		t.Errorf("Expected error to contain 'failed to unmarshal JSON', got '%s'", err.Error())
	}
}

// TestArticulateParser_LoadCourseFromFile tests the LoadCourseFromFile method.
func TestArticulateParser_LoadCourseFromFile(t *testing.T) {
	// Create a temporary test file
	testCourse := &models.Course{
		ShareID: "file-test-share-id",
		Author:  "File Test Author",
		Course: models.CourseInfo{
			ID:          "file-test-course-id",
			Title:       "File Test Course",
			Description: "File Test Description",
		},
	}

	// Create temporary directory and file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-course.json")

	// Write test data to file
	data, err := json.Marshal(testCourse)
	if err != nil {
		t.Fatalf("Failed to marshal test course: %v", err)
	}

	if err := os.WriteFile(tempFile, data, 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewArticulateParser(nil, "", 0)

	tests := []struct {
		name          string
		filePath      string
		expectedError string
	}{
		{
			name:     "valid file",
			filePath: tempFile,
		},
		{
			name:          "nonexistent file",
			filePath:      filepath.Join(tempDir, "nonexistent.json"),
			expectedError: "failed to read file",
		},
		{
			name:          "empty path",
			filePath:      "",
			expectedError: "failed to read file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			course, err := parser.LoadCourseFromFile(tt.filePath)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}
				if course == nil {
					t.Fatal("Expected course, got nil")
				}
				if course.ShareID != testCourse.ShareID {
					t.Errorf("Expected ShareID '%s', got '%s'", testCourse.ShareID, course.ShareID)
				}
			}
		})
	}
}

// TestArticulateParser_LoadCourseFromFile_InvalidJSON tests invalid JSON file handling.
func TestArticulateParser_LoadCourseFromFile_InvalidJSON(t *testing.T) {
	// Create temporary file with invalid JSON
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "invalid.json")

	if err := os.WriteFile(tempFile, []byte("invalid json content"), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewArticulateParser(nil, "", 0)
	_, err := parser.LoadCourseFromFile(tempFile)

	if err == nil {
		t.Fatal("Expected JSON parsing error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal JSON") {
		t.Errorf("Expected error to contain 'failed to unmarshal JSON', got '%s'", err.Error())
	}
}

// TestExtractShareID tests the extractShareID method.
func TestExtractShareID(t *testing.T) {
	parser := &ArticulateParser{}

	tests := []struct {
		name     string
		uri      string
		expected string
		hasError bool
	}{
		{
			name:     "standard articulate rise URI with fragment",
			uri:      "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/",
			expected: "N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO",
		},
		{
			name:     "standard articulate rise URI without fragment",
			uri:      "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO",
			expected: "N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO",
		},
		{
			name:     "URI with trailing slash",
			uri:      "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO/",
			expected: "N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO",
		},
		{
			name:     "short share ID",
			uri:      "https://rise.articulate.com/share/abc123",
			expected: "abc123",
		},
		{
			name:     "share ID with hyphens and underscores",
			uri:      "https://rise.articulate.com/share/test_ID-123_abc",
			expected: "test_ID-123_abc",
		},
		{
			name:     "invalid URI - no share path",
			uri:      "https://rise.articulate.com/",
			hasError: true,
		},
		{
			name:     "invalid URI - wrong domain",
			uri:      "https://example.com/share/test123",
			hasError: true,
		},
		{
			name:     "invalid URI - no share ID",
			uri:      "https://rise.articulate.com/share/",
			hasError: true,
		},
		{
			name:     "empty URI",
			uri:      "",
			hasError: true,
		},
		{
			name:     "malformed URI",
			uri:      "not-a-uri",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.extractShareID(tt.uri)

			if tt.hasError {
				if err == nil {
					t.Fatalf("Expected error for URI '%s', got nil", tt.uri)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error for URI '%s', got: %v", tt.uri, err)
				}
				if result != tt.expected {
					t.Errorf("Expected share ID '%s', got '%s'", tt.expected, result)
				}
			}
		})
	}
}

// TestBuildAPIURL tests the buildAPIURL method.
func TestBuildAPIURL(t *testing.T) {
	parser := &ArticulateParser{
		BaseURL: "https://rise.articulate.com",
	}

	tests := []struct {
		name     string
		shareID  string
		expected string
	}{
		{
			name:     "standard share ID",
			shareID:  "N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO",
			expected: "https://rise.articulate.com/api/rise-runtime/boot/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO",
		},
		{
			name:     "short share ID",
			shareID:  "abc123",
			expected: "https://rise.articulate.com/api/rise-runtime/boot/share/abc123",
		},
		{
			name:     "share ID with special characters",
			shareID:  "test_ID-123_abc",
			expected: "https://rise.articulate.com/api/rise-runtime/boot/share/test_ID-123_abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.buildAPIURL(tt.shareID)
			if result != tt.expected {
				t.Errorf("Expected URL '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestBuildAPIURL_DifferentBaseURL tests buildAPIURL with different base URLs.
func TestBuildAPIURL_DifferentBaseURL(t *testing.T) {
	parser := &ArticulateParser{
		BaseURL: "https://custom.domain.com",
	}

	shareID := "test123"
	expected := "https://custom.domain.com/api/rise-runtime/boot/share/test123"
	result := parser.buildAPIURL(shareID)

	if result != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, result)
	}
}

// BenchmarkExtractShareID benchmarks the extractShareID method.
func BenchmarkExtractShareID(b *testing.B) {
	parser := &ArticulateParser{}
	uri := "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/"

	for b.Loop() {
		_, _ = parser.extractShareID(uri)
	}
}

// BenchmarkBuildAPIURL benchmarks the buildAPIURL method.
func BenchmarkBuildAPIURL(b *testing.B) {
	parser := &ArticulateParser{
		BaseURL: "https://rise.articulate.com",
	}
	shareID := "N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO"

	for b.Loop() {
		_ = parser.buildAPIURL(shareID)
	}
}
