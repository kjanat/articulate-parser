# Contributing to Articulate Rise Parser

Thank you for your interest in contributing to the Articulate Rise Parser! We welcome contributions from the community.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find that the issue has already been reported. When creating a bug report, include as many details as possible:

- Use the bug report template
- Include sample Articulate Rise content that reproduces the issue
- Provide your environment details (OS, Go version, etc.)
- Include error messages and stack traces

### Suggesting Enhancements

Enhancement suggestions are welcome! Please use the feature request template and include:

- A clear description of the enhancement
- Your use case and why this would be valuable
- Any implementation ideas you might have

### Pull Requests

1. **Fork the repository** and create your branch from `master`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Ensure all tests pass** by running `go test ./...`
5. **Run `go fmt`** to format your code
6. **Run `go vet`** to check for common issues
7. **Update documentation** if needed
8. **Create a pull request** with a clear title and description

## Development Setup

1. **Prerequisites:**

- Go 1.21 or later
- Git

2. **Clone and setup:**

   ```bash
   git clone https://github.com/your-username/articulate-parser.git
   cd articulate-parser
   go mod download
   ```

3. **Run tests:**

   ```bash
   go test -v ./...
   ```

4. **Build:**

   ```bash
   go build main.go
   ```

## Coding Standards

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

### Testing

- Write tests for new functionality
- Use table-driven tests where appropriate
- Aim for good test coverage
- Test error cases and edge conditions

### Commit Messages

Use clear and meaningful commit messages:

```txt
Add support for new content type: interactive timeline

- Parse timeline content blocks
- Export timeline data to markdown
- Add tests for timeline parsing
- Update documentation

Fixes #123
```

## Project Structure

```txt
articulate-parser/
├── main.go                 # Entry point and CLI handling
├── parser/                 # Core parsing logic
├── exporters/              # Output format handlers
├── types/                  # Data structures
├── utils/                  # Utility functions
├── tests/                  # Test files and data
└── docs/                   # Documentation
```

## Adding New Features

### New Content Types

1. Add the content type definition to `types/`
2. Implement parsing logic in `parser/`
3. Add export handling in `exporters/`
4. Write comprehensive tests
5. Update documentation

### New Export Formats

1. Create a new exporter in `exporters/`
2. Implement the `Exporter` interface
3. Add CLI support in `main.go`
4. Add tests with sample output
5. Update README with usage examples

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestSpecificFunction ./...
```

### Test Data

- Add sample Articulate Rise JSON files to `tests/data/`
- Include both simple and complex content examples
- Test edge cases and error conditions

## Documentation

- Update the README for user-facing changes
- Add inline code comments for complex logic
- Update examples when adding new features
- Keep the feature list current

## Release Process

Releases are handled by maintainers:

1. Version bumping follows semantic versioning
2. Releases are created from the `master` branch
3. GitHub Actions automatically builds and publishes releases
4. Release notes are auto-generated from commits

## Questions?

- Open a discussion for general questions
- Use the question issue template for specific help
- Check existing issues and documentation first

## Recognition

Contributors will be recognized in release notes and the project README. Thank you for helping make this project better!
