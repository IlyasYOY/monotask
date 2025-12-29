package filestest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func RenderDir(t *testing.T, description string) string {
	tmpDir := t.TempDir()
	lines := make([]string, 0)
	for line := range strings.SplitSeq(description, "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if !strings.HasPrefix(line, "#") {
			continue
		}

		fileName := strings.TrimSpace(line[1:])
		var content strings.Builder
		for j := i + 1; j < len(lines); j++ {
			contentLine := lines[j]
			if strings.HasPrefix(contentLine, "#") {
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
