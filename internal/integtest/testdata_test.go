package integtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

// TestIntegration test walks the `testdata` directory for `*.txt` files.
// Each file describes a temporary directory layout (using the `-#file:` marker)
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
	files, err := filepath.Glob(filepath.Join(testdataDir, "*.txt"))
	if err != nil {
		t.Fatalf("Failed to find test files: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
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
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to run monotask binary: %v", err)
			}
			stdout := stdoutBuilder.String()

			actualSet := make(map[string]bool)
			for line := range strings.SplitSeq(strings.TrimSpace(stdout), "\n") {
				if strings.TrimSpace(line) != "" {
					shortenPath := strings.TrimPrefix(line, tempDir)
					shortenPath = strings.TrimPrefix(shortenPath, "/")
					actualSet[shortenPath] = true
				}
			}

			expectedSet := make(map[string]bool)
			for line := range strings.SplitSeq(strings.TrimSpace(header), "\n") {
				if strings.TrimSpace(line) != "" {
					expectedSet[line] = true
				}
			}
			if diff := cmp.Diff(expectedSet, actualSet); diff != "" {
				t.Errorf("Error matching stdout: %s", diff)
			}
		})
	}
}
