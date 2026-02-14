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
	adocLineCommentRegex  = regexp.MustCompile(`//\s*` + taskRegexCore)
	adocBlockCommentRegex = regexp.MustCompile(`////\s*` + taskRegexCore)
	adocInBlockRegex      = regexp.MustCompile(`\s*` + taskRegexCore)
)

func NewAsciiDocExtractor(filePath string) Extractor {
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

		inBlockComment := false
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if inBlockComment {
				if idx := strings.Index(line, "////"); idx >= 0 {
					inBlockComment = false
					rest := strings.TrimSpace(line[idx+4:])
					if rest != "" {
						if matches := adocInBlockRegex.FindStringSubmatch(rest); len(matches) > 0 {
							col := idx + 5 + strings.Index(rest, matches[0])
							task := ParseTask(matches, filePath, lineNum, col)
							tasks = append(tasks, task)
						}
					}
				} else {
					if matches := adocInBlockRegex.FindStringSubmatch(line); len(matches) > 0 {
						col := strings.Index(line, matches[0]) + 1
						task := ParseTask(matches, filePath, lineNum, col)
						tasks = append(tasks, task)
					}
				}
				continue
			}

			if strings.Contains(line, "////") {
				if matches := adocBlockCommentRegex.FindStringSubmatch(line); len(matches) > 0 {
					col := strings.Index(line, matches[0]) + 1
					task := ParseTask(matches, filePath, lineNum, col)
					tasks = append(tasks, task)

					if !strings.Contains(line, "////") || strings.Index(line, "////") == strings.LastIndex(line, "////") {
						// If there's only one //// or the second //// comes after the match,
						// we might still be in a block
						if !strings.Contains(line[strings.Index(line, matches[0])+len(matches[0]):], "////") {
							inBlockComment = true
						}
					}
				} else {
					inBlockComment = true
				}
				continue
			}

			if matches := adocLineCommentRegex.FindStringSubmatch(line); len(matches) > 0 {
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
