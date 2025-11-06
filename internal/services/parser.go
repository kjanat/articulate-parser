// Package services provides the core functionality for the articulate-parser application.
// It implements the interfaces defined in the interfaces package.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/kjanat/articulate-parser/internal/interfaces"
	"github.com/kjanat/articulate-parser/internal/models"
)

// ArticulateParser implements the CourseParser interface specifically for Articulate Rise courses.
// It can fetch courses from the Articulate Rise API or load them from local JSON files.
type ArticulateParser struct {
	// BaseURL is the root URL for the Articulate Rise API
	BaseURL string
	// Client is the HTTP client used to make requests to the API
	Client *http.Client
	// Logger for structured logging
	Logger interfaces.Logger
}

// NewArticulateParser creates a new ArticulateParser instance.
// If baseURL is empty, uses the default Articulate Rise API URL.
// If timeout is zero, uses a 30-second timeout.
func NewArticulateParser(logger interfaces.Logger, baseURL string, timeout time.Duration) interfaces.CourseParser {
	if logger == nil {
		logger = NewNoOpLogger()
	}
	if baseURL == "" {
		baseURL = "https://rise.articulate.com"
	}
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &ArticulateParser{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: timeout,
		},
		Logger: logger,
	}
}

// FetchCourse fetches a course from the given URI and returns the parsed course data.
// The URI should be an Articulate Rise share URL (e.g., https://rise.articulate.com/share/SHARE_ID).
// The context can be used for cancellation and timeout control.
func (p *ArticulateParser) FetchCourse(ctx context.Context, uri string) (*models.Course, error) {
	shareID, err := p.extractShareID(uri)
	if err != nil {
		return nil, err
	}

	apiURL := p.buildAPIURL(shareID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch course data: %w", err)
	}
	// Ensure response body is closed even if ReadAll fails. Close errors are logged
	// but not fatal since the body content has already been read and parsed. In the
	// context of HTTP responses, the body must be closed to release the underlying
	// connection, but a close error doesn't invalidate the data already consumed.
	defer func() {
		if err := resp.Body.Close(); err != nil {
			p.Logger.Warn("failed to close response body", "error", err, "url", apiURL)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var course models.Course
	if err := json.Unmarshal(body, &course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &course, nil
}

// LoadCourseFromFile loads an Articulate Rise course from a local JSON file.
func (p *ArticulateParser) LoadCourseFromFile(filePath string) (*models.Course, error) {
	// #nosec G304 - File path is provided by user via CLI argument, which is expected behavior
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var course models.Course
	if err := json.Unmarshal(data, &course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &course, nil
}

// extractShareID extracts the share ID from a Rise URI.
// It uses a regular expression to find the share ID in URIs like:
// https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/
//
// Parameters:
//   - uri: The Articulate Rise share URL
//
// Returns:
//   - The share ID string if found
//   - An error if the share ID can't be extracted from the URI
func (p *ArticulateParser) extractShareID(uri string) (string, error) {
	// Parse the URL to validate the domain
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("invalid URI: %s", uri)
	}

	// Validate that it's an Articulate Rise domain
	if parsedURL.Host != "rise.articulate.com" {
		return "", fmt.Errorf("invalid domain for Articulate Rise URI: %s", parsedURL.Host)
	}

	re := regexp.MustCompile(`/share/([a-zA-Z0-9_-]+)`)
	matches := re.FindStringSubmatch(uri)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract share ID from URI: %s", uri)
	}
	return matches[1], nil
}

// buildAPIURL constructs the API URL for fetching course data.
// It combines the base URL with the API path and the share ID.
//
// Parameters:
//   - shareID: The extracted share ID from the course URI
//
// Returns:
//   - The complete API URL string for fetching the course data
func (p *ArticulateParser) buildAPIURL(shareID string) string {
	return fmt.Sprintf("%s/api/rise-runtime/boot/share/%s", p.BaseURL, shareID)
}
