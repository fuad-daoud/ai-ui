# Site Components and Guides pages: implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: use superpowers:subagent-driven-development (recommended)
> or superpowers:executing-plans to carry out this plan task by task. Steps use checkbox (`- [ ]`) syntax
> for tracking.

**Goal:** Replace the nav's GitHub links with a Components page and Guides pages built from
`DESIGN.md` and `README.md`.

**Architecture:**
- The Go generator in `site/` gains `findSection`, which pulls one heading's section out of a Markdown
  file with goldmark's AST.
- Guide pages are lists of `{file, heading}` parts.
- The Components page is a list of `{name, DESIGN.md heading}`, each paired with a demo template
  rendered with the real partials.
- Everything renders before anything is written, as now.

**Tech stack:** Go 1.25, goldmark v1.8.7 with `extension.Table`, `html/template` with `ui.Parse`,
`golang.org/x/net/html` in tests, Playwright with axe on contabo via `dev run`.

**Spec:** `docs/superpowers/specs/2026-11-07-site-pages-design.md`, which extends
`docs/superpowers/specs/2026-11-07-site-design.md`.

## Global constraints

- Site copy never states how many components there are, because the number will change.
- Every internal link and asset path is relative, built from `.Root`: `./`, `../` or `../../`.
- Raw HTML in Markdown is dropped, never passed through. This is goldmark's default; don't enable
  `html.WithUnsafe`.
- A missing source file or heading fails the build, and nothing is written.
- No library change apart from the README edit in Task 2. No version bump, no tag.
- Node, Playwright and tests that need a browser run on contabo: `dev run "<cmd>"` from the repo root.
  Go runs locally.
- Demo element IDs start with `demo-`.
- Commits end with the session's attribution lines.

## Review focus

1. **Deep links:** a fragment link such as `getting-started/#nextjs-apps`, or a side-list `#button-ui-btn`,
   must land on an element with that `id`. The Task 4 link test checks fragments.
2. **ID collisions:** a guide that mixes README and DESIGN sections must not repeat an `id`, and demo IDs
   must not collide with heading IDs. Task 4 adds a Go duplicate-ID test over every page, and Task 5
   adds a browser one.
3. **Template actions in the contract text:** `{{template "ui.button" …}}` inside `DESIGN.md` must show
   as text, never run. Task 3 tests that the Components page shows it literally.
4. **Phones:** a contract table wider than a phone scrolls inside its own frame, reachable by keyboard.
   It must not push the page sideways. Task 5 runs the no-sideways-scroll test at 320 and 360 px and
   axe on every new page.
5. **Scripts off:** the dialog demo must still open, because `command="show-modal"` needs no JavaScript.
   Task 5 runs the dialog test with scripts on and off.

---

### Task 1: Markdown sections

**Files:**
- Create: `site/markdown.go`
- Create: `site/markdown_test.go`
- Modify: `site/main.go`, where `renderChangelog` switches to the shared renderer

**Interfaces:**
- Produces:
  - `var markdown goldmark.Markdown`
  - `type section struct { ID, Title string; Heading, Body template.HTML }`
  - `func findSection(src []byte, file, heading string, toAnyHeading bool) (section, error)`
  - `func renderMD(doc ast.Node, src []byte) (template.HTML, error)`

- [ ] **Step 1: Write the failing tests**

`site/markdown_test.go`:

```go
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

func TestSectionWrapsTables(t *testing.T) {
	s, err := findSection([]byte("## T\n\n| a | b |\n|---|---|\n| 1 | 2 |\n"), "x.md", "T", false)
	if err != nil {
		t.Fatal(err)
	}
	body := string(s.Body)
	if !strings.Contains(body, `<div class="ui-table-wrap" tabindex="0"><table class="ui-table">`) || !strings.Contains(body, "</table></div>") {
		t.Errorf("table isn't wrapped:\n%s", body)
	}
}
```

- [ ] **Step 2: Run the tests and confirm they fail**

Run: `cd site && go test -run Section ./...`
Expected: a build failure, because `findSection` and `section` are undefined.

- [ ] **Step 3: Implement**

`site/markdown.go`:

```go
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// markdown renders every Markdown source on the site: tables and heading IDs
// on, and raw HTML dropped (goldmark's default), so a source file can't put
// script on the page.
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.Table),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

// section is one heading of a Markdown file and what follows it.
type section struct {
	ID, Title string        // the heading's id and plain text
	Heading   template.HTML // the heading element, as an h2
	Body      template.HTML
}

// findSection returns the section under the top-level heading whose plain
// text is heading. It runs to the next heading at the same level or higher,
// or, with toAnyHeading, to the next heading of any level. Headings move so
// the section's own heading is an h2, and relative links point at the file
// on GitHub.
func findSection(src []byte, file, heading string, toAnyHeading bool) (section, error) {
	doc := markdown.Parser().Parse(text.NewReader(src))
	var h *ast.Heading
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if c, ok := n.(*ast.Heading); ok && plain(c, src) == heading {
			h = c
			break
		}
	}
	if h == nil {
		return section{}, fmt.Errorf("%s has no heading %q", file, heading)
	}
	body := ast.NewDocument()
	for n := h.NextSibling(); n != nil; {
		next := n.NextSibling()
		if c, ok := n.(*ast.Heading); ok && (toAnyHeading || c.Level <= h.Level) {
			break
		}
		body.AppendChild(body, n) // moves n out of doc
		n = next
	}
	shift := 2 - h.Level
	h.Level = 2
	_ = ast.Walk(body, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch c := n.(type) {
		case *ast.Heading:
			c.Level += shift
		case *ast.Link:
			if d := string(c.Destination); !strings.Contains(d, ":") && !strings.HasPrefix(d, "#") {
				c.Destination = []byte(repoURL + "/blob/main/" + d)
			}
		}
		return ast.WalkContinue, nil
	})
	s := section{Title: plain(h, src)}
	if id, ok := h.AttributeString("id"); ok {
		if b, ok := id.([]byte); ok {
			s.ID = string(b)
		}
	}
	head := ast.NewDocument()
	head.AppendChild(head, h)
	var err error
	if s.Heading, err = renderMD(head, src); err != nil {
		return section{}, fmt.Errorf("%s: %w", file, err)
	}
	if s.Body, err = renderMD(body, src); err != nil {
		return section{}, fmt.Errorf("%s: %w", file, err)
	}
	return s, nil
}

// plain is a node's text with the Markdown taken out: "Button: `.ui-btn`"
// reads "Button: .ui-btn".
func plain(n ast.Node, src []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := c.(type) {
			case *ast.Text:
				b.Write(t.Segment.Value(src))
			case *ast.String:
				b.Write(t.Value)
			}
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

// renderMD renders doc and frames each table as a library table. Raw HTML is
// dropped, so the renderer is the only source of a <table> tag.
func renderMD(doc ast.Node, src []byte) (template.HTML, error) {
	var b bytes.Buffer
	if err := markdown.Renderer().Render(&b, src, doc); err != nil {
		return "", err
	}
	out := strings.ReplaceAll(b.String(), "<table>", `<div class="ui-table-wrap" tabindex="0"><table class="ui-table">`)
	out = strings.ReplaceAll(out, "</table>", "</table></div>")
	return template.HTML(out), nil
}
```

