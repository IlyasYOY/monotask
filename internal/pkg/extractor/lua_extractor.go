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
	singleLinePattern  = regexp.MustCompile(`--\s*(TODO|BUG|NOTE):\s*(.+)`)
	inBlockLinePattern = regexp.MustCompile(`\s*(TODO|BUG|NOTE):\s*(.+)`)
)

type luaExtractor struct {
	filePath string
}

func NewLuaExtractor(filePath string) Extractor {
	return &luaExtractor{filePath: filePath}
}

func (e *luaExtractor) Extract(ctx context.Context) ([]Task, error) {
	file, err := os.Open(e.filePath)
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
			if before, after, ok := strings.Cut(line, "]]"); ok {
				inBlockComment = false

				// Check content before ]] for tasks
				beforeEnd := before
				if matches := regexp.MustCompile(`\s*(TODO|BUG|NOTE):\s*(.+)`).FindStringSubmatch(beforeEnd); len(matches) > 0 {
					task := Task{
						File:    e.filePath,
						Line:    lineNum,
						Column:  strings.Index(line, matches[0]) + 1,
						Type:    matches[1],
						Message: strings.TrimSpace(matches[2]),
					}
					tasks = append(tasks, task)
				}

				// Check after ]] for single-line comments
				afterEnd := strings.TrimSpace(after)
				if strings.HasPrefix(afterEnd, "--") {
					if matches := singleLinePattern.FindStringSubmatch(afterEnd); len(matches) > 0 {
						task := Task{
							File:    e.filePath,
							Line:    lineNum,
							Column:  strings.Index(line, afterEnd) + 1,
							Type:    matches[1],
							Message: strings.TrimSpace(matches[2]),
						}
						tasks = append(tasks, task)
					}
				}
			} else {
				// Still inside block comment, check this line for tasks
				if matches := inBlockLinePattern.FindStringSubmatch(line); len(matches) > 0 {
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

		// Check for start of block comment
		if _, after, ok := strings.Cut(line, "--[["); ok {
			inBlockComment = true

			// Check content after --[[ for tasks
			if before, after0, ok0 := strings.Cut(after, "]]"); ok0 {
				// Block comment ends on same line
				blockContent := before
				if matches := regexp.MustCompile(`\s*(TODO|BUG|NOTE):\s*(.+)`).FindStringSubmatch(blockContent); len(matches) > 0 {
					task := Task{
						File:    e.filePath,
						Line:    lineNum,
						Column:  strings.Index(line, matches[0]) + strings.Index(after, matches[0]) + 1,
						Type:    matches[1],
						Message: strings.TrimSpace(matches[2]),
					}
					tasks = append(tasks, task)
				}
				inBlockComment = false
				// Check after ]] for single-line comments
				afterEnd := strings.TrimSpace(after0)
				if strings.HasPrefix(afterEnd, "--") {
					if matches := singleLinePattern.FindStringSubmatch(afterEnd); len(matches) > 0 {
						task := Task{
							File:    e.filePath,
							Line:    lineNum,
							Column:  strings.Index(line, afterEnd) + 1,
							Type:    matches[1],
							Message: strings.TrimSpace(matches[2]),
						}
						tasks = append(tasks, task)
					}
				}
			} else {
				// Block comment continues to next line
				if matches := inBlockLinePattern.FindStringSubmatch(after); len(matches) > 0 {
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

		// Check for single-line comments
		if matches := singleLinePattern.FindStringSubmatch(line); len(matches) > 0 {
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
