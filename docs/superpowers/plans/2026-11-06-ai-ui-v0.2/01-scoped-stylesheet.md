# ai-ui v0.2, round 1: `ui-scoped.css`

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Goal:** Publish `ui-scoped.css`, which styles library markup and leaves everything else on the page untouched.

**Architecture:** `internal/build` renders the same sources twice. For the scoped file it rewrites `css/reset.css` and `css/base.css` so each selector carries a zero-specificity `:where()` limiting it to `ui-*` markup, and it fails the build on anything it can't rewrite. Tokens, components and utilities are copied unchanged. `dist/ui.css` must stay byte-identical in this round.

**Tech Stack:** Go (`embed`, `net/http`, `testing/fstest`), CSS layers, Playwright.

**Spec:** `docs/superpowers/specs/2026-11-06-v0.2-compact-and-scoped-design.md`, §4 and §7. Read it first.

## Global Constraints

- Go module `github.com/fuad-daoud/ai-ui`, `go 1.25`. No new Go or npm dependency.
- Modify **only** the files this plan names. Need another? **STOP and say so.**
- **Every code block in this plan was run on a trial copy and passed.** Apply it exactly. If a step's result differs from what it says, find out why before changing anything, and report it. Don't change a test to make it pass.
- `dist/ui.css` and `dist/ui-scoped.css` are generated. After any change under `css/` or to the generator, run `go run ./cmd/build` and commit `dist/` in the same commit. Never edit `dist/` by hand. `TestDistIsFresh`, in the gate, fails if you forget.
- **Never run `--update-snapshots` and never add, change or delete anything in `e2e/__screenshots__`.** Baselines depend on the machine's fonts, and the planner makes them on contabo. Always pass `--ignore-snapshots`.
- Run `npm ci` first if `node_modules/` is missing. Playwright uses `/usr/bin/chromium` when it exists (`playwright.config.ts`).
- Never run git except `git add <named files>`, `git commit`, `git status --short`, `git diff`, `git log`, `git show`. No amend, `add -A`, `add .`, stash, reset, checkout, branch, push, merge, rebase or tag. Never `git stash` in a relevo worktree: worktrees share one stash stack.
- No scratch files anywhere in the repository; use `/tmp`.
- A rule's selector never targets a bare element in `css/components/`. Tokens are never renamed.

**Gate, on every round:** `test -z "$(gofmt -l .)" && go vet ./... && go test -count=1 ./... && npx tsc --noEmit && npx vitest run && npx playwright test --ignore-snapshots`

**Report:** every file touched; every step's real result, including each failure you were told to expect; and every place the plan's wording did not fit the code and what you did. `git status --short` is empty after the commit.

## Review Focus

1. **A new reset rule that can't be scoped** (an at-rule, or `body a`) must fail the build, not reach an app's page. Pinned by `TestScopedFailsOnWhatItCannotRewrite`.
2. **Plain elements on a page using the scoped file**, such as a bare `<button>`, `<input>`, `<a>`, `<mark>` or `<code>`, must keep the browser's or the document's look. Pinned by `fixture_document.html`, in light and dark.
3. **Markup deep inside a library element**, such as a link or list inside a callout or disclosure, must look as it does under `ui.css`. Pinned by the library comparison, which covers every descendant.
4. **The fixture going blind:** if it stopped detecting leaks, "leaves it untouched" would pass on anything. Pinned by `ui.css restyles it`.
5. **A focused control,** because the focus ring lives in `base.css`, not in the components. Pinned by `a focused control draws the same ring`.

---

### Task 1: The generator

**Files:**
- Create: `internal/build/scoped.go`, `internal/build/scoped_test.go`
- Modify: `internal/build/build.go`

**Interfaces:**
- Produces: `build.ScopedCSS(fs.FS) ([]byte, error)`, beside the existing `build.CSS(fs.FS) ([]byte, error)`.

- [ ] **Step 1: Write the failing test**

`internal/build/scoped_test.go`:

```go
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
```

- [ ] **Step 2: Run it**

Run: `go test ./internal/build`
Expected: a compile failure, `undefined: build.ScopedCSS`.

