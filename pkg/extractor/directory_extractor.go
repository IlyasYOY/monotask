package extractor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type directoryExtractor struct {
	dirPath string
}

func NewDirectoryExtractor(dirPath string) Extractor {
	return &directoryExtractor{dirPath: dirPath}
}

func (e *directoryExtractor) Extract(ctx context.Context) ([]Task, error) {
	var allTasks []Task

	err := filepath.Walk(e.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if info.IsDir() {
			return nil
		}

		extractor := NewFileExtractor(path)
		tasks, err := extractor.Extract(ctx)
		if err != nil {
			return fmt.Errorf("failed to extract from %s: %w", path, err)
		}

		allTasks = append(allTasks, tasks...)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return allTasks, nil
}
