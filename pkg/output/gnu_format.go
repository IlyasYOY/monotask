package output

import (
	"fmt"
	"os"

	"github.com/ilyasyoy/monotask/pkg/extractor"
)

func PrintGNUFormat(tasks []extractor.Task) {
	for _, task := range tasks {
		fmt.Printf("%s:%d:%d: %s: %s\n", task.File, task.Line, task.Column, task.Type, task.Message)
	}
}

func PrintGNUFormatTo(tasks []extractor.Task, file *os.File) {
	for _, task := range tasks {
		fmt.Fprintf(file, "%s:%d:%d: %s: %s\n", task.File, task.Line, task.Column, task.Type, task.Message)
	}
}
