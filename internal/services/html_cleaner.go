// Package services provides the core functionality for the articulate-parser application.
// It implements the interfaces defined in the interfaces package.
package services

import (
	"regexp"
	"strings"
)

var (
	// htmlTagRegex matches HTML tags for removal
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	// whitespaceRegex matches multiple whitespace characters for normalization
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// HTMLCleaner provides utilities for converting HTML content to plain text.
// It removes HTML tags while preserving their content and converts HTML entities
// to their plain text equivalents.
type HTMLCleaner struct{}

// NewHTMLCleaner creates a new HTML cleaner instance.
// This service is typically injected into exporters that need to handle
// HTML content from Articulate Rise courses.
func NewHTMLCleaner() *HTMLCleaner {
	return &HTMLCleaner{}
}

// CleanHTML removes HTML tags and converts entities, returning clean plain text.
// The function preserves the textual content of the HTML while removing markup.
// It handles common HTML entities like &nbsp;, &amp;, etc., and normalizes whitespace.
//
// Parameters:
//   - html: The HTML content to clean
//
// Returns:
//   - A plain text string with all HTML elements and entities removed/converted
func (h *HTMLCleaner) CleanHTML(html string) string {
	// Remove HTML tags but preserve content
	cleaned := htmlTagRegex.ReplaceAllString(html, "")

	// Replace common HTML entities with their character equivalents
	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")
	cleaned = strings.ReplaceAll(cleaned, "&#39;", "'")
	cleaned = strings.ReplaceAll(cleaned, "&iuml;", "ï")
	cleaned = strings.ReplaceAll(cleaned, "&euml;", "ë")
	cleaned = strings.ReplaceAll(cleaned, "&eacute;", "é")

	// Clean up extra whitespace by replacing multiple spaces, tabs, and newlines
	// with a single space, then trim any leading/trailing whitespace
	cleaned = whitespaceRegex.ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
