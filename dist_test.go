package ui_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/fuad-daoud/ai-ui/internal/build"
)

func TestDistIsFresh(t *testing.T) {
	want, err := build.CSS(os.DirFS("."))
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("dist/ui.css")
	if err != nil {
		t.Fatalf("%v (run: go run ./cmd/build)", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("dist/ui.css is stale: run `go run ./cmd/build` and commit it")
	}
}
