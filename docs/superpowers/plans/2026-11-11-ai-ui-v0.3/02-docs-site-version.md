# ai-ui v0.3, round 2: the contract, the version and the site

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Precondition (STOP if unmet):** `git log --oneline -5` shows round 1's four commits, the last `test: measure every v0.3 option; show them in the gallery`, and `git status --short` is empty.

**Goal:** `DESIGN.md` documents every v0.3 option, `CHANGELOG.md` has 0.3.0 with the two changed looks first, the version is 0.3.0 everywhere, and the site's demos show the options.

**Architecture:** docs and the site only; no library file changes. The site names each Markdown table's region by its columns, because the Button section now has a second table and axe forbids two regions with one name on a page.

**Tech Stack:** Markdown, Go (`site/`), Playwright with axe.

**Spec:** `docs/superpowers/specs/2026-11-11-v0.3-options-design.md`, §8.

## Global Constraints

- Go module `github.com/fuad-daoud/ai-ui`, `go 1.25`. No new Go or npm dependency.
- Modify **only** the files this plan names. Need another? **STOP and say so.**
- **Every patch and file in this plan was applied to a trial copy and passed the gate.** Apply each one exactly: write the patch to the named file under `/tmp` byte for byte, apply it with `git apply`, then delete it. If a step's result differs from what it says, find out why before changing anything, and report it. Don't change a test to make it pass.
- `dist/ui.css` and `dist/ui-scoped.css` are generated. After any change under `css/`, run `go run ./cmd/build` and commit `dist/` in the same commit. Never edit `dist/` by hand. `TestDistIsFresh`, in the gate, fails if you forget.
- **Never run `--update-snapshots` and never add, change or delete anything in `e2e/__screenshots__`.** Baselines depend on the machine's fonts, and the planner makes them on contabo. Always pass `--ignore-snapshots`.
- Run `npm ci` first if `node_modules/` is missing. Playwright uses `/usr/bin/chromium` when it exists.
- Never run git except `git apply <patch>`, `git add <named files>`, `git commit`, `git status --short`, `git diff`, `git log`, `git show`. No amend, `add -A`, `add .`, stash, reset, checkout, branch, push, merge, rebase or tag. Never `git stash` in a relevo worktree: worktrees share one stash stack.
- No scratch files in the repository; patches go in `/tmp` and are deleted after applying.
- A rule's selector never targets a bare element in `css/components/`. Tokens are never renamed.

**Gate, on every round:** `test -z "$(gofmt -l .)" && go vet ./... && go test -count=1 ./... && (cd site && go vet ./... && go test -count=1 ./...) && npx tsc --noEmit && npx vitest run && npx playwright test --ignore-snapshots && npx playwright test -c playwright.site.config.ts`

**Report:** every file touched; every step's real result, including each failure you were told to expect; and every place the plan's wording did not fit the code and what you did. `git status --short` is empty after the last commit.

## Review Focus

1. **Every option in the spec is in `DESIGN.md`,** under its component, with Go and React names.
2. **`TestAgentsNamesCurrentVersion`:** the install lines name `v0.3.0`, matching `package.json`.
3. **axe on the Components page:** each table region has its own name.

---

### Task 1: Contract, changelog and version

**Files:** `DESIGN.md`, `CHANGELOG.md`, `AGENTS.md`, `README.md`, `package.json`, `package-lock.json`