In `site/main.go`, `renderChangelog` uses the shared renderer. Replace its last five lines:

```go
	var b bytes.Buffer
	if err := markdown.Convert(raw, &b); err != nil {
		return "", fmt.Errorf("CHANGELOG.md: %w", err)
	}
	return template.HTML(b.String()), nil
```

Then remove the now-unused `goldmark` and `parser` imports from `main.go`.

- [ ] **Step 4: Run the tests and confirm they pass**

Run: `cd site && go vet ./... && go test ./...`
Expected: PASS, including the existing tests.

- [ ] **Step 5: Commit**

```bash
git add site/markdown.go site/markdown_test.go site/main.go
git commit -m "feat(site): pull one section out of a Markdown file"
```

---

### Task 2: Guides

**Files:**
- Modify: `README.md`, the Next.js section
- Modify: `site/main.go`, adding guides, sources, `page` fields and `buildGuides`
- Create: `site/templates/guides.html`, which defines `guides.html` and `guide.html`
- Modify: `site/assets/site.css`, adding the docs layout and rendered Markdown
- Modify: `site/main_test.go`, extending `pages` and `fakeRepo` and adding a broken-guide case

**Interfaces:**
- Consumes: `findSection` (Task 1).
- Produces:
  - `type part struct{ File, Heading string }`
  - `type guide struct{ Slug, Title, Blurb string; Parts []part }`
  - `var guides []guide`
  - `func readSources(repo string) (map[string][]byte, error)`
  - `func buildGuides(set *template.Template, base page, src, files map[string][]byte) error`
  - `page` fields `Section string`, `Guides []guide`, `Guide guide`, `Sections []section`

- [ ] **Step 1: README**

In `README.md`, replace the Next.js section's lines from the dependency through the import, which run
from the line after `## Next.js apps` to the `import { Button, Field }` block, with:

````markdown
This needs React 19.

```json
"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.2.0" }
```

```ts
// next.config.ts
export default { transpilePackages: ["ai-ui"] };
```

With Tailwind v4:

```css
/* app/globals.css: this order matters */
@import "ai-ui/layers.css";
@import "tailwindcss";
@import "ai-ui/theme.css";
@import "ai-ui/ui.css";
```

Without Tailwind, import only `ai-ui/ui.css`. `layers.css` and `theme.css` exist only for Tailwind v4.

```tsx
import { Button, Field } from "ai-ui/react";
```
````

Run `go test . -run TestAgentsNamesCurrentVersion` from the repo root. Expected: PASS, because the
version strings are unchanged.

- [ ] **Step 2: Write the failing tests**

In `site/main_test.go`:

```go
var pages = []string{
	"index.html",
	"releases/index.html",
	"guides/index.html",
	"guides/getting-started/index.html",
	"guides/theming/index.html",
	"guides/density/index.html",
	"guides/existing-apps/index.html",
}
```

Replace `fakeRepo` so a fake repository also holds the real `DESIGN.md` and `README.md`:

```go
// fakeRepo holds only what build reads from the repository. DESIGN.md and
// README.md are the real ones unless a test changes them.
func fakeRepo(t *testing.T, pkg, changelog string) string {
	t.Helper()
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if pkg != "" {
		write("package.json", pkg)
	}
	if changelog != "" {
		write("CHANGELOG.md", changelog)
	}
	for _, name := range []string{"DESIGN.md", "README.md"} {
		raw, err := os.ReadFile(filepath.Join("..", name))
		if err != nil {
			t.Fatal(err)
		}
		write(name, string(raw))
	}
	return dir
}
```

Add:

```go
// A guide is built from headings in DESIGN.md and README.md, so renaming one
// must fail the build rather than publish a guide with a hole in it.
func TestBrokenGuideSourcesFailAndWriteNothing(t *testing.T) {
	for name, change := range map[string]func(dir string) error{
		"no DESIGN.md": func(dir string) error { return os.Remove(filepath.Join(dir, "DESIGN.md")) },
		"heading renamed": func(dir string) error {
			p := filepath.Join(dir, "DESIGN.md")
			raw, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			return os.WriteFile(p, []byte(strings.Replace(string(raw), "\n## Density\n", "\n## Sizes\n", 1)), 0o644)
		},
	} {
		t.Run(name, func(t *testing.T) {
			repo := fakeRepo(t, `{"version":"1.0.0"}`, "# Changelog\n\n## 1.0.0\n")
			if err := change(repo); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(t.TempDir(), "out")
			if err := build(out, repo); err == nil {
				t.Fatal("build succeeded")
			}
			if _, err := os.Stat(out); !os.IsNotExist(err) {
				t.Errorf("build left %s behind", out)
			}
		})
	}
}

func TestGuidesShowTheirSections(t *testing.T) {
	out := built(t)
	for _, g := range guides {
		page := parse(t, filepath.Join(out, "guides", g.Slug, "index.html"))
		ids := map[string]bool{}
		walk(page, func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "h2" {
				ids[strings.TrimSpace(text(n))] = true
			}
		})
		for _, p := range g.Parts {
			if !ids[p.Heading] {
				t.Errorf("%s: no h2 %q", g.Slug, p.Heading)
			}
		}
	}
}
```

