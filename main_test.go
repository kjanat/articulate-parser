// Package main_test provides tests for the main package utility functions.
package main

import (
	"bytes"
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

// BenchmarkIsURI benchmarks the isURI function performance.
func BenchmarkIsURI(b *testing.B) {
	testStr := "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/"

	for b.Loop() {
		isURI(testStr)
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

			// Restore stdout. Close errors are ignored: we've already captured the
			// output before closing, and any close error doesn't affect test validity.
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output. Copy errors are ignored: in this test context,
			// reading from a pipe that was just closed is not expected to fail, and
			// we're verifying the captured output regardless.
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
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

			// Restore stdout. Close errors are ignored: the pipe write end is already
			// closed before reading, and any close error doesn't affect the test.
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output. Copy errors are ignored: we successfully wrote
			// the help output to the pipe and can verify it regardless of close semantics.
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
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

// TestRunWithVersionFlags tests the run function with version flag arguments.
func TestRunWithVersionFlags(t *testing.T) {
	versionFlags := []string{"--version", "-v"}

	for _, flag := range versionFlags {
		t.Run("version_flag_"+flag, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run with version flag
			args := []string{"articulate-parser", flag}
			exitCode := run(args)

			// Restore stdout. Close errors are ignored: the version output has already
			// been written and we're about to read it; close semantics don't affect correctness.
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output. Copy errors are ignored: the output was successfully
			// produced and we can verify its contents regardless of any I/O edge cases.
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			// Verify exit code is 0 (success)
			if exitCode != 0 {
				t.Errorf("Expected exit code 0 for version flag %s, got %d", flag, exitCode)
			}

			// Verify version content is displayed
			expectedContent := []string{
				"articulate-parser version",
				"Build time:",
				"Git commit:",
			}

			for _, expected := range expectedContent {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected version output to contain %q when using flag %s, got: %s", expected, flag, output)
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

	// Restore stdout/stderr and log output. Close errors are ignored: we've already
	// written all error messages to these pipes before closing them, and the test
	// only cares about verifying the captured output.
	_ = stdoutW.Close()
	_ = stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	log.SetOutput(oldLogOutput)

	// Read captured output. Copy errors are ignored: the error messages have been
	// successfully written to the pipes, and we can verify the output content
	// regardless of any edge cases in pipe closure or I/O completion.
	var stdoutBuf, stderrBuf bytes.Buffer
	_, _ = io.Copy(&stdoutBuf, stdoutR)
	_, _ = io.Copy(&stderrBuf, stderrR)

	// Close read ends of pipes. Errors ignored: we've already consumed all data
	// from these pipes, and close errors don't affect test assertions.
	_ = stdoutR.Close()
	_ = stderrR.Close()

	// Verify exit code
	if exitCode != 1 {
		t.Errorf("Expected exit code 1 for non-existent file, got %d", exitCode)
	}

	// Should have error output in structured log format
	output := stdoutBuf.String()
	if !strings.Contains(output, "level=ERROR") && !strings.Contains(output, "failed to process course") {
		t.Errorf("Expected error message about processing course, got: %s", output)
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

	// Restore stdout/stderr and log output. Close errors are ignored: we've already
	// written all error messages about the invalid URI to these pipes before closing,
	// and test correctness only depends on verifying the captured error output.
	_ = stdoutW.Close()
	_ = stderrW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	log.SetOutput(oldLogOutput)

	// Read captured output. Copy errors are ignored: the error messages have been
	// successfully written and we can verify the failure output content regardless
	// of any edge cases in pipe lifecycle or I/O synchronization.
	var stdoutBuf, stderrBuf bytes.Buffer
	_, _ = io.Copy(&stdoutBuf, stdoutR)
	_, _ = io.Copy(&stderrBuf, stderrR)

	// Close read ends of pipes. Errors ignored: we've already consumed all data
	// and close errors don't affect the validation of the error output.
	_ = stdoutR.Close()
	_ = stderrR.Close()

	// Should fail because the URI is invalid/unreachable
	if exitCode != 1 {
		t.Errorf("Expected failure (exit code 1) for invalid URI, got %d", exitCode)
	}

	// Should have error output in structured log format
	output := stdoutBuf.String()
	if !strings.Contains(output, "level=ERROR") && !strings.Contains(output, "failed to process course") {
		t.Errorf("Expected error message about processing course, got: %s", output)
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
	// Ensure temporary test file is cleaned up. Remove errors are ignored because
	// the test has already used the file for its purpose, and cleanup failures don't
	// invalidate the test results (the OS will eventually clean up temp files).
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	// Close the temporary file. Errors are ignored because we've already written
	// the test content and the main test logic (loading the file) doesn't depend
	// on the success of closing this file descriptor.
	_ = tmpFile.Close()

	// Test successful run with valid file
	outputFile := "test-output.md"
	// Ensure test output file is cleaned up. Remove errors are ignored because the
	// test has already verified the export succeeded; cleanup failures don't affect
	// the test assertions.
	defer func() {
		_ = os.Remove(outputFile)
	}()

	// Save original stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	// Capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{"articulate-parser", tmpFile.Name(), "markdown", outputFile}
	exitCode := run(args)

	// Close write end and restore stdout. Close errors are ignored: we've already
	// written the success message before closing, and any close error doesn't affect
	// the validity of the captured output or the test assertions.
	_ = w.Close()
	os.Stdout = originalStdout

	// Read captured output. Copy errors are ignored: the success message was
	// successfully written to the pipe, and we can verify it regardless of any
	// edge cases in pipe closure or I/O synchronization.
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Verify successful execution
	if exitCode != 0 {
		t.Errorf("Expected successful execution (exit code 0), got %d", exitCode)
	}

	// Verify success message in structured log format
	if !strings.Contains(output, "level=INFO") || !strings.Contains(output, "successfully exported course") {
		t.Errorf("Expected success message in output, got: %s", output)
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

			// Restore stdout. Close errors are ignored: the export success message
			// has already been written and we're about to read it; close semantics
			// don't affect the validity of the captured output.
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output. Copy errors are ignored: the output was successfully
			// produced and we can verify its contents regardless of any I/O edge cases.
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			// Clean up test file. Remove errors are ignored because the test has
			// already verified the export succeeded; cleanup failures don't affect
			// the test assertions.
			defer func() {
				_ = os.Remove(format.output)
			}()

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
