package exporters

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kjanat/articulate-parser/internal/models"
	"github.com/kjanat/articulate-parser/internal/services"
)

// Item type constants.
const (
	itemTypeText           = "text"
	itemTypeList           = "list"
	itemTypeKnowledgeCheck = "knowledgecheck"
	itemTypeMultimedia     = "multimedia"
	itemTypeImage          = "image"
	itemTypeInteractive    = "interactive"
	itemTypeDivider        = "divider"
)

// templateData represents the data structure passed to the HTML template.
type templateData struct {
	Course   models.CourseInfo
	ShareID  string
	Sections []templateSection
	CSS      string
}

// templateSection represents a course section or lesson.
type templateSection struct {
	Type        string
	Title       string
	Number      int
	Description string
	Items       []templateItem
}

// templateItem represents a course item with preprocessed data.
type templateItem struct {
	Type      string
	TypeTitle string
	Items     []templateSubItem
}

// templateSubItem represents a sub-item with preprocessed data.
type templateSubItem struct {
	Heading   string
	Paragraph string
	Title     string
	Caption   string
	CleanText string
	Answers   []models.Answer
	Feedback  string
	Media     *models.Media
}

// prepareTemplateData converts a Course model into template-friendly data.
func prepareTemplateData(course *models.Course, htmlCleaner *services.HTMLCleaner) *templateData {
	data := &templateData{
		Course:   course.Course,
		ShareID:  course.ShareID,
		Sections: make([]templateSection, 0, len(course.Course.Lessons)),
		CSS:      defaultCSS,
	}

	lessonCounter := 0
	for _, lesson := range course.Course.Lessons {
		section := templateSection{
			Type:        lesson.Type,
			Title:       lesson.Title,
			Description: lesson.Description,
		}

		if lesson.Type != "section" {
			lessonCounter++
			section.Number = lessonCounter
			section.Items = prepareItems(lesson.Items, htmlCleaner)
		}

		data.Sections = append(data.Sections, section)
	}

	return data
}

// prepareItems converts model Items to template Items.
func prepareItems(items []models.Item, htmlCleaner *services.HTMLCleaner) []templateItem {
	result := make([]templateItem, 0, len(items))

	for _, item := range items {
		tItem := templateItem{
			Type:  strings.ToLower(item.Type),
			Items: make([]templateSubItem, 0, len(item.Items)),
		}

		// Set type title for unknown items
		if tItem.Type != itemTypeText && tItem.Type != itemTypeList && tItem.Type != itemTypeKnowledgeCheck &&
			tItem.Type != itemTypeMultimedia && tItem.Type != itemTypeImage && tItem.Type != itemTypeInteractive &&
			tItem.Type != itemTypeDivider {
			caser := cases.Title(language.English)
			tItem.TypeTitle = caser.String(item.Type)
		}

		// Process sub-items
		for _, subItem := range item.Items {
			tSubItem := templateSubItem{
				Heading:   subItem.Heading,
				Paragraph: subItem.Paragraph,
				Title:     subItem.Title,
				Caption:   subItem.Caption,
				Answers:   subItem.Answers,
				Feedback:  subItem.Feedback,
				Media:     subItem.Media,
			}

			// Clean HTML for list items
			if tItem.Type == itemTypeList && subItem.Paragraph != "" {
				tSubItem.CleanText = htmlCleaner.CleanHTML(subItem.Paragraph)
			}

			tItem.Items = append(tItem.Items, tSubItem)
		}

		result = append(result, tItem)
	}

	return result
}
