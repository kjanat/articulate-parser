// Package services provides the core functionality for the articulate-parser application.
// It implements the interfaces defined in the interfaces package.
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

// NewArticulateParser creates a new ArticulateParser instance with default settings.
// The default configuration uses the standard Articulate Rise API URL and a
// HTTP client with a 30-second timeout.
func NewArticulateParser() interfaces.CourseParser {
	return &ArticulateParser{
		BaseURL: "https://rise.articulate.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchCourse fetches a course from the given URI.
// It extracts the share ID from the URI, constructs an API URL, and fetches the course data.
// The course data is then unmarshalled into a Course model.
//
// Parameters:
//   - uri: The Articulate Rise share URL (e.g., https://rise.articulate.com/share/SHARE_ID)
//
// Returns:
//   - A parsed Course model if successful
//   - An error if the fetch fails, if the share ID can't be extracted,
//     or if the response can't be parsed
func (p *ArticulateParser) FetchCourse(uri string) (*models.Course, error) {
	shareID, err := p.extractShareID(uri)
	if err != nil {
		return nil, err
	}

	apiURL := p.buildAPIURL(shareID)

	resp, err := p.Client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch course data: %w", err)
	}
	defer resp.Body.Close()

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
// The file should contain a valid JSON representation of an Articulate Rise course.
//
// Parameters:
//   - filePath: The path to the JSON file containing the course data
//
// Returns:
//   - A parsed Course model if successful
//   - An error if the file can't be read or the JSON can't be parsed
func (p *ArticulateParser) LoadCourseFromFile(filePath string) (*models.Course, error) {
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