- [ ] **1. Modify `DESIGN.md`.** Write this patch to `/tmp/v03-DESIGN.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-DESIGN.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/DESIGN.md b/DESIGN.md
index c124034..789f093 100644
--- a/DESIGN.md
+++ b/DESIGN.md
@@ -123,15 +123,20 @@ Controls are 40px at 16px by default. Working tools can ask for compact, on the
 | Button `data-size="sm"` | 32px, 14px | 28px, 13px, 10px side padding |
 | Button `data-size="lg"` | 48px, 18px | 44px, 16px, 20px side padding |
 | Input, select | 40px, 16px | 32px, 14px |
+| Input, select `data-size="sm"` | 32px, 14px | 28px, 13px |
+| Input, select `data-size="lg"` | 48px, 16px | 44px, 16px |
 | Textarea | 80px, 16px | 80px, 14px |
 | Field label, hint, error | 14px | 13px |
 | Field gap | 8px | 4px |
 | Controls, callout, table wrap | 6px corners | 4px corners |
-| Card, dialog | 10px corners | 8px corners |
+| Card, dialog | 10px corners, a soft shadow | 8px corners, no shadow |
 | Pill | fully round | 4px corners |
 | Table text | 14px | 14px |
 | Body text (attribute on the page) | 16px | 14px |
 
+On a touch screen (`pointer: coarse`), compact controls grow to finger sizes: buttons, inputs and
+selects are 40px, 34px at `sm` and 44px at `lg`. The default density is already 40px there.
+
 Headings and the spacing tokens don't change. `compact` is the only value, and a default region
 nested inside a compact one isn't supported. Inside the attribute, `--text-md`, `--text-sm`,
 `--radius-md` and `--radius-lg` take compact values, so an app rule that reads them follows along.
@@ -217,8 +222,30 @@ Knobs: `--ui-btn-bg`, `--ui-btn-fg`, `--ui-btn-border`, `--ui-btn-bg-hover`, `--
 
 Variants: `primary` (default), `secondary`, `ghost`, `danger`. Sizes: `sm`, `md` (default), `lg`.
 
-Go args: `Label` (required), `Variant`, `Size`, `Type` (`submit`, `reset`; default `button`), `Name`,
-`Value`, `Command`, `CommandFor`, `Disabled`, and `Href`, which makes it a link. **A button that
+| Option | Markup | Go | React |
+|---|---|---|---|
+| Outlined | `data-tone="outline"` | `Tone "outline"` | `tone="outline"` |
+| Icon-only | `data-shape="square"` | `Shape "square"` | `shape="square"` |
+| Full width | `data-width="full"` | `Width "full"` | `width="full"` |
+| Icon | an `<svg>` before the label | `Icon` | `icon` |
+
+- **Outlined** draws a primary or danger button on the sheet with a `--rule` border, its text in
+  `--primary` or `--bad`. Use it for a destructive action in a row, and keep the filled danger
+  button for a dialog's confirm. Secondary and ghost are unchanged by it.
+- **Icon-only** is as wide as it is tall. It needs both `Icon` and `Label`: the label becomes its
+  `aria-label` and the icon its only content. In Go, a square button without `Icon` is an error; in
+  React, a square button's children must be text.
+- **Full width** fills its container, and a long label wraps.
+- **Icon** is the app's own markup. In Go, pass a `template.HTML` value, from the app's own
+  templates or code, never from user input; a plain string is escaped and shows as text. An `<svg>`
+  in a button is sized to the text.
+
+A disabled button is drawn in `--band` with `--ink-3` text and a `--rule` border, at full opacity,
+so its label stays readable. A disabled ghost button keeps no fill.
+
+Go args: `Label` (required), `Variant`, `Size`, `Tone`, `Shape`, `Width`, `Icon`, `Type` (`submit`,
+`reset`; default `button`), `Name`, `Value`, `Command`, `CommandFor`, `Disabled`, and `Href`, which
+makes it a link. **A button that
 submits its form needs `"Type" "submit"`:** the default is `button`, which submits nothing.
 React: `Button` passes every other `<button>` prop through, and its `type` also defaults to
 `button`. `ButtonLink` passes every `<a>` prop through.
@@ -231,10 +258,11 @@ Labels are text.
 
 ### Field: `.ui-field`
 
-Knobs: `--ui-field-gap`; on controls, `--ui-input-bg`, `--ui-input-fg`, `--ui-input-border`, `--ui-input-radius`, `--ui-input-height`.
+Knobs: `--ui-field-gap`; on controls, `--ui-input-bg`, `--ui-input-fg`, `--ui-input-border`, `--ui-input-radius`, `--ui-input-height`, `--ui-input-width` (default 100%).
 
