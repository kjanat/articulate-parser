// Package testdata provides shared test fixtures for Articulate Rise course testing.
// It contains builders for creating Course, Lesson, and Item structures with various
// configurations to support unit tests, integration tests, and benchmarks.
package testdata

import (
	"fmt"

	"github.com/kjanat/articulate-parser/internal/models"
)

// SampleAnswers provides a reusable set of quiz answers for testing.
var SampleAnswers = []models.Answer{
	{ID: "ans-1", Title: "Option A", Correct: false},
	{ID: "ans-2", Title: "Option B", Correct: true},
	{ID: "ans-3", Title: "Option C", Correct: false},
	{ID: "ans-4", Title: "Option D", Correct: false},
}

// SampleSubItems provides a reusable set of sub-items for testing.
var SampleSubItems = []models.SubItem{
	{
		ID:        "sub-1",
		Heading:   "Sample Heading",
		Paragraph: "<p>Sample paragraph with <strong>formatting</strong>.</p>",
	},
	{
		ID:        "sub-2",
		Paragraph: "<p>Another paragraph with <em>emphasis</em>.</p>",
	},
}

// MinimalCourse returns a course with just ShareID and Title set.
// Useful for testing edge cases with minimal data.
func MinimalCourse() *models.Course {
	return &models.Course{
		ShareID: "minimal-share-id",
		Course: models.CourseInfo{
			Title: "Minimal Course",
		},
	}
}

// BasicCourse returns a course with 1 section, 1 lesson, and 1 text item.
// Suitable for basic functionality tests.
func BasicCourse() *models.Course {
	return &models.Course{
		ShareID: "basic-share-id",
		Author:  "Test Author",
		Course: models.CourseInfo{
			ID:             "basic-course-id",
			Title:          "Basic Test Course",
			Description:    "<p>A basic course for testing.</p>",
			NavigationMode: "menu",
			Lessons: []models.Lesson{
				{
					ID:    "section-1",
					Title: "Introduction Section",
					Type:  "section",
				},
				{
					ID:          "lesson-1",
					Title:       "First Lesson",
					Type:        "lesson",
					Description: "<p>Lesson introduction text.</p>",
					Items: []models.Item{
						TextItem("Welcome", "<p>Welcome to the course!</p>"),
					},
				},
			},
		},
	}
}

// CompleteCourse returns a course with all item types represented.
// Useful for comprehensive export and parsing tests.
func CompleteCourse() *models.Course {
	return &models.Course{
		ShareID: "complete-share-id",
		Author:  "Complete Test Author",
		Course: models.CourseInfo{
			ID:             "complete-course-id",
			Title:          "Complete Test Course",
			Description:    "<p>A comprehensive course with <strong>all</strong> item types.</p>",
			Color:          "#3366CC",
			NavigationMode: "menu",
			ExportSettings: &models.ExportSettings{
				Title:  "Exported Complete Course",
				Format: "scorm",
			},
			CoverImage: &models.Media{
				Image: &models.ImageMedia{
					Key:         "cover-img-key",
					Type:        "jpg",
					Width:       1920,
					Height:      1080,
					OriginalURL: "https://example.com/cover.jpg",
				},
			},
			Lessons: []models.Lesson{
				{
					ID:    "section-1",
					Title: "Course Overview",
					Type:  "section",
				},
				{
					ID:          "lesson-1",
					Title:       "Introduction",
					Type:        "lesson",
					Description: "<p>Introduction to the course content.</p>",
					Icon:        "book",
					Items: []models.Item{
						TextItem("Welcome", "<p>Welcome to this comprehensive course!</p>"),
						TextItem("Objectives", "<p>In this course you will learn:</p>"),
						ListItem("First objective", "Second objective", "Third objective"),
						DividerItem(),
						MediaItem("video", "https://example.com/intro.mp4"),
						MediaItem("image", "https://example.com/diagram.png"),
					},
				},
				{
					ID:          "lesson-2",
					Title:       "Core Concepts",
					Type:        "lesson",
					Description: "<p>Core concepts and fundamentals.</p>",
					Items: []models.Item{
						TextItem("Key Concepts", "<p>Let's explore the key concepts.</p>"),
						KnowledgeCheckItem("What is the primary purpose of this course?", []models.Answer{
							{ID: "kc-1-a", Title: "Learning basics", Correct: false},
							{ID: "kc-1-b", Title: "Comprehensive understanding", Correct: true},
							{ID: "kc-1-c", Title: "Quick overview", Correct: false},
						}),
						InteractiveItem("Interactive Exercise"),
					},
				},
				{
					ID:          "lesson-3",
					Title:       "Advanced Topics",
					Type:        "lesson",
					Description: "<p>Deep dive into advanced topics.</p>",
					Items: []models.Item{
						TextItem("Advanced Section", "<p>This section covers advanced material.</p>"),
						FlipCardItem("Front content", "Back content"),
						QuoteItem("Learning is not attained by chance.", "Abigail Adams"),
					},
				},
				{
					ID:    "section-2",
					Title: "Summary",
					Type:  "section",
				},
				{
					ID:          "lesson-4",
					Title:       "Conclusion",
					Type:        "lesson",
					Description: "<p>Course wrap-up and next steps.</p>",
					Items: []models.Item{
						TextItem("Summary", "<p>Thank you for completing this course!</p>"),
						ListItem("Review materials regularly", "Practice learned skills", "Share knowledge"),
					},
				},
			},
		},
		LabelSet: models.LabelSet{
			ID:   "default-labels",
			Name: "Default",
			Labels: map[string]string{
				"next":     "Next",
				"previous": "Previous",
				"complete": "Mark Complete",
			},
		},
	}
}

