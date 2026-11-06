package ui_test

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"github.com/fuad-daoud/ai-ui/internal/build"
)

func TestDistIsFresh(t *testing.T) {
	for file, build := range map[string]func(fs.FS) ([]byte, error){
		"dist/ui.css":        build.CSS,
		"dist/ui-scoped.css": build.ScopedCSS,
	} {
		want, err := build(os.DirFS("."))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("%v (run: go run ./cmd/build)", err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s is stale: run `go run ./cmd/build` and commit it", file)
		}
	}
}