`text` puts a space after every text node, so `strings.TrimSpace(text(n))` of
`<h2>Density</h2>` is `Density`. None of the guides' headings contain inline code, so no inner
spaces are doubled.

- [ ] **Step 3: Run the tests and confirm they fail**

Run: `cd site && go test ./...`
Expected:
- a build failure, because `guides` is undefined
- after Step 4's types exist, the new pages are missing and the broken-guide cases succeed

- [ ] **Step 4: Implement the generator**

In `site/main.go`:

1. Leave `contractURL`, `readmeURL` and their `page` fields alone. The home page and nav still use
   them until Task 4 removes them.
2. Add the types and the guide list after the constants:

```go
// part is one heading's section of a Markdown file in the repository.
type part struct{ File, Heading string }

// guide is a page made of parts, in order.
type guide struct {
	Slug, Title, Blurb string
	Parts              []part
}

var guides = []guide{
	{"getting-started", "Getting started", "Install it in a Go app or a Next.js app.",
		[]part{{"README.md", "Go apps"}, {"README.md", "Next.js apps"}}},
	{"theming", "Theming", "Tokens, knobs, light and dark, palettes, and why your CSS always wins.",
		[]part{{"DESIGN.md", "Override model"}, {"DESIGN.md", "Tokens"}, {"DESIGN.md", "Modes and palettes"}}},
	{"density", "Density", "A compact size for working tools, for a whole page or one part of it.",
		[]part{{"DESIGN.md", "Density"}}},
	{"existing-apps", "Adding to an existing app", "Bring it into an app that already has its own CSS.",
		[]part{{"DESIGN.md", "Adopting ai-ui in an app that has its own CSS"}, {"DESIGN.md", "Choosing a stylesheet"}}},
}
```

3. Add to `page`:

```go
	Section  string    // the nav link this page sits under: "components", "guides", "releases", or "" for home
	Guides   []guide   // every guide, for the index and the side list
	Guide    guide     // the guide this page shows
	Sections []section // that guide's parts, rendered
```

4. Add `readSources` and `buildGuides`:

```go
// readSources reads the Markdown the guides and the Components page are built from.
func readSources(repo string) (map[string][]byte, error) {
	src := map[string][]byte{}
	for _, name := range []string{"DESIGN.md", "README.md"} {
		raw, err := os.ReadFile(filepath.Join(repo, name))
		if err != nil {
			return nil, err
		}
		src[name] = raw
	}
	return src, nil
}

func buildGuides(set *template.Template, base page, src, files map[string][]byte) error {
	index := base
	index.Title = "Guides: ai-ui"
	index.Description = "Install ai-ui, theme it, size it, and bring it into an app that has its own CSS."
	index.Root = "../"
	index.Section = "guides"
	var err error
	if files["guides/index.html"], err = render(set, "guides.html", index); err != nil {
		return err
	}
	for _, g := range guides {
		p := index
		p.Title = g.Title + ": ai-ui"
		p.Description = g.Blurb
		p.Root = "../../"
		p.Guide = g
		p.Sections = nil
		for _, pt := range g.Parts {
			s, err := findSection(src[pt.File], pt.File, pt.Heading, false)
			if err != nil {
				return err
			}
			p.Sections = append(p.Sections, s)
		}
		if files["guides/"+g.Slug+"/index.html"], err = render(set, "guide.html", p); err != nil {
			return err
		}
	}
	return nil
}
```

`files` is `map[string][]byte`, and so is `src`, so the signature
`buildGuides(set *template.Template, base page, src, files map[string][]byte) error` holds.

5. In `build`:
   - read the sources right after the changelog: `sources, err := readSources(repo)` with the usual
     `if err != nil { return err }`
   - set `Guides: guides` in `base`
   - set `releases.Section = "releases"`
   - call `buildGuides(set, base, sources, files)` after the releases page, returning its error

- [ ] **Step 5: Templates**

`site/templates/guides.html`:

```html
{{define "guides.html" -}}
{{template "top" .}}
<main id="main" class="wrap prose">
<h1>Guides</h1>
<ul class="guide-list">
{{- range .Guides}}
  <li><a href="{{.Slug}}/">{{.Title}}</a><p>{{.Blurb}}</p></li>
{{- end}}
</ul>
</main>
{{template "bottom" .}}
{{- end}}

{{define "guide.html" -}}
{{template "top" .}}
<div class="wrap doc">
  <nav class="side" aria-label="Guides">
    <ul>
    {{- $current := .Guide.Slug}}
    {{- range .Guides}}
      <li><a href="../{{.Slug}}/"{{if eq .Slug $current}} aria-current="page"{{end}}>{{.Title}}</a></li>
    {{- end}}
    </ul>
  </nav>
  <main id="main" class="prose">
    <h1>{{.Guide.Title}}</h1>
    {{- range .Sections}}
    {{.Heading}}
    {{.Body}}
    {{- end}}
  </main>
</div>
{{template "bottom" .}}
{{- end}}
```

- [ ] **Step 6: CSS**

Append to `site/assets/site.css`, before the `@media (prefers-reduced-motion)` block. Also change the
`.prose p, .prose ul` rule to `.prose p, .prose :is(ul, ol)`.

