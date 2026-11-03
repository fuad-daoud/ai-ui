package build_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/fuad-daoud/ai-ui/internal/build"
)

func sources() fstest.MapFS {
	return fstest.MapFS{
		"css/reset.css":            {Data: []byte("/*R*/\n")},
		"css/tokens.css":           {Data: []byte("/*T*/\n")},
		"css/base.css":             {Data: []byte("/*B*/\n")},
		"css/components/zeta.css":  {Data: []byte("/*Z*/\n")},
		"css/components/alpha.css": {Data: []byte("/*A*/\n")},
		"css/utilities.css":        {Data: []byte("/*U*/\n")},
	}
}

func TestCSSLayersAndOrder(t *testing.T) {
	out, err := build.CSS(sources())
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	first := strings.Index(s, "@layer")
	if first < 0 || !strings.HasPrefix(s[first:], "@layer ui.reset, ui.base, ui.components, ui.utilities;") {
		t.Fatalf("the first @layer is not the order statement:\n%s", s)
	}
	last := -1
	for _, m := range []string{
		"@layer ui.reset {", "/*R*/",
		"@layer ui.base {", "/*T*/", "/*B*/",
		"@layer ui.components {", "/*A*/", "/*Z*/",
		"@layer ui.utilities {", "/*U*/",
	} {
		i := strings.Index(s, m)
		if i < 0 {
			t.Fatalf("missing %q in:\n%s", m, s)
		}
		if i <= last {
			t.Fatalf("%q is out of order in:\n%s", m, s)
		}
		last = i
	}
}

func TestCSSMissingSource(t *testing.T) {
	fsys := sources()
	delete(fsys, "css/base.css")
	if _, err := build.CSS(fsys); err == nil || !strings.Contains(err.Error(), "css/base.css") {
		t.Fatalf("err = %v, want it to name css/base.css", err)
	}
}

func TestCSSWithoutComponents(t *testing.T) {
	fsys := sources()
	delete(fsys, "css/components/zeta.css")
	delete(fsys, "css/components/alpha.css")
	out, err := build.CSS(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "@layer ui.components {") {
		t.Error("an empty ui.components block was written")
	}
}
