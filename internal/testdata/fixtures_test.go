package testdata

import (
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
)

func TestMinimalCourse(t *testing.T) {
	course := MinimalCourse()
	if course == nil {
		t.Fatal("MinimalCourse() returned nil")
	}
	if course.ShareID == "" {
		t.Error("MinimalCourse() should have ShareID set")
	}
	if course.Course.Title == "" {
		t.Error("MinimalCourse() should have Title set")
	}
}

func TestBasicCourse(t *testing.T) {
	course := BasicCourse()
	if course == nil {
		t.Fatal("BasicCourse() returned nil")
	}
	if len(course.Course.Lessons) != 2 {
		t.Errorf("BasicCourse() should have 2 lessons, got %d", len(course.Course.Lessons))
	}
	if course.Course.Lessons[0].Type != "section" {
		t.Error("First lesson should be a section")
	}
	if course.Course.Lessons[1].Type != "lesson" {
		t.Error("Second lesson should be a lesson")
	}
	if len(course.Course.Lessons[1].Items) != 1 {
		t.Errorf("Lesson should have 1 item, got %d", len(course.Course.Lessons[1].Items))
	}
}

func TestCompleteCourse(t *testing.T) {
	course := CompleteCourse()
	if course == nil {
		t.Fatal("CompleteCourse() returned nil")
	}
	if len(course.Course.Lessons) == 0 {
		t.Error("CompleteCourse() should have lessons")
	}

	// Verify all expected item types exist
	itemTypes := make(map[string]bool)
	for _, lesson := range course.Course.Lessons {
		for _, item := range lesson.Items {
			itemTypes[item.Type] = true
		}
	}

	expectedTypes := []string{"text", "list", "divider", "multimedia", "knowledgeCheck", "interactive", "flashcard", "quote"}
	for _, et := range expectedTypes {
		if !itemTypes[et] {
			t.Errorf("CompleteCourse() should have item type %q", et)
		}
	}
}

func TestLargeCourse(t *testing.T) {
	tests := []struct {
		lessons  int
		expected int
	}{
		{1, 2},   // 1 lesson + 1 section
		{5, 6},   // 5 lessons + 1 section
		{10, 12}, // 10 lessons + 2 sections
		{50, 60}, // 50 lessons + 10 sections
	}

	for _, tt := range tests {
		course := LargeCourse(tt.lessons)
		if course == nil {
			t.Fatalf("LargeCourse(%d) returned nil", tt.lessons)
		}
		if len(course.Course.Lessons) != tt.expected {
			t.Errorf("LargeCourse(%d) should have %d lessons+sections, got %d",
				tt.lessons, tt.expected, len(course.Course.Lessons))
		}
	}
}

func TestTextItem(t *testing.T) {
	item := TextItem("Test Heading", "<p>Test paragraph</p>")
	if item.Type != "text" {
		t.Errorf("TextItem Type should be 'text', got %q", item.Type)
	}
	if len(item.Items) != 1 {
		t.Fatalf("TextItem should have 1 sub-item, got %d", len(item.Items))
	}
	if item.Items[0].Heading != "Test Heading" {
		t.Errorf("TextItem heading mismatch: got %q", item.Items[0].Heading)
	}
}

func TestListItem(t *testing.T) {
	item := ListItem("Item 1", "Item 2", "Item 3")
	if item.Type != "list" {
		t.Errorf("ListItem Type should be 'list', got %q", item.Type)
	}
	if len(item.Items) != 3 {
		t.Errorf("ListItem should have 3 sub-items, got %d", len(item.Items))
	}
}

func TestKnowledgeCheckItem(t *testing.T) {
	answers := []models.Answer{
		{Title: "A", Correct: false},
		{Title: "B", Correct: true},
	}
	item := KnowledgeCheckItem("Test question?", answers)
	if item.Type != "knowledgeCheck" {
		t.Errorf("KnowledgeCheckItem Type should be 'knowledgeCheck', got %q", item.Type)
	}
	if len(item.Items) != 1 {
		t.Fatalf("KnowledgeCheckItem should have 1 sub-item, got %d", len(item.Items))
	}
	if len(item.Items[0].Answers) != 2 {
		t.Errorf("KnowledgeCheckItem should have 2 answers, got %d", len(item.Items[0].Answers))
	}
}

