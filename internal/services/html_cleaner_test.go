// Package services_test provides tests for the HTML cleaner service.
package services

import (
	"strings"
	"testing"
)

// TestNewHTMLCleaner tests the NewHTMLCleaner constructor.
func TestNewHTMLCleaner(t *testing.T) {
	cleaner := NewHTMLCleaner()

	if cleaner == nil {
		t.Fatal("NewHTMLCleaner() returned nil")
	}
}

// TestHTMLCleaner_CleanHTML tests the CleanHTML method with various HTML inputs.
func TestHTMLCleaner_CleanHTML(t *testing.T) {
	cleaner := NewHTMLCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text (no HTML)",
			input:    "This is plain text",
			expected: "This is plain text",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple HTML tag",
			input:    "<p>Hello world</p>",
			expected: "Hello world",
		},
		{
			name:     "multiple HTML tags",
			input:    "<h1>Title</h1><p>Paragraph text</p>",
			expected: "TitleParagraph text",
		},
		{
			name:     "nested HTML tags",
			input:    "<div><h1>Title</h1><p>Paragraph with <strong>bold</strong> text</p></div>",
			expected: "TitleParagraph with bold text",
		},
		{
			name:     "HTML with attributes",
			input:    "<p class=\"test\" id=\"para1\">Text with attributes</p>",
			expected: "Text with attributes",
		},
		{
			name:     "self-closing tags",
			input:    "Line 1<br/>Line 2<hr/>End",
			expected: "Line 1Line 2End",
		},
		{
			name:     "HTML entities - basic",
			input:    "AT&amp;T &lt;company&gt; &quot;quoted&quot; &nbsp; text",
			expected: "AT&T <company> \"quoted\" text",
		},
		{
			name:     "HTML entities - apostrophe",
			input:    "It&#39;s a test",
			expected: "It's a test",
		},
		{
			name:     "HTML entities - special characters",
			input:    "&iuml;ber &euml;lite &eacute;cart&eacute;",
			expected: "ïber ëlite écarté",
		},
		{
			name:     "HTML entities - nbsp",
			input:    "Word1&nbsp;&nbsp;&nbsp;Word2",
			expected: "Word1 Word2",
		},
		{
			name:     "mixed HTML and entities",
			input:    "<p>Hello &amp; welcome to <strong>our</strong> site!</p>",
			expected: "Hello & welcome to our site!",
		},
		{
			name:     "multiple whitespace",
			input:    "Text   with\t\tmultiple\n\nspaces",
			expected: "Text with multiple spaces",
		},
		{
			name:     "whitespace with HTML",
			input:    "<p>  Text  with  </p>  <div>  spaces  </div>  ",
			expected: "Text with spaces",
		},
		{
			name:     "complex content",
			input:    "<div class=\"content\"><h1>Course Title</h1><p>This is a <em>great</em> course about &amp; HTML entities like &nbsp; and &quot;quotes&quot;.</p></div>",
			expected: "Course TitleThis is a great course about & HTML entities like and \"quotes\".",
		},
		{
			name:     "malformed HTML",
			input:    "<p>Unclosed paragraph<div>Another <span>tag</p></div>",
			expected: "Unclosed paragraphAnother tag",
		},
		{
			name:     "HTML comments (should be removed)",
			input:    "Text before<!-- This is a comment -->Text after",
			expected: "Text beforeText after",
		},
		{
			name:     "script and style tags content",
			input:    "<script>alert('test');</script>Content<style>body{color:red;}</style>",
			expected: "Content", // Script and style tags are correctly skipped
		},
		{
			name:     "line breaks and formatting",
			input:    "<p>Line 1</p>\n<p>Line 2</p>\n<p>Line 3</p>",
			expected: "Line 1 Line 2 Line 3",
		},
		{
			name:     "only whitespace",
			input:    "   \t\n   ",
			expected: "",
		},
		{
			name:     "only HTML tags",
			input:    "<div><p></p></div>",
			expected: "",
		},
		{
			name:     "HTML with newlines",
			input:    "<p>\n  Paragraph with\n  line breaks\n</p>",
			expected: "Paragraph with line breaks",
		},
		{
			name:     "complex nested structure",
			input:    "<article><header><h1>Title</h1></header><section><p>First paragraph with <a href=\"#\">link</a>.</p><ul><li>Item 1</li><li>Item 2</li></ul></section></article>",
			expected: "TitleFirst paragraph with link.Item 1Item 2",
		},
		{
			name:     "entities in attributes (should still be processed)",
			input:    "<p title=\"AT&amp;T\">Content</p>",
			expected: "Content",
		},
		{
			name:     "special HTML5 entities",
			input:    "Left arrow &larr; Right arrow &rarr;",
			expected: "Left arrow ← Right arrow →", // HTML5 entities are properly handled by the parser
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.CleanHTML(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHTMLCleaner_CleanHTML_LargeContent tests the CleanHTML method with large content.
func TestHTMLCleaner_CleanHTML_LargeContent(t *testing.T) {
	cleaner := NewHTMLCleaner()

	// Create a large HTML string
	var builder strings.Builder
	builder.WriteString("<html><body>")
	for i := range 1000 {
		builder.WriteString("<p>Paragraph ")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(" with some content &amp; entities.</p>")
	}
	builder.WriteString("</body></html>")

	input := builder.String()
	result := cleaner.CleanHTML(input)

	// Check that HTML tags are removed
	if strings.Contains(result, "<") || strings.Contains(result, ">") {
		t.Error("Result should not contain HTML tags")
	}

	// Check that content is preserved
	if !strings.Contains(result, "Paragraph") {
		t.Error("Result should contain paragraph content")
	}

	// Check that entities are converted
	if strings.Contains(result, "&amp;") {
		t.Error("Result should not contain unconverted HTML entities")
	}
	if !strings.Contains(result, "&") {
		t.Error("Result should contain converted ampersand")
	}
}

// TestHTMLCleaner_CleanHTML_EdgeCases tests edge cases for the CleanHTML method.
func TestHTMLCleaner_CleanHTML_EdgeCases(t *testing.T) {
	cleaner := NewHTMLCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "only entities",
			input:    "&amp;&lt;&gt;&quot;&#39;&nbsp;",
			expected: "&<>\"'",
		},
		{
			name:     "repeated entities",
			input:    "&amp;&amp;&amp;",
			expected: "&&&",
		},
		{
			name:     "entities without semicolon (properly converted)",
			input:    "&amp test &lt test",
			expected: "& test < test", // Parser handles entities even without semicolons in some cases
		},
		{
			name:     "mixed valid and invalid entities",
			input:    "&amp; &invalid; &lt; &fake;",
			expected: "& &invalid; < &fake;",
		},
		{
			name:     "unclosed tag at end",
			input:    "Content <p>with unclosed",
			expected: "Content with unclosed",
		},
		{
			name:     "tag with no closing bracket",
			input:    "Content <p class='test' with no closing bracket",
			expected: "Content", // Parser handles malformed HTML gracefully
		},
		{
			name:     "extremely nested tags",
			input:    "<div><div><div><div><div>Deep content</div></div></div></div></div>",
			expected: "Deep content",
		},
		{
			name:     "empty tags with whitespace",
			input:    "<p>   </p><div>\t\n</div>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.CleanHTML(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHTMLCleaner_CleanHTML_Unicode tests Unicode content handling.
func TestHTMLCleaner_CleanHTML_Unicode(t *testing.T) {
	cleaner := NewHTMLCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unicode characters",
			input:    "<p>Hello 世界! Café naïve résumé</p>",
			expected: "Hello 世界! Café naïve résumé",
		},
		{
			name:     "unicode with entities",
			input:    "<p>Unicode: 你好 &amp; emoji: 🌍</p>",
			expected: "Unicode: 你好 & emoji: 🌍",
		},
		{
			name:     "mixed scripts",
			input:    "<div>English العربية русский 日本語</div>",
			expected: "English العربية русский 日本語",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.CleanHTML(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkHTMLCleaner_CleanHTML benchmarks the CleanHTML method.
func BenchmarkHTMLCleaner_CleanHTML(b *testing.B) {
	cleaner := NewHTMLCleaner()
	input := "<div class=\"content\"><h1>Course Title</h1><p>This is a <em>great</em> course about &amp; HTML entities like &nbsp; and &quot;quotes&quot;.</p><ul><li>Item 1</li><li>Item 2</li></ul></div>"

	for b.Loop() {
		cleaner.CleanHTML(input)
	}
}

// BenchmarkHTMLCleaner_CleanHTML_Large benchmarks the CleanHTML method with large content.
func BenchmarkHTMLCleaner_CleanHTML_Large(b *testing.B) {
	cleaner := NewHTMLCleaner()

	// Create a large HTML string
	var builder strings.Builder
	for i := range 100 {
		builder.WriteString("<p>Paragraph ")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(" with some content &amp; entities &lt;test&gt;.</p>")
	}
	input := builder.String()

	for b.Loop() {
		cleaner.CleanHTML(input)
	}
}
