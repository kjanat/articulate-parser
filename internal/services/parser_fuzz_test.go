package services

import (
	"encoding/json"
	"testing"

	"github.com/kjanat/articulate-parser/internal/models"
)

// FuzzArticulateParser_ParseJSON fuzzes JSON parsing to ensure no panics on malformed input.
// It tests that the parser gracefully handles any byte sequence without crashing.
func FuzzArticulateParser_ParseJSON(f *testing.F) {
	// Add seed corpus with valid JSON structures
	f.Add([]byte(`{"shareId":"test","course":{"id":"1","title":"Test"}}`))
	f.Add([]byte(`{"shareId":"abc123","author":"John Doe","course":{"id":"c1","title":"Course Title","description":"A test course"}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"course":{}}`))
	f.Add([]byte(`{"shareId":"","course":{"id":"","title":"","lessons":[]}}`))

	// Add seed corpus with invalid/edge case inputs
	f.Add([]byte(`invalid`))
	f.Add([]byte(``))
	f.Add([]byte(`null`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`true`))
	f.Add([]byte(`{`))
	f.Add([]byte(`}`))
	f.Add([]byte(`{"unclosed": `))
	f.Add([]byte(`{"key": "value with \x00 null byte"}`))

	// Add seed corpus with nested structures
	f.Add([]byte(`{"course":{"lessons":[{"id":"l1","title":"Lesson 1","items":[{"id":"i1","type":"text","data":"content"}]}]}}`))
	f.Add([]byte(`{"course":{"lessons":[{},{},{}]}}`))

	// Add seed corpus with special characters
	f.Add([]byte(`{"shareId":"test<script>alert('xss')</script>"}`))
	f.Add([]byte(`{"course":{"title":"Test\n\r\t\u0000Title"}}`))
	f.Add([]byte(`{"course":{"description":"<p>HTML content</p>"}}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Test that unmarshaling doesn't panic
		var course models.Course
		_ = json.Unmarshal(data, &course)

		// Also test that re-marshaling doesn't panic
		_, _ = json.Marshal(&course)

		// Test accessing fields doesn't panic
		_ = course.ShareID
		_ = course.Author
		_ = course.Course.ID
		_ = course.Course.Title
		_ = course.Course.Description

		// Test iterating over lessons doesn't panic
		for _, lesson := range course.Course.Lessons {
			_ = lesson.ID
			_ = lesson.Title
			for _, item := range lesson.Items {
				_ = item.ID
				_ = item.Type
			}
		}
	})
}

// FuzzArticulateParser_LoadCourseFromFile_JSON tests JSON parsing behavior used by LoadCourseFromFile.
// This fuzz test ensures the JSON unmarshaling logic handles arbitrary input gracefully.
func FuzzArticulateParser_LoadCourseFromFile_JSON(f *testing.F) {
	// Seed corpus mimicking file contents
	f.Add([]byte(`{"shareId":"file-test","course":{"id":"1","title":"From File"}}`))
	f.Add([]byte(`{
		"shareId": "multiline",
		"course": {
			"id": "2",
			"title": "Multiline JSON"
		}
	}`))
	f.Add([]byte("\xef\xbb\xbf{\"shareId\":\"bom\"}")) // UTF-8 BOM prefix
	f.Add([]byte(`{"shareId": "spaces",  "course":  {}}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var course models.Course
		err := json.Unmarshal(data, &course)

		// If unmarshaling succeeded, verify the result is usable
		if err == nil {
			// Access all fields to ensure no panics on zero values
			_ = course.ShareID
			_ = course.Author
			_ = course.Course.ID
			_ = course.Course.Title

			// Iterate through nested structures
			for _, lesson := range course.Course.Lessons {
				_ = lesson.ID
				_ = lesson.Title
				_ = lesson.Type
				for _, item := range lesson.Items {
					_ = item.ID
					_ = item.Type
					_ = item.Data
				}
			}
		}
	})
}

