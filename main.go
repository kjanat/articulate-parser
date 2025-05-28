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
	"github.com/kjanat/articulate-parser/internal/version"
)

// main is the entry point of the application.
// It handles command-line arguments, sets up dependencies,
// and coordinates the parsing and exporting of courses.
func main() {
	os.Exit(run(os.Args))
}

// run contains the main application logic and returns an exit code.
// This function is testable as it doesn't call os.Exit directly.
func run(args []string) int {
	// Dependency injection setup
	htmlCleaner := services.NewHTMLCleaner()
	parser := services.NewArticulateParser()
	exporterFactory := exporters.NewFactory(htmlCleaner)
	app := services.NewApp(parser, exporterFactory)

	// Check for version flag
	if len(args) > 1 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Printf("%s version %s\n", args[0], version.Version)
		fmt.Printf("Build time: %s\n", version.BuildTime)
		fmt.Printf("Git commit: %s\n", version.GitCommit)
		return 0
	}

	// Check for help flag
	if len(args) > 1 && (args[1] == "--help" || args[1] == "-h" || args[1] == "help") {
		printUsage(args[0], app.GetSupportedFormats())
		return 0
	}

	// Check for required command-line arguments
	if len(args) < 4 {
		printUsage(args[0], app.GetSupportedFormats())
		return 1
	}

	source := args[1]
	format := args[2]
	output := args[3]

	var err error

	// Determine if source is a URI or file path
	if isURI(source) {
		err = app.ProcessCourseFromURI(source, format, output)
	} else {
		err = app.ProcessCourseFromFile(source, format, output)
	}

	if err != nil {
		log.Printf("Error processing course: %v", err)
		return 1
	}

	fmt.Printf("Successfully exported course to %s\n", output)
	return 0
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

// printUsage prints the command-line usage information.
//
// Parameters:
//   - programName: The name of the program (args[0])
//   - supportedFormats: Slice of supported export formats
func printUsage(programName string, supportedFormats []string) {
	fmt.Printf("Usage: %s <source> <format> <output>\n", programName)
	fmt.Printf("  source: URI or file path to the course\n")
	fmt.Printf("  format: export format (%s)\n", joinStrings(supportedFormats, ", "))
	fmt.Printf("  output: output file path\n")
	fmt.Println("\nExample:")
	fmt.Printf("  %s articulate-sample.json markdown output.md\n", programName)
	fmt.Printf("  %s https://rise.articulate.com/share/xyz docx output.docx\n", programName)
}
