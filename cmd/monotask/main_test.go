package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrintsVersionWithLongFlag(t *testing.T) {
	withVersion("v0.2.0+1f4b76a", func() {
		var stdout bytes.Buffer

		exitCode := run([]string{"--version", "/path/that/does/not/exist"}, &stdout)

		if exitCode != 0 {
			t.Fatalf("run() exit code = %d, want 0", exitCode)
		}
		if got, want := stdout.String(), "v0.2.0+1f4b76a\n"; got != want {
			t.Fatalf("stdout = %q, want %q", got, want)
		}
	})
}

func TestRunPrintsVersionWithShortFlag(t *testing.T) {
	withVersion("v0.2.0+1f4b76a", func() {
		var stdout bytes.Buffer

		exitCode := run([]string{"-v", "/path/that/does/not/exist"}, &stdout)

		if exitCode != 0 {
			t.Fatalf("run() exit code = %d, want 0", exitCode)
		}
		if got, want := stdout.String(), "v0.2.0+1f4b76a\n"; got != want {
			t.Fatalf("stdout = %q, want %q", got, want)
		}
	})
}

func TestRunScansPathArgument(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(filePath, []byte("// TODO: Implement error handling\n"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var stdout bytes.Buffer

	exitCode := run([]string{dir}, &stdout)

	if exitCode != 0 {
		t.Fatalf("run() exit code = %d, want 0", exitCode)
	}
	if got := stdout.String(); !strings.Contains(got, "TODO: Implement error handling") {
		t.Fatalf("stdout = %q, want extracted TODO", got)
	}
}

func withVersion(version string, fn func()) {
	oldVersionString := versionString
	versionString = func() string {
		return version
	}
	defer func() {
		versionString = oldVersionString
	}()

	fn()
}
