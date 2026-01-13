package output

import (
	"fmt"
	"io"

	"github.com/IlyasYOY/monotask/internal/pkg/extractor"
)

func PrintGNUFormatTo(tasks []extractor.Task, writer io.Writer) {
	for _, task := range tasks {
		if task.Assignee != "" {
			fmt.Fprintf(writer, "%s:%d:%d: %s(%s): %s\n", task.File, task.Line, task.Column, task.Type, task.Assignee, task.Message)
		} else {
			fmt.Fprintf(writer, "%s:%d:%d: %s: %s\n", task.File, task.Line, task.Column, task.Type, task.Message)
		}
	}
}
