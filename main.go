package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/unidoc/unioffice/document"
)

// Core data structures based on the Articulate Rise JSON format
type Course struct {
	ShareID  string     `json:"shareId"`
	Author   string     `json:"author"`
	Course   CourseInfo `json:"course"`
	LabelSet LabelSet   `json:"labelSet"`
}

type CourseInfo struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Color          string          `json:"color"`
	NavigationMode string          `json:"navigationMode"`
	Lessons        []Lesson        `json:"lessons"`
	CoverImage     *Media          `json:"coverImage,omitempty"`
	ExportSettings *ExportSettings `json:"exportSettings,omitempty"`
}

type Lesson struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Icon        string      `json:"icon"`
	Items       []Item      `json:"items"`
	Position    interface{} `json:"position"`
	Ready       bool        `json:"ready"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type Item struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Family   string      `json:"family"`
	Variant  string      `json:"variant"`
	Items    []SubItem   `json:"items"`
	Settings interface{} `json:"settings"`
	Data     interface{} `json:"data"`
	Media    *Media      `json:"media,omitempty"`
}

type SubItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type,omitempty"`
	Title     string    `json:"title,omitempty"`
	Heading   string    `json:"heading,omitempty"`
	Paragraph string    `json:"paragraph,omitempty"`
	Caption   string    `json:"caption,omitempty"`
	Media     *Media    `json:"media,omitempty"`
	Answers   []Answer  `json:"answers,omitempty"`
	Feedback  string    `json:"feedback,omitempty"`
	Front     *CardSide `json:"front,omitempty"`
	Back      *CardSide `json:"back,omitempty"`
}

type Answer struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Correct    bool   `json:"correct"`
	MatchTitle string `json:"matchTitle,omitempty"`
}

type CardSide struct {
	Media       *Media `json:"media,omitempty"`
	Description string `json:"description,omitempty"`
}

type Media struct {
	Image *ImageMedia `json:"image,omitempty"`
	Video *VideoMedia `json:"video,omitempty"`
}

type ImageMedia struct {
	Key           string `json:"key"`
	Type          string `json:"type"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
	CrushedKey    string `json:"crushedKey,omitempty"`
	OriginalUrl   string `json:"originalUrl"`
	UseCrushedKey bool   `json:"useCrushedKey,omitempty"`
}

type VideoMedia struct {
	Key         string `json:"key"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Poster      string `json:"poster,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	InputKey    string `json:"inputKey,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	OriginalUrl string `json:"originalUrl"`
}

type ExportSettings struct {
	Title  string `json:"title"`
	Format string `json:"format"`
}

type LabelSet struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

// Parser main struct
type ArticulateParser struct {
	BaseURL string
	Client  *http.Client
}

func NewArticulateParser() *ArticulateParser {
	return &ArticulateParser{
		BaseURL: "https://rise.articulate.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *ArticulateParser) ExtractShareID(uri string) (string, error) {
	// Extract share ID from URI like: https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/
	re := regexp.MustCompile(`/share/([a-zA-Z0-9_-]+)`)
	matches := re.FindStringSubmatch(uri)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract share ID from URI: %s", uri)
	}
	return matches[1], nil
}

func (p *ArticulateParser) BuildAPIURL(shareID string) string {
	return fmt.Sprintf("%s/api/rise-runtime/boot/share/%s", p.BaseURL, shareID)
}