- [ ] **Step 3: Write the generator**

`internal/build/scoped.go`:

```go
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
```

`internal/build/build.go`:

```diff
--- a/internal/build/build.go
+++ b/internal/build/build.go
@@ -1,4 +1,5 @@
-// Package build concatenates the css/ sources into dist/ui.css, wrapping each
-// group in its cascade layer (spec §4.1). Sources contain no @layer of their own.
+// Package build concatenates the css/ sources into dist/ui.css and
+// dist/ui-scoped.css, wrapping each group in its cascade layer (spec §4.1).
+// Sources contain no @layer of their own.
 package build
 
@@ -16,4 +17,18 @@
 // CSS returns the content of dist/ui.css built from fsys (the repository root).
 func CSS(fsys fs.FS) ([]byte, error) {
+	return render(fsys, nil)
+}
+
+// ScopedCSS returns the content of dist/ui-scoped.css: the same sources, with the
+// reset and base limited to library markup. Tokens, components and utilities are
+// unchanged, so library markup looks the same under either file.
+func ScopedCSS(fsys fs.FS) ([]byte, error) {
+	return render(fsys, map[string]func(string) (string, error){
+		"css/reset.css": scopeRules,
+		"css/base.css":  scopeRules,
+	})
+}
+
+func render(fsys fs.FS, rewrite map[string]func(string) (string, error)) ([]byte, error) {
 	comps, err := fs.Glob(fsys, "css/components/*.css")
 	if err != nil {
@@ -40,4 +55,11 @@
 				return nil, fmt.Errorf("build: %w", err)
 			}
+			if fn := rewrite[f]; fn != nil {
+				out, err := fn(string(src))
+				if err != nil {
+					return nil, fmt.Errorf("build: %s: %w", f, err)
+				}
+				src = []byte(out)
+			}
 			fmt.Fprintf(&b, "/* %s */\n", f)
 			b.Write(bytes.TrimRight(src, "\n"))
```

- [ ] **Step 4: Run it**

Run: `go test ./internal/build`
Expected: PASS.

### Task 2: Publish the file

**Files:**
- Modify: `cmd/build/main.go`, `dist_test.go`, `ui.go`, `ui_test.go`, `package.json`
- Create: `dist/ui-scoped.css` (generated)

**Interfaces:**
- Consumes: `build.ScopedCSS`.
- Produces: `ui.Href("ui-scoped.css")`, `ui.Assets` entry `ui-scoped.css`, npm export `ai-ui/ui-scoped.css`.

- [ ] **Step 1: Write the failing tests**

`ui_test.go`:

```diff
--- a/ui_test.go
+++ b/ui_test.go
@@ -125,2 +125,29 @@
 	}
 }
+
+func TestScopedStylesheetIsPublished(t *testing.T) {
+	rec := get(t, mux(), ui.Href("ui-scoped.css"), nil)
+	if rec.Code != http.StatusOK {
+		t.Fatalf("status %d", rec.Code)
+	}
+	if ct := rec.Header().Get("Content-Type"); ct != "text/css; charset=utf-8" {
+		t.Errorf("Content-Type = %q", ct)
+	}
+	if cc := rec.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
+		t.Errorf("Cache-Control = %q", cc)
+	}
+	want, err := os.ReadFile("dist/ui-scoped.css")
+	if err != nil {
+		t.Fatal(err)
+	}
+	if rec.Body.String() != string(want) {
+		t.Error("body differs from dist/ui-scoped.css")
+	}
+	got, err := fs.ReadFile(ui.Assets, "ui-scoped.css")
+	if err != nil {
+		t.Fatal(err)
+	}
+	if string(got) != string(want) {
+		t.Error("Assets ui-scoped.css differs from dist/ui-scoped.css")
+	}
+}
```

`dist_test.go`:

