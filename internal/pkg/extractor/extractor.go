package extractor

import "context"

type Task struct {
	File    string
	Line    int
	Column  int
	Type    string
	Message string
}

type Extractor interface {
	Extract(ctx context.Context) ([]Task, error)
}

type ExtractorFunc func(ctx context.Context) ([]Task, error)

func (f ExtractorFunc) Extract(ctx context.Context) ([]Task, error) {
	return f(ctx)
}
