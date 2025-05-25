// Package models_test provides tests for the data models.
package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestCourse_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of Course.
func TestCourse_JSONMarshalUnmarshal(t *testing.T) {
	original := Course{
		ShareID: "test-share-id",
		Author:  "Test Author",
		Course: CourseInfo{
			ID:             "course-123",
			Title:          "Test Course",
			Description:    "A test course description",
			Color:          "#FF5733",
			NavigationMode: "menu",
			Lessons: []Lesson{
				{
					ID:          "lesson-1",
					Title:       "First Lesson",
					Description: "Lesson description",
					Type:        "lesson",
					Icon:        "icon-1",
					Ready:       true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-02T00:00:00Z",
				},
			},
			ExportSettings: &ExportSettings{
				Title:  "Export Title",
				Format: "scorm",
			},
		},
		LabelSet: LabelSet{
			ID:   "labelset-1",
			Name: "Test Labels",
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Course to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled Course
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Course from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled Course structs do not match")
		t.Logf("Original: %+v", original)
		t.Logf("Unmarshaled: %+v", unmarshaled)
	}
}

// TestCourseInfo_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of CourseInfo.
func TestCourseInfo_JSONMarshalUnmarshal(t *testing.T) {
	original := CourseInfo{
		ID:             "course-456",
		Title:          "Another Test Course",
		Description:    "Another test description",
		Color:          "#33FF57",
		NavigationMode: "linear",
		Lessons: []Lesson{
			{
				ID:    "lesson-2",
				Title: "Second Lesson",
				Type:  "section",
				Items: []Item{
					{
						ID:      "item-1",
						Type:    "text",
						Family:  "text",
						Variant: "paragraph",
						Items: []SubItem{
							{
								Title:     "Sub Item Title",
								Heading:   "Sub Item Heading",
								Paragraph: "Sub item paragraph content",
							},
						},
					},
				},
			},
		},
		CoverImage: &Media{
			Image: &ImageMedia{
				Key:         "img-123",
				Type:        "jpg",
				Width:       800,
				Height:      600,
				OriginalUrl: "https://example.com/image.jpg",
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal CourseInfo to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled CourseInfo
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal CourseInfo from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled CourseInfo structs do not match")
	}
}

// TestLesson_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of Lesson.
func TestLesson_JSONMarshalUnmarshal(t *testing.T) {
	original := Lesson{
		ID:          "lesson-test",
		Title:       "Test Lesson",
		Description: "Test lesson description",
		Type:        "lesson",
		Icon:        "lesson-icon",
		Ready:       true,
		CreatedAt:   "2023-06-01T12:00:00Z",
		UpdatedAt:   "2023-06-01T13:00:00Z",
		Position:    map[string]interface{}{"x": 1, "y": 2},
		Items: []Item{
			{
				ID:      "item-test",
				Type:    "multimedia",
				Family:  "media",
				Variant: "video",
				Items: []SubItem{
					{
						Caption: "Video caption",
						Media: &Media{
							Video: &VideoMedia{
								Key:         "video-123",
								URL:         "https://example.com/video.mp4",
								Type:        "mp4",
								Duration:    120,
								OriginalUrl: "https://example.com/video.mp4",
							},
						},
					},
				},
				Settings: map[string]interface{}{"autoplay": false},
				Data:     map[string]interface{}{"metadata": "test"},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Lesson to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled Lesson
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Lesson from JSON: %v", err)
	}

	// Compare structures
	if !compareLessons(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled Lesson structs do not match")
	}
}

// TestItem_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of Item.
func TestItem_JSONMarshalUnmarshal(t *testing.T) {
	original := Item{
		ID:      "item-json-test",
		Type:    "knowledgeCheck",
		Family:  "assessment",
		Variant: "multipleChoice",
		Items: []SubItem{
			{
				Title: "What is the answer?",
				Answers: []Answer{
					{Title: "Option A", Correct: false},
					{Title: "Option B", Correct: true},
					{Title: "Option C", Correct: false},
				},
				Feedback: "Well done!",
			},
		},
		Settings: map[string]interface{}{
			"allowRetry": true,
			"showAnswer": true,
		},
		Data: map[string]interface{}{
			"points": 10,
			"weight": 1.5,
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Item to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled Item
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Item from JSON: %v", err)
	}

	// Compare structures
	if !compareItem(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled Item structs do not match")
	}
}

// TestSubItem_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of SubItem.
func TestSubItem_JSONMarshalUnmarshal(t *testing.T) {
	original := SubItem{
		Title:     "Test SubItem Title",
		Heading:   "Test SubItem Heading",
		Paragraph: "Test paragraph with content",
		Caption:   "Test caption",
		Feedback:  "Test feedback message",
		Answers: []Answer{
			{Title: "First answer", Correct: true},
			{Title: "Second answer", Correct: false},
		},
		Media: &Media{
			Image: &ImageMedia{
				Key:           "subitem-img",
				Type:          "png",
				Width:         400,
				Height:        300,
				OriginalUrl:   "https://example.com/subitem.png",
				CrushedKey:    "crushed-123",
				UseCrushedKey: true,
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal SubItem to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled SubItem
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SubItem from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled SubItem structs do not match")
	}
}

// TestAnswer_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of Answer.
func TestAnswer_JSONMarshalUnmarshal(t *testing.T) {
	original := Answer{
		Title:   "Test answer text",
		Correct: true,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Answer to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled Answer
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Answer from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled Answer structs do not match")
	}
}

// TestMedia_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of Media.
func TestMedia_JSONMarshalUnmarshal(t *testing.T) {
	// Test with Image
	originalImage := Media{
		Image: &ImageMedia{
			Key:           "media-img-test",
			Type:          "jpeg",
			Width:         1200,
			Height:        800,
			OriginalUrl:   "https://example.com/media.jpg",
			CrushedKey:    "crushed-media",
			UseCrushedKey: false,
		},
	}

	jsonData, err := json.Marshal(originalImage)
	if err != nil {
		t.Fatalf("Failed to marshal Media with Image to JSON: %v", err)
	}

	var unmarshaledImage Media
	err = json.Unmarshal(jsonData, &unmarshaledImage)
	if err != nil {
		t.Fatalf("Failed to unmarshal Media with Image from JSON: %v", err)
	}

	if !reflect.DeepEqual(originalImage, unmarshaledImage) {
		t.Errorf("Marshaled and unmarshaled Media with Image do not match")
	}

	// Test with Video
	originalVideo := Media{
		Video: &VideoMedia{
			Key:         "media-video-test",
			URL:         "https://example.com/media.mp4",
			Type:        "mp4",
			Duration:    300,
			Poster:      "https://example.com/poster.jpg",
			Thumbnail:   "https://example.com/thumb.jpg",
			InputKey:    "input-123",
			OriginalUrl: "https://example.com/original.mp4",
		},
	}

	jsonData, err = json.Marshal(originalVideo)
	if err != nil {
		t.Fatalf("Failed to marshal Media with Video to JSON: %v", err)
	}

	var unmarshaledVideo Media
	err = json.Unmarshal(jsonData, &unmarshaledVideo)
	if err != nil {
		t.Fatalf("Failed to unmarshal Media with Video from JSON: %v", err)
	}

	if !reflect.DeepEqual(originalVideo, unmarshaledVideo) {
		t.Errorf("Marshaled and unmarshaled Media with Video do not match")
	}
}

// TestImageMedia_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of ImageMedia.
func TestImageMedia_JSONMarshalUnmarshal(t *testing.T) {
	original := ImageMedia{
		Key:           "image-media-test",
		Type:          "gif",
		Width:         640,
		Height:        480,
		OriginalUrl:   "https://example.com/image.gif",
		CrushedKey:    "crushed-gif",
		UseCrushedKey: true,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal ImageMedia to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled ImageMedia
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ImageMedia from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled ImageMedia structs do not match")
	}
}

// TestVideoMedia_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of VideoMedia.
func TestVideoMedia_JSONMarshalUnmarshal(t *testing.T) {
	original := VideoMedia{
		Key:         "video-media-test",
		URL:         "https://example.com/video.webm",
		Type:        "webm",
		Duration:    450,
		Poster:      "https://example.com/poster.jpg",
		Thumbnail:   "https://example.com/thumbnail.jpg",
		InputKey:    "upload-456",
		OriginalUrl: "https://example.com/original.webm",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal VideoMedia to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled VideoMedia
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VideoMedia from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled VideoMedia structs do not match")
	}
}

// TestExportSettings_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of ExportSettings.
func TestExportSettings_JSONMarshalUnmarshal(t *testing.T) {
	original := ExportSettings{
		Title:  "Custom Export Title",
		Format: "xAPI",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal ExportSettings to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled ExportSettings
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ExportSettings from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled ExportSettings structs do not match")
	}
}

// TestLabelSet_JSONMarshalUnmarshal tests JSON marshaling and unmarshaling of LabelSet.
func TestLabelSet_JSONMarshalUnmarshal(t *testing.T) {
	original := LabelSet{
		ID:   "labelset-test",
		Name: "Test Label Set",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal LabelSet to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled LabelSet
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal LabelSet from JSON: %v", err)
	}

	// Compare structures
	if !reflect.DeepEqual(original, unmarshaled) {
		t.Errorf("Marshaled and unmarshaled LabelSet structs do not match")
	}
}

// TestEmptyStructures tests marshaling and unmarshaling of empty structures.
func TestEmptyStructures(t *testing.T) {
	testCases := []struct {
		name string
		data interface{}
	}{
		{"Empty Course", Course{}},
		{"Empty CourseInfo", CourseInfo{}},
		{"Empty Lesson", Lesson{}},
		{"Empty Item", Item{}},
		{"Empty SubItem", SubItem{}},
		{"Empty Answer", Answer{}},
		{"Empty Media", Media{}},
		{"Empty ImageMedia", ImageMedia{}},
		{"Empty VideoMedia", VideoMedia{}},
		{"Empty ExportSettings", ExportSettings{}},
		{"Empty LabelSet", LabelSet{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tc.data)
			if err != nil {
				t.Fatalf("Failed to marshal %s to JSON: %v", tc.name, err)
			}

			// Unmarshal from JSON
			result := reflect.New(reflect.TypeOf(tc.data)).Interface()
			err = json.Unmarshal(jsonData, result)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s from JSON: %v", tc.name, err)
			}

			// Basic validation that no errors occurred
			if len(jsonData) == 0 {
				t.Errorf("%s should produce some JSON output", tc.name)
			}
		})
	}
}

// TestNilPointerSafety tests that nil pointers in optional fields are handled correctly.
func TestNilPointerSafety(t *testing.T) {
	course := Course{
		ShareID: "nil-test",
		Course: CourseInfo{
			ID:             "nil-course",
			Title:          "Nil Pointer Test",
			CoverImage:     nil, // Test nil pointer
			ExportSettings: nil, // Test nil pointer
			Lessons: []Lesson{
				{
					ID:    "lesson-nil",
					Title: "Lesson with nil media",
					Items: []Item{
						{
							ID:   "item-nil",
							Type: "text",
							Items: []SubItem{
								{
									Title: "SubItem with nil media",
									Media: nil, // Test nil pointer
								},
							},
							Media: nil, // Test nil pointer
						},
					},
				},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(course)
	if err != nil {
		t.Fatalf("Failed to marshal Course with nil pointers to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled Course
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Course with nil pointers from JSON: %v", err)
	}

	// Basic validation
	if unmarshaled.ShareID != "nil-test" {
		t.Error("ShareID should be preserved")
	}
	if unmarshaled.Course.Title != "Nil Pointer Test" {
		t.Error("Course title should be preserved")
	}
}

// TestJSONTagsPresence tests that JSON tags are properly defined.
func TestJSONTagsPresence(t *testing.T) {
	// Test that important fields have JSON tags
	courseType := reflect.TypeOf(Course{})
	if courseType.Kind() == reflect.Struct {
		field, found := courseType.FieldByName("ShareID")
		if !found {
			t.Error("ShareID field not found")
		} else {
			tag := field.Tag.Get("json")
			if tag == "" {
				t.Error("ShareID should have json tag")
			}
			if tag != "shareId" {
				t.Errorf("ShareID json tag should be 'shareId', got '%s'", tag)
			}
		}
	}

	// Test CourseInfo
	courseInfoType := reflect.TypeOf(CourseInfo{})
	if courseInfoType.Kind() == reflect.Struct {
		field, found := courseInfoType.FieldByName("NavigationMode")
		if !found {
			t.Error("NavigationMode field not found")
		} else {
			tag := field.Tag.Get("json")
			if tag == "" {
				t.Error("NavigationMode should have json tag")
			}
		}
	}
}

// BenchmarkCourse_JSONMarshal benchmarks JSON marshaling of Course.
func BenchmarkCourse_JSONMarshal(b *testing.B) {
	course := Course{
		ShareID: "benchmark-id",
		Author:  "Benchmark Author",
		Course: CourseInfo{
			ID:    "benchmark-course",
			Title: "Benchmark Course",
			Lessons: []Lesson{
				{
					ID:    "lesson-1",
					Title: "Lesson 1",
					Items: []Item{
						{
							ID:   "item-1",
							Type: "text",
							Items: []SubItem{
								{Title: "SubItem 1"},
							},
						},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(course)
	}
}

// BenchmarkCourse_JSONUnmarshal benchmarks JSON unmarshaling of Course.
func BenchmarkCourse_JSONUnmarshal(b *testing.B) {
	course := Course{
		ShareID: "benchmark-id",
		Author:  "Benchmark Author",
		Course: CourseInfo{
			ID:    "benchmark-course",
			Title: "Benchmark Course",
			Lessons: []Lesson{
				{
					ID:    "lesson-1",
					Title: "Lesson 1",
					Items: []Item{
						{
							ID:   "item-1",
							Type: "text",
							Items: []SubItem{
								{Title: "SubItem 1"},
							},
						},
					},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(course)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Course
		_ = json.Unmarshal(jsonData, &result)
	}
}

// compareMaps compares two interface{} values that should be maps
func compareMaps(original, unmarshaled interface{}) bool {
	origMap, origOk := original.(map[string]interface{})
	unMap, unOk := unmarshaled.(map[string]interface{})

	if !origOk || !unOk {
		// If not maps, use deep equal
		return reflect.DeepEqual(original, unmarshaled)
	}

	if len(origMap) != len(unMap) {
		return false
	}

	for key, origVal := range origMap {
		unVal, exists := unMap[key]
		if !exists {
			return false
		}

		// Handle numeric type conversion from JSON
		switch origVal := origVal.(type) {
		case int:
			if unFloat, ok := unVal.(float64); ok {
				if float64(origVal) != unFloat {
					return false
				}
			} else {
				return false
			}
		case float64:
			if unFloat, ok := unVal.(float64); ok {
				if origVal != unFloat {
					return false
				}
			} else {
				return false
			}
		default:
			if !reflect.DeepEqual(origVal, unVal) {
				return false
			}
		}
	}
	return true
}

// compareLessons compares two Lesson structs accounting for JSON type conversion
func compareLessons(original, unmarshaled Lesson) bool {
	// Compare all fields except Position and Items
	if original.ID != unmarshaled.ID ||
		original.Title != unmarshaled.Title ||
		original.Description != unmarshaled.Description ||
		original.Type != unmarshaled.Type ||
		original.Icon != unmarshaled.Icon ||
		original.Ready != unmarshaled.Ready ||
		original.CreatedAt != unmarshaled.CreatedAt ||
		original.UpdatedAt != unmarshaled.UpdatedAt {
		return false
	}

	// Compare Position
	if !compareMaps(original.Position, unmarshaled.Position) {
		return false
	}

	// Compare Items
	return compareItems(original.Items, unmarshaled.Items)
}

// compareItems compares two Item slices accounting for JSON type conversion
func compareItems(original, unmarshaled []Item) bool {
	if len(original) != len(unmarshaled) {
		return false
	}

	for i := range original {
		if !compareItem(original[i], unmarshaled[i]) {
			return false
		}
	}
	return true
}

// compareItem compares two Item structs accounting for JSON type conversion
func compareItem(original, unmarshaled Item) bool {
	// Compare basic fields
	if original.ID != unmarshaled.ID ||
		original.Type != unmarshaled.Type ||
		original.Family != unmarshaled.Family ||
		original.Variant != unmarshaled.Variant {
		return false
	}

	// Compare Settings and Data
	if !compareMaps(original.Settings, unmarshaled.Settings) {
		return false
	}

	if !compareMaps(original.Data, unmarshaled.Data) {
		return false
	}

	// Compare Items (SubItems)
	if len(original.Items) != len(unmarshaled.Items) {
		return false
	}

	for i := range original.Items {
		if !reflect.DeepEqual(original.Items[i], unmarshaled.Items[i]) {
			return false
		}
	}

	// Compare Media
	if !reflect.DeepEqual(original.Media, unmarshaled.Media) {
		return false
	}

	return true
}
