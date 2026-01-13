package extractor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func NewPythonExtractor(filePath string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer file.Close()

		var tasks []Task
		scanner := bufio.NewScanner(file)
		lineNum := 0

		// Pattern for # comments (like shell)
		hashCommentPattern := regexp.MustCompile(`#\s*(TODO|BUG|NOTE):\s*(.+)`)

		// Patterns for single-line docstrings
		tripleDoublePattern := regexp.MustCompile(`""".*?(TODO|BUG|NOTE):\s*(.+?)"""`)
		tripleSinglePattern := regexp.MustCompile(`'''.*?(TODO|BUG|NOTE):\s*(.+?)'''`)

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Check for # comments
			if matches := hashCommentPattern.FindStringSubmatch(line); len(matches) > 0 {
				task := Task{
					File:    filePath,
					Line:    lineNum,
					Column:  strings.Index(line, matches[0]) + 1,
					Type:    matches[1],
					Message: strings.TrimSpace(matches[2]),
				}
				tasks = append(tasks, task)
			}

			// Check for triple double quote docstrings
			if matches := tripleDoublePattern.FindStringSubmatch(line); len(matches) > 0 {
				col := strings.Index(line, matches[0]) + strings.Index(matches[0], matches[1]+":") + 1
				task := Task{
					File:    filePath,
					Line:    lineNum,
					Column:  col,
					Type:    matches[1],
					Message: strings.TrimSpace(matches[2]),
				}
				tasks = append(tasks, task)
			}

			// Check for triple single quote docstrings
			if matches := tripleSinglePattern.FindStringSubmatch(line); len(matches) > 0 {
				col := strings.Index(line, matches[0]) + strings.Index(matches[0], matches[1]+":") + 1
				task := Task{
					File:    filePath,
					Line:    lineNum,
					Column:  col,
					Type:    matches[1],
					Message: strings.TrimSpace(matches[2]),
				}
				tasks = append(tasks, task)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		return tasks, nil
	})
}
