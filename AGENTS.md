# AGENTS.md - Monotask Development Guidelines

This document provides guidelines for coding agents working on the monotask codebase. Monotask is a CLI tool written in Go that extracts tasks (TODO, BUG, NOTE markers from comments and unchecked checkboxes from markdown files) and outputs them in GNU error format.

## Build, Test, and Development Commands

### Building
```bash
# Build the monotask binary
go build -o ./bin/monotask ./cmd/monotask

# Clean build artifacts
make clean
```

### Testing
```bash
# Run all tests (builds binary first, includes race detection and parallel execution)
make test

# Run tests for a specific package
go test -v ./pkg/extractor
go test -v ./pkg/output

# Run a single test function
go test -run TestCCommentsExtractor -v ./pkg/extractor

# Run integration tests only
go test -run TestIntegration -v ./internal/integtest
```

### Linting and Formatting
```bash
# Format Go code
gofmt -w .

# Check for formatting issues (non-zero exit if issues found)
gofmt -d . | tee /dev/stderr | [ "$(cat)" = "" ]

# Vet for suspicious constructs
go vet ./...

# Check for unused dependencies
go mod tidy
```

### Development Workflow
```bash
# Full development cycle
make clean && go build -o ./bin/monotask ./cmd/monotask && make test

# Quick iteration (rebuild and test)
go build -o ./bin/monotask ./cmd/monotask && MONOTASK_BINARY=./bin/monotask go test -race -v ./...
```

## Code Style Guidelines

### Go Version Requirements
- **Minimum Go version**: 1.25.5
- Use Go modules for dependency management
- Keep dependencies minimal; prefer standard library when possible

### Import Organization
```go
import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/extractor"
)
```
- Standard library imports first
- Third-party imports second
- Local project imports last
- Blank lines between groups
- No unused imports (go mod tidy will catch this)

### Naming Conventions
- **Exported types/functions**: PascalCase (`NewDirectoryExtractor`, `Task`)
- **Unexported types/functions**: camelCase (`directoryExtractor`, `extractTasks`)
- **Variables**: camelCase (`filePath`, `lineNum`)
- **Constants**: PascalCase if exported, camelCase if unexported
- **Struct fields**: PascalCase for exported structs, camelCase for unexported
- **Interface methods**: PascalCase

### Error Handling
```go
// Good: Wrap errors with context
func (e *extractor) Extract(ctx context.Context) ([]Task, error) {
	file, err := os.Open(e.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", e.filePath, err)
	}
	defer file.Close()

	// ... processing ...

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return tasks, nil
}
```
- Use `fmt.Errorf` with `%w` verb for error wrapping
- Return errors early rather than nested if-else chains
- Include relevant context in error messages
- Use `defer` for resource cleanup

### Logging
```go
import (
	"log"
)

// Good: Use log package for error logging
func processDirectory(dirPath string) error {
	// ... processing ...
	if err != nil {
		log.Printf("Error processing directory %s: %v", dirPath, err)
		return err
	}
	return nil
}
```
- Use `log` package for logging errors to stderr
- `log.Printf` automatically adds timestamps and newlines
- Prefer `log` over `fmt.Fprintf(os.Stderr, ...)` for consistency

### Context Usage
- Always include `context.Context` in public API methods
- Pass context through to underlying operations
- Use `context.Background()` in main functions and tests

### File and I/O Operations
```go
// Good: Use bufio.Scanner for line-by-line reading
func processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Process line
	}

	return scanner.Err()
}
```
- Use `bufio.Scanner` for line-by-line file reading
- Always check `scanner.Err()` after scanning
- Use `defer` for file cleanup
- Use `io.Writer` interfaces for flexible output

