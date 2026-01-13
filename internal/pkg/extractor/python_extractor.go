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
	hashCommentRegex  = regexp.MustCompile(`#\s*` + taskRegexCore)
	tripleDoubleRegex = regexp.MustCompile(`""".*?` + taskRegexCore + `"""`)
	tripleSingleRegex = regexp.MustCompile(`'''.*?` + taskRegexCore + `'''`)
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

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Check for # comments
			if matches := hashCommentRegex.FindStringSubmatch(line); len(matches) > 0 {
				col := strings.Index(line, matches[0]) + 1
				task := ParseTask(matches, filePath, lineNum, col)
				tasks = append(tasks, task)
			}

			// Check for triple double quote docstrings
			if matches := tripleDoubleRegex.FindStringSubmatch(line); len(matches) > 0 {
				col := strings.Index(line, matches[0]) + strings.Index(matches[0], matches[1]+":") + 1
				task := ParseTask(matches, filePath, lineNum, col)
				tasks = append(tasks, task)
			}

			// Check for triple single quote docstrings
			if matches := tripleSingleRegex.FindStringSubmatch(line); len(matches) > 0 {
				col := strings.Index(line, matches[0]) + strings.Index(matches[0], matches[1]+":") + 1
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
