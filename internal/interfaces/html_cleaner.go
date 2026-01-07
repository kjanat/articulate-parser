package interfaces

// HTMLCleaner defines the interface for cleaning HTML content.
// Implementations convert HTML to plain text, stripping tags and decoding entities.
type HTMLCleaner interface {
	// CleanHTML removes HTML tags and decodes HTML entities from a string.
	// Returns plain text suitable for non-HTML output formats.
	CleanHTML(htmlStr string) string

	// CleanAndTrim removes HTML tags, decodes entities, and trims whitespace.
	// Convenience method combining CleanHTML with strings.TrimSpace.
	CleanAndTrim(htmlStr string) string
}
