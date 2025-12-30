package extractor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type cCommentsExtractor struct {
	filePath string
}

func NewCCommentsExtractor(filePath string) Extractor {
	return &cCommentsExtractor{filePath: filePath}
}

func (e *cCommentsExtractor) Extract(ctx context.Context) ([]Task, error) {
	file, err := os.Open(e.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var tasks []Task
	scanner := bufio.NewScanner(file)
	lineNum := 0

	lineCommentPattern := regexp.MustCompile(`//\s*(TODO|BUG|NOTE):\s*(.+)`)

	inBlockComment := false
	blockCommentPattern := regexp.MustCompile(`/\*\s*(TODO|BUG|NOTE):\s*(.+)`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if inBlockComment {
			if idx := strings.Index(line, "*/"); idx >= 0 {
				inBlockComment = false
				rest := strings.TrimSpace(line[idx+2:])
				if rest != "" {
					if matches := regexp.MustCompile(`\s*(TODO|BUG|NOTE):\s*(.+)`).FindStringSubmatch(rest); len(matches) > 0 {
						task := Task{
							File:    e.filePath,
							Line:    lineNum,
							Column:  idx + 3,
							Type:    matches[1],
							Message: strings.TrimSpace(matches[2]),
						}
						tasks = append(tasks, task)
					}
				}
			} else {
				if matches := regexp.MustCompile(`\s*(TODO|BUG|NOTE):\s*(.+)`).FindStringSubmatch(line); len(matches) > 0 {
					task := Task{
						File:    e.filePath,
						Line:    lineNum,
						Column:  strings.Index(line, matches[0]) + 1,
						Type:    matches[1],
						Message: strings.TrimSpace(matches[2]),
					}
					tasks = append(tasks, task)
				}
			}
			continue
		}

		if strings.Contains(line, "/*") {
			if matches := blockCommentPattern.FindStringSubmatch(line); len(matches) > 0 {
				task := Task{
					File:    e.filePath,
					Line:    lineNum,
					Column:  strings.Index(line, matches[0]) + 1,
					Type:    matches[1],
					Message: strings.TrimSpace(matches[2]),
				}
				tasks = append(tasks, task)

				if !strings.Contains(line, "*/") {
					inBlockComment = true
				}
			} else {
				inBlockComment = true
			}
			continue
		}

		if matches := lineCommentPattern.FindStringSubmatch(line); len(matches) > 0 {
			task := Task{
				File:    e.filePath,
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
}
