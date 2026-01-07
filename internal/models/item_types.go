package models

// Item type constants used throughout the application.
// These match the type field values in Articulate Rise course data.
const (
	// ItemTypeText represents text content blocks with headings and paragraphs.
	ItemTypeText = "text"

	// ItemTypeList represents bullet or numbered list content.
	ItemTypeList = "list"

	// ItemTypeKnowledgeCheck represents quiz or assessment items.
	ItemTypeKnowledgeCheck = "knowledgecheck"

	// ItemTypeMultimedia represents video or audio content.
	ItemTypeMultimedia = "multimedia"

	// ItemTypeImage represents standalone image content.
	ItemTypeImage = "image"

	// ItemTypeInteractive represents interactive elements like accordions or tabs.
	ItemTypeInteractive = "interactive"

	// ItemTypeDivider represents visual separator elements.
	ItemTypeDivider = "divider"

	// ItemTypeQuote represents quotation blocks.
	ItemTypeQuote = "quote"

	// ItemTypeFlipcard represents flashcard-style content.
	ItemTypeFlipcard = "flipcard"
)

// Lesson type constants.
const (
	// LessonTypeSection represents a section header (not a real lesson).
	LessonTypeSection = "section"

	// LessonTypeLesson represents a regular lesson.
	LessonTypeLesson = "lesson"
)
