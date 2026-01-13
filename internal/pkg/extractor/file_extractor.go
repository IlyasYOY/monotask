package extractor

import (
	"context"
	"path/filepath"
	"strings"
)

func NewFileExtractor(filePath string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		ext := strings.ToLower(filepath.Ext(filePath))

		switch ext {
		case ".md":
			return NewMarkdownExtractor(filePath).Extract(ctx)
		case ".lua":
			return NewLuaExtractor(filePath).Extract(ctx)
		case ".sh", ".bash":
			return NewShellExtractor(filePath).Extract(ctx)
		case ".py":
			return NewPythonExtractor(filePath).Extract(ctx)
		case ".c", ".h", ".java", ".go", ".js", ".mjs", ".ts", ".mts", ".cpp", ".hpp", ".cxx", ".cc":
			return NewCCommentsExtractor(filePath).Extract(ctx)
		default:
			return []Task{}, nil
		}
	})
}
