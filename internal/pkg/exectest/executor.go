// Package exectest is [os/exec] package testing facilities.
//
// The main goal of the package: declarative testing of any executable..
package exectest

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	// TODO: replace -- with something more unique.
	// This is the comment start in Lua, so the tests look odd.
	filePrefix       = "--file:"
	stdoutPrefix     = "--stdout"
	stderrPrefix     = "--stderr"
	argPrefix        = "--arg:"
	returnCodePrefix = "--return-code:"
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
	wantStdout, wantStderr, wantReturnCode, args, dir := prepareScheme(t, scheme)

	gotStdout, gotStderr, gotReturnCode := executeCommand(t, binary, dir, args, opts)

	assertReturnCode(t, wantReturnCode, gotReturnCode)
	assertNoDiff(t, "stderr", wantStdout, gotStdout)
	assertNoDiff(t, "stdout", wantStderr, gotStderr)
}

func assertReturnCode(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("Failed to match return code: want %d, got %d", want, got)
	}
}

func assertNoDiff(t *testing.T, name string, want string, got string) {
	t.Helper()
	wantLines := slices.Collect(strings.Lines(want))
	gotLines := slices.Collect(strings.Lines(got))
	if diff := cmp.Diff(wantLines, gotLines); diff != "" {
		t.Errorf("Failed matching %s (-missing line, +extra line): %s", name, diff)
	}
}

func executeCommand(t *testing.T, binary string, dir string, args []string, opts []cmdOption) (string, string, int) {
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

	// this is intentional, we will assert exit code manually
	_ = cmd.Run()

	return stdoutBuilder.String(), stderrBuilder.String(), cmd.ProcessState.ExitCode()
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

func prepareScheme(t *testing.T, scheme string) (string, string, int, []string, string) {
	t.Helper()

	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("Test scheme: %s", scheme)
		}
	})

	// TODO: Make test fail if the same field defined twice.
	var stdout strings.Builder
	var stderr strings.Builder
	var returnCode int
	var args []string

	dir := t.TempDir()
	files := make(map[string]string)

	// TODO: Replace with enum.
	// FSM is always good for readability.
	var isStdout bool
	var isStderr bool
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

	for line := range strings.Lines(scheme) {
		if strings.HasPrefix(line, stderrPrefix) {
			saveFile("")
			isFile = false
			isStderr = true
			isStdout = false
			continue
		}
		if strings.HasPrefix(line, stdoutPrefix) {
			saveFile("")
			isFile = false
			isStderr = false
			isStdout = true
			continue
		}
		if fileName, ok := strings.CutPrefix(line, filePrefix); ok {
			saveFile(strings.TrimSpace(fileName))
			isFile = true
			isStderr = false
			isStdout = false
			continue
		}

		if rtCodeText, ok := strings.CutPrefix(line, returnCodePrefix); ok {
			rtCodeText = strings.TrimSpace(rtCodeText)
			var err error
			returnCode, err = strconv.Atoi(rtCodeText)
			if err != nil {
				t.Fatalf("Failed to convert return code %q to int: %s", rtCodeText, err)
			}
			continue
		}
		if arg, ok := strings.CutPrefix(line, argPrefix); ok {
			arg = strings.TrimSpace(arg)
			arg = evaluateVariables(arg, dir)
			args = append(args, arg)
			continue
		}
		if isStderr {
			line = evaluateVariables(line, dir)
			stderr.WriteString(line)
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
		fileDir := filepath.Dir(path)
		if err := os.MkdirAll(fileDir, 0o755); err != nil {
			t.Fatalf("Failed to create directory (%q) for test file: %s", fileDir, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to write file (%v): %s", path, err)
		}
	}

	return stdout.String(), stderr.String(), returnCode, args, dir
}

func evaluateVariables(data string, dir string) string {
	data = strings.ReplaceAll(data, "{dir}", dir)
	return data
}
