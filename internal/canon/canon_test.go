package canon_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fuad-daoud/ai-ui/internal/canon"
)

func TestSharedVectors(t *testing.T) {
	b, err := os.ReadFile("../../testdata/canon/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct{ Name, In, Out string }
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no vectors")
	}
	for _, c := range cases {
		got, err := canon.HTML(c.In)
		if err != nil {
			t.Errorf("%s: %v", c.Name, err)
			continue
		}
		if got != c.Out {
			t.Errorf("%s:\n got %q\nwant %q", c.Name, got, c.Out)
		}
	}
}
