package filestest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/internal/pkg/filestest"
)

func TestRenderDir_Header(t *testing.T) {
	_, header := filestest.RenderDir(t, `This is header.
--file:file.txt
line1
line2`)

	if diff := cmp.Diff("This is header.\n", header); diff != "" {
		t.Errorf("header mismatch (-want +got):\n%s", diff)
	}
}

func TestRenderDir_BasicFile(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `--file:file.txt
line1
line2`)

	filePath := filepath.Join(dir, "file.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "line1\nline2\n"
	if diff := cmp.Diff(expected, string(content)); diff != "" {
		t.Errorf("content mismatch (-want +got):\n%s", diff)
	}
}

func TestRenderDir_MultipleFiles(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `--file:file1.txt
content1
--file:file2.txt
content2`)

	file1Path := filepath.Join(dir, "file1.txt")
	content1, err := os.ReadFile(file1Path)
	if err != nil {
		t.Fatalf("Failed to read file file1.txt: %v", err)
	}
	file2Path := filepath.Join(dir, "file2.txt")
	content2, err := os.ReadFile(file2Path)
	if err != nil {
		t.Fatalf("Failed to read file file2.txt: %v", err)
	}

	if diff := cmp.Diff("content1\n", string(content1)); diff != "" {
		t.Errorf("file1.txt content mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff("content2\n", string(content2)); diff != "" {
		t.Errorf("file2.txt content mismatch (-want +got):\n%s", diff)
	}
}

func TestRenderDir_NestedDirectories(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `--file:dir/subdir/file.txt
nested content`)

	filePath := filepath.Join(dir, "dir", "subdir", "file.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "nested content\n"
	if diff := cmp.Diff(expected, string(content)); diff != "" {
		t.Errorf("content mismatch (-want +got):\n%s", diff)
	}
}

func TestRenderDir_EmptyContent(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `--file:empty.txt
`)

	filePath := filepath.Join(dir, "empty.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if diff := cmp.Diff([]byte{}, content); diff != "" {
		t.Errorf("expected empty content (-want +got):\n%s", diff)
	}
}

func TestRenderDir_BlankLines(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `--file:file.txt

line

`)

	filePath := filepath.Join(dir, "file.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "line\n"
	if diff := cmp.Diff(expected, string(content)); diff != "" {
		t.Errorf("content mismatch (-want +got):\n%s", diff)
	}
}

func TestRenderDir_NoMarkers(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `plain text
no markers`)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			t.Errorf("Unexpected file found: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk directory: %v", err)
	}
}