```css
/* Docs: a side list beside the text when there's room, a row above it when there isn't. */
.doc { display: grid; gap: var(--space-5); padding-block: var(--space-5) var(--space-7); }
.doc > main { min-inline-size: 0; }
.doc .prose, .doc.prose { padding-block: 0; }
.side ul { display: flex; flex-wrap: wrap; gap: var(--space-2) var(--space-4); padding: 0; list-style: none; font-size: var(--text-sm); }
.side a { color: var(--ink-2); text-decoration: none; }
.side a:hover, .side a[aria-current="page"] { color: var(--ink); text-decoration: underline; }
@media (min-width: 60rem) {
  .doc { grid-template-columns: 12rem minmax(0, 1fr); gap: var(--space-7); }
  .side { position: sticky; inset-block-start: var(--space-5); align-self: start; }
  .side ul { flex-direction: column; }
}

/* Rendered Markdown. Code wraps rather than scrolls, so it never widens the page. */
.prose :is(ul, ol) { padding-inline-start: 1.5rem; }
.prose h3 { margin-block-start: var(--space-6); }
.prose :is(h2, h3)[id] { scroll-margin-block-start: var(--space-5); }
.prose pre { margin-block-start: var(--space-4); white-space: pre-wrap; overflow-wrap: anywhere; }
.prose pre code { padding: 0; background: none; }
.prose .ui-table-wrap { margin-block-start: var(--space-4); }
.prose td code { overflow-wrap: anywhere; }

.guide-list { padding: 0; list-style: none; }
.guide-list li { margin-block-start: var(--space-5); }
.guide-list a { font-size: var(--text-xl); font-weight: 700; }
.guide-list p { margin-block-start: var(--space-1); }
```

`.prose`'s existing `padding-block` stays for the release notes and the guides index. A guide's
`main` is `.prose` inside `.doc`, so the `.doc .prose` rule zeroes it there.

- [ ] **Step 7: Run the tests and confirm they pass**

Run: `cd site && go vet ./... && go test ./...`
Expected: PASS. The link test covers the new pages. The Guides nav link still goes to GitHub until
Task 4.

- [ ] **Step 8: Commit**

```bash
git add README.md site/main.go site/main_test.go site/templates/guides.html site/assets/site.css
git commit -m "feat(site): guides built from DESIGN.md and README.md"
```

---

### Task 3: Components page

**Files:**
- Modify: `site/main.go`, adding components, `demoSection`, `buildComponents` and `page` fields
- Create: `site/templates/components.html`, the page and one `demo <Name>` template per component
- Modify: `site/assets/site.css`, adding `.component` and `.demo`
- Modify: `site/main_test.go`, adding the coverage tests

**Interfaces:**
- Consumes: `findSection` (Task 1), `readSources` and `page.Section` (Task 2).
- Produces:
  - `type component struct{ Name, Heading string }`
  - `var components []component`
  - `type demoSection struct{ Name string; section; Demo template.HTML }`
  - `func buildComponents(set *template.Template, base page, src, files map[string][]byte) error`
  - `page` fields `Lead template.HTML`, `Components []demoSection`
  - the output file `components/index.html`

- [ ] **Step 1: Write the failing tests**

In `site/main_test.go`, add `"components/index.html"` to `pages`, then factor the class lookup out of
`TestHomeShowsEveryComponent`:

```go
// componentClasses is the first .ui-* class of each file in css/components/, except density.
func componentClasses(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("../css/components/*.css")
	if err != nil || len(files) == 0 {
		t.Fatalf("no component CSS found: %v", err)
	}
	first := regexp.MustCompile(`(?m)^\.(ui-[a-z]+)`)
	var classes []string
	for _, f := range files {
		if filepath.Base(f) == "density.css" {
			continue
		}
		css, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		m := first.FindSubmatch(css)
		if m == nil {
			t.Errorf("%s: no component class found", f)
			continue
		}
		classes = append(classes, string(m[1]))
	}
	return classes
}

// classesIn collects the classes used under every element that has class root.
func classesIn(n *html.Node, root string) map[string]bool {
	found := map[string]bool{}
	walk(n, func(c *html.Node) {
		if c.Type != html.ElementNode || !strings.Contains(" "+attr(c, "class")+" ", " "+root+" ") {
			return
		}
		walk(c, func(d *html.Node) {
			for _, cl := range strings.Fields(attr(d, "class")) {
				found[cl] = true
			}
		})
	})
	return found
}
```

Rewrite `TestHomeShowsEveryComponent` to use them:

```go
// A new component has to be placed in the pile before this passes.
func TestHomeShowsEveryComponent(t *testing.T) {
	found := classesIn(parse(t, filepath.Join(built(t), "index.html")), "pile")
	for _, c := range componentClasses(t) {
		if !found[c] {
			t.Errorf("home page has no .%s", c)
		}
	}
}
```

Then add:

```go
// A new component needs a demo on the Components page before this passes.
func TestComponentsPageDemosEveryComponent(t *testing.T) {
	found := classesIn(parse(t, filepath.Join(built(t), "components/index.html")), "demo")
	for _, c := range componentClasses(t) {
		if !found[c] {
			t.Errorf("Components page has no demo using .%s", c)
		}
	}
}

// Every ### under ## Components in the contract must be on the Components page.
func TestEveryContractComponentIsListed(t *testing.T) {
	raw, err := os.ReadFile("../DESIGN.md")
	if err != nil {
		t.Fatal(err)
	}
	_, rest, ok := strings.Cut(string(raw), "\n## Components\n")
	if !ok {
		t.Fatal(`DESIGN.md has no "## Components"`)
	}
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	listed := map[string]bool{}
	for _, c := range components {
		listed[c.Heading] = true
	}
	for _, m := range regexp.MustCompile(`(?m)^### (.+)$`).FindAllStringSubmatch(rest, -1) {
		if h := strings.ReplaceAll(m[1], "`", ""); !listed[h] {
			t.Errorf("DESIGN.md's %q isn't on the Components page", h)
		}
	}
}

