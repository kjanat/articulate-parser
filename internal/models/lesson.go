package models

// Lesson represents a single lesson or section within an Articulate Rise course.
// Lessons are the main organizational units and contain various content items.
type Lesson struct {
	// ID is the unique identifier for the lesson
	ID string `json:"id"`
	// Title is the name of the lesson
	Title string `json:"title"`
	// Description is the introductory text for the lesson
	Description string `json:"description"`
	// Type indicates whether this is a regular lesson or a section header
	Type string `json:"type"`
	// Icon is the identifier for the icon displayed with this lesson
	Icon string `json:"icon"`
	// Items is an ordered array of content items within the lesson
	Items []Item `json:"items"`
	// Position stores the ordering information for the lesson
	Position any `json:"position"`
	// Ready indicates whether the lesson is marked as complete
	Ready bool `json:"ready"`
	// CreatedAt is the timestamp when the lesson was created
	CreatedAt string `json:"createdAt"`
	// UpdatedAt is the timestamp when the lesson was last modified
	UpdatedAt string `json:"updatedAt"`
}

// Item represents a content block within a lesson.
// Items can be of various types such as text, multimedia, knowledge checks, etc.
type Item struct {
	// ID is the unique identifier for the item
	ID string `json:"id"`
	// Type indicates the kind of content (text, image, knowledge check, etc.)
	Type string `json:"type"`
	// Family groups similar item types together
	Family string `json:"family"`
	// Variant specifies a sub-type within the main type
	Variant string `json:"variant"`
	// Items contains the actual content elements (sub-items) of this item
	Items []SubItem `json:"items"`
	// Settings contains configuration options specific to this item type
	Settings any `json:"settings"`
	// Data contains additional structured data for the item
	Data any `json:"data"`
	// Media contains any associated media for the item
	Media *Media `json:"media,omitempty"`
}

// SubItem represents a specific content element within an Item.
// SubItems are the most granular content units like paragraphs, headings, or answers.
type SubItem struct {
	// ID is the unique identifier for the sub-item
	ID string `json:"id"`
	// Type indicates the specific kind of sub-item
	Type string `json:"type,omitempty"`
	// Title is the name or label of the sub-item
	Title string `json:"title,omitempty"`
	// Heading is a heading text for this sub-item
	Heading string `json:"heading,omitempty"`
	// Paragraph contains regular text content
	Paragraph string `json:"paragraph,omitempty"`
	// Caption is text associated with media elements
	Caption string `json:"caption,omitempty"`
	// Media contains any associated images or videos
	Media *Media `json:"media,omitempty"`
	// Answers contains possible answers for question-type sub-items
	Answers []Answer `json:"answers,omitempty"`
	// Feedback is the response shown after user interaction
	Feedback string `json:"feedback,omitempty"`
	// Front contains content for the front side of a card-type sub-item
	Front *CardSide `json:"front,omitempty"`
	// Back contains content for the back side of a card-type sub-item
	Back *CardSide `json:"back,omitempty"`
}

// Answer represents a possible response in a knowledge check or quiz item.
type Answer struct {
	// ID is the unique identifier for the answer
	ID string `json:"id"`
	// Title is the text of the answer option
	Title string `json:"title"`
	// Correct indicates whether this is the right answer
	Correct bool `json:"correct"`
	// MatchTitle is used in matching-type questions to pair answers
	MatchTitle string `json:"matchTitle,omitempty"`
}

// CardSide represents one side of a flipcard-type content element.
type CardSide struct {
	// Media is the image or video associated with this side of the card
	Media *Media `json:"media,omitempty"`
	// Description is the text content for this side of the card
	Description string `json:"description,omitempty"`
}
