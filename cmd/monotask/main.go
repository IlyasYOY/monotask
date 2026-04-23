package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/IlyasYOY/monotask/internal/pkg/extractor"
	"github.com/IlyasYOY/monotask/internal/pkg/output"
	"github.com/IlyasYOY/monotask/internal/pkg/version"
)

var versionString = version.Current

func main() {
	os.Exit(run(os.Args[1:], os.Stdout))
}

func run(args []string, stdout io.Writer) int {
	// I don't need time here:
	// - makes testing harder,
	// - doesn't add benefits.
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	flags := flag.NewFlagSet("monotask", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var showVersion bool
	flags.BoolVar(&showVersion, "v", false, "print version")
	flags.BoolVar(&showVersion, "version", false, "print version")
	if err := flags.Parse(args); err != nil {
		log.Printf("Error parsing flags: %v", err)
		return 2
	}

	if showVersion {
		fmt.Fprintln(stdout, versionString())
		return 0
	}

	path := "."
	if flags.NArg() > 0 {
		path = flags.Arg(0)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
		return 1
	}

	ctx := context.Background()
	taskExtractor := extractor.NewDirectoryExtractor(absPath)
	if fileInfo, err := os.Stat(absPath); err == nil && !fileInfo.IsDir() {
		taskExtractor = extractor.NewFileExtractor(absPath)
	}

	tasks, err := taskExtractor.Extract(ctx)
	if err != nil {
		log.Printf("Error extracting tasks: %v", err)
		return 1
	}

	output.PrintGNUFormatTo(tasks, stdout)
	return 0
}
