package filestest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const separatorPrefix = "-#file:"

// RenderDir creates a temporary directory structure based on the provided description.
// The description is a multi-line string where file contents are defined using
// the format: "-#file:<filename>" followed by the file's content lines.
// Files can be nested in subdirectories by including slashes in the filename.
// It returns the path to the created temporary directory.
// This is primarily used for testing purposes to set up file structures.
//
// Example:
//
//	description := `-#file:hello.txt
//	Hello, World!
//
//	-#file:subdir/goodbye.txt
//	Goodbye!`
//	tmpDir := RenderDir(t, description)
//
// Creates hello.txt with "Hello, World!\n" and subdir/goodbye.txt with "Goodbye!\n"
func RenderDir(t *testing.T, description string) string {
	t.Helper()

	tmpDir := t.TempDir()
	lines := make([]string, 0)
	for line := range strings.SplitSeq(description, "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}

	for i := range len(lines) {
		line := lines[i]
		if !strings.HasPrefix(line, separatorPrefix) {
			continue
		}

		fileName := strings.TrimSpace(line[len(separatorPrefix):])
		var content strings.Builder
		for j := i + 1; j < len(lines); j++ {
			contentLine := lines[j]
			if strings.HasPrefix(contentLine, separatorPrefix) {
				break
			}
			if contentLine != "" {
				content.WriteString(contentLine + "\n")
			}
		}

		filePath := filepath.Join(tmpDir, fileName)
		os.MkdirAll(filepath.Dir(filePath), 0o755)
		os.WriteFile(filePath, []byte(content.String()), 0o644)
	}

	return tmpDir
}
