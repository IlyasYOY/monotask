package extractor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type markdownExtractor struct {
	filePath string
}

func NewMarkdownExtractor(filePath string) Extractor {
	return &markdownExtractor{filePath: filePath}
}

func (e *markdownExtractor) Extract(ctx context.Context) ([]Task, error) {
	file, err := os.Open(e.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var tasks []Task
	scanner := bufio.NewScanner(file)
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
			File:    e.filePath,
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
}
