package integtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/IlyasYOY/monotask/internal/pkg/exectest"
)

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
