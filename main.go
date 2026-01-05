// Package main provides the entry point for the articulate-parser application.
// This application fetches Articulate Rise courses from URLs or local files and
// exports them to different formats such as Markdown or DOCX.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kjanat/articulate-parser/internal/config"
	"github.com/kjanat/articulate-parser/internal/exporters"
	"github.com/kjanat/articulate-parser/internal/interfaces"
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
	// Load configuration
	cfg := config.Load()

	// Dependency injection setup with configuration
	var logger interfaces.Logger
	if cfg.LogFormat == "json" {
		logger = services.NewSlogLogger(cfg.LogLevel)
	} else {
		logger = services.NewTextLogger(cfg.LogLevel)
	}

	htmlCleaner := services.NewHTMLCleaner()
	parser := services.NewArticulateParser(logger, cfg.BaseURL, cfg.RequestTimeout)
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
		printUsage(args[0], app.SupportedFormats())
		return 0
	}

	// Check for required command-line arguments
	if len(args) < 4 {
		printUsage(args[0], app.SupportedFormats())
		return 1
	}

	source := args[1]
	format := args[2]
	output := args[3]

	var err error

	// Determine if source is a URI or file path
	if isURI(source) {
		err = app.ProcessCourseFromURI(context.Background(), source, format, output)
	} else {
		err = app.ProcessCourseFromFile(source, format, output)
	}

	if err != nil {
		logger.Error("failed to process course", "error", err, "source", source)
		return 1
	}

	logger.Info("successfully exported course", "output", output, "format", format)
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
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

// printUsage prints the command-line usage information.
//
// Parameters:
//   - programName: The name of the program (args[0])
//   - supportedFormats: Slice of supported export formats
func printUsage(programName string, supportedFormats []string) {
	fmt.Printf("Usage: %s <source> <format> <output>\n", programName)
	fmt.Printf("  source: URI or file path to the course\n")
	fmt.Printf("  format: export format (%s)\n", strings.Join(supportedFormats, ", "))
	fmt.Printf("  output: output file path\n")
	fmt.Println("\nExample:")
	fmt.Printf("  %s articulate-sample.json markdown output.md\n", programName)
	fmt.Printf("  %s https://rise.articulate.com/share/xyz docx output.docx\n", programName)
}