```diff
--- a/dist_test.go
+++ b/dist_test.go
@@ -3,4 +3,5 @@
 import (
 	"bytes"
+	"io/fs"
 	"os"
 	"testing"
@@ -10,14 +11,19 @@
 
 func TestDistIsFresh(t *testing.T) {
-	want, err := build.CSS(os.DirFS("."))
-	if err != nil {
-		t.Fatal(err)
-	}
-	got, err := os.ReadFile("dist/ui.css")
-	if err != nil {
-		t.Fatalf("%v (run: go run ./cmd/build)", err)
-	}
-	if !bytes.Equal(got, want) {
-		t.Fatal("dist/ui.css is stale: run `go run ./cmd/build` and commit it")
+	for file, build := range map[string]func(fs.FS) ([]byte, error){
+		"dist/ui.css":        build.CSS,
+		"dist/ui-scoped.css": build.ScopedCSS,
+	} {
+		want, err := build(os.DirFS("."))
+		if err != nil {
+			t.Fatal(err)
+		}
+		got, err := os.ReadFile(file)
+		if err != nil {
+			t.Fatalf("%v (run: go run ./cmd/build)", err)
+		}
+		if !bytes.Equal(got, want) {
+			t.Errorf("%s is stale: run `go run ./cmd/build` and commit it", file)
+		}
 	}
 }
```

- [ ] **Step 2: Run them**

Run: `go test -run 'TestScopedStylesheetIsPublished|TestDistIsFresh' .`
Expected: FAIL. `ui.Href` panics with `unknown asset "ui-scoped.css"`, or `dist/ui-scoped.css` is missing.

- [ ] **Step 3: Implement**

`cmd/build/main.go`:

```diff
--- a/cmd/build/main.go
+++ b/cmd/build/main.go
@@ -1,3 +1,3 @@
-// Command build writes dist/ui.css from the sources under css/.
+// Command build writes dist/ui.css and dist/ui-scoped.css from the sources under css/.
 // Run it from the repository root: go run ./cmd/build
 package main
@@ -5,4 +5,5 @@
 import (
 	"fmt"
+	"io/fs"
 	"os"
 
@@ -14,15 +15,23 @@
 		fail("run from the ai-ui repository root (no go.mod here)")
 	}
-	out, err := build.CSS(os.DirFS("."))
-	if err != nil {
-		fail(err.Error())
-	}
 	if err := os.MkdirAll("dist", 0o755); err != nil {
 		fail(err.Error())
 	}
-	if err := os.WriteFile("dist/ui.css", out, 0o644); err != nil {
-		fail(err.Error())
+	for _, t := range []struct {
+		file  string
+		build func(fs.FS) ([]byte, error)
+	}{
+		{"dist/ui.css", build.CSS},
+		{"dist/ui-scoped.css", build.ScopedCSS},
+	} {
+		out, err := t.build(os.DirFS("."))
+		if err != nil {
+			fail(err.Error())
+		}
+		if err := os.WriteFile(t.file, out, 0o644); err != nil {
+			fail(err.Error())
+		}
+		fmt.Printf("wrote %s (%d bytes)\n", t.file, len(out))
 	}
-	fmt.Printf("wrote dist/ui.css (%d bytes)\n", len(out))
 }
 
```

`ui.go`:

```diff
--- a/ui.go
+++ b/ui.go
@@ -14,5 +14,5 @@
 )
 
-//go:embed dist/ui.css js/theme.js
+//go:embed dist/ui.css dist/ui-scoped.css js/theme.js
 var raw embed.FS
 
@@ -24,6 +24,7 @@
 // published maps each public asset name to its file in the repository.
 var published = map[string]string{
-	"ui.css":   "dist/ui.css",
-	"theme.js": "js/theme.js",
+	"ui.css":        "dist/ui.css",
+	"ui-scoped.css": "dist/ui-scoped.css",
+	"theme.js":      "js/theme.js",
 }
 
@@ -55,5 +56,5 @@
 
 // Assets holds the published files, opened by their public names
-// ("ui.css", "theme.js").
+// ("ui.css", "ui-scoped.css", "theme.js").
 var Assets fs.FS = publicFS{}
 
@@ -74,5 +75,5 @@
 	a, ok := assets[name]
 	if !ok {
-		panic(fmt.Sprintf("ui: unknown asset %q (published: ui.css, theme.js)", name))
+		panic(fmt.Sprintf("ui: unknown asset %q (published: ui.css, ui-scoped.css, theme.js)", name))
 	}
 	return Prefix + name + "?v=" + a.hash
```

