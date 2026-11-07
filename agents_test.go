package ui_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// A release must not leave AGENTS.md telling agents to install an old version.
func TestAgentsNamesCurrentVersion(t *testing.T) {
	raw, err := os.ReadFile("package.json")
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct{ Version string }
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatal(err)
	}
	agents, err := os.ReadFile("AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"ai-ui@v" + pkg.Version, "ai-ui#v" + pkg.Version} {
		if !strings.Contains(string(agents), want) {
			t.Errorf("AGENTS.md doesn't contain %q", want)
		}
	}
}