### Regular Expressions
```go
// Good: Compile regex patterns once
type extractor struct {
	todoPattern *regexp.Regexp
}

func NewExtractor() *extractor {
	return &extractor{
		todoPattern: regexp.MustCompile(`//\s*TODO:\s*(.+)`),
	}
}
```
- Compile regex patterns once (in constructor or init)
- Use `regexp.MustCompile` for patterns that are known to be valid
- Use raw strings for regex patterns to avoid escaping issues

### Struct Initialization
```go
// Good: Use field names for clarity
task := Task{
    File:    filePath,
    Line:    lineNum,
    Column:  colNum,
    Type:    "TODO",
    Message: message,
}
```
- Use field names when initializing structs with multiple fields
- Keep field order consistent with struct definition
- Use `&` for pointer returns from constructors

### Testing Patterns
```go
func TestExtractor(t *testing.T) {
	t.Parallel() // Run tests in parallel

	tests := []struct {
		name     string
		input    string
		expected []Task
	}{
		{"simple todo", "// TODO: fix bug", []Task{{Type: "TODO", Message: "fix bug"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
```
- Use `t.Parallel()` for parallel test execution
- Use table-driven tests with subtests (`t.Run()`)
- Use `github.com/google/go-cmp/cmp` for comparing complex types
- Use descriptive test names and clear error messages

### Interface Design
```go
type Extractor interface {
	Extract(ctx context.Context) ([]Task, error)
}

type Task struct {
	File    string
	Line    int
	Column  int
	Type    string
	Message string
}
```
- Keep interfaces small and focused
- Use context in interface methods
- Design structs with exported fields for JSON serialization compatibility
- Use meaningful field names that clearly indicate their purpose

### Project-Specific Conventions

#### Task Types
- **TODO**: General tasks that need to be completed
- **BUG**: Known bugs that need fixing
- **NOTE**: Important notes or reminders
- **CHECKBOX**: Unchecked markdown checkboxes (`- [ ]`)

#### Assignee Support
Tasks can optionally include an assignee in parentheses after the type:
- `// TODO(ilyasyoy): fix this bug`
- `# BUG(user): handle error case`
- `-- NOTE(team): review implementation`

Empty parentheses `TODO():` are supported and treated as tasks without assignee.

#### Supported File Types
- **C/C++**: `.c`, `.h`, `.cpp`, `.hpp`, `.cxx`, `.cc` (C-style comments)
- **Java**: `.java` (C-style comments)
- **Go**: `.go` (C-style comments)
- **JavaScript/TypeScript**: `.js`, `.mjs`, `.ts`, `.mts` (C-style comments)
- **Lua**: `.lua` (single-line `--` and multi-line `--[[ ]]`)
- **Shell**: `.sh`, `.bash` (single-line `#`)
- **Python**: `.py` (# comments and single-line docstrings)
- **Markdown**: `.md` (unchecked checkboxes)

#### Output Format
Always use GNU error format for consistency:
```
file:line:column: type: message
```
Or with optional assignee:
```
file:line:column: type(assignee): message
```
Examples:
- `src/main.go:15:3: TODO: implement error handling`
- `src/main.go:15:3: TODO(ilyasyoy): implement error handling`

#### Architecture Patterns
- Use the Extractor interface for all extraction logic
- Separate concerns: extractors for different file types, output formatters
- Directory traversal should be recursive but skip hidden directories
- File processing should be streaming (not loading entire files into memory)

### Code Organization
- **cmd/**: Main application entry points
- **pkg/**: Library code that can be imported by other projects
- **internal/**: Private application code
- Keep packages focused on single responsibilities
- Use internal packages for code that shouldn't be imported externally

### Performance Considerations
- Compile regex patterns once, not per file
- Use streaming I/O for large files
- Avoid unnecessary string allocations
- Use `strings.Builder` for string concatenation in loops

### Security Best Practices
- Validate file paths before opening
- Don't follow symlinks in directory traversal
- Use context for cancellation in long-running operations
- Sanitize output to prevent injection in error messages

Remember: When in doubt, follow the existing patterns in the codebase. Consistency is more important than personal preference.

