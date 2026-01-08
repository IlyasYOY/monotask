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
