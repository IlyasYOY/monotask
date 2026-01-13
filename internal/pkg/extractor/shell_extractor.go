package extractor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	commentRegex = regexp.MustCompile(`#\s*` + taskRegexCore)
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
		scanner.Buffer(nil, 1024*1024) // Set max token size to 1MB for long lines
		lineNum := 0

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if matches := commentRegex.FindStringSubmatch(line); len(matches) > 0 {
				col := strings.Index(line, matches[0]) + 1
				task := ParseTask(matches, filePath, lineNum, col)
				tasks = append(tasks, task)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		return tasks, nil
	})
}
