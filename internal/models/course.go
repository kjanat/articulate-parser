// Package models defines the data structures representing Articulate Rise courses.
// These structures closely match the JSON format used by Articulate Rise.
package models

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
type LabelSet struct {
	// ID is the unique identifier for this label set
	ID string `json:"id"`
	// Name is the descriptive name of the label set
	Name string `json:"name"`
	// Labels is a mapping of label keys to their customized values
	Labels map[string]string `json:"labels"`
}
