package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/IlyasYOY/monotask/internal/pkg/extractor"
	"github.com/IlyasYOY/monotask/internal/pkg/output"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	dirExtractor := extractor.NewDirectoryExtractor(absPath)

	tasks, err := dirExtractor.Extract(ctx)
	if err != nil {
		log.Printf("Error extracting tasks: %v", err)
		os.Exit(1)
	}

	output.PrintGNUFormatTo(tasks, os.Stdout)
}
