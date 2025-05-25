// Package main provides the entry point for the articulate-parser application.
// This application fetches Articulate Rise courses from URLs or local files and
// exports them to different formats such as Markdown or DOCX.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kjanat/articulate-parser/internal/exporters"
	"github.com/kjanat/articulate-parser/internal/services"
)

// main is the entry point of the application.
// It handles command-line arguments, sets up dependencies,
// and coordinates the parsing and exporting of courses.
func main() {
	// Dependency injection setup
	htmlCleaner := services.NewHTMLCleaner()
	parser := services.NewArticulateParser()
	exporterFactory := exporters.NewFactory(htmlCleaner)
	app := services.NewApp(parser, exporterFactory)

	// Check for required command-line arguments
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <source> <format> <output>\n", os.Args[0])
		fmt.Printf("  source: URI or file path to the course\n")
		fmt.Printf("  format: export format (%s)\n", joinStrings(app.GetSupportedFormats(), ", "))
		fmt.Printf("  output: output file path\n")
		fmt.Println("\nExample:")
		fmt.Printf("  %s articulate-sample.json markdown output.md\n", os.Args[0])
		fmt.Printf("  %s https://rise.articulate.com/share/xyz docx output.docx\n", os.Args[0])
		os.Exit(1)
	}

	source := os.Args[1]
	format := os.Args[2]
	output := os.Args[3]

	var err error

	// Determine if source is a URI or file path
	if isURI(source) {
		err = app.ProcessCourseFromURI(source, format, output)
	} else {
		err = app.ProcessCourseFromFile(source, format, output)
	}

	if err != nil {
		log.Fatalf("Error processing course: %v", err)
	}

	fmt.Printf("Successfully exported course to %s\n", output)
}

// isURI checks if a string is a URI by looking for http:// or https:// prefixes.
//
// Parameters:
//   - str: The string to check
//
// Returns:
//   - true if the string appears to be a URI, false otherwise
func isURI(str string) bool {
	return len(str) > 7 && (str[:7] == "http://" || str[:8] == "https://")
}

// joinStrings concatenates a slice of strings using the specified separator.
//
// Parameters:
//   - strs: The slice of strings to join
//   - sep: The separator to insert between each string
//
// Returns:
//   - A single string with all elements joined by the separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
