// Command build writes dist/ui.css from the sources under css/.
// Run it from the repository root: go run ./cmd/build
package main

import (
	"fmt"
	"os"

	"github.com/fuad-daoud/ai-ui/internal/build"
)

func main() {
	if _, err := os.Stat("go.mod"); err != nil {
		fail("run from the ai-ui repository root (no go.mod here)")
	}
	out, err := build.CSS(os.DirFS("."))
	if err != nil {
		fail(err.Error())
	}
	if err := os.MkdirAll("dist", 0o755); err != nil {
		fail(err.Error())
	}
	if err := os.WriteFile("dist/ui.css", out, 0o644); err != nil {
		fail(err.Error())
	}
	fmt.Printf("wrote dist/ui.css (%d bytes)\n", len(out))
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "build:", msg)
	os.Exit(1)
}
