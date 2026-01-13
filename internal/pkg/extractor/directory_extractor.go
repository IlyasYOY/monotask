package extractor

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func NewDirectoryExtractor(dirPath string, ignores ...string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		var allIgnores []string
		allIgnores = append(allIgnores, ignores...)

		mtignores, err := readMtignores(dirPath)
		if err != nil {
			return nil, err
		}
		allIgnores = append(allIgnores, mtignores...)

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("error reading directory: %w", err)
		}

		var allTasks []Task
		for _, entry := range entries {
			if entry.Name() == mtIgnoreFilename {
				continue
			}

			fullPath := filepath.Join(dirPath, entry.Name())
			fullPath, err = filepath.Abs(fullPath)
			if err != nil {
				log.Printf("Error getting absolute path to %s: %v", fullPath, err)
				continue
			}

			if slices.Contains(allIgnores, fullPath) {
				continue
			}

			if entry.IsDir() {
				subExtractor := NewDirectoryExtractor(fullPath, allIgnores...)
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

const mtIgnoreFilename = ".mtignore"

func readMtignores(atDir string) ([]string, error) {
	mtignorePath := filepath.Join(atDir, mtIgnoreFilename)
	f, err := os.Open(mtignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var mtignores []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ignore := strings.TrimSpace(scanner.Text())
		// otherwise empty strings turn to current dir ignores
		if len(ignore) == 0 {
			continue
		}

		ignore = filepath.Join(atDir, ignore)
		ignore, err := filepath.Abs(ignore)
		if err != nil {
			log.Printf(
				"Error transform ignore path %s from %s: %v",
				ignore, mtignorePath, err,
			)
			return nil, err
		}
		mtignores = append(mtignores, ignore)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mtignores, nil
}
