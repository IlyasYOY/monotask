package extractor_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/extractor"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestCCommentsExtractor_SingleLine(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
// TODO: fix the bug
int main() { return 0; }
`)
	extr := extractor.NewCCommentsExtractor(filepath.Join(dir, "test.c"))

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
			Message: "fix the bug",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_BlockComment(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
/* BUG: crash here */
int main() { return 0; }
`)
	extr := extractor.NewCCommentsExtractor(filepath.Join(dir, "test.c"))

	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff([]extractor.Task{
		{
			File:    filepath.Join(dir, "test.c"),
			Line:    1,
			Column:  1,
			Type:    "BUG",
			Message: "crash here */",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_MultiLineBlock(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
/*
NOTE: important note
*/
int main() { return 0; }
`)
	extr := extractor.NewCCommentsExtractor(filepath.Join(dir, "test.c"))

	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff([]extractor.Task{
		{
			File:    filepath.Join(dir, "test.c"),
			Line:    2,
			Column:  1,
			Type:    "NOTE",
			Message: "important note",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestCCommentsExtractor_EmptyFile(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
`)
	extractor := extractor.NewCCommentsExtractor(filepath.Join(dir, "test.c"))

	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty tasks but was: %v", tasks)
	}
}

func TestCCommentsExtractor_NoMarkers(t *testing.T) {
	dir := filestest.RenderDir(t, `
-#file:test.c
// this is a comment
/* another comment */
int main() { return 0; }
`)
	extractor := extractor.NewCCommentsExtractor(filepath.Join(dir, "test.c"))

	tasks, err := extractor.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty tasks but was: %v", tasks)
	}
}