// The contract shows Go template calls as code; the site must not run them.
func TestContractCodeShowsAsText(t *testing.T) {
	page := text(parse(t, filepath.Join(built(t), "components/index.html")))
	if !strings.Contains(page, `{{template "ui.button" (dict "Label" "Save"`) {
		t.Error("the Components page doesn't show the Go usage from DESIGN.md as text")
	}
}
```

The pile's container is `<div class="pile run-a" id="pile">`, so `classesIn(…, "pile")` finds it.

- [ ] **Step 2: Run the tests and confirm they fail**

Run: `cd site && go test ./...`
Expected: a build failure, because `components` is undefined. After the types exist, the page is
missing.

- [ ] **Step 3: Implement the generator**

In `site/main.go`:

```go
// component is one ### section under ## Components in DESIGN.md. Its demo is
// the template "demo <Name>" in templates/components.html.
type component struct{ Name, Heading string }

var components = []component{
	{"Button", "Button: .ui-btn"},
	{"Field", "Field: .ui-field"},
	{"Check", "Check: .ui-check"},
	{"Pill", "Pill: .ui-pill"},
	{"Dialog", "Dialog: .ui-dialog"},
	{"Containers", "Containers (plain markup in Go and React)"},
	{"Utilities", "Utilities"},
}

// demoSection is a component's contract section with its live demo.
type demoSection struct {
	Name string
	section
	Demo template.HTML
}
```

Add to `page`:

```go
	Lead       template.HTML // the Components page's opening text
	Components []demoSection
```

```go
func buildComponents(set *template.Template, base page, src, files map[string][]byte) error {
	p := base
	p.Title = "Components: ai-ui"
	p.Description = "Every ai-ui component, working, with its knobs and its Go and React usage."
	p.Root = "../"
	p.Section = "components"
	lead, err := findSection(src["DESIGN.md"], "DESIGN.md", "Components", true)
	if err != nil {
		return err
	}
	p.Lead = lead.Body
	for _, c := range components {
		s, err := findSection(src["DESIGN.md"], "DESIGN.md", c.Heading, false)
		if err != nil {
			return err
		}
		demo, err := render(set, "demo "+c.Name, p)
		if err != nil {
			return err
		}
		p.Components = append(p.Components, demoSection{Name: c.Name, section: s, Demo: template.HTML(demo)})
	}
	files["components/index.html"], err = render(set, "components.html", p)
	return err
}
```

In `build`:
- move `home.Currencies` into `base`, so every page, and every demo, has it
- call `buildComponents(set, base, sources, files)` after `buildGuides`, returning its error

- [ ] **Step 4: Templates**

`site/templates/components.html`:

```html
{{define "components.html" -}}
{{template "top" .}}
<div class="wrap doc">
  <nav class="side" aria-label="Components">
    <ul>
    {{- range .Components}}
      <li><a href="#{{.ID}}">{{.Name}}</a></li>
    {{- end}}
    </ul>
  </nav>
  <main id="main">
    <div class="prose">
      <h1>Components</h1>
      {{.Lead}}
    </div>
    {{- range .Components}}
    <section class="component" aria-labelledby="{{.ID}}">
      <div class="prose">{{.Heading}}</div>
      <div class="demo">{{.Demo}}</div>
      <div class="prose">{{.Body}}</div>
    </section>
    {{- end}}
  </main>
</div>
{{template "bottom" .}}
{{- end}}

{{define "demo Button" -}}
<div class="ui-stack">
  <div class="ui-cluster">
    {{template "ui.button" (dict "Label" "Save")}}
    {{template "ui.button" (dict "Label" "Preview" "Variant" "secondary")}}
    {{template "ui.button" (dict "Label" "Cancel" "Variant" "ghost")}}
    {{template "ui.button" (dict "Label" "Delete" "Variant" "danger")}}
  </div>
  <div class="ui-cluster">
    {{template "ui.button" (dict "Label" "Small" "Size" "sm")}}
    {{template "ui.button" (dict "Label" "Medium")}}
    {{template "ui.button" (dict "Label" "Large" "Size" "lg")}}
  </div>
  <div class="ui-cluster">
    {{template "ui.button" (dict "Label" "Saving…" "Disabled" true)}}
    {{template "ui.button" (dict "Label" "Read the guides" "Variant" "secondary" "Href" (print .Root "guides/"))}}
    {{template "ui.button" (dict "Label" "Unavailable link" "Variant" "secondary" "Href" "#" "Disabled" true)}}
  </div>
</div>
{{- end}}

{{define "demo Field" -}}
<div class="ui-stack">
  {{template "ui.field" (dict "Label" "Email" "ID" "demo-email" "Type" "email" "Hint" "We never share it.")}}
  {{template "ui.field" (dict "Label" "Work email" "ID" "demo-work-email" "Type" "email" "Value" "ada@example" "Error" "Enter a full email address.")}}
  <div class="ui-form-row">
    {{template "ui.field" (dict "Label" "First name" "ID" "demo-first" "Autocomplete" "given-name")}}
    {{template "ui.field" (dict "Label" "Last name" "ID" "demo-last" "Autocomplete" "family-name")}}
  </div>
  {{template "ui.field" (dict "Label" "Currency" "ID" "demo-currency" "Control" "select" "Options" .Currencies)}}
  {{template "ui.field" (dict "Label" "Notes" "ID" "demo-notes" "Control" "textarea" "Rows" "3")}}
  {{template "ui.field" (dict "Label" "Account number" "ID" "demo-account" "Value" "DE89 3704 0044" "Disabled" true)}}
</div>
{{- end}}