func (p *ArticulateParser) FetchCourse(uri string) (*Course, error) {
	shareID, err := p.ExtractShareID(uri)
	if err != nil {
		return nil, err
	}

	apiURL := p.BuildAPIURL(shareID)

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

	var course Course
	if err := json.Unmarshal(body, &course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &course, nil
}

func (p *ArticulateParser) LoadCourseFromFile(filePath string) (*Course, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var course Course
	if err := json.Unmarshal(data, &course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &course, nil
}

// HTML cleaner utility
func cleanHTML(html string) string {
	// Remove HTML tags but preserve content
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(html, "")

	// Replace HTML entities
	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")
	cleaned = strings.ReplaceAll(cleaned, "&#39;", "'")
	cleaned = strings.ReplaceAll(cleaned, "&iuml;", "ï")
	cleaned = strings.ReplaceAll(cleaned, "&euml;", "ë")
	cleaned = strings.ReplaceAll(cleaned, "&eacute;", "é")

	// Clean up extra whitespace
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

// Markdown export functions
func (p *ArticulateParser) ExportToMarkdown(course *Course, outputPath string) error {
	var buf bytes.Buffer

	// Write course header
	buf.WriteString(fmt.Sprintf("# %s\n\n", course.Course.Title))

	if course.Course.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", cleanHTML(course.Course.Description)))
	}

	// Add metadata
	buf.WriteString("## Course Information\n\n")
	buf.WriteString(fmt.Sprintf("- **Course ID**: %s\n", course.Course.ID))
	buf.WriteString(fmt.Sprintf("- **Share ID**: %s\n", course.ShareID))
	buf.WriteString(fmt.Sprintf("- **Navigation Mode**: %s\n", course.Course.NavigationMode))
	if course.Course.ExportSettings != nil {
		buf.WriteString(fmt.Sprintf("- **Export Format**: %s\n", course.Course.ExportSettings.Format))
	}
	buf.WriteString("\n---\n\n")

	// Process lessons
	for i, lesson := range course.Course.Lessons {
		if lesson.Type == "section" {
			buf.WriteString(fmt.Sprintf("# %s\n\n", lesson.Title))
			continue
		}

		buf.WriteString(fmt.Sprintf("## Lesson %d: %s\n\n", i+1, lesson.Title))

		if lesson.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", cleanHTML(lesson.Description)))
		}

		// Process lesson items
		for _, item := range lesson.Items {
			p.processItemToMarkdown(&buf, item, 3)
		}

		buf.WriteString("\n---\n\n")
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func (p *ArticulateParser) processItemToMarkdown(buf *bytes.Buffer, item Item, level int) {
	headingPrefix := strings.Repeat("#", level)

	switch item.Type {
	case "text":
		for _, subItem := range item.Items {
			if subItem.Heading != "" {
				heading := cleanHTML(subItem.Heading)
				if heading != "" {
					buf.WriteString(fmt.Sprintf("%s %s\n\n", headingPrefix, heading))
				}
			}
			if subItem.Paragraph != "" {
				paragraph := cleanHTML(subItem.Paragraph)
				if paragraph != "" {
					buf.WriteString(fmt.Sprintf("%s\n\n", paragraph))
				}
			}
		}

	case "list":
		for _, subItem := range item.Items {
			if subItem.Paragraph != "" {
				paragraph := cleanHTML(subItem.Paragraph)
				if paragraph != "" {
					buf.WriteString(fmt.Sprintf("- %s\n", paragraph))
				}
			}
		}
		buf.WriteString("\n")

	case "multimedia":
		buf.WriteString(fmt.Sprintf("%s Media Content\n\n", headingPrefix))
		for _, subItem := range item.Items {
			if subItem.Media != nil {
				if subItem.Media.Video != nil {
					buf.WriteString(fmt.Sprintf("**Video**: %s\n", subItem.Media.Video.OriginalUrl))
					if subItem.Media.Video.Duration > 0 {
						buf.WriteString(fmt.Sprintf("- Duration: %d seconds\n", subItem.Media.Video.Duration))
					}
				}
				if subItem.Media.Image != nil {
					buf.WriteString(fmt.Sprintf("**Image**: %s\n", subItem.Media.Image.OriginalUrl))
				}
			}
			if subItem.Caption != "" {
				caption := cleanHTML(subItem.Caption)
				buf.WriteString(fmt.Sprintf("*%s*\n", caption))
			}
		}
		buf.WriteString("\n")

	case "image":
		buf.WriteString(fmt.Sprintf("%s Image\n\n", headingPrefix))
		for _, subItem := range item.Items {
			if subItem.Media != nil && subItem.Media.Image != nil {
				buf.WriteString(fmt.Sprintf("**Image**: %s\n", subItem.Media.Image.OriginalUrl))
			}
			if subItem.Caption != "" {
				caption := cleanHTML(subItem.Caption)
				buf.WriteString(fmt.Sprintf("*%s*\n", caption))
			}
		}
		buf.WriteString("\n")

	case "knowledgeCheck":
		buf.WriteString(fmt.Sprintf("%s Knowledge Check\n\n", headingPrefix))
		for _, subItem := range item.Items {
			if subItem.Title != "" {
				title := cleanHTML(subItem.Title)
				buf.WriteString(fmt.Sprintf("**Question**: %s\n\n", title))
			}

			buf.WriteString("**Answers**:\n")
			for i, answer := range subItem.Answers {
				answerText := cleanHTML(answer.Title)
				correctMark := ""
				if answer.Correct {
					correctMark = " ✓"
				}
				buf.WriteString(fmt.Sprintf("%d. %s%s\n", i+1, answerText, correctMark))
			}

			if subItem.Feedback != "" {
				feedback := cleanHTML(subItem.Feedback)
				buf.WriteString(fmt.Sprintf("\n**Feedback**: %s\n", feedback))
			}
		}
		buf.WriteString("\n")

	case "interactive":
		buf.WriteString(fmt.Sprintf("%s Interactive Content\n\n", headingPrefix))
		for _, subItem := range item.Items {
			if subItem.Front != nil && subItem.Front.Description != "" {
				desc := cleanHTML(subItem.Front.Description)
				buf.WriteString(fmt.Sprintf("**Front**: %s\n", desc))
			}
			if subItem.Back != nil && subItem.Back.Description != "" {
				desc := cleanHTML(subItem.Back.Description)
				buf.WriteString(fmt.Sprintf("**Back**: %s\n", desc))
			}
		}
		buf.WriteString("\n")

	case "divider":
		buf.WriteString("---\n\n")

	default:
		// Handle unknown types
		if len(item.Items) > 0 {
			buf.WriteString(fmt.Sprintf("%s %s Content\n\n", headingPrefix, strings.Title(item.Type)))
			for _, subItem := range item.Items {
				if subItem.Title != "" {
					title := cleanHTML(subItem.Title)
					buf.WriteString(fmt.Sprintf("- %s\n", title))
				}
			}
			buf.WriteString("\n")
		}
	}
}

// DOCX export functions
func (p *ArticulateParser) ExportToDocx(course *Course, outputPath string) error {
	doc := document.New()

	// Add title
	title := doc.AddParagraph()
	titleRun := title.AddRun()
	titleRun.AddText(course.Course.Title)
	titleRun.Properties().SetSize(20)
	titleRun.Properties().SetBold(true)

	// Add description
	if course.Course.Description != "" {
		desc := doc.AddParagraph()
		descRun := desc.AddRun()
		descRun.AddText(cleanHTML(course.Course.Description))
	}

	// Add course metadata
	metadata := doc.AddParagraph()
	metadataRun := metadata.AddRun()
	metadataRun.Properties().SetBold(true)
	metadataRun.AddText("Course Information")

	courseInfo := doc.AddParagraph()
	courseInfoRun := courseInfo.AddRun()
	courseInfoText := fmt.Sprintf("Course ID: %s\nShare ID: %s\nNavigation Mode: %s",
		course.Course.ID, course.ShareID, course.Course.NavigationMode)
	courseInfoRun.AddText(courseInfoText)

	// Process lessons
	for i, lesson := range course.Course.Lessons {
		if lesson.Type == "section" {
			section := doc.AddParagraph()
			sectionRun := section.AddRun()
			sectionRun.AddText(lesson.Title)
			sectionRun.Properties().SetSize(18)
			sectionRun.Properties().SetBold(true)
			continue
		}

		// Lesson title
		lessonTitle := doc.AddParagraph()
		lessonTitleRun := lessonTitle.AddRun()
		lessonTitleRun.AddText(fmt.Sprintf("Lesson %d: %s", i+1, lesson.Title))
		lessonTitleRun.Properties().SetSize(16)
		lessonTitleRun.Properties().SetBold(true)

		// Lesson description
		if lesson.Description != "" {
			lessonDesc := doc.AddParagraph()
			lessonDescRun := lessonDesc.AddRun()
			lessonDescRun.AddText(cleanHTML(lesson.Description))
		}

		// Process lesson items
		for _, item := range lesson.Items {
			p.processItemToDocx(doc, item)
		}
	}

	return doc.SaveToFile(outputPath)
}

func (p *ArticulateParser) processItemToDocx(doc *document.Document, item Item) {
	switch item.Type {
	case "text":
		for _, subItem := range item.Items {
			if subItem.Heading != "" {
				heading := cleanHTML(subItem.Heading)
				if heading != "" {
					para := doc.AddParagraph()
					run := para.AddRun()
					run.AddText(heading)
					run.Properties().SetBold(true)
				}
			}
			if subItem.Paragraph != "" {
				paragraph := cleanHTML(subItem.Paragraph)
				if paragraph != "" {
					para := doc.AddParagraph()
					run := para.AddRun()
					run.AddText(paragraph)
				}
			}
		}

	case "list":
		for _, subItem := range item.Items {
			if subItem.Paragraph != "" {
				paragraph := cleanHTML(subItem.Paragraph)
				if paragraph != "" {
					para := doc.AddParagraph()
					run := para.AddRun()
					run.AddText("• " + paragraph)
				}
			}
		}

	case "multimedia", "image":
		para := doc.AddParagraph()
		run := para.AddRun()
		run.AddText("[Media Content]")
		run.Properties().SetItalic(true)

		for _, subItem := range item.Items {
			if subItem.Media != nil {
				if subItem.Media.Video != nil {
					mediaPara := doc.AddParagraph()
					mediaRun := mediaPara.AddRun()
					mediaRun.AddText(fmt.Sprintf("Video: %s", subItem.Media.Video.OriginalUrl))
				}
				if subItem.Media.Image != nil {
					mediaPara := doc.AddParagraph()
					mediaRun := mediaPara.AddRun()
					mediaRun.AddText(fmt.Sprintf("Image: %s", subItem.Media.Image.OriginalUrl))
				}
			}
			if subItem.Caption != "" {
				caption := cleanHTML(subItem.Caption)
				captionPara := doc.AddParagraph()
				captionRun := captionPara.AddRun()
				captionRun.AddText(caption)
				captionRun.Properties().SetItalic(true)
			}
		}

	case "knowledgeCheck":
		for _, subItem := range item.Items {
			if subItem.Title != "" {
				title := cleanHTML(subItem.Title)
				questionPara := doc.AddParagraph()
				questionRun := questionPara.AddRun()
				questionRun.AddText("Question: " + title)
				questionRun.Properties().SetBold(true)
			}

			for i, answer := range subItem.Answers {
				answerText := cleanHTML(answer.Title)
				correctMark := ""
				if answer.Correct {
					correctMark = " [CORRECT]"
				}
				answerPara := doc.AddParagraph()
				answerRun := answerPara.AddRun()
				answerRun.AddText(fmt.Sprintf("%d. %s%s", i+1, answerText, correctMark))
			}

			if subItem.Feedback != "" {
				feedback := cleanHTML(subItem.Feedback)
				feedbackPara := doc.AddParagraph()
				feedbackRun := feedbackPara.AddRun()
				feedbackRun.AddText("Feedback: " + feedback)
				feedbackRun.Properties().SetItalic(true)
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: articulate-parser <input_uri_or_file> <output_format> [output_path]")
		fmt.Println("  input_uri_or_file: Articulate Rise URI or local JSON file path")
		fmt.Println("  output_format: md (Markdown) or docx (Word Document)")
		fmt.Println("  output_path: Optional output file path")
		os.Exit(1)
	}

	input := os.Args[1]
	format := strings.ToLower(os.Args[2])

	if format != "md" && format != "docx" {
		log.Fatal("Output format must be 'md' or 'docx'")
	}

	parser := NewArticulateParser()
	var course *Course
	var err error

	// Determine if input is a URI or file path
	if strings.HasPrefix(input, "http") {
		course, err = parser.FetchCourse(input)
	} else {
		course, err = parser.LoadCourseFromFile(input)
	}

	if err != nil {
		log.Fatalf("Failed to load course: %v", err)
	}

	// Determine output path
	var outputPath string
	if len(os.Args) > 3 {
		outputPath = os.Args[3]
	} else {
		baseDir := "output"
		os.MkdirAll(baseDir, 0755)

		// Create safe filename from course title
		safeTitle := regexp.MustCompile(`[^a-zA-Z0-9\-_]`).ReplaceAllString(course.Course.Title, "_")
		if safeTitle == "" {
			safeTitle = "articulate_course"
		}

		outputPath = filepath.Join(baseDir, fmt.Sprintf("%s.%s", safeTitle, format))
	}

	// Export based on format
	switch format {
	case "md":
		err = parser.ExportToMarkdown(course, outputPath)
	case "docx":
		err = parser.ExportToDocx(course, outputPath)
	}

	if err != nil {
		log.Fatalf("Failed to export course: %v", err)
	}

	fmt.Printf("Course successfully exported to: %s\n", outputPath)
	fmt.Printf("Course: %s (%d lessons)\n", course.Course.Title, len(course.Course.Lessons))
}
