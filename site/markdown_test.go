package main

import (
	"strings"
	"testing"
)

const sample = "# Title\n\n## Has `code` in it\n\nsee [x](CHANGELOG.md), [y](#has-code-in-it) and [z](https://example.com)\n\n<script>alert(1)</script>\n\n### Deeper\n\nkept\n\n## Next\n\nnot kept\n"

func mustSection(t *testing.T, heading string, toAny bool) section {
	t.Helper()
	s, err := findSection([]byte(sample), "x.md", heading, toAny)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSectionFindsHeadingWithInlineCode(t *testing.T) {
	s := mustSection(t, "Has code in it", false)
	if s.ID != "has-code-in-it" || s.Title != "Has code in it" {
		t.Errorf("got ID %q, Title %q", s.ID, s.Title)
	}
	if want := `<h2 id="has-code-in-it">Has <code>code</code> in it</h2>`; !strings.Contains(string(s.Heading), want) {
		t.Errorf("heading is %s, want %s", s.Heading, want)
	}
}

func TestSectionStopsAtSameLevelAndKeepsDeeper(t *testing.T) {
	body := string(mustSection(t, "Has code in it", false).Body)
	if !strings.Contains(body, `<h3 id="deeper">Deeper</h3>`) || !strings.Contains(body, "<p>kept</p>") {
		t.Errorf("deeper heading or its text is missing:\n%s", body)
	}
	if strings.Contains(body, "not kept") || strings.Contains(body, "Next") {
		t.Errorf("ran past the next heading at the same level:\n%s", body)
	}
}

func TestSectionToAnyHeadingStopsAtTheFirst(t *testing.T) {
	body := string(mustSection(t, "Has code in it", true).Body)
	if !strings.Contains(body, "see ") || strings.Contains(body, "Deeper") {
		t.Errorf("want only the text before the first heading:\n%s", body)
	}
}

func TestSectionMovesHeadingsSoItStartsAtH2(t *testing.T) {
	s := mustSection(t, "Title", false)
	if !strings.HasPrefix(string(s.Heading), `<h2 id="title">`) {
		t.Errorf("h1 wasn't moved to h2: %s", s.Heading)
	}
	body := string(s.Body)
	for _, want := range []string{`<h3 id="has-code-in-it">`, `<h4 id="deeper">`, `<h3 id="next">`} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %s:\n%s", want, body)
		}
	}
}

func TestSectionRewritesRelativeLinksOnly(t *testing.T) {
	body := string(mustSection(t, "Has code in it", false).Body)
	for _, want := range []string{
		`href="https://github.com/fuad-daoud/ai-ui/blob/main/CHANGELOG.md"`,
		`href="#has-code-in-it"`,
		`href="https://example.com"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %s:\n%s", want, body)
		}
	}
}

func TestSectionDropsRawHTML(t *testing.T) {
	if body := string(mustSection(t, "Has code in it", false).Body); strings.Contains(body, "alert(") {
		t.Errorf("raw HTML reached the output:\n%s", body)
	}
}

func TestSectionMissingHeadingIsAnError(t *testing.T) {
	_, err := findSection([]byte(sample), "x.md", "Nope", false)
	if err == nil || !strings.Contains(err.Error(), `x.md has no heading "Nope"`) {
		t.Errorf("got %v", err)
	}
}

// DESIGN.md's tables leave the corner cell empty; an empty <th> names no
// column, so it renders as an empty <td>.
func TestEmptyCornerHeaderIsNotAHeader(t *testing.T) {
	s, err := findSection([]byte("## T\n\n| | Go |\n|---|---|\n| Button | x |\n"), "x.md", "T", false)
	if err != nil {
		t.Fatal(err)
	}
	body := string(s.Body)
	if strings.Contains(body, "<th></th>") || !strings.Contains(body, "<td></td>\n<th>Go</th>") {
		t.Errorf("empty header cell is still a <th>:\n%s", body)
	}
}

func TestSectionWrapsTables(t *testing.T) {
	s, err := findSection([]byte("## T\n\n| a | b |\n|---|---|\n| 1 | 2 |\n"), "x.md", "T", false)
	if err != nil {
		t.Fatal(err)
	}
	body := string(s.Body)
	if !strings.Contains(body, `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="Table"><table class="ui-table">`) || !strings.Contains(body, "</table></div>") {
		t.Errorf("table isn't wrapped:\n%s", body)
	}
}
