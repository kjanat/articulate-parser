// Package main_test provides tests for the main package utility functions.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

// TestRunWithInsufficientArgs tests the run function with insufficient command-line arguments.
func TestRunWithInsufficientArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{"articulate-parser"},
		},
		{
			name: "one argument",
			args: []string{"articulate-parser", "source"},
		},
		{
			name: "two arguments",
			args: []string{"articulate-parser", "source", "format"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the function
			exitCode := run(tt.args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify exit code
			if exitCode != 1 {
				t.Errorf("Expected exit code 1, got %d", exitCode)
			}

			// Verify usage message is displayed
			if !strings.Contains(output, "Usage:") {
				t.Errorf("Expected usage message in output, got: %s", output)
			}

			if !strings.Contains(output, "export format") {
				t.Errorf("Expected format information in output, got: %s", output)
			}
		})
	}
}

// TestRunWithHelpFlags tests the run function with help flag arguments.
func TestRunWithHelpFlags(t *testing.T) {
	helpFlags := []string{"--help", "-h", "help"}
	
	for _, flag := range helpFlags {
		t.Run("help_flag_"+flag, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run with help flag
			args := []string{"articulate-parser", flag}
			exitCode := run(args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify exit code is 0 (success)
			if exitCode != 0 {
				t.Errorf("Expected exit code 0 for help flag %s, got %d", flag, exitCode)
			}

			// Verify help content is displayed
			expectedContent := []string{
				"Usage:",
				"source: URI or file path to the course",
				"format: export format",
				"output: output file path",
				"Example:",
				"articulate-sample.json markdown output.md",
				"https://rise.articulate.com/share/xyz docx output.docx",
			}

			for _, expected := range expectedContent {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected help output to contain %q when using flag %s, got: %s", expected, flag, output)
				}
			}
		})
	}
}

// TestRunWithInvalidFile tests the run function with a non-existent file.
func TestRunWithInvalidFile(t *testing.T) {
	// Capture stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()

	os.Stdout = stdoutW
	os.Stderr = stderrW

	// Also need to redirect log output
	oldLogOutput := log.Writer()
	log.SetOutput(stderrW)

	// Run with non-existent file
	args := []string{"articulate-parser", "nonexistent-file.json", "markdown", "output.md"}
	exitCode := run(args)

	// Restore stdout/stderr and log output
	stdoutW.Close()
	stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	log.SetOutput(oldLogOutput)

	// Read captured output
	var stdoutBuf, stderrBuf bytes.Buffer
	io.Copy(&stdoutBuf, stdoutR)
	io.Copy(&stderrBuf, stderrR)

	stdoutR.Close()
	stderrR.Close()

	// Verify exit code
	if exitCode != 1 {
		t.Errorf("Expected exit code 1 for non-existent file, got %d", exitCode)
	}

	// Should have error output
	errorOutput := stderrBuf.String()
	if !strings.Contains(errorOutput, "Error processing course") {
		t.Errorf("Expected error message about processing course, got: %s", errorOutput)
	}
}

// TestRunWithInvalidURI tests the run function with an invalid URI.
func TestRunWithInvalidURI(t *testing.T) {
	// Capture stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()

	os.Stdout = stdoutW
	os.Stderr = stderrW

	// Also need to redirect log output
	oldLogOutput := log.Writer()
	log.SetOutput(stderrW)

	// Run with invalid URI (will fail because we can't actually fetch)
	args := []string{"articulate-parser", "https://example.com/invalid", "markdown", "output.md"}
	exitCode := run(args)

	// Restore stdout/stderr and log output
	stdoutW.Close()
	stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	log.SetOutput(oldLogOutput)

	// Read captured output
	var stdoutBuf, stderrBuf bytes.Buffer
	io.Copy(&stdoutBuf, stdoutR)
	io.Copy(&stderrBuf, stderrR)

	stdoutR.Close()
	stderrR.Close()

	// Should fail because the URI is invalid/unreachable
	if exitCode != 1 {
		t.Errorf("Expected failure (exit code 1) for invalid URI, got %d", exitCode)
	}

	// Should have error output
	errorOutput := stderrBuf.String()
	if !strings.Contains(errorOutput, "Error processing course") {
		t.Errorf("Expected error message about processing course, got: %s", errorOutput)
	}
}

// TestRunWithValidJSONFile tests the run function with a valid JSON file.
func TestRunWithValidJSONFile(t *testing.T) {
	// Create a temporary test JSON file
	testContent := `{
		"title": "Test Course",
		"lessons": [
			{
				"id": "lesson1",
				"title": "Test Lesson",
				"blocks": [
					{
						"type": "text",
						"id": "block1",
						"data": {
							"text": "Test content"
						}
					}
				]
			}
		]
	}`

	tmpFile, err := os.CreateTemp("", "test-course-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	tmpFile.Close()

	// Test successful run with valid file
	outputFile := "test-output.md"
	defer os.Remove(outputFile)

	// Save original stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	// Capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"articulate-parser", tmpFile.Name(), "markdown", outputFile}
	exitCode := run(args)

	// Close write end and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify successful execution
	if exitCode != 0 {
		t.Errorf("Expected successful execution (exit code 0), got %d", exitCode)
	}

	// Verify success message
	expectedMsg := fmt.Sprintf("Successfully exported course to %s", outputFile)
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("Expected success message '%s' in output, got: %s", expectedMsg, output)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to be created", outputFile)
	}
}

// TestRunIntegration tests the run function with different output formats using sample file.
func TestRunIntegration(t *testing.T) {
	// Skip if sample file doesn't exist
	if _, err := os.Stat("articulate-sample.json"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: articulate-sample.json not found")
	}

	formats := []struct {
		format string
		output string
	}{
		{"markdown", "test-output.md"},
		{"html", "test-output.html"},
		{"docx", "test-output.docx"},
	}

	for _, format := range formats {
		t.Run("format_"+format.format, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the function
			args := []string{"articulate-parser", "articulate-sample.json", format.format, format.output}
			exitCode := run(args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Clean up test file
			defer os.Remove(format.output)

			// Verify successful execution
			if exitCode != 0 {
				t.Errorf("Expected successful execution (exit code 0), got %d", exitCode)
			}

			// Verify success message
			expectedMsg := "Successfully exported course to " + format.output
			if !strings.Contains(output, expectedMsg) {
				t.Errorf("Expected success message '%s' in output, got: %s", expectedMsg, output)
			}

			// Verify output file was created
			if _, err := os.Stat(format.output); os.IsNotExist(err) {
				t.Errorf("Expected output file %s to be created", format.output)
			}
		})
	}
}

// TestMainFunction tests that the main function exists and is properly structured.
// We can't test os.Exit behavior directly, but we can verify the main function
// calls the run function correctly by testing run function behavior.
func TestMainFunction(t *testing.T) {
	// Test that insufficient args return exit code 1
	exitCode := run([]string{"program"})
	if exitCode != 1 {
		t.Errorf("Expected run to return exit code 1 for insufficient args, got %d", exitCode)
	}

	// Test that main function exists (this is mainly for coverage)
	// The main function just calls os.Exit(run(os.Args)), which we can't test directly
	// but we've tested the run function thoroughly above.
}