-Go args: `Label` (required), `ID` or `Name` (at least one), `Control` (`textarea`, `select`; default input), `Type` (`email`, `password`, `number`, `search`, `tel`, `url`, `date`; default `text`), `Value`, `Placeholder`, `Autocomplete`, `Rows`, `Options` (`[]ui.Option`), `Hint`, `Error`, `Required`, `Disabled`.
-React: `<Field label name id control type options defaultValue hint error … />`. Every other prop, including `ref`, goes to the control,
+Go args: `Label` (required), `ID` or `Name` (at least one), `Control` (`textarea`, `select`; default input), `Type` (`email`, `password`, `number`, `search`, `tel`, `url`, `date`; default `text`), `Size` (`sm`, `lg`), `Icon` (input only), `Value`, `Placeholder`, `Autocomplete`, `Rows`, `Options` (`[]ui.Option`), `Hint`, `Error`, `Required`, `Disabled`.
+React: `<Field label name id control type size icon options defaultValue hint error … />`. React's
+`size` is the library's size, `sm` or `lg`, in place of the native numeric `size` attribute. Every other prop, including `ref`, goes to the control,
 so a controlled field takes `value` and `onChange`. `options` is `{ value, label }[]` in React and
 `[]ui.Option{{Value: "m", Label: "Monthly"}}` in Go; the option whose value equals `Value` is selected.
 
@@ -250,7 +278,23 @@ height at least, and `Rows` (`rows`) makes it taller.
 </div>
 ```
 
-Lay fields side by side with `<div class="ui-form-row">`.
+**Sizes.** `data-size="sm"` or `"lg"` on a control (`Size`, `size`) makes it as tall as the
+button of that size, at either density.
+
+**A leading icon.** `Icon` (`icon`) puts an icon inside an input, before its text:
+
+```html
+<div class="ui-input-group">
+  <svg class="ui-input-icon" aria-hidden="true">…</svg>
+  <input class="ui-input" type="search" aria-label="Search merchants">
+</div>
+```
+
+The icon is 1rem in `--ink-3` and takes no clicks. Knob on the group: `--ui-input-icon-inset`.
+
+Lay fields side by side with `<div class="ui-form-row">`. Each child grows from 14rem; set
+`--ui-row-grow` and `--ui-row-basis` on one child to change that, for a narrow field beside a wide
+one: `<div style="--ui-row-grow: 0; --ui-row-basis: 6rem">…port…</div>`.
 
 ### Check: `.ui-check`
 
@@ -270,11 +314,19 @@ Go args: `Label` (required), `Type` (`radio`; default `checkbox`), `ID`, `Name`,
 Knobs: `--ui-pill-bg`, `--ui-pill-fg`, `--ui-pill-border`. Variants: neutral (default), `info`, `good`, `warning`, `bad`.
 Go: `{{template "ui.pill" (dict "Label" "Paid" "Variant" "good")}}`. React: `<Pill variant="good">Paid</Pill>`.
 
+- **Dot:** `data-dot` (Go `Dot true`, React `dot`) draws a small dot in the pill's text colour
+  before the label, for a status.
+- **Outlined:** `data-tone="outline"` (`Tone "outline"`, `tone="outline"`) drops the fill; the
+  border takes the text colour, or `--rule-field` on a neutral pill. `info` is already outlined.
+
+Go args: `Label` (required), `Variant`, `Tone`, `Dot`.
+
 ### Dialog: `.ui-dialog`
 
-Knobs: `--ui-dialog-width`, `--ui-dialog-bg`, `--ui-dialog-radius`, `--ui-dialog-pad`. It needs no JavaScript:
+Knobs: `--ui-dialog-width`, `--ui-dialog-bg`, `--ui-dialog-radius`, `--ui-dialog-pad`, `--ui-dialog-shadow`, `--ui-dialog-border`, `--ui-dialog-eyebrow-fg`. It needs no JavaScript:
 open it with a button carrying `command="show-modal" commandfor="<id>"`. React: `<Dialog id="confirm" title="Delete resume?">…</Dialog>`;
-`closeLabel` (default `Close`) is the × button's accessible name.
+`closeLabel` (default `Close`) is the × button's accessible name, `closeButton={false}` leaves the ×
+out, and `eyebrow` adds a small line above the title.
 In Go, write this markup (it is pinned by `testdata/parity/dialog-basic.html`):
 
 ```html