// FuzzHTMLCleaner_CleanHTML fuzzes the HTML cleaner to ensure no panics on malformed HTML.
// It tests that arbitrary HTML-like input is handled gracefully.
func FuzzHTMLCleaner_CleanHTML(f *testing.F) {
	cleaner := NewHTMLCleaner()

	// Add seed corpus with valid HTML patterns
	f.Add("<p>Simple paragraph</p>")
	f.Add("<div><p>Nested</p></div>")
	f.Add("<h1>Title</h1><p>Content</p>")
	f.Add("<strong>Bold</strong> and <em>italic</em>")
	f.Add("<a href=\"https://example.com\">Link</a>")
	f.Add("<img src=\"image.png\" alt=\"Image\"/>")
	f.Add("<br/><hr/>")
	f.Add("<ul><li>Item 1</li><li>Item 2</li></ul>")
	f.Add("<table><tr><td>Cell</td></tr></table>")

	// Add seed corpus with HTML entities
	f.Add("AT&amp;T")
	f.Add("&lt;script&gt;alert('xss')&lt;/script&gt;")
	f.Add("&nbsp;&nbsp;&nbsp;")
	f.Add("&#39;apostrophe&#39;")
	f.Add("&quot;quoted&quot;")
	f.Add("&eacute;&iuml;&ouml;")
	f.Add("&#x27;hex entity&#x27;")
	f.Add("&#8217;numeric entity&#8217;")

	// Add seed corpus with malformed/edge case HTML
	f.Add("")
	f.Add("plain text no HTML")
	f.Add("<p>Unclosed paragraph")
	f.Add("</p>Close tag first<p>")
	f.Add("<div><div><div>Deep nesting")
	f.Add("<script>alert('xss')</script>")
	f.Add("<style>body{color:red}</style>")
	f.Add("<!--comment-->")
	f.Add("<!DOCTYPE html>")
	f.Add("<![CDATA[data]]>")
	f.Add("<p class=\"test\" id='test' data-attr=unquoted>Attributes</p>")
	f.Add("<p\n\t>Whitespace in tag</p\n>")
	f.Add("<>empty tag</>")
	f.Add("<<<nested angles>>>")
	f.Add("<p onclick=\"evil()\">Event handler</p>")

	// Add seed corpus with special characters
	f.Add("<p>\x00null byte</p>")
	f.Add("<p>\xffhigh byte</p>")
	f.Add("<p>emoji: \xf0\x9f\x98\x80</p>")
	f.Add("<p>unicode: \u2603</p>")
	f.Add("<p>newlines\n\r\nand\rtabs\there</p>")

	// Add seed corpus with potentially problematic patterns
	f.Add(string(make([]byte, 0)))                     // Empty
	f.Add(string(make([]byte, 1)))                     // Single null byte
	f.Add("<" + string(make([]byte, 100)) + ">")       // Tag with null bytes
	f.Add("<p>" + string(make([]byte, 1000)) + "</p>") // Large content

	f.Fuzz(func(t *testing.T, input string) {
		// Test that CleanHTML doesn't panic on any input
		result := cleaner.CleanHTML(input)

		// Verify the result is a valid string (doesn't panic on access)
		_ = len(result)

		// Test that the result can be used in common string operations without panic
		_ = result == ""
		_ = result == input

		// Test that calling CleanHTML on its own output doesn't panic (idempotency check)
		_ = cleaner.CleanHTML(result)
	})
}

// FuzzHTMLCleaner_CleanHTML_Bytes tests HTML cleaner with raw byte sequences.
// This catches issues with invalid UTF-8 and binary data in HTML content.
func FuzzHTMLCleaner_CleanHTML_Bytes(f *testing.F) {
	cleaner := NewHTMLCleaner()

	// Add seed corpus with byte sequences
	f.Add([]byte("<p>normal</p>"))
	f.Add([]byte{0x3c, 0x70, 0x3e, 0xff, 0xfe, 0x3c, 0x2f, 0x70, 0x3e})       // <p>..invalid..</p>
	f.Add([]byte{0x00})                                                       // Null byte
	f.Add([]byte{0xff, 0xfe})                                                 // Invalid UTF-8
	f.Add([]byte{0xef, 0xbb, 0xbf, 0x3c, 0x70, 0x3e, 0x3c, 0x2f, 0x70, 0x3e}) // UTF-8 BOM + <p></p>
	f.Add([]byte("<p>\xc0\xaf</p>"))                                          // Overlong UTF-8

	f.Fuzz(func(t *testing.T, data []byte) {
		// Convert bytes to string and test
		input := string(data)

		// Test that CleanHTML doesn't panic
		result := cleaner.CleanHTML(input)

		// Verify result is usable
		_ = len(result)
		_ = []byte(result)
	})
}
