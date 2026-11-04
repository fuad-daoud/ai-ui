package color_test

import (
	"math"
	"testing"

	"github.com/fuad-daoud/ai-ui/internal/color"
)

func TestContrastKnownPairs(t *testing.T) {
	cases := []struct {
		a, b string
		want float64
	}{
		{"#000000", "#ffffff", 21},
		{"#ffffff", "#ffffff", 1},
		{"#2563eb", "#ffffff", 5.17},
		{"oklch(100% 0 0)", "oklch(0% 0 0)", 21},
		{"oklch(54.6% 0.215 263)", "#ffffff", 5.17},
		{"oklch(49% 0.108 151)", "oklch(94.8% 0.015 158)", 5.14},
	}
	for _, c := range cases {
		a, err := color.Parse(c.a)
		if err != nil {
			t.Fatal(err)
		}
		b, err := color.Parse(c.b)
		if err != nil {
			t.Fatal(err)
		}
		if got := color.Contrast(a, b); math.Abs(got-c.want) > 0.02 {
			t.Errorf("Contrast(%s, %s) = %.3f, want %.2f", c.a, c.b, got, c.want)
		}
		if ab, ba := color.Contrast(a, b), color.Contrast(b, a); ab != ba {
			t.Errorf("Contrast is not symmetric for %s / %s: %v vs %v", c.a, c.b, ab, ba)
		}
	}
}

func TestParseIgnoresAlpha(t *testing.T) {
	a, err := color.Parse("oklch(21% 0.03 262 / 0.12)")
	if err != nil {
		t.Fatal(err)
	}
	b, err := color.Parse("oklch(21% 0.03 262)")
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("alpha changed the colour: %v vs %v", a, b)
	}
}

func TestParseRejects(t *testing.T) {
	for _, s := range []string{"", "rgb(0 0 0)", "#fff", "#gggggg", "oklch(50 0.1 20)", "oklch(50% 0.1)", "light-dark(#000000, #ffffff)"} {
		if _, err := color.Parse(s); err == nil {
			t.Errorf("Parse(%q) succeeded, want an error", s)
		}
	}
}
