package extractor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func NewMarkdownExtractor(filePath string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		var tasks []Task
		scanner := bufio.NewScanner(file)
		scanner.Buffer(nil, 1024*1024) // Set max token size to 1MB for long lines
		lineNum := 0

		checkboxPattern := regexp.MustCompile(`^- \[ \] (.+)`)
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			matches := checkboxPattern.FindStringSubmatch(line)
			if matches == nil {
				continue
			}

			tasks = append(tasks, Task{
				File:    filePath,
				Line:    lineNum,
				Column:  strings.Index(line, "- [ ]") + 1,
				Type:    "CHECKBOX",
				Message: strings.TrimSpace(matches[1]),
			})
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		return tasks, nil
	})
}
