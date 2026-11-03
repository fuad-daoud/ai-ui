package ui_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	ui "github.com/fuad-daoud/ai-ui"
	"github.com/fuad-daoud/ai-ui/internal/canon"
)

var update = flag.Bool("update", false, "rewrite testdata/parity/*.html from the Go partials")

type parityCase struct {
	Go *struct {
		Template string         `json:"template"`
		Args     map[string]any `json:"args"`
	} `json:"go"`
	React *struct {
		Component string         `json:"component"`
		Props     map[string]any `json:"props"`
	} `json:"react"`
	Error string `json:"error"`
}

// parityCases loads every fixture, keyed by its path without ".json".
func parityCases(t *testing.T) (map[string]parityCase, []string) {
	t.Helper()
	files, err := filepath.Glob("testdata/parity/*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no fixtures in testdata/parity")
	}
	cases := make(map[string]parityCase, len(files))
	var keys []string
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var c parityCase
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&c); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if c.Go == nil && c.React == nil {
			t.Fatalf("%s: needs a go part, a react part, or both", f)
		}
		key := strings.TrimSuffix(f, ".json")
		cases[key] = c
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return cases, keys
}

func mustCanon(t *testing.T, s string) string {
	t.Helper()
	c, err := canon.HTML(s)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestParityGo(t *testing.T) {
	tmpl, err := ui.Parse(template.New("parity"))
	if err != nil {
		t.Fatal(err)
	}
	cases, keys := parityCases(t)
	for _, key := range keys {
		c := cases[key]
		if c.Go == nil {
			continue // container markup: React side only, golden written by hand
		}
		t.Run(filepath.Base(key), func(t *testing.T) {
			var b strings.Builder
			err := tmpl.ExecuteTemplate(&b, c.Go.Template, c.Go.Args)
			if c.Error != "" {
				if err == nil || !strings.Contains(err.Error(), c.Error) {
					t.Fatalf("err = %v, want it to contain %q", err, c.Error)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(b.String(), "no value") {
				t.Fatalf("a missing key leaked into the markup: %s", b.String())
			}
			got := mustCanon(t, b.String())
			golden := key + ".html"
			if *update {
				if err := os.WriteFile(golden, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("%v (run: go test . -run TestParityGo -update)", err)
			}
			if w := mustCanon(t, string(want)); got != w {
				t.Errorf("markup differs from %s\n got: %s\nwant: %s", golden, got, w)
			}
		})
	}
}
