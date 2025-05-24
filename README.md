# Articulate Rise Parser

A Go-based parser that converts Articulate Rise e-learning content to various formats including Markdown and Word documents.

[![Go version](https://img.shields.io/github/go-mod/go-version/kjanat/articulate-parser?logo=Go&logoColor=white)][gomod]
[![Go Doc](https://godoc.org/github.com/kjanat/articulate-parser?status.svg)][Package documentation]
[![Go Report Card](https://goreportcard.com/badge/github.com/kjanat/articulate-parser)][Go report]
[![Tag](https://img.shields.io/github/v/tag/kjanat/articulate-parser?sort=semver&label=Tag)][Tags]
[![Release Date](https://img.shields.io/github/release-date/kjanat/articulate-parser?label=Release%20date)][Latest release]
[![License](https://img.shields.io/github/license/kjanat/articulate-parser?label=License)](LICENSE)
[![Commit activity](https://img.shields.io/github/commit-activity/m/kjanat/articulate-parser?label=Commit%20activity)][Commits]
[![Last commit](https://img.shields.io/github/last-commit/kjanat/articulate-parser?label=Last%20commit)][Commits]
[![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/kjanat/articulate-parser?label=Issues)][Issues]
[![CI](https://img.shields.io/github/actions/workflow/status/kjanat/articulate-parser/ci.yml?logo=github&label=CI)][Build]
[![Codecov](https://img.shields.io/codecov/c/gh/kjanat/articulate-parser?token=eHhaHY8nut&logo=codecov&logoColor=%23F01F7A&label=Codecov)][Codecov]

## Features

-   Parse Articulate Rise JSON data from URLs or local files
-   Export to Markdown (.md) format
-   Export to Word Document (.docx) format
-   Support for various content types:
    -   Text content with headings and paragraphs
    -   Lists and bullet points
    -   Multimedia content (videos and images)
    -   Knowledge checks and quizzes
    -   Interactive content (flashcards)
    -   Course structure and metadata

## Installation

1.  Ensure you have Go 1.21 or later installed
2.  Clone or download the parser code
3.  Initialize the Go module:

```bash
go mod init articulate-parser
go mod tidy
```

## Dependencies

The parser uses the following external library:

-   `github.com/unidoc/unioffice` - For creating Word documents

## Usage

### Command Line Interface

```bash
go run main.go <input_uri_or_file> <output_format> [output_path]
```

#### Parameters

-   `input_uri_or_file`: Either an Articulate Rise share URL or path to a local JSON file
-   `output_format`: `md` for Markdown or `docx` for Word Document
-   `output_path`: Optional. If not provided, files are saved to `./output/` directory

#### Examples

1.  **Parse from URL and export to Markdown:**

```bash
go run main.go "https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/" md
```

2.  **Parse from local file and export to Word:**

```bash
go run main.go "articulate-sample.json" docx "my-course.docx"
```

3.  **Parse from local file and export to Markdown:**

```bash
go run main.go "C:\Users\kjana\Projects\articulate-parser\articulate-sample.json" md
```

### Building the Executable

To build a standalone executable:

```bash
go build -o articulate-parser main.go
```

Then run:

```bash
./articulate-parser input.json md output.md
```

## Output Formats

### Markdown (.md)

-   Hierarchical structure with proper heading levels
-   Clean text content with HTML tags removed
-   Lists and bullet points preserved
-   Quiz questions with correct answers marked
-   Media references included
-   Course metadata at the top

### Word Document (.docx)

-   Professional document formatting
-   Bold headings and proper typography
-   Bulleted lists
-   Quiz questions with answers
-   Media content references
-   Maintains course structure

## Supported Content Types

The parser handles the following Articulate Rise content types:

-   **Text blocks**: Headings and paragraphs
-   **Lists**: Bullet points and numbered lists
-   **Multimedia**: Videos and images (references only)
-   **Knowledge Checks**: Multiple choice, multiple response, fill-in-the-blank, matching
-   **Interactive Content**: Flashcards and interactive scenarios
-   **Dividers**: Section breaks
-   **Sections**: Course organization

## Data Structure

The parser works with the standard Articulate Rise JSON format which includes:

-   Course metadata (title, description, settings)
-   Lesson structure
-   Content items with various types
-   Media references
-   Quiz/assessment data
-   Styling and layout information

## URL Pattern Recognition

The parser automatically extracts share IDs from Articulate Rise URLs:

-   Input: `https://rise.articulate.com/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO#/`
-   API URL: `https://rise.articulate.com/api/rise-runtime/boot/share/N_APNg40Vr2CSH2xNz-ZLATM5kNviDIO`

## Error Handling

The parser includes error handling for:

-   Invalid URLs or share IDs
-   Network connection issues
-   Malformed JSON data
-   File I/O errors
-   Unsupported content types

<!-- ## Code coverage

![Sunburst](https://codecov.io/gh/kjanat/articulate-parser/graphs/tree.svg?token=eHhaHY8nut)

![Grid](https://codecov.io/gh/kjanat/articulate-parser/graphs/tree.svg?token=eHhaHY8nut)

![Icicle](https://codecov.io/gh/kjanat/articulate-parser/graphs/icicle.svg?token=eHhaHY8nut) -->

## Limitations

-   Media files (videos, images) are referenced but not downloaded
-   Complex interactive elements may be simplified in export
-   Styling and visual formatting is not preserved
-   Assessment logic and interactivity is lost in static exports

## Future Enhancements

Potential improvements could include:

-   PDF export support
-   Media file downloading
-   HTML export with preserved styling
-   SCORM package support
-   Batch processing capabilities
-   Custom template support

## License

This is a utility tool for educational content conversion. Please ensure you have appropriate rights to the Articulate Rise content you're parsing.

[Build]: https://github.com/kjanat/articulate-parser/actions/workflows/ci.yml
[Codecov]: https://codecov.io/gh/kjanat/articulate-parser
[Commits]: https://github.com/kjanat/articulate-parser/commits/master/
[Go report]: https://goreportcard.com/report/github.com/kjanat/articulate-parser
[gomod]: go.mod
[Issues]: https://github.com/kjanat/articulate-parser/issues
[Latest release]: https://github.com/kjanat/articulate-parser/releases/latest
[Package documentation]: https://godoc.org/github.com/kjanat/articulate-parser
[Tags]: https://github.com/kjanat/articulate-parser/tags