// LargeCourse returns a parameterized course for benchmarking with the specified number of lessons.
// Each lesson contains multiple item types for realistic performance testing.
func LargeCourse(lessons int) *models.Course {
	courseLessons := make([]models.Lesson, 0, lessons+lessons/5) // +sections

	for i := range lessons {
		// Add section header every 5 lessons
		if i%5 == 0 {
			courseLessons = append(courseLessons, models.Lesson{
				ID:    fmt.Sprintf("section-%d", i/5+1),
				Title: fmt.Sprintf("Section %d", i/5+1),
				Type:  "section",
			})
		}

		lesson := models.Lesson{
			ID:          fmt.Sprintf("lesson-%d", i+1),
			Title:       fmt.Sprintf("Lesson %d", i+1),
			Type:        "lesson",
			Description: fmt.Sprintf("<p>Description for lesson %d with <em>formatting</em>.</p>", i+1),
			Items: []models.Item{
				TextItem(
					fmt.Sprintf("Heading %d", i+1),
					fmt.Sprintf("<p>Content for lesson %d with <strong>bold</strong> and <em>italic</em> text.</p>", i+1),
				),
				ListItem("Point 1", "Point 2", "Point 3"),
				KnowledgeCheckItem(
					fmt.Sprintf("Quiz Question %d", i+1),
					[]models.Answer{
						{ID: fmt.Sprintf("q%d-a", i), Title: "Answer A", Correct: false},
						{ID: fmt.Sprintf("q%d-b", i), Title: "Answer B", Correct: true},
						{ID: fmt.Sprintf("q%d-c", i), Title: "Answer C", Correct: false},
					},
				),
			},
		}
		courseLessons = append(courseLessons, lesson)
	}

	return &models.Course{
		ShareID: "large-benchmark-id",
		Author:  "Benchmark Author",
		Course: models.CourseInfo{
			ID:             "large-benchmark-course",
			Title:          fmt.Sprintf("Large Benchmark Course (%d lessons)", lessons),
			Description:    "<p>Large course for performance testing.</p>",
			NavigationMode: "menu",
			Lessons:        courseLessons,
		},
	}
}

// TextItem creates a text item with the specified heading and paragraph content.
func TextItem(heading, paragraph string) models.Item {
	return models.Item{
		Type:   "text",
		Family: "text",
		Items: []models.SubItem{
			{
				ID:        "text-sub-1",
				Heading:   heading,
				Paragraph: paragraph,
			},
		},
	}
}

// ListItem creates a list item with the specified items as list entries.
func ListItem(items ...string) models.Item {
	subItems := make([]models.SubItem, len(items))
	for i, item := range items {
		subItems[i] = models.SubItem{
			ID:        fmt.Sprintf("list-sub-%d", i+1),
			Paragraph: fmt.Sprintf("<p>%s</p>", item),
		}
	}

	return models.Item{
		Type:   "list",
		Family: "list",
		Items:  subItems,
	}
}

// KnowledgeCheckItem creates a knowledge check item with the specified question and answers.
func KnowledgeCheckItem(question string, answers []models.Answer) models.Item {
	return models.Item{
		Type:   "knowledgeCheck",
		Family: "knowledge",
		Items: []models.SubItem{
			{
				ID:       "kc-sub-1",
				Title:    fmt.Sprintf("<p>%s</p>", question),
				Answers:  answers,
				Feedback: "<p>Thank you for your response.</p>",
			},
		},
	}
}

