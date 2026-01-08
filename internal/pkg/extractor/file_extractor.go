package extractor

import (
	"context"
	"path/filepath"
	"strings"
)

type fileExtractor struct {
	filePath string
}

func NewFileExtractor(filePath string) Extractor {
	return &fileExtractor{filePath: filePath}
}

func (e *fileExtractor) Extract(ctx context.Context) ([]Task, error) {
	ext := strings.ToLower(filepath.Ext(e.filePath))

	switch ext {
	case ".md":
		return NewMarkdownExtractor(e.filePath).Extract(ctx)
	case ".c", ".h", ".java", ".go", ".js", ".mjs", ".ts", ".mts", ".cpp", ".hpp", ".cxx", ".cc":
		return NewCCommentsExtractor(e.filePath).Extract(ctx)
	default:
		return []Task{}, nil
	}
}
