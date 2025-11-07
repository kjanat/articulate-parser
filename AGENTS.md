# Agent Guidelines for articulate-parser

## Build/Test Commands
- **Build**: `task build` or `go build -o bin/articulate-parser main.go`
- **Run tests**: `task test` or `go test -race -timeout 5m ./...`
- **Run single test**: `go test -v -race -run ^TestName$ ./path/to/package`
- **Test with coverage**: 
  - `task test:coverage` or 
  - `go test -race -coverprofile=coverage/coverage.out -covermode=atomic ./...`
- **Lint**: `task lint` (runs vet, fmt check, staticcheck, golangci-lint)
- **Format**: `task fmt` or `gofmt -s -w .`
- **CI checks**: `task ci` (deps, lint, test with coverage, build)

## Code Style Guidelines

### Imports
- Use `goimports` with local prefix: `github.com/kjanat/articulate-parser`
- Order: stdlib, external, internal packages
- Group related imports together

### Formatting
- Use `gofmt -s` (simplify) and `gofumpt` with extra rules
- Function length: max 100 lines, 50 statements
- Cyclomatic complexity: max 15
- Cognitive complexity: max 20

### Types & Naming
- Use interface-based design (see `internal/interfaces/`)
- Export types/functions with clear godoc comments ending with period
- Use descriptive names: `ArticulateParser`, `MarkdownExporter`
- Receiver names: short (1-2 chars), consistent per type

### Error Handling
- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use `%w` verb for error wrapping to preserve error chain
- Check all error returns (enforced by `errcheck`)
- Document error handling rationale in defer blocks when ignoring close errors

### Comments
- All exported types/functions require godoc comments
- End sentences with periods (`godot` linter enforced)
- Mark known issues with TODO/FIXME/HACK/BUG/XXX

### Security
- Use `#nosec` with justification for deliberate security exceptions (G304 for CLI file paths, G306 for export file permissions)
- Run `gosec` and `govulncheck` for security audits

### Testing
- Enable race detection: `-race` flag
- Use table-driven tests where applicable
- Mark test helpers with `t.Helper()`
- Benchmarks in `*_bench_test.go`, examples in `*_example_test.go`

### Dependencies
- Minimal external dependencies (currently: go-docx, golang.org/x/net, golang.org/x/text)
- Run `task deps:tidy` after adding/removing dependencies
