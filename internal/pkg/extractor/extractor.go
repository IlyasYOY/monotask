package extractor

import (
	"context"
	"strings"
)

const taskRegexCore = `(TODO|BUG|NOTE)(\([^)]*\))?:\s*(.+)`

type Task struct {
	File     string
	Line     int
	Column   int
	Type     string
	Assignee string
	Message  string
}

func ParseTask(matches []string, filePath string, lineNum int, column int) Task {
	typ := matches[1]
	assignee := ""
	if len(matches) > 2 && matches[2] != "" {
		assignee = strings.Trim(matches[2], "()")
	}
	message := strings.TrimSpace(matches[len(matches)-1])
	return Task{
		File:     filePath,
		Line:     lineNum,
		Column:   column,
		Type:     typ,
		Assignee: assignee,
		Message:  message,
	}
}

type Extractor interface {
	Extract(ctx context.Context) ([]Task, error)
}

type ExtractorFunc func(ctx context.Context) ([]Task, error)

func (f ExtractorFunc) Extract(ctx context.Context) ([]Task, error) {
	return f(ctx)
}
