package output_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/extractor"
	"github.com/ilyasyoy/monotask/pkg/output"
)

func TestPrintGNUFormatTo(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []extractor.Task
		expected string
	}{
		{
			name:     "empty tasks",
			tasks:    []extractor.Task{},
			expected: "",
		},
		{
			name: "single task",
			tasks: []extractor.Task{
				{File: "main.go", Line: 10, Column: 5, Type: "error", Message: "syntax error"},
			},
			expected: "main.go:10:5: error: syntax error\n",
		},
		{
			name: "multiple tasks",
			tasks: []extractor.Task{
				{File: "main.go", Line: 10, Column: 5, Type: "error", Message: "syntax error"},
				{File: "utils.go", Line: 25, Column: 12, Type: "warning", Message: "unused variable"},
				{File: "test.go", Line: 1, Column: 1, Type: "info", Message: "test message"},
			},
			expected: "main.go:10:5: error: syntax error\nutils.go:25:12: warning: unused variable\ntest.go:1:1: info: test message\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			output.PrintGNUFormatTo(tt.tasks, &buf)

			if diff := cmp.Diff(tt.expected, buf.String()); diff != "" {
				t.Errorf("(-want +got):\\n%s", diff)
			}
		})
	}
}