@@ -290,7 +342,17 @@ In Go, write this markup (it is pinned by `testdata/parity/dialog-basic.html`):
 </dialog>
 ```
 
-`.ui-dialog-actions` is the row of buttons at the end, and it can be the `<form>` itself. A button
+`.ui-dialog-actions` is the row of buttons at the end, and it can be the `<form>` itself.
+
+**An eyebrow** is a `<p class="ui-dialog-eyebrow">` as the dialog's first child, before the head:
+small uppercase text in `--ink-3`, such as "Can't be undone". Set `--ui-dialog-eyebrow-fg: var(--bad)`
+on the dialog for red.
+
+**The × is optional.** Leave its button out of the head. The dialog must still close, so give it a
+Cancel button with `command="close" commandfor="<id>"`.
+
+**Shown without script.** A `<dialog open>` lays out as a dialog in the page's flow, with no
+backdrop, so a confirm can be served open when script is off. Placing it is the app's. A button
 with `command="close" commandfor="<id>"` closes the dialog without script.
 
 **A confirm per row.** In Go, give each row its own dialog with its own id
@@ -331,7 +393,14 @@ confirm button acts on it and closes the dialog with `command="close"`.
 <section class="ui-card">…</section>
 ```
 
-Knobs: `--ui-callout-bg`, `--ui-callout-fg`, `--ui-callout-accent`; `--ui-table-pad-y`, `--ui-table-pad-x`; `--ui-card-pad`, `--ui-card-radius`.
+Knobs: `--ui-callout-bg`, `--ui-callout-fg`, `--ui-callout-accent`; `--ui-table-pad-y`, `--ui-table-pad-x`; `--ui-card-pad`, `--ui-card-radius`, `--ui-card-shadow`, `--ui-card-border`; `--ui-disclosure-rule`.
+
+A card has a soft shadow by default and none in compact density. `--ui-card-shadow: none` makes any
+card flat.
+
+`data-marker="start"` on a disclosure puts its chevron before the summary text, pointing to the
+end when closed and down when open, for a list of folders or groups. `--ui-disclosure-rule:
+transparent` drops the line under it.
 
 A callout that reports what just happened takes `role="status"`, and one that reports an error
 takes `role="alert"`, so a screen reader announces it. An error about one field belongs in that
````

- [ ] **2. Modify `CHANGELOG.md`.** Write this patch to `/tmp/v03-CHANGELOG.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-CHANGELOG.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/CHANGELOG.md b/CHANGELOG.md
index 1b3ec0c..e2cf3c1 100644
--- a/CHANGELOG.md
+++ b/CHANGELOG.md
@@ -4,6 +4,26 @@ Breaking changes come first in each entry. Renaming or removing a token,
 class, data-attribute value, knob, partial name or argument, or React prop
 is breaking.
 
+## 0.3.0 (2026-11-11)
+
+No breaking changes. Two looks change:
+
+- A disabled button is drawn in `--band` with `--ink-3` text and a `--rule` border, at full opacity, at every density. It used to fade to half opacity. A disabled ghost button keeps no fill, and a disabled `<button>` now shows the not-allowed cursor.
+- In compact density, cards and dialogs have no shadow. Their border stays.
+
+New:
+
+- Button: `data-tone="outline"` for an outlined primary or danger button, `data-shape="square"` for an icon-only button, `data-width="full"` for a full-width button whose label wraps, and an icon before the label. Go `Tone`, `Shape`, `Width`, `Icon`; React `tone`, `shape`, `width`, `icon`.
+- Inputs, selects and textareas: `data-size="sm"` and `"lg"`, as tall as the button of that size. Go `Size`, React `size`.
+- A leading icon inside an input: `.ui-input-group` and `.ui-input-icon`. Go `Icon`, React `icon`.
+- Knobs `--ui-input-width`, and `--ui-row-grow` and `--ui-row-basis` on a child of `.ui-form-row`.
+- In compact density, on a touch screen, controls grow to 40px (34px small, 44px large).
+- Card: knobs `--ui-card-shadow` and `--ui-card-border`.
+- Dialog: knobs `--ui-dialog-shadow`, `--ui-dialog-border` and `--ui-dialog-eyebrow-fg`; `.ui-dialog-eyebrow` above the title; the × is optional. React `eyebrow` and `closeButton`.
+- Pill: `data-dot` for a status dot and `data-tone="outline"`. Go `Dot`, `Tone`; React `dot`, `tone`.
+- Disclosure: `data-marker="start"` puts the chevron before the summary, and the knob `--ui-disclosure-rule`.
+- `DESIGN.md`: every option above, and a dialog shown without script.
+
 ## 0.2.1 (2026-11-09)
 
 No breaking changes.
```

