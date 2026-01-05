# Agent Guidelines for articulate-parser

A Go CLI tool that parses Articulate Rise courses from URLs or local JSON files and exports them to Markdown, HTML, or DOCX formats.

## Repository Info

- **GitHub**: https://github.com/kjanat/articulate-parser
- **Default branch**: `master` (not `main`)

## Build/Test Commands

### Primary Commands (using Taskfile)

```bash
task build              # Build binary to bin/articulate-parser
task test               # Run all tests with race detection
task lint               # Run all linters (vet, fmt, staticcheck, golangci-lint)
task fmt                # Format all Go files
task ci                 # Full CI pipeline: deps, lint, test with coverage, build
task qa                 # Quick QA: fmt + lint + test
```

### Direct Go Commands

```bash
# Build
go build -o bin/articulate-parser main.go

# Run all tests
go test -race -timeout 5m ./...

# Run single test by name
go test -v -race -run ^TestMarkdownExporter_Export$ ./internal/exporters

# Run tests in specific package
go test -v -race ./internal/services

# Run tests matching pattern
go test -v -race -run "TestParser" ./...

# Test with coverage
go test -race -coverprofile=coverage/coverage.out -covermode=atomic ./...
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Benchmarks
go test -bench=. -benchmem ./...
go test -bench=BenchmarkMarkdownExporter ./internal/exporters
```

### Security & Auditing

```bash
task security:check     # Run gosec security scanner
task security:audit     # Run govulncheck for vulnerabilities
```

## Code Style Guidelines

### Imports

- Use `goimports` with local prefix: `github.com/kjanat/articulate-parser`
- Order: stdlib, blank line, external packages, blank line, internal packages

```go
import (
    "context"
    "fmt"

    "github.com/fumiama/go-docx"

    "github.com/kjanat/articulate-parser/internal/interfaces"
)
```

### Formatting

- Use `gofmt -s` (simplify) and `gofumpt` with extra rules
- Function length: max 100 lines, 50 statements
- Cyclomatic complexity: max 15; Cognitive complexity: max 20

### Types & Naming

- Use interface-based design (see `internal/interfaces/`)
- Exported types/functions require godoc comments ending with period
- Use descriptive names: `ArticulateParser`, `MarkdownExporter`
- Receiver names: short (1-2 chars), consistent per type

### Error Handling

- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use `%w` verb for error wrapping to preserve error chain
- Check all error returns (enforced by `errcheck`)
- Document error handling rationale in defer blocks when ignoring close errors

```go
// Good: Error wrapping with context
if err := json.Unmarshal(body, &course); err != nil {
    return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
}

// Good: Documented defer with error handling
defer func() {
    if err := resp.Body.Close(); err != nil {
        p.Logger.Warn("failed to close response body", "error", err)
    }
}()
```

### Comments

- All exported types/functions require godoc comments
- End sentences with periods (`godot` linter enforced)
- Mark known issues with TODO/FIXME/HACK/BUG/XXX

### Security

- Use `#nosec` with justification for deliberate security exceptions
- G304: File paths from CLI args; G306: Export file permissions

```go
// #nosec G304 - File path provided by user via CLI argument
data, err := os.ReadFile(filePath)
```

### Testing

- Enable race detection: `-race` flag always
- Use table-driven tests where applicable
- Mark test helpers with `t.Helper()`
- Use `t.TempDir()` for temporary files
- Benchmarks in `*_bench_test.go`, examples in `*_example_test.go`
- Test naming: `Test<Type>_<Method>` or `Test<Function>`

```go
func TestMarkdownExporter_ProcessItemToMarkdown_AllTypes(t *testing.T) {
    tests := []struct {
        name, itemType, expectedText string
    }{
        {"text item", "text", ""},
        {"divider item", "divider", "---"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Dependencies

- Minimal external dependencies (go-docx, golang.org/x/net, golang.org/x/text)
- Run `task deps:tidy` after adding/removing dependencies
- CGO disabled by default (`CGO_ENABLED=0`)

## Project Structure

```
articulate-parser/
  internal/
    config/              # Configuration loading
    exporters/           # Export implementations (markdown, html, docx)
    interfaces/          # Core interfaces (Exporter, CourseParser, Logger)
    models/              # Data models (Course, Lesson, Item, Media)
    services/            # Core services (parser, html cleaner, app, logger)
    version/             # Version information
  main.go                # Application entry point
```

## Common Patterns

### Creating a new exporter

1. Implement `interfaces.Exporter` interface
2. Add factory method to `internal/exporters/factory.go`
3. Register format in `NewFactory()`
4. Add tests following existing patterns

### Adding configuration options

1. Add field to `Config` struct in `internal/config/config.go`
2. Load from environment variable with sensible default
3. Document in config struct comments
