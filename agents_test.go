package ui_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// A release must not leave the install steps naming an old version.
func TestAgentsNamesCurrentVersion(t *testing.T) {
	raw, err := os.ReadFile("package.json")
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct{ Version string }
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"AGENTS.md", "README.md"} {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"ai-ui@v" + pkg.Version, "ai-ui#v" + pkg.Version} {
			if !strings.Contains(string(raw), want) {
				t.Errorf("%s doesn't contain %q", file, want)
			}
		}
	}
}