`package.json`:

```diff
--- a/package.json
+++ b/package.json
@@ -8,4 +8,5 @@
   "exports": {
     "./ui.css": "./dist/ui.css",
+    "./ui-scoped.css": "./dist/ui-scoped.css",
     "./layers.css": "./react/layers.css",
     "./theme.css": "./react/theme.css",
```

Then run `go run ./cmd/build`.

- [ ] **Step 4: Check**

Run: `go test -count=1 ./... && git diff --stat dist/ui.css`
Expected: PASS, and **no change to `dist/ui.css`**. `dist/ui-scoped.css` is 19,889 bytes.

### Task 3: Prove it in a browser

**Files:**
- Create: `gallery/fixture_library.html`, `gallery/fixture_document.html`, `e2e/support.ts`, `e2e/scoped.spec.ts`
- Modify: `gallery/main.go`, `gallery/main_test.go`

**Interfaces:**
- Produces: routes `/fixture/library/{sheet}` (with `?density=root|container`) and `/fixture/document/{sheet}`, where sheet is `none`, `ui.css` or `ui-scoped.css`. Helpers `snapshot(page, root)` and `differences(a, b, ignore)` in `e2e/support.ts`. Round 2 uses all of these.

- [ ] **Step 1: Add the fixtures, routes and their Go test**

`gallery/fixture_library.html`:

```html
<!doctype html>
<html lang="en"{{if eq .Density "root"}} data-density="compact"{{end}}>
<head>
<meta charset="utf-8">
<title>ai-ui fixture: library markup</title>
{{if ne .Sheet "none"}}<link rel="stylesheet" href="{{href .Sheet}}">{{end}}
<!-- What ui.css's body rule does, restated so the two sheets inherit alike. -->
<style>body { margin: 0; min-height: 100vh; font-family: var(--font-sans); font-size: var(--text-md); line-height: 1.6; color: var(--ink); }</style>
</head>
<body>
<main id="library" class="ui-stack"{{if eq .Density "container"}} data-density="compact"{{end}}>
  <div class="ui-cluster">
    <button class="ui-btn" id="btn">Primary</button>
    <button class="ui-btn" id="sec" data-variant="secondary">Secondary</button>
    <button class="ui-btn" data-variant="ghost">Ghost</button>
    <button class="ui-btn" data-variant="danger">Danger</button>
    <button class="ui-btn" id="btn-sm" data-size="sm">Small</button>
    <button class="ui-btn" id="btn-lg" data-size="lg">Large</button>
    <a class="ui-btn" href="#library">Link</a>
    <button class="ui-btn" disabled>Off</button>
  </div>
  <div class="ui-field" id="field">
    <label class="ui-field-label" id="label" for="input">Label</label>
    <input class="ui-input" id="input" value="text">
    <p class="ui-field-hint" id="hint">Hint</p>
    <p class="ui-field-error" id="error">Error</p>
  </div>
  <div class="ui-field"><label class="ui-field-label" for="select">Select</label><select class="ui-select" id="select"><option>One</option></select></div>
  <div class="ui-field"><label class="ui-field-label" for="multi">Several</label><select class="ui-select" id="multi" multiple size="2"><option>One</option><option>Two</option></select></div>
  <div class="ui-field"><label class="ui-field-label" for="area">Area</label><textarea class="ui-textarea" id="area">x</textarea></div>
  <label class="ui-check"><input class="ui-check-input" type="checkbox" checked> Check</label>
  <div class="ui-cluster"><span class="ui-pill" id="pill">Neutral</span><span class="ui-pill" data-variant="good">Good</span><span class="ui-pill" data-variant="bad">Bad</span></div>
  <div class="ui-callout" id="callout" data-variant="bad"><p class="ui-callout-title">Title</p><p>Body with <a id="callout-link" href="#library">a link</a> and <strong>bold</strong>.</p></div>
  <div class="ui-table-wrap" id="wrap"><table class="ui-table"><caption>Caption</caption><thead><tr><th scope="col">A</th><th scope="col" class="ui-num">B</th></tr></thead><tbody><tr><td id="td">x</td><td class="ui-num">1</td></tr></tbody></table></div>
  <details class="ui-disclosure" open><summary class="ui-disclosure-summary">More</summary><div class="ui-disclosure-body"><p>Body</p><ul><li>Item</li></ul></div></details>
  <section class="ui-card" id="card"><h2>Card head</h2><p>Text</p></section>
  <dialog class="ui-dialog" id="dialog" aria-labelledby="dialog-title">
    <div class="ui-dialog-head"><h2 class="ui-dialog-title" id="dialog-title">Title</h2><button class="ui-dialog-x" aria-label="Close">x</button></div>
    <div class="ui-dialog-body"><p>Body</p></div>
    <div class="ui-dialog-actions"><button class="ui-btn">OK</button></div>
  </dialog>
</main>
</body>
</html>
```

