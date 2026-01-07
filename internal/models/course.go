// Package models defines the data structures representing Articulate Rise courses.
// These structures closely match the JSON format used by Articulate Rise.
package models

import "fmt"

// Course represents the top-level structure of an Articulate Rise course.
// It contains metadata and the actual course content.
type Course struct {
	// ShareID is the unique identifier used in public sharing URLs
	ShareID string `json:"shareId"`
	// Author is the name of the course creator
	Author string `json:"author"`
	// Course contains the detailed course information and content
	Course CourseInfo `json:"course"`
	// LabelSet contains customized labels used in the course
	LabelSet LabelSet `json:"labelSet"`
}

// CourseInfo contains the main details and content of an Articulate Rise course.
type CourseInfo struct {
	// ID is the internal unique identifier for the course
	ID string `json:"id"`
	// Title is the name of the course
	Title string `json:"title"`
	// Description is the course summary or introduction text
	Description string `json:"description"`
	// Color is the theme color of the course
	Color string `json:"color"`
	// NavigationMode specifies how users navigate through the course
	NavigationMode string `json:"navigationMode"`
	// Lessons is an ordered array of all lessons in the course
	Lessons []Lesson `json:"lessons"`
	// CoverImage is the main image displayed for the course
	CoverImage *Media `json:"coverImage,omitempty"`
	// ExportSettings contains configuration for exporting the course
	ExportSettings *ExportSettings `json:"exportSettings,omitempty"`
}

// ExportSettings defines configuration options for exporting a course.
type ExportSettings struct {
	// Title specifies the export title which might differ from course title
	Title string `json:"title"`
	// Format indicates the preferred export format
	Format string `json:"format"`
}

// LabelSet contains customized labels used throughout the course.
// This allows course creators to modify standard terminology.
//
// NOTE: Currently parsed but not used in export logic.
type LabelSet struct {
	// ID is the unique identifier for this label set
	ID string `json:"id"`
	// Name is the descriptive name of the label set
	Name string `json:"name"`
	// Labels is a mapping of label keys to their customized values
	Labels map[string]string `json:"labels"`
}

// Validate checks whether the Course structure contains valid data.
// It returns an error describing the first validation failure found,
// or nil if the course is valid.
//
// Validation rules:
//   - ShareID must not be empty
//   - Course.ID must not be empty
//   - Course.Title must not be empty
//   - Each Lesson must have a non-empty ID and Title
func (c *Course) Validate() error {
	if c.ShareID == "" {
		return &ValidationError{Field: "ShareID", Message: "share ID is required"}
	}
	if c.Course.ID == "" {
		return &ValidationError{Field: "Course.ID", Message: "course ID is required"}
	}
	if c.Course.Title == "" {
		return &ValidationError{Field: "Course.Title", Message: "course title is required"}
	}
	for i, lesson := range c.Course.Lessons {
		if lesson.ID == "" {
			return &ValidationError{Field: "Course.Lessons", Message: fmt.Sprintf("lesson at index %d has empty ID", i)}
		}
		if lesson.Title == "" {
			return &ValidationError{Field: "Course.Lessons", Message: fmt.Sprintf("lesson at index %d has empty title", i)}
		}
	}
	return nil
}

// ValidationError represents a validation failure for a specific field.
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string
	// Message describes the validation failure
	Message string
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
