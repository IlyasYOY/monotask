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
	// I don't need time here:
	// - makes testing harder,
	// - doesn't add benefits.
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

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
	taskExtractor := extractor.NewDirectoryExtractor(absPath)
	if fileInfo, err := os.Stat(absPath); err == nil && !fileInfo.IsDir() {
		taskExtractor = extractor.NewFileExtractor(absPath)
	}

	tasks, err := taskExtractor.Extract(ctx)
	if err != nil {
		log.Printf("Error extracting tasks: %v", err)
		os.Exit(1)
	}

	output.PrintGNUFormatTo(tasks, os.Stdout)
}
