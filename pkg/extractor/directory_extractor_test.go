package extractor

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestDirectoryExtractor_SingleFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
// TODO: test
`)
	extractor := NewDirectoryExtractor(dir)
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.c"), Line: 1, Column: 1, Type: "TODO", Message: "test"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestDirectoryExtractor_MultipleFiles(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
// TODO: c task
#test.md
- [ ] md task
#test.txt
some text
`)
	extractor := NewDirectoryExtractor(dir)
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.c"), Line: 1, Column: 1, Type: "TODO", Message: "c task"},
		{File: filepath.Join(dir, "test.md"), Line: 1, Column: 1, Type: "CHECKBOX", Message: "md task"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestDirectoryExtractor_NestedDirectories(t *testing.T) {
	dir := filestest.RenderDir(t, `
#subdir/test.c
// BUG: nested
`)
	extractor := NewDirectoryExtractor(dir)
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "subdir", "test.c"), Line: 1, Column: 1, Type: "BUG", Message: "nested"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestDirectoryExtractor_EmptyDirectory(t *testing.T) {
	dir := filestest.RenderDir(t, ``)
	extractor := NewDirectoryExtractor(dir)
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTasks := ([]Task)(nil)
	if diff := cmp.Diff(expectedTasks, tasks, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("expected no tasks (-want +got):\n%s", diff)
	}
}