func TestMediaItem(t *testing.T) {
	tests := []struct {
		name      string
		mediaType string
		hasVideo  bool
		hasImage  bool
	}{
		{"video", "video", true, false},
		{"image", "image", false, true},
		{"unknown", "unknown", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := MediaItem(tt.mediaType, "https://example.com/media")
			if item.Type != "multimedia" {
				t.Errorf("MediaItem Type should be 'multimedia', got %q", item.Type)
			}
			if len(item.Items) != 1 {
				t.Fatalf("MediaItem should have 1 sub-item, got %d", len(item.Items))
			}
			if tt.hasVideo && item.Items[0].Media.Video == nil {
				t.Error("Expected video media to be set")
			}
			if tt.hasImage && item.Items[0].Media.Image == nil {
				t.Error("Expected image media to be set")
			}
		})
	}
}

func TestInteractiveItem(t *testing.T) {
	item := InteractiveItem("Interactive Title")
	if item.Type != "interactive" {
		t.Errorf("InteractiveItem Type should be 'interactive', got %q", item.Type)
	}
}

func TestDividerItem(t *testing.T) {
	item := DividerItem()
	if item.Type != "divider" {
		t.Errorf("DividerItem Type should be 'divider', got %q", item.Type)
	}
}

func TestFlipCardItem(t *testing.T) {
	item := FlipCardItem("Front text", "Back text")
	if item.Type != "flashcard" {
		t.Errorf("FlipCardItem Type should be 'flashcard', got %q", item.Type)
	}
	if len(item.Items) != 1 {
		t.Fatalf("FlipCardItem should have 1 sub-item, got %d", len(item.Items))
	}
	if item.Items[0].Front == nil || item.Items[0].Back == nil {
		t.Error("FlipCardItem should have Front and Back set")
	}
}

func TestQuoteItem(t *testing.T) {
	item := QuoteItem("Test quote", "Test author")
	if item.Type != "quote" {
		t.Errorf("QuoteItem Type should be 'quote', got %q", item.Type)
	}
}

func TestImageItem(t *testing.T) {
	item := ImageItem("https://example.com/img.jpg", "Image caption")
	if item.Type != "image" {
		t.Errorf("ImageItem Type should be 'image', got %q", item.Type)
	}
	if len(item.Items) != 1 || item.Items[0].Media == nil || item.Items[0].Media.Image == nil {
		t.Error("ImageItem should have image media set")
	}
}

func TestVideoItem(t *testing.T) {
	item := VideoItem("https://example.com/video.mp4", 120, "Video caption")
	if item.Type != "video" {
		t.Errorf("VideoItem Type should be 'video', got %q", item.Type)
	}
	if len(item.Items) != 1 || item.Items[0].Media == nil || item.Items[0].Media.Video == nil {
		t.Error("VideoItem should have video media set")
	}
	if item.Items[0].Media.Video.Duration != 120 {
		t.Errorf("VideoItem duration should be 120, got %d", item.Items[0].Media.Video.Duration)
	}
}

func TestMatchingItem(t *testing.T) {
	pairs := map[string]string{
		"Term 1": "Match 1",
		"Term 2": "Match 2",
	}
	item := MatchingItem(pairs)
	if item.Type != "matching" {
		t.Errorf("MatchingItem Type should be 'matching', got %q", item.Type)
	}
	if len(item.Items) != 1 {
		t.Fatalf("MatchingItem should have 1 sub-item, got %d", len(item.Items))
	}
	if len(item.Items[0].Answers) != 2 {
		t.Errorf("MatchingItem should have 2 answers, got %d", len(item.Items[0].Answers))
	}
}

func TestEmptyLesson(t *testing.T) {
	lesson := EmptyLesson("test-id", "Test Title")
	if lesson.ID != "test-id" {
		t.Errorf("EmptyLesson ID mismatch: got %q", lesson.ID)
	}
	if lesson.Type != "lesson" {
		t.Errorf("EmptyLesson Type should be 'lesson', got %q", lesson.Type)
	}
	if len(lesson.Items) != 0 {
		t.Errorf("EmptyLesson should have no items, got %d", len(lesson.Items))
	}
}

func TestSectionLesson(t *testing.T) {
	lesson := SectionLesson("section-id", "Section Title")
	if lesson.Type != "section" {
		t.Errorf("SectionLesson Type should be 'section', got %q", lesson.Type)
	}
}

func TestSampleVariables(t *testing.T) {
	if len(SampleAnswers) == 0 {
		t.Error("SampleAnswers should not be empty")
	}
	if len(SampleSubItems) == 0 {
		t.Error("SampleSubItems should not be empty")
	}

	// Verify one correct answer exists in SampleAnswers
	hasCorrect := false
	for _, a := range SampleAnswers {
		if a.Correct {
			hasCorrect = true
			break
		}
	}
	if !hasCorrect {
		t.Error("SampleAnswers should have at least one correct answer")
	}
}
