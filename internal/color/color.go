// Package color parses the colour syntaxes css/tokens.css uses and computes
// WCAG 2 contrast. It exists for the contrast test (spec §8.3).
package color

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB is a colour in linear-light sRGB, each channel clamped to [0, 1].
type RGB struct{ R, G, B float64 }

// Parse accepts "#rrggbb" and "oklch(L% C H)" with an optional " / alpha",
// which is ignored. Anything else is an error.
func Parse(s string) (RGB, error) {
	s = strings.TrimSpace(s)
	switch {
	case strings.HasPrefix(s, "#"):
		return parseHex(s)
	case strings.HasPrefix(s, "oklch(") && strings.HasSuffix(s, ")"):
		return parseOKLCH(strings.TrimSuffix(strings.TrimPrefix(s, "oklch("), ")"))
	}
	return RGB{}, fmt.Errorf("color: unsupported syntax %q", s)
}

func parseHex(s string) (RGB, error) {
	h := strings.TrimPrefix(s, "#")
	if len(h) != 6 {
		return RGB{}, fmt.Errorf("color: %q is not #rrggbb", s)
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("color: %q: %w", s, err)
	}
	ch := func(shift uint) float64 { return toLinear(float64((v>>shift)&0xff) / 255) }
	return RGB{ch(16), ch(8), ch(0)}, nil
}

func toLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func parseOKLCH(body string) (RGB, error) {
	if i := strings.Index(body, "/"); i >= 0 {
		body = body[:i]
	}
	f := strings.Fields(body)
	if len(f) != 3 || !strings.HasSuffix(f[0], "%") {
		return RGB{}, fmt.Errorf("color: oklch(%s) wants \"L%% C H\"", body)
	}
	l, errL := strconv.ParseFloat(strings.TrimSuffix(f[0], "%"), 64)
	c, errC := strconv.ParseFloat(f[1], 64)
	h, errH := strconv.ParseFloat(f[2], 64)
	if err := errors.Join(errL, errC, errH); err != nil {
		return RGB{}, fmt.Errorf("color: oklch(%s): %w", body, err)
	}
	return fromOKLCH(l/100, c, h), nil
}

// fromOKLCH converts via Oklab to linear sRGB (Björn Ottosson's matrices),
// clamping out-of-gamut channels the way a browser's gamut mapping roughly does.
func fromOKLCH(L, C, H float64) RGB {
	a := C * math.Cos(H*math.Pi/180)
	b := C * math.Sin(H*math.Pi/180)
	l := cube(L + 0.3963377774*a + 0.2158037573*b)
	m := cube(L - 0.1055613458*a - 0.0638541728*b)
	s := cube(L - 0.0894841775*a - 1.2914855480*b)
	return RGB{
		clamp(4.0767416621*l - 3.3077115913*m + 0.2309699292*s),
		clamp(-1.2684380046*l + 2.6097574011*m - 0.3413193965*s),
		clamp(-0.0041960863*l - 0.7034186147*m + 1.7076147010*s),
	}
}

func cube(x float64) float64  { return x * x * x }
func clamp(x float64) float64 { return math.Min(1, math.Max(0, x)) }

// Luminance is the WCAG relative luminance of c.
func (c RGB) Luminance() float64 { return 0.2126*c.R + 0.7152*c.G + 0.0722*c.B }

// Contrast is the WCAG 2 contrast ratio between a and b, in [1, 21].
func Contrast(a, b RGB) float64 {
	la, lb := a.Luminance(), b.Luminance()
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}