// MediaItem creates a media item (video or image) with the specified type and URL.
// The mediaType should be "video" or "image".
func MediaItem(mediaType, url string) models.Item {
	var media *models.Media

	switch mediaType {
	case "video":
		media = &models.Media{
			Video: &models.VideoMedia{
				Key:         "video-key",
				URL:         url,
				Type:        "mp4",
				Duration:    120,
				OriginalURL: url,
			},
		}
	case "image":
		media = &models.Media{
			Image: &models.ImageMedia{
				Key:         "image-key",
				Type:        "png",
				Width:       800,
				Height:      600,
				OriginalURL: url,
			},
		}
	default:
		media = &models.Media{}
	}

	return models.Item{
		Type:   "multimedia",
		Family: "media",
		Items: []models.SubItem{
			{
				ID:      "media-sub-1",
				Title:   fmt.Sprintf("<p>%s content</p>", mediaType),
				Caption: fmt.Sprintf("<p>Caption for %s</p>", mediaType),
				Media:   media,
			},
		},
	}
}

// InteractiveItem creates an interactive content item with the specified title.
func InteractiveItem(title string) models.Item {
	return models.Item{
		Type:   "interactive",
		Family: "interactive",
		Items: []models.SubItem{
			{
				ID:    "interactive-sub-1",
				Title: fmt.Sprintf("<p>%s</p>", title),
			},
		},
	}
}

// DividerItem creates a divider/separator item.
func DividerItem() models.Item {
	return models.Item{
		Type:   "divider",
		Family: "divider",
	}
}

// FlipCardItem creates a flashcard-style item with front and back content.
func FlipCardItem(front, back string) models.Item {
	return models.Item{
		Type:   "flashcard",
		Family: "interactive",
		Items: []models.SubItem{
			{
				ID: "flipcard-sub-1",
				Front: &models.CardSide{
					Description: fmt.Sprintf("<p>%s</p>", front),
				},
				Back: &models.CardSide{
					Description: fmt.Sprintf("<p>%s</p>", back),
				},
			},
		},
	}
}

// QuoteItem creates a quote/blockquote item with text and attribution.
func QuoteItem(quote, author string) models.Item {
	return models.Item{
		Type:   "quote",
		Family: "text",
		Items: []models.SubItem{
			{
				ID:        "quote-sub-1",
				Paragraph: fmt.Sprintf("<blockquote>%s</blockquote>", quote),
				Title:     author,
			},
		},
	}
}

// ImageItem creates a standalone image item with caption.
func ImageItem(url, caption string) models.Item {
	return models.Item{
		Type:   "image",
		Family: "media",
		Items: []models.SubItem{
			{
				ID:      "image-sub-1",
				Caption: fmt.Sprintf("<p>%s</p>", caption),
				Media: &models.Media{
					Image: &models.ImageMedia{
						Key:         "standalone-img-key",
						Type:        "jpg",
						Width:       1200,
						Height:      800,
						OriginalURL: url,
					},
				},
			},
		},
	}
}

// VideoItem creates a standalone video item with caption.
func VideoItem(url string, duration int, caption string) models.Item {
	return models.Item{
		Type:   "video",
		Family: "media",
		Items: []models.SubItem{
			{
				ID:      "video-sub-1",
				Caption: fmt.Sprintf("<p>%s</p>", caption),
				Media: &models.Media{
					Video: &models.VideoMedia{
						Key:         "standalone-video-key",
						URL:         url,
						Type:        "mp4",
						Duration:    duration,
						OriginalURL: url,
					},
				},
			},
		},
	}
}

// MatchingItem creates a matching/pairing knowledge check item.
func MatchingItem(pairs map[string]string) models.Item {
	answers := make([]models.Answer, 0, len(pairs))
	i := 0
	for term, match := range pairs {
		answers = append(answers, models.Answer{
			ID:         fmt.Sprintf("match-%d", i),
			Title:      term,
			MatchTitle: match,
			Correct:    true,
		})
		i++
	}

	return models.Item{
		Type:   "matching",
		Family: "knowledge",
		Items: []models.SubItem{
			{
				ID:       "matching-sub-1",
				Title:    "<p>Match the following items:</p>",
				Answers:  answers,
				Feedback: "<p>Check your matches!</p>",
			},
		},
	}
}

// EmptyLesson returns a lesson with no items for edge case testing.
func EmptyLesson(id, title string) models.Lesson {
	return models.Lesson{
		ID:    id,
		Title: title,
		Type:  "lesson",
		Items: []models.Item{},
	}
}

// SectionLesson returns a section header (not a content lesson).
func SectionLesson(id, title string) models.Lesson {
	return models.Lesson{
		ID:    id,
		Title: title,
		Type:  "section",
	}
}
