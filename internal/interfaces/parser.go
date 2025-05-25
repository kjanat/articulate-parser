// Package interfaces provides the core contracts for the articulate-parser application.
// It defines interfaces for parsing and exporting Articulate Rise courses.
package interfaces

import "github.com/kjanat/articulate-parser/internal/models"

// CourseParser defines the interface for loading course data.
// It provides methods to fetch course content either from a remote URI
// or from a local file path.
type CourseParser interface {
	// FetchCourse loads a course from a URI (typically an Articulate Rise share URL).
	// It retrieves the course data from the remote location and returns a parsed Course model.
	// Returns an error if the fetch operation fails or if the data cannot be parsed.
	FetchCourse(uri string) (*models.Course, error)

	// LoadCourseFromFile loads a course from a local file.
	// It reads and parses the course data from the specified file path.
	// Returns an error if the file cannot be read or if the data cannot be parsed.
	LoadCourseFromFile(filePath string) (*models.Course, error)
}
