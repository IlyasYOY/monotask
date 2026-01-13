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
	lineCommentRegex  = regexp.MustCompile(`//\s*` + taskRegexCore)
	blockCommentRegex = regexp.MustCompile(`/\*\s*` + taskRegexCore)
	inBlockRegex      = regexp.MustCompile(`\s*` + taskRegexCore)
)

func NewCCommentsExtractor(filePath string) Extractor {
	return ExtractorFunc(func(ctx context.Context) ([]Task, error) {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		var tasks []Task
		scanner := bufio.NewScanner(file)
		lineNum := 0

		inBlockComment := false
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if inBlockComment {
				if idx := strings.Index(line, "*/"); idx >= 0 {
					inBlockComment = false
					rest := strings.TrimSpace(line[idx+2:])
					if rest != "" {
						if matches := inBlockRegex.FindStringSubmatch(rest); len(matches) > 0 {
							col := idx + 3 + strings.Index(rest, matches[0])
							task := ParseTask(matches, filePath, lineNum, col)
							tasks = append(tasks, task)
						}
					}
				} else {
					if matches := inBlockRegex.FindStringSubmatch(line); len(matches) > 0 {
						col := strings.Index(line, matches[0]) + 1
						task := ParseTask(matches, filePath, lineNum, col)
						tasks = append(tasks, task)
					}
				}
				continue
			}

			if strings.Contains(line, "/*") {
				if matches := blockCommentRegex.FindStringSubmatch(line); len(matches) > 0 {
					col := strings.Index(line, matches[0]) + 1
					task := ParseTask(matches, filePath, lineNum, col)
					tasks = append(tasks, task)

					if !strings.Contains(line, "*/") {
						inBlockComment = true
					}
				} else {
					inBlockComment = true
				}
				continue
			}

			if matches := lineCommentRegex.FindStringSubmatch(line); len(matches) > 0 {
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
