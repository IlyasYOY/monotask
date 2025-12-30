package extractor_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ilyasyoy/monotask/pkg/extractor"
	"github.com/ilyasyoy/monotask/pkg/filestest"
)

func TestMarkdownExtractor_SingleCheckbox(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.md
- [ ] incomplete task
`)
	extr := extractor.NewMarkdownExtractor(filepath.Join(dir, "test.md"))
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
			Message: "incomplete task",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestMarkdownExtractor_MultipleCheckboxes(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.md
- [ ] task 1
- [x] task 2
- [ ] task 3
`)
	extr := extractor.NewMarkdownExtractor(filepath.Join(dir, "test.md"))
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
			Message: "task 1",
		},
		{
			File:    filepath.Join(dir, "test.md"),
			Line:    3,
			Column:  1,
			Type:    "CHECKBOX",
			Message: "task 3",
		},
	}, tasks); diff != "" {
		t.Errorf("tasks mismatch (-want +got):\n%s", diff)
	}
}

func TestMarkdownExtractor_NoCheckboxes(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.md
- item
- [x] done
`)
	extr := extractor.NewMarkdownExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty tasks but was: %v", tasks)
	}
}

func TestMarkdownExtractor_EmptyFile(t *testing.T) {
	dir, _ := filestest.RenderDir(t, `
-#file:test.md
`)
	extr := extractor.NewMarkdownExtractor(filepath.Join(dir, "test.md"))
	tasks, err := extr.Extract(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty tasks but was: %v", tasks)
	}
}