- [ ] **3. Modify `AGENTS.md`.** Write this patch to `/tmp/v03-AGENTS.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-AGENTS.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/AGENTS.md b/AGENTS.md
index 5eb04e6..98d8393 100644
--- a/AGENTS.md
+++ b/AGENTS.md
@@ -27,7 +27,7 @@ Check which stack this project uses, then follow one of these.
 ### Go with `html/template`
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.2.1
+go get github.com/fuad-daoud/ai-ui@v0.3.0
 ```
 
 It needs Go 1.25, and `go get` raises the `go` line in an older `go.mod` to match.
@@ -60,13 +60,13 @@ read it too, then install:
 
 ```bash
 echo "allow-git=root" >> .npmrc
-npm install github:fuad-daoud/ai-ui#v0.2.1
+npm install github:fuad-daoud/ai-ui#v0.3.0
 ```
 
 That leaves this in `package.json`:
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.2.1" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.0" }
 ```
 
 ```ts
````

- [ ] **4. Modify `README.md`.** Write this patch to `/tmp/v03-README.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-README.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/README.md b/README.md
index 230f126..67e5332 100644
--- a/README.md
+++ b/README.md
@@ -8,7 +8,7 @@ an agent starts. The site is at https://ui.fuad-daoud.com/.
 ## Go apps
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.2.1
+go get github.com/fuad-daoud/ai-ui@v0.3.0
 ```
 
 It needs Go 1.25.
@@ -41,11 +41,11 @@ allows it, so allow it in a `.npmrc` (which `npm ci` and CI read too), then inst
 
 ```bash
 echo "allow-git=root" >> .npmrc
-npm install github:fuad-daoud/ai-ui#v0.2.1
+npm install github:fuad-daoud/ai-ui#v0.3.0
 ```
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.2.1" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.0" }
 ```
 
 ```ts
````

- [ ] **5. Modify `package.json`.** Write this patch to `/tmp/v03-package.json.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-package.json.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/package.json b/package.json
index be649d6..87649d9 100644
--- a/package.json
+++ b/package.json
@@ -1,6 +1,6 @@
 {
   "name": "ai-ui",
-  "version": "0.2.1",
+  "version": "0.3.0",
   "description": "A small UI library written to be used by coding agents: tokens, layered CSS, Go partials and React components.",
   "license": "MIT",
   "type": "module",
```

- [ ] **6. Modify `package-lock.json`.** Write this patch to `/tmp/v03-package-lock.json.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-package-lock.json.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/package-lock.json b/package-lock.json
index f4ff51e..fb14ee6 100644
--- a/package-lock.json
+++ b/package-lock.json
@@ -1,12 +1,12 @@
 {
   "name": "ai-ui",
-  "version": "0.2.1",
+  "version": "0.3.0",
   "lockfileVersion": 3,
   "requires": true,
   "packages": {
     "": {
       "name": "ai-ui",
-      "version": "0.2.1",
+      "version": "0.3.0",
       "license": "MIT",
       "devDependencies": {
         "@axe-core/playwright": "^4.14.0",
```

