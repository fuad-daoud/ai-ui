// Command build writes dist/ui.css and dist/ui-scoped.css from the sources under css/.
// Run it from the repository root: go run ./cmd/build
package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/fuad-daoud/ai-ui/internal/build"
)

func main() {
	if _, err := os.Stat("go.mod"); err != nil {
		fail("run from the ai-ui repository root (no go.mod here)")
	}
	if err := os.MkdirAll("dist", 0o755); err != nil {
		fail(err.Error())
	}
	for _, t := range []struct {
		file  string
		build func(fs.FS) ([]byte, error)
	}{
		{"dist/ui.css", build.CSS},
		{"dist/ui-scoped.css", build.ScopedCSS},
	} {
		out, err := t.build(os.DirFS("."))
		if err != nil {
			fail(err.Error())
		}
		if err := os.WriteFile(t.file, out, 0o644); err != nil {
			fail(err.Error())
		}
		fmt.Printf("wrote %s (%d bytes)\n", t.file, len(out))
	}
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "build:", msg)
	os.Exit(1)
}
