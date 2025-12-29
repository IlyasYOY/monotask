package extractor_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/extractor"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestDirectoryExtractor_SingleFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
// TODO: test
`)
	extr := extractor.NewDirectoryExtractor(dir)
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
			Message: "test",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestDirectoryExtractor_MultipleFiles(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
// TODO: c task
-#file:test.md
- [ ] md task
-#file:test.txt
some text
`)
	extr := extractor.NewDirectoryExtractor(dir)
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
			Message: "c task",
		},
		{
			File:    filepath.Join(dir, "test.md"),
			Line:    1,
			Column:  1,
			Type:    "CHECKBOX",
			Message: "md task",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestDirectoryExtractor_NestedDirectories(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:subdir/test.c
// BUG: nested
`)
	extr := extractor.NewDirectoryExtractor(dir)
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff([]extractor.Task{
		{
			File:    filepath.Join(dir, "subdir", "test.c"),
			Line:    1,
			Column:  1,
			Type:    "BUG",
			Message: "nested",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestDirectoryExtractor_EmptyDirectory(t *testing.T) {
	dir := filestest.RenderDir(t, ``)
	extr := extractor.NewDirectoryExtractor(dir)
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty tasks but was: %v", tasks)
	}
}
