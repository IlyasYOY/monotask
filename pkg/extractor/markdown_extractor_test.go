package extractor

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestMarkdownExtractor_SingleCheckbox(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.md
- [ ] incomplete task
`)
	extractor := NewMarkdownExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.md"), Line: 1, Column: 1, Type: "CHECKBOX", Message: "incomplete task"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestMarkdownExtractor_MultipleCheckboxes(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.md
- [ ] task 1
- [x] task 2
- [ ] task 3
`)
	extractor := NewMarkdownExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.md"), Line: 1, Column: 1, Type: "CHECKBOX", Message: "task 1"},
		{File: filepath.Join(dir, "test.md"), Line: 3, Column: 1, Type: "CHECKBOX", Message: "task 3"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestMarkdownExtractor_NoCheckboxes(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.md
- item
- [x] done
`)
	extractor := NewMarkdownExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTasks := ([]Task)(nil)
	if diff := cmp.Diff(expectedTasks, tasks, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("expected no tasks (-want +got):\n%s", diff)
	}
}

func TestMarkdownExtractor_EmptyFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.md
`)
	extractor := NewMarkdownExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTasks := ([]Task)(nil)
	if diff := cmp.Diff(expectedTasks, tasks, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("expected no tasks (-want +got):\n%s", diff)
	}
}
