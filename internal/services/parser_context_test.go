// Package services_test provides context-aware tests for the parser service.
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kjanat/articulate-parser/internal/models"
)

// TestArticulateParser_FetchCourse_ContextCancellation tests that FetchCourse
// respects context cancellation.
func TestArticulateParser_FetchCourse_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep to give time for context cancellation
		time.Sleep(100 * time.Millisecond)

		testCourse := &models.Course{
			ShareID: "test-id",
			Course: models.CourseInfo{
				Title: "Test Course",
			},
		}
		// Encode errors are ignored in test setup; httptest.ResponseWriter is reliable
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Logger: NewNoOpLogger(),
	}

	// Create a context that we'll cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/test-id")

	// Should get a context cancellation error
	if err == nil {
		t.Fatal("Expected error due to context cancellation, got nil")
	}

	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

// TestArticulateParser_FetchCourse_ContextTimeout tests that FetchCourse
// respects context timeout.
func TestArticulateParser_FetchCourse_ContextTimeout(t *testing.T) {
	// Create a server that delays response longer than timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep longer than the context timeout
		time.Sleep(200 * time.Millisecond)

		testCourse := &models.Course{
			ShareID: "test-id",
			Course: models.CourseInfo{
				Title: "Test Course",
			},
		}
		// Encode errors are ignored in test setup; httptest.ResponseWriter is reliable
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Logger: NewNoOpLogger(),
	}

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/test-id")

	// Should get a context deadline exceeded error
	if err == nil {
		t.Fatal("Expected error due to context timeout, got nil")
	}

	if !strings.Contains(err.Error(), "deadline exceeded") &&
		!strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context timeout error, got: %v", err)
	}
}

// TestArticulateParser_FetchCourse_ContextDeadline tests that FetchCourse
// respects context deadline.
func TestArticulateParser_FetchCourse_ContextDeadline(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)

		testCourse := &models.Course{
			ShareID: "test-id",
			Course: models.CourseInfo{
				Title: "Test Course",
			},
		}
		// Encode errors are ignored in test setup; httptest.ResponseWriter is reliable
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Logger: NewNoOpLogger(),
	}

	// Create a context with a deadline in the past
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Millisecond))
	defer cancel()

	_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/test-id")

	// Should get a deadline exceeded error
	if err == nil {
		t.Fatal("Expected error due to context deadline, got nil")
	}

	if !strings.Contains(err.Error(), "deadline exceeded") &&
		!strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected deadline exceeded error, got: %v", err)
	}
}

// TestArticulateParser_FetchCourse_ContextSuccess tests that FetchCourse
// succeeds when context is not cancelled.
func TestArticulateParser_FetchCourse_ContextSuccess(t *testing.T) {
	testCourse := &models.Course{
		ShareID: "test-id",
		Course: models.CourseInfo{
			Title: "Test Course",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond quickly
		// Encode errors are ignored in test setup; httptest.ResponseWriter is reliable
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Logger: NewNoOpLogger(),
	}

	// Create a context with generous timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	course, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/test-id")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if course == nil {
		t.Fatal("Expected course, got nil")
	}

	if course.Course.Title != testCourse.Course.Title {
		t.Errorf("Expected title '%s', got '%s'", testCourse.Course.Title, course.Course.Title)
	}
}

// TestArticulateParser_FetchCourse_CancellationDuringRequest tests cancellation
// during an in-flight request.
func TestArticulateParser_FetchCourse_CancellationDuringRequest(t *testing.T) {
	requestStarted := make(chan bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStarted <- true
		// Keep the handler running to simulate slow response
		time.Sleep(300 * time.Millisecond)

		testCourse := &models.Course{
			ShareID: "test-id",
		}
		// Encode errors are ignored in test setup; httptest.ResponseWriter is reliable
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Logger: NewNoOpLogger(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start the request in a goroutine
	errChan := make(chan error, 1)
	go func() {
		_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/test-id")
		errChan <- err
	}()

	// Wait for request to start
	<-requestStarted

	// Cancel after request has started
	cancel()

	// Get the error
	err := <-errChan

	if err == nil {
		t.Fatal("Expected error due to context cancellation, got nil")
	}

	// Should contain context canceled somewhere in the error chain
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context canceled error, got: %v", err)
	}
}

// TestArticulateParser_FetchCourse_MultipleTimeouts tests behavior with
// multiple concurrent requests and timeouts.
func TestArticulateParser_FetchCourse_MultipleTimeouts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		testCourse := &models.Course{ShareID: "test"}
		// Encode errors are ignored in test setup; httptest.ResponseWriter is reliable
		_ = json.NewEncoder(w).Encode(testCourse)
	}))
	defer server.Close()

	parser := &ArticulateParser{
		BaseURL: server.URL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Logger: NewNoOpLogger(),
	}

	// Launch multiple requests with different timeouts
	tests := []struct {
		name          string
		timeout       time.Duration
		shouldSucceed bool
	}{
		{"very short timeout", 10 * time.Millisecond, false},
		{"short timeout", 50 * time.Millisecond, false},
		{"adequate timeout", 500 * time.Millisecond, true},
		{"long timeout", 2 * time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			_, err := parser.FetchCourse(ctx, "https://rise.articulate.com/share/test-id")

			if tt.shouldSucceed && err != nil {
				t.Errorf("Expected success with timeout %v, got error: %v", tt.timeout, err)
			}

			if !tt.shouldSucceed && err == nil {
				t.Errorf("Expected timeout error with timeout %v, got success", tt.timeout)
			}
		})
	}
}
