package build

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// libraryScope matches an element with a ui-* class and everything inside one.
// It is a :where(), so scoping a selector adds no specificity.
const libraryScope = `:where([class^="ui-"], [class*=" ui-"], [class^="ui-"] *, [class*=" ui-"] *)`

// pageOnly finds a selector that names the page itself, which scoping can't express.
var pageOnly = regexp.MustCompile(`\b(html|body)\b|:root`)

// scopeRules rewrites flat CSS so every rule applies to library markup only.
// Rules on html and body are dropped, and a bare :root (the tokens' reduced-motion
// rule) is kept. It returns an error for anything else it can't rewrite, so a new
// rule in the reset or base can't reach an app's page unnoticed.
func scopeRules(src string) (string, error) {
	rest, err := stripComments(src)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for {
		rest = strings.TrimSpace(rest)
		if rest == "" {
			return b.String(), nil
		}
		open := strings.IndexByte(rest, '{')
		if open < 0 {
			return "", fmt.Errorf("text outside a rule: %q", rest)
		}
		end, err := closingBrace(rest, open)
		if err != nil {
			return "", err
		}
		prelude, body := strings.TrimSpace(rest[:open]), rest[open+1:end]
		rest = rest[end+1:]
		switch {
		case strings.HasPrefix(prelude, "@media"):
			inner, err := scopeRules(body)
			if err != nil {
				return "", err
			}
			if inner != "" {
				fmt.Fprintf(&b, "%s {\n%s}\n", prelude, inner)
			}
		case strings.HasPrefix(prelude, "@"):
			return "", fmt.Errorf("can't scope %q", prelude)
		default:
			sel, err := scopeSelectors(prelude)
			if err != nil {
				return "", err
			}
			if sel != "" {
				fmt.Fprintf(&b, "%s {%s}\n", sel, body)
			}
		}
	}
}

func scopeSelectors(list string) (string, error) {
	var out []string
	for _, sel := range splitSelectors(list) {
		switch {
		case sel == "html" || sel == "body":
			continue
		case sel == ":root":
			out = append(out, sel)
			continue
		case sel == "" || pageOnly.MatchString(sel):
			return "", fmt.Errorf("can't scope selector %q", sel)
		}
		head, pseudo := sel, ""
		if i := strings.Index(sel, "::"); i >= 0 {
			head, pseudo = sel[:i], sel[i:]
		}
		out = append(out, head+libraryScope+pseudo)
	}
	return strings.Join(out, ", "), nil
}

// splitSelectors splits on commas outside parentheses.
func splitSelectors(list string) []string {
	var out []string
	depth, start := 0, 0
	for i, r := range list {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(list[start:i]))
				start = i + 1
			}
		}
	}
	return append(out, strings.TrimSpace(list[start:]))
}

func closingBrace(s string, open int) (int, error) {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return 0, errors.New("unbalanced braces")
}

func stripComments(s string) (string, error) {
	var b strings.Builder
	for {
		i := strings.Index(s, "/*")
		if i < 0 {
			b.WriteString(s)
			return b.String(), nil
		}
		b.WriteString(s[:i])
		j := strings.Index(s[i+2:], "*/")
		if j < 0 {
			return "", errors.New("unterminated comment")
		}
		s = s[i+2+j+2:]
	}
}
