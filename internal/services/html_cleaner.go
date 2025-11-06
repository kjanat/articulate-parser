// Package services provides the core functionality for the articulate-parser application.
// It implements the interfaces defined in the interfaces package.
package services

import (
	"bytes"
	stdhtml "html"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// HTMLCleaner provides utilities for converting HTML content to plain text.
// It removes HTML tags while preserving their content and converts HTML entities
// to their plain text equivalents using proper HTML parsing instead of regex.
type HTMLCleaner struct{}

// NewHTMLCleaner creates a new HTML cleaner instance.
// This service is typically injected into exporters that need to handle
// HTML content from Articulate Rise courses.
func NewHTMLCleaner() *HTMLCleaner {
	return &HTMLCleaner{}
}

// CleanHTML removes HTML tags and converts entities, returning clean plain text.
// It parses the HTML into a node tree and extracts only text content,
// skipping script and style tags. HTML entities are automatically handled
// by the parser, and whitespace is normalized.
func (h *HTMLCleaner) CleanHTML(htmlStr string) string {
	// Parse the HTML into a node tree
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		// If parsing fails, return empty string
		// This maintains backward compatibility with the test expectations
		return ""
	}

	// Extract text content from the node tree
	var buf bytes.Buffer
	extractText(&buf, doc)

	// Unescape any remaining HTML entities
	unescaped := stdhtml.UnescapeString(buf.String())

	// Normalize whitespace: replace multiple spaces, tabs, and newlines with a single space
	cleaned := strings.Join(strings.Fields(unescaped), " ")
	return strings.TrimSpace(cleaned)
}

// extractText recursively traverses the HTML node tree and extracts text content.
// It skips script and style tags to avoid including their content in the output.
func extractText(w io.Writer, n *html.Node) {
	// Skip script and style tags entirely
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		return
	}

	// If this is a text node, write its content
	if n.Type == html.TextNode {
		// Write errors are ignored because we're writing to an in-memory buffer
		// which cannot fail in normal circumstances
		_, _ = w.Write([]byte(n.Data))
	}

	// Recursively process all child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(w, c)
	}
}
