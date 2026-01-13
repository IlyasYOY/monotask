package extractor

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func NewDirectoryExtractor(dirPath string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		var allTasks []Task

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("error reading directory: %w", err)
		}

		for _, entry := range entries {
			fullPath := filepath.Join(dirPath, entry.Name())

			if entry.IsDir() {
				subExtractor := NewDirectoryExtractor(fullPath)
				subTasks, err := subExtractor.Extract(ctx)
				if err != nil {
					log.Printf("Error extracting from directory %s: %v", fullPath, err)
					continue
				}
				allTasks = append(allTasks, subTasks...)
			} else {
				extractor := NewFileExtractor(fullPath)
				tasks, err := extractor.Extract(ctx)
				if err != nil {
					log.Printf("Error extracting from %s: %v", fullPath, err)
					continue
				}
				allTasks = append(allTasks, tasks...)
			}
		}

		return allTasks, nil
	})
}
