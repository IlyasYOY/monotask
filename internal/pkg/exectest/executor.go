// Package exectest is [os/exec] package testing facilities.
//
// The main goal of the package: declarative testing of any executable..
package exectest

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	// TODO: replace -- with something more unique.
	// This is the comment start in Lua, so the tests look odd.
	filePrefix   = "--file:"
	stdoutPrefix = "--stdout"
	argPrefix    = "--arg:"
)

type cmdOption func(*exec.Cmd)

// Execute is the main testing facility of the package.
//
// Function consists of the steps:
//
//   - [prepareScheme]: create directory, parse args and required stdout.
//   - execute given binary in the prepared conditions.
//   - assert results of the binary evaluation.
//
// Examples:
//
//	--file:a.txt
//	--file:.b.txt
//	--arg:-a
//	--stdout
//	.
//	..
//	.b.txt
//	a.txt
//
// This is a desciption of the command `ls -a` run in the
// directory with a.txt and .b.txt files.
func Execute(t *testing.T, binary, scheme string, opts ...cmdOption) {
	t.Helper()
	want, args, dir := prepareScheme(t, scheme)

	got := executeCommand(t, binary, dir, args, opts)

	wantLines := slices.Collect(strings.Lines(want))
	gotLines := slices.Collect(strings.Lines(got))
	if diff := cmp.Diff(wantLines, gotLines); diff != "" {
		t.Fatalf("Error matching stdout (-missing line, +extra line): %s\n\nFor schema: \n%s", diff, scheme)
	}
}

func executeCommand(t *testing.T, binary string, dir string, args []string, opts []cmdOption) string {
	t.Helper()

	cmd := exec.Command(binary)
	var stdoutBuilder strings.Builder
	cmd.Stdout = &stdoutBuilder
	var stderrBuilder strings.Builder
	cmd.Stderr = &stderrBuilder
	cmd.Dir = dir
	cmd.Args = append(cmd.Args, args...)
	for _, opt := range opts {
		opt(cmd)
	}
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run monotask binary: %v\nError message:\n%s", err, stderrBuilder.String())
	}
	return stdoutBuilder.String()
}

// ExecuteForFile the same as the [Execute] but uses a file (path) with a scheme.
func ExecuteForFile(t *testing.T, binary string, file string, opts ...cmdOption) {
	t.Helper()
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", file, err)
	}
	Execute(t, binary, string(content), opts...)
}

func prepareScheme(t *testing.T, header string) (string, []string, string) {
	t.Helper()

	dir := t.TempDir()

	var stdout strings.Builder
	var args []string
	files := make(map[string]string)

	// TODO: Replace with enum.
	// FSM is always good for readability.
	var isStdout bool
	var isFile bool

	var currentFileName string
	var currentFile strings.Builder

	saveFile := func(name string) {
		if isFile {
			resultPath := filepath.Join(dir, currentFileName)
			files[resultPath] = currentFile.String()
		}
		currentFileName = name
		currentFile.Reset()
	}

	for line := range strings.Lines(header) {
		if strings.HasPrefix(line, stdoutPrefix) {
			saveFile("")
			isStdout = true
			isFile = false
			continue
		}
		if fileName, ok := strings.CutPrefix(line, filePrefix); ok {
			saveFile(strings.TrimSpace(fileName))
			isFile = true
			isStdout = false
			continue
		}

		if arg, ok := strings.CutPrefix(line, argPrefix); ok {
			arg = strings.TrimSpace(arg)
			arg = evaluateVariables(arg, dir)
			args = append(args, arg)
			continue
		}
		if isStdout {
			line = evaluateVariables(line, dir)
			stdout.WriteString(line)
			continue
		}
		if isFile {
			currentFile.WriteString(line)
		}
	}
	if isFile {
		saveFile("")
	}

	for path, content := range files {
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, []byte(content), 0o644)
	}

	return stdout.String(), args, dir
}

func evaluateVariables(data string, dir string) string {
	data = strings.ReplaceAll(data, "{dir}", dir)
	return data
}