- [ ] **7.** Run `go test -count=1 .`. Expect `ok`, `TestAgentsNamesCurrentVersion` among them.

- [ ] **8. Commit:** `git add DESIGN.md CHANGELOG.md AGENTS.md README.md package.json package-lock.json` then `git commit -m "docs: v0.3 options in DESIGN.md; 0.3.0"`.

### Task 2: The site

**Files:** `site/markdown.go`, `site/markdown_test.go`, `site/main.go`, `site/templates/components.html`

- [ ] **9. Modify `site/markdown.go`.** Write this patch to `/tmp/v03-site-markdown.go.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-site-markdown.go.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/site/markdown.go b/site/markdown.go
index 403a392..c1a08f4 100644
--- a/site/markdown.go
+++ b/site/markdown.go
@@ -3,6 +3,7 @@ package main
 import (
 	"bytes"
 	"fmt"
+	"html"
 	"html/template"
 	"net/url"
 	"path"
@@ -22,6 +23,13 @@ const rawURL = "https://raw.githubusercontent.com/fuad-daoud/ai-ui/main/"
 // emptyHeader is a header cell with nothing in it, aligned or not.
 var emptyHeader = regexp.MustCompile(`<th([^>]*)></th>`)
 
+// tableHead, headCell and anyTag find a table's column names, which name its region.
+var (
+	tableHead = regexp.MustCompile(`(?s)<table>\s*<thead>(.*?)</thead>`)
+	headCell  = regexp.MustCompile(`(?s)<th[^>]*>(.*?)</th>`)
+	anyTag    = regexp.MustCompile(`<[^>]+>`)
+)
+
 // markdown renders every Markdown source on the site: tables and heading IDs
 // on, and raw HTML dropped (goldmark's default), so a source file can't put
 // script on the page.
@@ -175,7 +183,8 @@ func plain(n ast.Node, src []byte) string {
 }
 
 // renderMD renders doc and frames each table as a library table, in a named
-// region that takes focus so a table wider than a phone can scroll. An empty
+// region that takes focus so a table wider than a phone can scroll. The region
+// is named by the table's columns, so two tables on one page differ. An empty
 // header cell, the corner of a table like DESIGN.md's, names no column, so it
 // becomes a plain cell. Raw HTML is dropped, so the renderer is the only
 // source of these tags.
@@ -184,7 +193,20 @@ func renderMD(doc ast.Node, src []byte) (template.HTML, error) {
 	if err := markdown.Renderer().Render(&b, src, doc); err != nil {
 		return "", err
 	}
-	out := strings.ReplaceAll(b.String(), "<table>", `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="Table"><table class="ui-table">`)
+	out := tableHead.ReplaceAllStringFunc(b.String(), func(m string) string {
+		var names []string
+		for _, c := range headCell.FindAllStringSubmatch(tableHead.FindStringSubmatch(m)[1], -1) {
+			if n := strings.TrimSpace(html.UnescapeString(anyTag.ReplaceAllString(c[1], ""))); n != "" {
+				names = append(names, n)
+			}
+		}
+		label := "Table"
+		if len(names) > 0 {
+			label += ": " + strings.Join(names, ", ")
+		}
+		return `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="` + html.EscapeString(label) + `"><table class="ui-table">` + strings.TrimPrefix(m, "<table>")
+	})
+	out = strings.ReplaceAll(out, "<table>", `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="Table"><table class="ui-table">`)
 	out = strings.ReplaceAll(out, "</table>", "</table></div>")
 	out = emptyHeader.ReplaceAllString(out, "<td$1></td>")
 	return template.HTML(out), nil
```

