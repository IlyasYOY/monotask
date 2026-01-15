package integtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/IlyasYOY/monotask/internal/pkg/filestest"
	"github.com/google/go-cmp/cmp"
)

// TestIntegration test walks the `testdata` directory recursively for `*.txt` files.
// Each file describes a temporary directory layout (using the `--file:` marker)
// and the expected monotask output in its header.
//
// For every test file the test:
//   - Renders the described directory structure into a temporary directory
//     via filestest.RenderDir.
//   - Executes the monotask binary (path supplied by the MONOTASK_BINARY env var)
//     with the temporary directory as its working directory.
//   - Captures the binary’s stdout, normalises the paths, and builds a set of
//     reported task lines.
//   - Builds a set of expected lines from the file header.
//   - Compares the two sets using go‑cmp, failing the test on any mismatch.
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
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", file, err)
			}
			tempDir, header := filestest.RenderDir(t, string(content))

			cmd := exec.Command(monotaskBinary)
			var stdoutBuilder strings.Builder
			cmd.Stdout = &stdoutBuilder
			cmd.Stderr = os.Stderr
			cmd.Dir = tempDir
			cmd.Args = append(cmd.Args, tempDir)
			cmd.Env = []string{"GOCOVERDIR="+os.Getenv("GOCOVERDIR")}
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to run monotask binary: %v", err)
			}
			stdout := stdoutBuilder.String()

			var actualLines []string
			for line := range strings.SplitSeq(strings.TrimSpace(stdout), "\n") {
				if strings.TrimSpace(line) != "" {
					shortenPath := strings.TrimPrefix(line, tempDir)
					shortenPath = strings.TrimPrefix(shortenPath, "/")
					actualLines = append(actualLines, shortenPath)
				}
			}

			var wantLines []string
			for line := range strings.SplitSeq(strings.TrimSpace(header), "\n") {
				if strings.TrimSpace(line) != "" {
					wantLines = append(wantLines, line)
				}
			}

			if diff := cmp.Diff(wantLines, actualLines); diff != "" {
				t.Errorf("Error matching stdout: %s\n\nFor schema: \n%s", diff, string(content))
			}
		})
	}
}