{{define "demo Check" -}}
<div class="ui-stack">
  {{template "ui.check" (dict "Label" "Email me receipts" "ID" "demo-receipts")}}
  {{template "ui.check" (dict "Label" "Remember this device" "ID" "demo-remember" "Checked" true)}}
  <div class="ui-cluster" role="radiogroup" aria-label="Billing period">
    {{template "ui.check" (dict "Label" "Monthly" "Type" "radio" "Name" "demo-period" "Value" "m" "Checked" true)}}
    {{template "ui.check" (dict "Label" "Yearly" "Type" "radio" "Name" "demo-period" "Value" "y")}}
  </div>
  {{template "ui.check" (dict "Label" "Two-factor sign-in (set by your admin)" "ID" "demo-2fa" "Checked" true "Disabled" true)}}
</div>
{{- end}}

{{define "demo Pill" -}}
<div class="ui-cluster">
  {{template "ui.pill" (dict "Label" "Draft")}}
  {{template "ui.pill" (dict "Label" "Sent" "Variant" "info")}}
  {{template "ui.pill" (dict "Label" "Paid" "Variant" "good")}}
  {{template "ui.pill" (dict "Label" "Due soon" "Variant" "warning")}}
  {{template "ui.pill" (dict "Label" "Overdue" "Variant" "bad")}}
</div>
{{- end}}

{{define "demo Dialog" -}}
{{template "ui.button" (dict "Label" "Delete invoice" "Variant" "danger" "Command" "show-modal" "CommandFor" "demo-dialog")}}
<dialog class="ui-dialog" id="demo-dialog" aria-labelledby="demo-dialog-title">
  <header class="ui-dialog-head">
    <h3 class="ui-dialog-title" id="demo-dialog-title">Delete invoice?</h3>
    <button class="ui-dialog-x" type="button" command="close" commandfor="demo-dialog" aria-label="Close">×</button>
  </header>
  <div class="ui-dialog-body">
    This cannot be undone.
    <div class="ui-dialog-actions">
      {{template "ui.button" (dict "Label" "Cancel" "Variant" "secondary" "Command" "close" "CommandFor" "demo-dialog")}}
      {{template "ui.button" (dict "Label" "Delete" "Variant" "danger" "Command" "close" "CommandFor" "demo-dialog")}}
    </div>
  </div>
</dialog>
{{- end}}

{{define "demo Containers" -}}
<div class="ui-stack">
  <div class="ui-callout"><p class="ui-callout-title">Export ready</p><p>Your October statement is ready to download.</p></div>
  <div class="ui-callout" data-variant="good"><p class="ui-callout-title">Paid</p><p>Acme paid invoice 1042.</p></div>
  <div class="ui-callout" data-variant="warning"><p class="ui-callout-title">Due soon</p><p>Northwind's invoice is due on Friday.</p></div>
  <div class="ui-callout" data-variant="bad"><p class="ui-callout-title">Build failed</p><p>The gate found unsupported claims.</p></div>
  <div class="ui-table-wrap">
    <table class="ui-table">
      <caption>October spending</caption>
      <thead><tr><th scope="col">Category</th><th scope="col" class="ui-num">Amount</th></tr></thead>
      <tbody>
        <tr><td>Groceries</td><td class="ui-num">212.40</td></tr>
        <tr><td>Transport</td><td class="ui-num">64.00</td></tr>
      </tbody>
    </table>
  </div>
  <div>
    <details class="ui-disclosure" name="demo-faq">
      <summary class="ui-disclosure-summary">Can I export it?</summary>
      <div class="ui-disclosure-body"><p>Yes, as CSV or JSON.</p></div>
    </details>
    <details class="ui-disclosure" name="demo-faq">
      <summary class="ui-disclosure-summary">Is there an API?</summary>
      <div class="ui-disclosure-body"><p>Not yet.</p></div>
    </details>
  </div>
  <section class="ui-card" aria-label="Sample card"><p>A card holds a group of related things.</p></section>
</div>
{{- end}}

{{define "demo Utilities" -}}
<div class="ui-stack">
  <p class="ui-muted">Muted text, for what matters less.</p>
  <p class="ui-mono">go get github.com/fuad-daoud/ai-ui</p>
  <p class="ui-num">1,040.00</p>
  <div class="ui-cluster">
    {{template "ui.pill" (dict "Label" "A cluster")}}
    {{template "ui.pill" (dict "Label" "wraps")}}
    {{template "ui.pill" (dict "Label" "when it runs out of room")}}
  </div>
</div>
{{- end}}
```

- [ ] **Step 5: CSS**

Append to `site/assets/site.css`, next to the docs rules:

```css
.component { padding-block-start: var(--space-3); }
.demo { max-inline-size: 48rem; margin-block-start: var(--space-4); padding: var(--space-5); background: var(--sheet); border: 1px solid var(--rule-soft); border-radius: var(--radius-lg); }
```

- [ ] **Step 6: Run the tests and confirm they pass**

Run: `cd site && go vet ./... && go test ./...`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add site/main.go site/main_test.go site/templates/components.html site/assets/site.css
git commit -m "feat(site): a Components page with a live demo per component"
```

---

### Task 4: Nav and home links

**Files:**
- Modify: `site/templates/layout.html`, the nav
- Modify: `site/templates/home.html`, the browse and install links
- Modify: `site/main.go`, removing `contractURL`, `readmeURL`, `ContractURL` and `ReadmeURL`
- Modify: `site/assets/site.css`, adding the current nav link
- Modify: `site/main_test.go`, adding fragment-aware links, the nav-current test, duplicate IDs and no
  counts on any page, and removing `TestLinkedHeadingsExist`

**Interfaces:**
- Consumes: `page.Section` (Task 2), and the pages from Tasks 2 and 3.

- [ ] **Step 1: Write the failing tests**

In `site/main_test.go`, add:

```go
func hasID(t *testing.T, path, id string) bool {
	t.Helper()
	found := false
	walk(parse(t, path), func(n *html.Node) {
		if n.Type == html.ElementNode && attr(n, "id") == id {
			found = true
		}
	})
	return found
}
```

In `TestLinksResolveAndStayOnGitHub`, replace the `case ref == "", …` line and the `default:` branch:

```go
			case ref == "", strings.HasPrefix(ref, "data:"):
			case strings.HasPrefix(ref, "#"):
				if !hasID(t, filepath.Join(out, p), ref[1:]) {
					t.Errorf("%s: %q has no target on the page", p, ref)
				}
```

```go
			default:
				file, frag, _ := strings.Cut(ref, "#")
				target := filepath.Join(out, filepath.Dir(p), file)
				if info, err := os.Stat(target); err == nil && info.IsDir() {
					target = filepath.Join(target, "index.html")
				}
				if _, err := os.Stat(target); err != nil {
					t.Errorf("%s: %q doesn't resolve", p, ref)
				} else if frag != "" && !hasID(t, target, frag) {
					t.Errorf("%s: %q has no #%s", p, ref, frag)
				}
```

Then add:

```go
// The nav's Components and Guides go to site pages, and each page marks the
// part of the site it belongs to, and nothing else.
func TestNavPointsAtSitePagesAndMarksTheCurrentOne(t *testing.T) {
	out := built(t)
	want := map[string]string{
		"index.html":                     "",
		"releases/index.html":            "Release notes",
		"components/index.html":          "Components",
		"guides/index.html":              "Guides",
		"guides/theming/index.html":      "Guides",
		"guides/existing-apps/index.html": "Guides",
	}
	for p, current := range want {
		var marked []string
		walk(parse(t, filepath.Join(out, p)), func(n *html.Node) {
			if n.Type != html.ElementNode || n.Data != "nav" || attr(n, "class") != "site-nav" {
				return
			}
			walk(n, func(a *html.Node) {
				if a.Type != html.ElementNode || a.Data != "a" {
					return
				}
				name := strings.TrimSpace(text(a))
				if (name == "Components" || name == "Guides") && strings.Contains(attr(a, "href"), "://") {
					t.Errorf("%s: nav %s leaves the site: %s", p, name, attr(a, "href"))
				}
				if attr(a, "aria-current") == "page" {
					marked = append(marked, name)
				}
			})
		})
		if got := strings.Join(marked, ","); got != current {
			t.Errorf("%s: nav marks %q, want %q", p, got, current)
		}
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	out := built(t)
	for _, p := range pages {
		seen := map[string]bool{}
		walk(parse(t, filepath.Join(out, p)), func(n *html.Node) {
			if id := attr(n, "id"); n.Type == html.ElementNode && id != "" {
				if seen[id] {
					t.Errorf("%s: id %q is used twice", p, id)
				}
				seen[id] = true
			}
		})
	}
}
```

Replace `TestHomeDoesNotCountComponents` with:

```go
func TestNoPageCountsComponents(t *testing.T) {
	count := regexp.MustCompile(`(?i)\b(\d+|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve)\s+components\b`)
	for _, p := range pages {
		if m := count.FindString(text(parse(t, filepath.Join(built(t), p)))); m != "" {
			t.Errorf("%s counts components: %q", p, m)
		}
	}
}
```

Delete `TestLinkedHeadingsExist`. The site no longer links into headings on GitHub, and the build now
fails on any heading it uses that's missing.

- [ ] **Step 2: Run the tests and confirm they fail**

Run: `cd site && go test ./...`
Expected: `TestNavPointsAtSitePagesAndMarksTheCurrentOne` fails, because the nav Components and Guides
links leave the site and nothing is marked.

- [ ] **Step 3: Implement**

The nav in `site/templates/layout.html`:

```html
  <nav class="site-nav" aria-label="Main">
    <a class="mark" href="{{.Root}}">ai-ui</a>
    <a href="{{.Root}}components/"{{if eq .Section "components"}} aria-current="page"{{end}}>Components</a>
    <a href="{{.Root}}guides/"{{if eq .Section "guides"}} aria-current="page"{{end}}>Guides</a>
    <a href="{{.Root}}releases/"{{if eq .Section "releases"}} aria-current="page"{{end}}>Release notes</a>
    <a class="end" href="{{.RepoURL}}">GitHub</a>
  </nav>
```

In `site/templates/home.html`:
- `<a href="{{.ContractURL}}#components">Browse the components</a>` becomes
  `<a href="{{.Root}}components/">Browse the components</a>`.
- The install line becomes:

```html
    <p class="ui-muted">Installing it yourself? Follow the steps for <a href="{{.Root}}guides/getting-started/#go-apps">Go</a> or <a href="{{.Root}}guides/getting-started/#nextjs-apps">Next.js</a>.</p>
```

In `site/main.go`, delete the `contractURL` and `readmeURL` constants, the `ContractURL` and `ReadmeURL`
fields and their assignments in `base`.

In `site/assets/site.css`, after `.site-nav a:hover`, add:

```css
.site-nav a[aria-current="page"] { color: var(--ink); text-decoration: underline; text-decoration-thickness: 2px; text-underline-offset: 0.4em; }
```

- [ ] **Step 4: Run the tests and confirm they pass**

Run: `cd site && go vet ./... && go test ./...`
Expected: PASS. If the duplicate-ID test fails on a guide, rename the clashing demo or heading ID. Don't
weaken the test.

- [ ] **Step 5: Commit**

```bash
git add site/
git commit -m "feat(site): the nav goes to the site's own Components and Guides"
```

---

### Task 5: Browser tests

**Files:**
- Modify: `site/e2e/site.spec.ts`

- [ ] **Step 1: Add the tests**

Near the top of `site/e2e/site.spec.ts`:

```ts
const components = "components/";
const guides = ["guides/", "guides/getting-started/", "guides/theming/", "guides/density/", "guides/existing-apps/"];
const docs = [components, ...guides];
```

Change the duplicate-ID loop to `for (const path of [home, releases, ...docs])`. Then add:

```ts
for (const path of docs) {
  test(`axe passes on ${path}`, async ({ page }) => {
    await page.goto(path);
    const { violations } = await new AxeBuilder({ page }).analyze();
    expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target).join(", ")}`)).toEqual([]);
  });

  for (const width of [320, 360]) {
    test(`${path} doesn't scroll sideways at ${width}px`, async ({ page }) => {
      await page.setViewportSize({ width, height: 800 });
      await page.goto(path);
      const { scroll, client } = await page.evaluate(() => ({ scroll: document.documentElement.scrollWidth, client: document.documentElement.clientWidth }));
      expect(scroll).toBeLessThanOrEqual(client);
    });
  }
}

test("the nav reaches Components and Guides on the site", async ({ page }) => {
  await page.goto(home);
  const nav = page.getByRole("navigation", { name: "Main" });
  await nav.getByRole("link", { name: "Components" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/components\/$/);
  await expect(page.locator("main h1")).toHaveText("Components");
  await nav.getByRole("link", { name: "Guides" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/guides\/$/);
  await page.getByRole("main").getByRole("link", { name: "Theming" }).click();
  await expect(page).toHaveURL(/\/ai-ui\/guides\/theming\/$/);
  await expect(page.locator("main h1")).toHaveText("Theming");
  await expect(page.getByRole("navigation", { name: "Guides" }).getByRole("link", { name: "Theming" })).toHaveAttribute("aria-current", "page");
});

for (const scripts of [true, false]) {
  test(`the dialog demo opens and closes${scripts ? "" : " without scripts"}`, async ({ browser, baseURL }) => {
    const context = await browser.newContext({ javaScriptEnabled: scripts, baseURL });
    const page = await context.newPage();
    await page.goto(components);
    const opener = page.getByRole("button", { name: "Delete invoice" });
    await opener.click();
    const dialog = page.getByRole("dialog", { name: "Delete invoice?" });
    await expect(dialog).toBeVisible();
    await dialog.getByRole("button", { name: "Cancel" }).click();
    await expect(dialog).toBeHidden();
    await expect(opener).toBeFocused();
    await context.close();
  });
}

test("Tab reaches the side list and every demo control on the Components page", async ({ page }) => {
  await page.goto(components);
  // A radio group is one tab stop, at its checked radio; disabled controls and a closed dialog's aren't stops.
  const total = await page.evaluate(() => {
    const els = [...document.querySelectorAll<HTMLElement>(".side a, .demo :is(a[href], button, input, select, textarea, summary)")]
      .filter((el) => !(el as HTMLInputElement).disabled && !el.closest("dialog")
        && !(el instanceof HTMLInputElement && el.type === "radio" && !el.checked));
    els.forEach((el, i) => { el.dataset.want = String(i); });
    return els.length;
  });
  const seen = new Set<string>();
  for (let i = 0; i < total * 2 + 20; i++) {
    await page.keyboard.press("Tab");
    const want = await page.evaluate(() => (document.activeElement as HTMLElement | null)?.dataset.want ?? "");
    if (want) seen.add(want);
  }
  expect(seen.size).toBe(total);
});
```

- [ ] **Step 2: Run the browser tests**

Run, from the repo root: `dev run "npm run e2e:site"`
Expected: every test passes, the original ones and the new ones.

If a test fails, fix the page, not the test. The likely failures:
- axe's `scrollable-region-focusable`, or colour contrast in `.side`
- sideways scroll from a long token in `.prose` at 320 px: add `overflow-wrap: anywhere` to the
  element that overflows

Re-run until green.

- [ ] **Step 3: Commit**

```bash
git add site/e2e/site.spec.ts site/
git commit -m "test(site): Components and Guides in a browser"
```

---

### Task 6: Look, review, ship

- [ ] **Step 1: Look at the pages.**
  - Take full-page screenshots of `components/`, `guides/` and `guides/theming/`, at 360 and 1280 px,
    from the preview server on contabo.
  - Read them.
  - Fix anything that reads badly, such as spacing, the side list, or a demo cramped at 360 px.
  - Re-run `go test` and `npm run e2e:site` after any fix.

- [ ] **Step 2: Whole-branch review.**
  - One fresh reviewer, on the most capable model, reads the spec, this plan and
    `git diff origin/main...HEAD`.
  - Fix what it confirms, with a test where one can fail first.
  - Add a "Changes after the whole-branch review" section to the end of this plan.

- [ ] **Step 3: PR.**
  - Push the branch and open a PR with a summary and the test evidence.
  - Wait for CI to report green, with `gh pr checks --watch` after the checks have registered.
  - Ask Fuad before merging.

- [ ] **Step 4: Deploy after the merge.**
  - On an up-to-date `main`, run the two commands in
    `~/projects/servers/contabo/services/ai-ui-site/NOTES.md`.
  - Check that `https://ui.fuad-daoud.com/components/` and `/guides/theming/` answer 200.
  - Check that `srv.fish probe now` is all up.
  - Update the `ai-ui-site` memory.

## Changes after the whole-branch review

A fresh reviewer read the spec, this plan and the whole diff. It found nothing critical and one important
issue, and one minor item was regraded to important by its effect. Both are fixed, each with a test
that failed first:

- **Guide pages named the wrong link as the current page.** On `guides/theming/` the main nav's Guides
  link had `aria-current="page"` while pointing at `guides/`. §4 and §7 of the spec asked for that, and
  it was wrong. A link now gets `"page"` only when it is the page itself, and the Guides link gets
  `"true"` on a guide. Both are underlined.
- **Table frames were unnamed tab stops.** Each frame now has `role="region"` and `aria-label="Table"`.

The minor findings were deferred:
- the link rewrite's handling of `../`, `//` and `/` paths, images and links in headings
- heading levels past h6
- aligned empty corner cells
- a test that a wide table scrolls inside its frame
- cross-file heading-ID collisions failing at build time rather than in tests
- framed tables in the changelog

All of these were fixed in a follow-up PR on 2026-11-07, `fix/site-minors`, each with a test that failed first.