- [ ] **10. Modify `site/markdown_test.go`.** Write this patch to `/tmp/v03-site-markdown_test.go.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-site-markdown_test.go.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/site/markdown_test.go b/site/markdown_test.go
index b6214a1..55763ef 100644
--- a/site/markdown_test.go
+++ b/site/markdown_test.go
@@ -172,7 +172,22 @@ func TestSectionWrapsTables(t *testing.T) {
 		t.Fatal(err)
 	}
 	body := string(s.Body)
-	if !strings.Contains(body, `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="Table"><table class="ui-table">`) || !strings.Contains(body, "</table></div>") {
+	if !strings.Contains(body, `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="Table: a, b"><table class="ui-table">`) || !strings.Contains(body, "</table></div>") {
 		t.Errorf("table isn't wrapped:\n%s", body)
 	}
 }
+
+// Two tables on one page are two regions, so each is named by its columns. An
+// empty corner cell names nothing, and markup in a heading is read as text.
+func TestTableRegionsAreNamedByTheirColumns(t *testing.T) {
+	s, err := findSection([]byte("## T\n\n|   | `Go` | **React** |\n|---|---|---|\n| x | 1 | 2 |\n\n| a & b |\n|---|\n| 1 |\n"), "x.md", "T", false)
+	if err != nil {
+		t.Fatal(err)
+	}
+	body := string(s.Body)
+	for _, want := range []string{`aria-label="Table: Go, React"`, `aria-label="Table: a &amp; b"`} {
+		if !strings.Contains(body, want) {
+			t.Errorf("want %s in:\n%s", want, body)
+		}
+	}
+}
```

