package ui_test

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/fuad-daoud/ai-ui/internal/color"
)

var (
	commentRE = regexp.MustCompile(`(?s)/\*.*?\*/`)
	blockRE   = regexp.MustCompile(`(?s)(?::root(?:,\s*\[data-palette="blue"\])?|\[data-palette="([a-z0-9-]+)"\])\s*\{(.*?)\}`)
	declRE    = regexp.MustCompile(`--([a-z0-9-]+)\s*:\s*([^;]+);`)
)

// palettes parses css/tokens.css into palette → token → [light, dark].
// "blue" is :root (and :root, [data-palette="blue"]); every other palette is
// blue overlaid with its [data-palette="…"] block. Blocks such as
// [data-theme="dark"] carry no tokens and are skipped.
func palettes(t *testing.T) map[string]map[string][2]string {
	t.Helper()
	src, err := os.ReadFile("css/tokens.css")
	if err != nil {
		t.Fatal(err)
	}
	clean := commentRE.ReplaceAllString(string(src), "")
	base := map[string][2]string{}
	overlays := map[string]map[string][2]string{}
	for _, m := range blockRE.FindAllStringSubmatch(clean, -1) {
		target := base
		if name := m[1]; name != "" {
			if overlays[name] == nil {
				overlays[name] = map[string][2]string{}
			}
			target = overlays[name]
		}
		for _, d := range declRE.FindAllStringSubmatch(m[2], -1) {
			v, err := lightDark(strings.TrimSpace(d[2]))
			if err != nil {
				t.Fatalf("--%s: %v", d[1], err)
			}
			target[d[1]] = v
		}
	}
	out := map[string]map[string][2]string{"blue": base}
	for name, o := range overlays {
		p := maps.Clone(base)
		maps.Copy(p, o)
		out[name] = p
	}
	return out
}

// lightDark splits "light-dark(A, B)" at its top-level comma. Any other value
// applies to both modes.
func lightDark(v string) ([2]string, error) {
	if !strings.HasPrefix(v, "light-dark(") {
		return [2]string{v, v}, nil
	}
	if !strings.HasSuffix(v, ")") {
		return [2]string{}, fmt.Errorf("unclosed light-dark: %q", v)
	}
	body := v[len("light-dark(") : len(v)-1]
	depth := 0
	for i, r := range body {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				return [2]string{strings.TrimSpace(body[:i]), strings.TrimSpace(body[i+1:])}, nil
			}
		}
	}
	return [2]string{}, fmt.Errorf("light-dark without a top-level comma: %q", v)
}

func TestLightDarkSplit(t *testing.T) {
	got, err := lightDark("light-dark(oklch(21% 0.03 262 / 0.12), oklch(0% 0 0 / 0.5))")
	if err != nil {
		t.Fatal(err)
	}
	want := [2]string{"oklch(21% 0.03 262 / 0.12)", "oklch(0% 0 0 / 0.5)"}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if _, err := lightDark("light-dark(oklch(1% 0 0))"); err == nil {
		t.Error("a single-argument light-dark parsed")
	}
}

var requiredTokens = []string{
	"paper", "band", "sheet", "ink", "ink-2", "ink-3", "rule", "rule-soft", "rule-field",
	"primary", "on-primary", "primary-hover", "link", "focus", "accent",
	"good", "good-bg", "warning", "warning-bg", "bad", "bad-bg", "mark-bg", "shadow", "scrim",
	"font-sans", "font-mono", "font-display",
	"text-xs", "text-sm", "text-md", "text-lg", "text-xl", "text-2xl", "text-3xl",
	"space-1", "space-2", "space-3", "space-4", "space-5", "space-6", "space-7", "space-8",
	"radius-sm", "radius-md", "radius-lg", "dur", "ease",
}

func TestTokensDefined(t *testing.T) {
	p := palettes(t)
	for _, name := range requiredTokens {
		if _, ok := p["blue"][name]; !ok {
			t.Errorf("--%s is missing from :root", name)
		}
	}
	if _, ok := p["mono"]; !ok {
		t.Error(`the mono palette ([data-palette="mono"]) is missing`)
	}
}

type contrastPair struct {
	fg, bg string
	min    float64
}

// contrastPairs is spec §8.3.
var contrastPairs = func() []contrastPair {
	var ps []contrastPair
	for _, fg := range []string{"ink", "ink-2", "ink-3"} {
		for _, bg := range []string{"paper", "band", "sheet"} {
			ps = append(ps, contrastPair{fg, bg, 4.5})
		}
	}
	return append(ps,
		contrastPair{"on-primary", "primary", 4.5},
		contrastPair{"on-primary", "primary-hover", 4.5},
		contrastPair{"on-primary", "bad", 4.5},
		contrastPair{"link", "paper", 4.5},
		contrastPair{"link", "band", 4.5},
		contrastPair{"link", "sheet", 4.5},
		contrastPair{"good", "good-bg", 4.5},
		contrastPair{"warning", "warning-bg", 4.5},
		contrastPair{"bad", "bad-bg", 4.5},
		contrastPair{"rule-field", "sheet", 3},
	)
}()

func TestTokensContrast(t *testing.T) {
	for name, p := range palettes(t) {
		for mode, label := range []string{"light", "dark"} {
			for _, c := range contrastPairs {
				fg, err := color.Parse(p[c.fg][mode])
				if err != nil {
					t.Errorf("%s %s --%s: %v", name, label, c.fg, err)
					continue
				}
				bg, err := color.Parse(p[c.bg][mode])
				if err != nil {
					t.Errorf("%s %s --%s: %v", name, label, c.bg, err)
					continue
				}
				if got := color.Contrast(fg, bg); got < c.min {
					t.Errorf("%s %s: --%s on --%s = %.2f, want ≥ %.1f", name, label, c.fg, c.bg, got, c.min)
				}
			}
		}
	}
}