`gallery/fixture_document.html`:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>ai-ui fixture: a document with bare CSS</title>
{{if ne .Sheet "none"}}<link rel="stylesheet" href="{{href .Sheet}}">{{end}}
<!-- Styled the way a résumé is: bare element rules that state only some properties. -->
<style>
  :root { color-scheme: light; }
  body { font-family: Georgia, serif; max-width: 40rem; }
  h1 { font-size: 22px; }
  h2 { font-size: 14px; border-bottom: 1px solid #333; }
  p, li { font-size: 12px; }
  table { width: 100%; }
  a { color: #123; }
</style>
</head>
<body>
<main id="document">
  <h1>Sam Invented</h1>
  <p>sam@example.org · <a href="#document">example.org</a></p>
  <h2>Experience</h2>
  <h3>Example Corp</h3>
  <p>Built things that <mark>mattered</mark>, in <code>Go</code>.</p>
  <ul><li>One thing</li><li>Another thing</li></ul>
  <table><tr><td>A</td><td>B</td></tr></table>
  <button>A plain button</button> <input value="a plain field">
  <img alt="" width="20" height="20" src="data:image/gif;base64,R0lGODlhAQABAAAAACw=">
</main>
</body>
</html>
```

`gallery/main.go`:

```diff
--- a/gallery/main.go
+++ b/gallery/main.go
@@ -20,4 +20,10 @@
 var pageSrc string
 
+//go:embed fixture_library.html
+var fixtureLibrarySrc string
+
+//go:embed fixture_document.html
+var fixtureDocumentSrc string
+
 var (
 	palettes   = []string{"blue", "mono"}
@@ -30,4 +36,10 @@
 	Current     bool
 }
+
+// fixtureSheets are the sheets a fixture page can load. "none" loads no library
+// CSS, for comparing a document with and without it.
+var fixtureSheets = []string{"none", "ui.css", "ui-scoped.css"}
+
+type fixture struct{ Sheet, Density string }
 
 type page struct {
@@ -54,6 +66,12 @@
 		return nil, err
 	}
-	if _, err := set.New("page").Parse(pageSrc); err != nil {
-		return nil, err
+	for name, src := range map[string]string{
+		"page":             pageSrc,
+		"fixture-library":  fixtureLibrarySrc,
+		"fixture-document": fixtureDocumentSrc,
+	} {
+		if _, err := set.New(name).Parse(src); err != nil {
+			return nil, err
+		}
 	}
 	mux := http.NewServeMux()
@@ -81,4 +99,23 @@
 	mux.HandleFunc("GET /p/{palette}/{mode}", serve(false))
 	mux.HandleFunc("GET /p/{palette}/{mode}/dialog", serve(true))
+	serveFixture := func(name string) http.HandlerFunc {
+		return func(w http.ResponseWriter, r *http.Request) {
+			sheet := r.PathValue("sheet")
+			if !slices.Contains(fixtureSheets, sheet) {
+				http.NotFound(w, r)
+				return
+			}
+			var b bytes.Buffer
+			data := fixture{Sheet: sheet, Density: r.URL.Query().Get("density")}
+			if err := set.ExecuteTemplate(&b, name, data); err != nil {
+				http.Error(w, err.Error(), http.StatusInternalServerError)
+				return
+			}
+			w.Header().Set("Content-Type", "text/html; charset=utf-8")
+			_, _ = b.WriteTo(w)
+		}
+	}
+	mux.HandleFunc("GET /fixture/library/{sheet}", serveFixture("fixture-library"))
+	mux.HandleFunc("GET /fixture/document/{sheet}", serveFixture("fixture-document"))
 	return mux, nil
 }
```

`gallery/main_test.go`:

```diff
--- a/gallery/main_test.go
+++ b/gallery/main_test.go
@@ -56,5 +56,5 @@
 		t.Fatal(err)
 	}
-	for _, path := range []string{"/p/rose/light", "/p/blue/sepia", "/p/blue"} {
+	for _, path := range []string{"/p/rose/light", "/p/blue/sepia", "/p/blue", "/fixture/library/other.css", "/fixture/document/theme.js"} {
 		if rec := fetch(t, h, path); rec.Code != http.StatusNotFound {
 			t.Errorf("%s: status %d, want 404", path, rec.Code)
@@ -95,2 +95,27 @@
 	}
 }
+
+func TestFixturePages(t *testing.T) {
+	h, err := newHandler()
+	if err != nil {
+		t.Fatal(err)
+	}
+	for path, want := range map[string]string{
+		"/fixture/library/ui.css":                   `/_ui/ui.css?v=`,
+		"/fixture/library/ui-scoped.css":            `/_ui/ui-scoped.css?v=`,
+		"/fixture/library/ui.css?density=root":      `<html lang="en" data-density="compact">`,
+		"/fixture/library/ui.css?density=container": `class="ui-stack" data-density="compact"`,
+		"/fixture/document/ui-scoped.css":           `/_ui/ui-scoped.css?v=`,
+	} {
+		rec := fetch(t, h, path)
+		if rec.Code != http.StatusOK {
+			t.Fatalf("%s: status %d: %s", path, rec.Code, rec.Body)
+		}
+		if !strings.Contains(rec.Body.String(), want) {
+			t.Errorf("%s: missing %s", path, want)
+		}
+	}
+	if body := fetch(t, h, "/fixture/document/none").Body.String(); strings.Contains(body, "/_ui/") {
+		t.Error("/fixture/document/none links a library sheet")
+	}
+}
```

- [ ] **Step 2: Add the browser tests**

`e2e/support.ts`:

```ts
import type { Page } from "@playwright/test";

export type Snapshot = { key: string; style: Record<string, string> }[];

// Every computed property of every element under root, and of its ::before and ::after.
export function snapshot(page: Page, root: string): Promise<Snapshot> {
  return page.evaluate((root) => {
    const out: { key: string; style: Record<string, string> }[] = [];
    const els = document.querySelectorAll(`${root}, ${root} *`);
    els.forEach((el, i) => {
      for (const pseudo of [null, "::before", "::after"]) {
        const c = getComputedStyle(el, pseudo);
        const style: Record<string, string> = {};
        for (let n = 0; n < c.length; n++) {
          const prop = c.item(n);
          if (!prop.startsWith("--")) style[prop] = c.getPropertyValue(prop);
        }
        out.push({ key: `${i}:${el.tagName.toLowerCase()}#${el.id}.${el.getAttribute("class") ?? ""}${pseudo ?? ""}`, style });
      }
    });
    return out;
  }, root);
}

export function differences(a: Snapshot, b: Snapshot, ignore: string[] = []): string[] {
  if (a.length !== b.length) return [`different element counts: ${a.length} and ${b.length}`];
  const out: string[] = [];
  a.forEach((row, i) => {
    for (const prop of Object.keys(row.style)) {
      if (ignore.includes(prop)) continue;
      if (row.style[prop] !== b[i].style[prop]) out.push(`${row.key} ${prop}: ${row.style[prop]} -> ${b[i].style[prop]}`);
    }
  });
  return out;
}
```

`e2e/scoped.spec.ts`:

```ts
import { expect, test, type Page } from "@playwright/test";
import { differences, snapshot, type Snapshot } from "./support";

// These two are set on html and body only, so they can't match on a page ui-scoped.css leaves alone.
const PAGE_LEVEL = ["text-size-adjust", "-webkit-font-smoothing"];

async function library(page: Page, sheet: string, query = ""): Promise<Snapshot> {
  await page.goto(`/fixture/library/${sheet}${query}`);
  await page.locator("#dialog").evaluate((d: HTMLDialogElement) => d.showModal());
  return snapshot(page, "#library");
}

for (const query of ["", "?density=container"]) {
  test(`library markup matches under both sheets (${query || "default"})`, async ({ page }) => {
    const full = await library(page, "ui.css", query);
    const scoped = await library(page, "ui-scoped.css", query);
    expect(full.length).toBeGreaterThan(100);
    expect(differences(full, scoped, PAGE_LEVEL)).toEqual([]);
  });
}

test("a focused control draws the same ring under both sheets", async ({ page }) => {
  const ring = async (sheet: string) => {
    await page.goto(`/fixture/library/${sheet}`);
    await page.keyboard.press("Tab");
    await page.keyboard.press("Tab");
    return page.evaluate(() => {
      const el = document.activeElement!;
      const c = getComputedStyle(el);
      return { id: el.id, visible: el.matches(":focus-visible"), style: c.outlineStyle, width: c.outlineWidth, offset: c.outlineOffset, color: c.outlineColor };
    });
  };
  const full = await ring("ui.css");
  expect(full).toMatchObject({ id: "sec", visible: true, style: "solid", width: "2px" });
  expect(await ring("ui-scoped.css")).toEqual(full);
});

for (const colorScheme of ["light", "dark"] as const) {
  test.describe(`a document beside the library, ${colorScheme}`, () => {
    test.use({ colorScheme });

    const load = async (page: Page, sheet: string) => {
      await page.goto(`/fixture/document/${sheet}`);
      return { style: await snapshot(page, "#document"), shot: await page.screenshot({ fullPage: true }) };
    };

    test("ui-scoped.css leaves it untouched", async ({ page }) => {
      const none = await load(page, "none");
      const scoped = await load(page, "ui-scoped.css");
      expect(differences(none.style, scoped.style)).toEqual([]);
      expect(Buffer.compare(none.shot, scoped.shot)).toBe(0);
    });

    // Proves the fixture can see a leak: if ui.css stopped changing it, the test above would prove nothing.
    test("ui.css restyles it", async ({ page }) => {
      const none = await load(page, "none");
      const full = await load(page, "ui.css");
      expect(differences(none.style, full.style).length).toBeGreaterThan(50);
      expect(Buffer.compare(none.shot, full.shot)).not.toBe(0);
    });
  });
}
```

- [ ] **Step 3: Run them**

Run: `go test -count=1 ./gallery && npx playwright test e2e/scoped.spec.ts --ignore-snapshots`
Expected: PASS, 7 browser tests.

- [ ] **Step 4: Prove the document test can fail**

In `internal/build/scoped.go`, temporarily change `libraryScope` to `` `:where(*)` ``, then run `go run ./cmd/build && npx playwright test e2e/scoped.spec.ts --ignore-snapshots`.
Expected: exactly `ui-scoped.css leaves it untouched` fails, in light and dark. Paste that in your report.
Restore `libraryScope` and run `go run ./cmd/build`.

- [ ] **Step 5: Run the gate, then commit**

```bash
git add internal/build/scoped.go internal/build/scoped_test.go internal/build/build.go cmd/build/main.go dist_test.go ui.go ui_test.go package.json dist/ui-scoped.css gallery/fixture_library.html gallery/fixture_document.html gallery/main.go gallery/main_test.go e2e/support.ts e2e/scoped.spec.ts
git commit -m "feat: ui-scoped.css styles library markup and nothing else (plan 01 of 2026-11-06-ai-ui-v0.2)"
```
