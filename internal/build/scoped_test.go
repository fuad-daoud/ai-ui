package build_test

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/fuad-daoud/ai-ui/internal/build"
)

const scope = `:where([class^="ui-"], [class*=" ui-"], [class^="ui-"] *, [class*=" ui-"] *)`

func scopedWith(reset, base string) (string, error) {
	fsys := sources()
	fsys["css/reset.css"] = &fstest.MapFile{Data: []byte(reset)}
	fsys["css/base.css"] = &fstest.MapFile{Data: []byte(base)}
	out, err := build.ScopedCSS(fsys)
	return string(out), err
}

func TestScopedRewritesSelectors(t *testing.T) {
	out, err := scopedWith(`
*, *::before { box-sizing: border-box; }
html { text-size-adjust: 100%; }
body { min-height: 100vh; }
h1, h2 { text-wrap: balance; }
`, `
/* a comment with { a brace */
body, a { color: red; }
:focus-visible { outline: 2px solid; }
::selection { color: blue; }
a:is(.x, .y) { color: green; }
@media (prefers-reduced-motion: reduce) { :root { --dur: 0ms; } }
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"*" + scope + ", *" + scope + "::before { box-sizing: border-box; }",
		"h1" + scope + ", h2" + scope + " { text-wrap: balance; }",
		"a" + scope + " { color: red; }",
		":focus-visible" + scope + " { outline: 2px solid; }",
		scope + "::selection { color: blue; }",
		"a:is(.x, .y)" + scope + " { color: green; }",
		":root { --dur: 0ms; }",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, gone := range []string{"text-size-adjust", "min-height", "body", "html"} {
		if strings.Contains(out, gone) {
			t.Errorf("%q should be dropped:\n%s", gone, out)
		}
	}
}

func TestScopedFailsOnWhatItCannotRewrite(t *testing.T) {
	for name, src := range map[string]string{
		"at-rule":         `@supports (display: grid) { a { color: red; } }`,
		"page in a chain": `body a { color: red; }`,
		"html in a chain": `html:has(a) p { color: red; }`,
		"root in a chain": `:root a { color: red; }`,
		"stray text":      `a { color: red; } b`,
		"unbalanced":      `a { color: red;`,
		"open comment":    `a { color: red; } /* never closed`,
		"empty selector":  `a, { color: red; }`,
	} {
		t.Run(name, func(t *testing.T) {
			if out, err := scopedWith(src, ""); err == nil {
				t.Fatalf("want an error, got:\n%s", out)
			} else if !strings.Contains(err.Error(), "css/reset.css") {
				t.Errorf("err = %v, want it to name css/reset.css", err)
			}
		})
	}
}

// The real sources: every rule left in the reset and base carries the scope.
func TestScopedSourcesAreAllScoped(t *testing.T) {
	full, err := build.ScopedCSS(os.DirFS("../.."))
	if err != nil {
		t.Fatal(err)
	}
	s := string(full)
	reset := s[strings.Index(s, "/* css/reset.css */"):strings.Index(s, "@layer ui.base {")]
	base := s[strings.Index(s, "/* css/base.css */"):strings.Index(s, "@layer ui.components {")]
	rule := regexp.MustCompile(`(?m)^([^\s@}/][^{]*)\{`)
	for name, part := range map[string]string{"reset": reset, "base": base} {
		rules := rule.FindAllStringSubmatch(part, -1)
		if len(rules) == 0 {
			t.Fatalf("%s: no rules found", name)
		}
		for _, m := range rules {
			sel := strings.TrimSpace(m[1])
			if sel == ":root" {
				continue
			}
			for _, one := range strings.Split(strings.ReplaceAll(sel, scope, "§"), ", ") {
				if !strings.Contains(one, "§") {
					t.Errorf("%s: selector %q is not scoped", name, one)
				}
			}
		}
	}
	for _, p := range []string{"css/tokens.css", "css/utilities.css", "css/components/button.css"} {
		src, err := os.ReadFile("../../" + p)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(s, strings.TrimRight(string(src), "\n")) {
			t.Errorf("%s was not copied unchanged", p)
		}
	}
}
