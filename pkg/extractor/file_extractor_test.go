package extractor

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestFileExtractor_CFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.c
// TODO: c
`)
	extractor := NewFileExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.c"), Line: 1, Column: 1, Type: "TODO", Message: "c"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestFileExtractor_HFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.h
// BUG: h
`)
	extractor := NewFileExtractor(filepath.Join(dir, "test.h"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.h"), Line: 1, Column: 1, Type: "BUG", Message: "h"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestFileExtractor_MarkdownFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.md
- [ ] md
`)
	extractor := NewFileExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Task{
		{File: filepath.Join(dir, "test.md"), Line: 1, Column: 1, Type: "CHECKBOX", Message: "md"},
	}
	if diff := cmp.Diff(expected, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestFileExtractor_UnsupportedFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
#test.txt
TODO: txt
`)
	extractor := NewFileExtractor(filepath.Join(dir, "test.txt"))
	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTasks := []Task{}
	if diff := cmp.Diff(expectedTasks, tasks, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("expected no tasks (-want +got):\n%s", diff)
	}
}
