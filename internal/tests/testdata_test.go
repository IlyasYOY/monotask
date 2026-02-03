// Package integtest provides integration tests for the monotask CLI.
//
// The tests use binary of the application. In case of missing binary they skip.
//
// Binary for tests must be provided using: MONOTASK_BINARY:
//
//   - specify binary manually.
//   - use dotenv files.
//   - use make goal
//
// User might set BINARY_GOCOVERDIR to pass to binary as GOCOVERDIR env variable, thus collecting coverage data.
package integtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/IlyasYOY/exectest"
)

// TestIntegration walks all testdata files and creates test-cases ([exectest]) for running binary of them.
func TestIntegration(t *testing.T) {
	monotaskBinary := os.Getenv("MONOTASK_BINARY")
	if monotaskBinary == "" {
		t.Skipf("MONOTASK_BINARY flag is required")
	}
	if _, err := os.Stat(monotaskBinary); os.IsNotExist(err) {
		t.Skipf("MONOTASK_BINARY does not exist: %s", monotaskBinary)
	}

	testdataDir := "testdata"
	var files []string
	err := filepath.WalkDir(testdataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".txt") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to find test files: %v", err)
	}

	for _, file := range files {
		relPath, _ := filepath.Rel(testdataDir, file)
		t.Run(relPath, func(t *testing.T) {
			exectest.ExecuteForFile(t, monotaskBinary, file, func(cmd *exec.Cmd) {
				cmd.Env = []string{"GOCOVERDIR=" + os.Getenv("BINARY_GOCOVERDIR")}
			})
		})
	}
}