- [ ] **11. Modify `site/main.go`.** Write this patch to `/tmp/v03-site-main.go.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-site-main.go.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/site/main.go b/site/main.go
index 81faf77..fa70c09 100644
--- a/site/main.go
+++ b/site/main.go
@@ -88,6 +88,7 @@ type page struct {
 	RepoURL            string
 	Changelog          template.HTML
 	Currencies         []ui.Option
+	Icons              demoIcons
 	Guides             []guide       // every guide, for the index and the side list
 	Guide              guide         // the guide this page shows
 	Sections           []section     // that guide's parts, rendered
@@ -95,6 +96,15 @@ type page struct {
 	Components         []demoSection
 }
 
+// demoIcons are the demos' own markup for the partials' Icon argument.
+type demoIcons struct{ Download, Theme, Search template.HTML }
+
+var icons = demoIcons{
+	Download: `<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M8 2v8m-3.5-3.5L8 10l3.5-3.5M3 13h10"/></svg>`,
+	Theme:    `<svg viewBox="0 0 16 16" aria-hidden="true"><circle cx="8" cy="8" r="6" fill="none" stroke="currentColor" stroke-width="1.6"/><path d="M8 2a6 6 0 0 1 0 12z" fill="currentColor"/></svg>`,
+	Search:   `<svg class="ui-input-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" aria-hidden="true"><circle cx="7" cy="7" r="4.5"/><path d="m10.5 10.5 3 3"/></svg>`,
+}
+
 func main() {
 	out := flag.String("out", "out", "directory to write the site to")
 	repo := flag.String("repo", "..", "the repository root")
@@ -141,6 +151,7 @@ func build(out, repo string) error {
 		RepoURL:    repoURL,
 		Guides:     guides,
 		Currencies: []ui.Option{{Value: "usd", Label: "US dollar"}, {Value: "eur", Label: "Euro"}},
+		Icons:      icons,
 	}
 
 	files := map[string][]byte{}
```

- [ ] **12. Modify `site/templates/components.html`.** Write this patch to `/tmp/v03-site-templates-components.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-site-templates-components.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/site/templates/components.html b/site/templates/components.html
index 605a145..cb14327 100644
--- a/site/templates/components.html
+++ b/site/templates/components.html
@@ -43,6 +43,14 @@
     {{template "ui.button" (dict "Label" "Read the guides" "Variant" "secondary" "Href" (print .Root "guides/"))}}
     {{template "ui.button" (dict "Label" "Unavailable link" "Variant" "secondary" "Href" "#" "Disabled" true)}}
   </div>
+  <div class="ui-cluster">
+    {{template "ui.button" (dict "Label" "Remove" "Variant" "danger" "Tone" "outline")}}
+    {{template "ui.button" (dict "Label" "Download" "Variant" "secondary" "Icon" .Icons.Download)}}
+    {{template "ui.button" (dict "Label" "Theme" "Variant" "secondary" "Shape" "square" "Icon" .Icons.Theme)}}
+  </div>
+  <div style="max-inline-size: 16rem">
+    {{template "ui.button" (dict "Label" "Use this category for every payment to Acme" "Variant" "secondary" "Width" "full")}}
+  </div>
 </div>
 {{- end}}
 
@@ -57,6 +65,11 @@
   {{template "ui.field" (dict "Label" "Currency" "ID" "demo-currency" "Control" "select" "Options" .Currencies)}}
   {{template "ui.field" (dict "Label" "Notes" "ID" "demo-notes" "Control" "textarea" "Rows" "3")}}
   {{template "ui.field" (dict "Label" "Account number" "ID" "demo-account" "Value" "DE89 3704 0044" "Disabled" true)}}
+  {{template "ui.field" (dict "Label" "Search invoices" "ID" "demo-search" "Type" "search" "Placeholder" "Customer or number" "Icon" .Icons.Search)}}
+  <div class="ui-form-row">
+    {{template "ui.field" (dict "Label" "Mail server" "ID" "demo-host" "Size" "sm" "Value" "imap.example.com")}}
+    <div style="--ui-row-grow: 0; --ui-row-basis: 6rem">{{template "ui.field" (dict "Label" "Port" "ID" "demo-port" "Type" "number" "Size" "sm" "Value" "993")}}</div>
+  </div>
 </div>
 {{- end}}
 
@@ -79,12 +92,15 @@
   {{template "ui.pill" (dict "Label" "Paid" "Variant" "good")}}
   {{template "ui.pill" (dict "Label" "Due soon" "Variant" "warning")}}
   {{template "ui.pill" (dict "Label" "Overdue" "Variant" "bad")}}
+  {{template "ui.pill" (dict "Label" "2 failing" "Variant" "bad" "Dot" true)}}
+  {{template "ui.pill" (dict "Label" "draft" "Tone" "outline")}}
 </div>
 {{- end}}
 
 {{define "demo Dialog" -}}
 {{template "ui.button" (dict "Label" "Delete invoice" "Variant" "danger" "Command" "show-modal" "CommandFor" "demo-dialog")}}
 <dialog class="ui-dialog" id="demo-dialog" aria-labelledby="demo-dialog-title">
+  <p class="ui-dialog-eyebrow">Can't be undone</p>
   <header class="ui-dialog-head">
     <h3 class="ui-dialog-title" id="demo-dialog-title">Delete invoice?</h3>
     <button class="ui-dialog-x" type="button" command="close" commandfor="demo-dialog" aria-label="Close">×</button>
@@ -125,6 +141,16 @@
       <div class="ui-disclosure-body"><p>Not yet.</p></div>
     </details>
   </div>
+  <div>
+    <details class="ui-disclosure" data-marker="start" open>
+      <summary class="ui-disclosure-summary">2026</summary>
+      <div class="ui-disclosure-body"><p>12 invoices</p></div>
+    </details>
+    <details class="ui-disclosure" data-marker="start">
+      <summary class="ui-disclosure-summary">2025</summary>
+      <div class="ui-disclosure-body"><p>31 invoices</p></div>
+    </details>
+  </div>
   <section class="ui-card" aria-label="Sample card"><p>A card holds a group of related things.</p></section>
 </div>
 {{- end}}
```

- [ ] **13.** Run `cd site && gofmt -l . && go vet ./... && go test -count=1 ./...`. Expect no gofmt output and `ok`, `TestTableRegionsAreNamedByTheirColumns` among them.

- [ ] **14.** Run `npx playwright test -c playwright.site.config.ts`. Expect 57 passed, `axe passes on components/` among them.

- [ ] **15. Commit:** `git add site/markdown.go site/markdown_test.go site/main.go site/templates/components.html` then `git commit -m "feat(site): demos of the v0.3 options; tables named by their columns"`.

### Task 3: The gate

- [ ] **16.** Run the gate above. Expect it all to pass.
