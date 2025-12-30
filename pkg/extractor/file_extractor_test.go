package extractor_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/extractor"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestFileExtractor_CFile(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.c
// TODO: c
`)
	extr := extractor.NewFileExtractor(filepath.Join(dir, "test.c"))
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff([]extractor.Task{
		{
			File:    filepath.Join(dir, "test.c"),
			Line:    1,
			Column:  1,
			Type:    "TODO",
			Message: "c",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestFileExtractor_HFile(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.h
// BUG: h
`)
	extr := extractor.NewFileExtractor(filepath.Join(dir, "test.h"))
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff([]extractor.Task{
		{
			File:    filepath.Join(dir, "test.h"),
			Line:    1,
			Column:  1,
			Type:    "BUG",
			Message: "h",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestFileExtractor_MarkdownFile(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.md
- [ ] md
`)
	extr := extractor.NewFileExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff([]extractor.Task{
		{
			File:    filepath.Join(dir, "test.md"),
			Line:    1,
			Column:  1,
			Type:    "CHECKBOX",
			Message: "md",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestFileExtractor_UnsupportedFile(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.txt
TODO: txt
`)
	extr := extractor.NewFileExtractor(filepath.Join(dir, "test.txt"))
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty tasks but was: %v", tasks)
	}
}
