// Package main_test provides tests for the main package utility functions.
package main

import (
	"testing"
)

// TestIsURI tests the isURI function with various input scenarios.
func TestIsURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid HTTP URI",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "valid HTTPS URI",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "valid Articulate Rise URI",
			input:    "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/",
			expected: true,
		},
		{
			name:     "local file path",
			input:    "C:\\Users\\test\\file.json",
			expected: false,
		},
		{
			name:     "relative file path",
			input:    "./sample.json",
			expected: false,
		},
		{
			name:     "filename only",
			input:    "sample.json",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "short string",
			input:    "http",
			expected: false,
		},
		{
			name:     "malformed URI",
			input:    "htp://example.com",
			expected: false,
		},
		{
			name:     "FTP URI",
			input:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "HTTP with extra characters",
			input:    "xhttp://example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isURI(tt.input)
			if result != tt.expected {
				t.Errorf("isURI(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestJoinStrings tests the joinStrings function with various input scenarios.
func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name      string
		strs      []string
		separator string
		expected  string
	}{
		{
			name:      "empty slice",
			strs:      []string{},
			separator: ", ",
			expected:  "",
		},
		{
			name:      "single string",
			strs:      []string{"hello"},
			separator: ", ",
			expected:  "hello",
		},
		{
			name:      "two strings with comma separator",
			strs:      []string{"markdown", "docx"},
			separator: ", ",
			expected:  "markdown, docx",
		},
		{
			name:      "three strings with comma separator",
			strs:      []string{"markdown", "md", "docx"},
			separator: ", ",
			expected:  "markdown, md, docx",
		},
		{
			name:      "multiple strings with pipe separator",
			strs:      []string{"option1", "option2", "option3"},
			separator: " | ",
			expected:  "option1 | option2 | option3",
		},
		{
			name:      "strings with no separator",
			strs:      []string{"a", "b", "c"},
			separator: "",
			expected:  "abc",
		},
		{
			name:      "strings with newline separator",
			strs:      []string{"line1", "line2", "line3"},
			separator: "\n",
			expected:  "line1\nline2\nline3",
		},
		{
			name:      "empty strings in slice",
			strs:      []string{"", "middle", ""},
			separator: "-",
			expected:  "-middle-",
		},
		{
			name:      "nil slice",
			strs:      nil,
			separator: ", ",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinStrings(tt.strs, tt.separator)
			if result != tt.expected {
				t.Errorf("joinStrings(%v, %q) = %q, want %q", tt.strs, tt.separator, result, tt.expected)
			}
		})
	}
}

// BenchmarkIsURI benchmarks the isURI function performance.
func BenchmarkIsURI(b *testing.B) {
	testStr := "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isURI(testStr)
	}
}

// BenchmarkJoinStrings benchmarks the joinStrings function performance.
func BenchmarkJoinStrings(b *testing.B) {
	strs := []string{"markdown", "md", "docx", "word", "pdf", "html"}
	separator := ", "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinStrings(strs, separator)
	}
}
