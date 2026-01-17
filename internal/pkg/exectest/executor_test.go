package exectest_test

import (
	"testing"

	"github.com/IlyasYOY/monotask/internal/pkg/exectest"
)

func TestExecute(t *testing.T) {
	t.Run("ls", func(t *testing.T) {
		t.Run("stdout after files", func(t *testing.T) {
			exectest.Execute(t, "ls", `
--file:a.txt
--file:b.txt
--stdout
a.txt
b.txt
`)
		})

		t.Run("stdout before files", func(t *testing.T) {
			exectest.Execute(t, "ls", `
--stdout
a.txt
b.txt
--file:a.txt
--file:b.txt
`)
		})

		t.Run("stdout between files", func(t *testing.T) {
			exectest.Execute(t, "ls", `
--file:a.txt
--stdout
a.txt
b.txt
--file:b.txt
`)
		})

		t.Run("no hidden files showed", func(t *testing.T) {
			exectest.Execute(t, "ls", `
--file:a.txt
--file:.b.txt
--stdout
a.txt
`)
		})

		t.Run("argument passed to show hidden file", func(t *testing.T) {
			exectest.Execute(t, "ls", `
--file:a.txt
--file:.b.txt
--arg:-a
--stdout
.
..
.b.txt
a.txt
`)
		})
	})
}
