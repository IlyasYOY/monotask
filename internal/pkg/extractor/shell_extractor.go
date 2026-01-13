package extractor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func NewShellExtractor(filePath string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer file.Close()

		var tasks []Task
		scanner := bufio.NewScanner(file)
		lineNum := 0

		commentPattern := regexp.MustCompile(`#\s*(TODO|BUG|NOTE):\s*(.+)`)

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if matches := commentPattern.FindStringSubmatch(line); len(matches) > 0 {
				task := Task{
					File:    filePath,
					Line:    lineNum,
					Column:  strings.Index(line, matches[0]) + 1,
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
