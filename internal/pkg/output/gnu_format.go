package output

import (
	"fmt"
	"io"

	"github.com/ilyasyoy/monotask/internal/pkg/extractor"
)

func PrintGNUFormatTo(tasks []extractor.Task, writer io.Writer) {
	for _, task := range tasks {
		fmt.Fprintf(writer, "%s:%d:%d: %s: %s\n", task.File, task.Line, task.Column, task.Type, task.Message)
	}
}
