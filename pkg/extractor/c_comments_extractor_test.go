package extractor

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestCCommentsExtractor_SingleLine(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
// TODO: fix the bug
int main() { return 0; }
`)
	extractor := NewCCommentsExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.c"), Line: 1, Column: 1, Type: "TODO", Message: "fix the bug"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_BlockComment(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
/* BUG: crash here */
int main() { return 0; }
`)
	extractor := NewCCommentsExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.c"), Line: 1, Column: 1, Type: "BUG", Message: "crash here */"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_MultiLineBlock(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
/*
NOTE: important note
*/
int main() { return 0; }
`)
	extractor := NewCCommentsExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.c"), Line: 2, Column: 1, Type: "NOTE", Message: "important note"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_EmptyFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
`)
	extractor := NewCCommentsExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTasks := ([]Task)(nil)
	if diff := cmp.Diff(expectedTasks, tasks, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("expected no tasks (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_NoMarkers(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
// this is a comment
/* another comment */
int main() { return 0; }
`)
	extractor := NewCCommentsExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTasks := ([]Task)(nil)
	if diff := cmp.Diff(expectedTasks, tasks, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("expected no tasks (-want +got):\n%s", diff)
	}
}
