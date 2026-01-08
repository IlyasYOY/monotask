package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyasyoy/monotask/internal/pkg/extractor"
	"github.com/ilyasyoy/monotask/internal/pkg/output"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	dirExtractor := extractor.NewDirectoryExtractor(absPath)

	tasks, err := dirExtractor.Extract(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting tasks: %v\n", err)
		os.Exit(1)
	}

	output.PrintGNUFormatTo(tasks, os.Stdout)
}
