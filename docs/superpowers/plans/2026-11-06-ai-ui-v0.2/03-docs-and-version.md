# ai-ui v0.2, round 3: docs and version

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Precondition (STOP if unmet):** `git log --oneline` shows the commit `feat: compact density, and fields match buttons`, and `git status --short` is empty.

**Goal:** `DESIGN.md` documents both features, and the repository says 0.2.0.

**Architecture:** Docs and version numbers only. `dist/` must not change. Tagging is the planner's, after the merge.

**Tech Stack:** Markdown, npm metadata.

**Spec:** `docs/superpowers/specs/2026-11-06-v0.2-compact-and-scoped-design.md`, §3, §4, §7 and §8.

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

1. **The docs table against the code:** a size in `DESIGN.md` that differs from `e2e/density.spec.ts` would mislead every app. Compare each row.
2. **The CHANGELOG's order:** breaking changes come first by its own rule, and this entry says there are none and puts the one visible change first.

---

### Task 1: Docs and version

**Files:**
- Modify: `DESIGN.md`, `CHANGELOG.md`, `README.md`, `package.json`, `package-lock.json`

- [ ] **Step 1: Apply**

`DESIGN.md`:

````diff
--- a/DESIGN.md
+++ b/DESIGN.md
@@ -88,4 +88,60 @@
 
 Delete the compat sheet once no legacy markup is left.
+
+## Density
+
+Controls are 40px at 16px by default. Working tools can ask for compact, on the page or on one part of it:
+
+```html
+<html data-density="compact">          <!-- the whole page -->
+<div data-density="compact">…</div>    <!-- one part of a page -->
+```
+
+| | Default | Compact |
+|---|---|---|
+| Button | 40px, 16px text | 32px, 14px text, 12px side padding |
+| Button `data-size="sm"` | 32px, 14px | 28px, 13px, 10px side padding |
+| Button `data-size="lg"` | 48px, 18px | 44px, 16px, 20px side padding |
+| Input, select | 40px, 16px | 32px, 14px |
+| Textarea | 80px, 16px | 80px, 14px |
+| Field label, hint, error | 14px | 13px |
+| Field gap | 8px | 4px |
+| Controls, callout, table wrap | 6px corners | 4px corners |
+| Card, dialog | 10px corners | 8px corners |
+| Pill | fully round | 4px corners |
+| Table text | 14px | 14px |
+| Body text (attribute on the page) | 16px | 14px |
+
+Headings and the spacing tokens don't change. `compact` is the only value, and a default region
+nested inside a compact one isn't supported. Inside the attribute, `--text-md`, `--text-sm`,
+`--radius-md` and `--radius-lg` take compact values, so an app rule that reads them follows along.
+An unlayered app rule still beats all of this.
+
+## Choosing a stylesheet
+
+Two files are published, and a page loads one of them:
+
+- **`ui.css`** is the normal choice. Its reset and base style the whole page.
+- **`ui-scoped.css`** styles library markup only, meaning an element with a `ui-*` class and
+  everything inside it. Use it on a page whose other content must stay exactly as it is, such as a
+  document drawn under your nav. Don't give your own classes the `ui-` prefix, because they would
+  count as library markup.
+
+```html
+<link rel="stylesheet" href="{{href "ui-scoped.css"}}">
+```
+
+Library markup looks the same under either file. With `ui-scoped.css`, the page gives up:
+
+- **Page-level styling.** `body` gets no background, font or colour, and links outside library
+  markup get no colour.
+- **An inherited font.** Library markup takes its font from around it, so the surrounding app sets
+  `font-family`.
+- **Page-level compact.** With no `body` rule, `<html data-density="compact">` doesn't set body
+  text, so put the attribute on the element that holds the library markup.
+
+Tokens still apply to `:root`, including `color-scheme: light dark`. A page that sets no
+`color-scheme` of its own follows the system's dark mode for its default colours, so set one if
+the page must not.
 
 ## Tokens
````

`CHANGELOG.md`:

```diff
--- a/CHANGELOG.md
+++ b/CHANGELOG.md
@@ -4,4 +4,15 @@
 class, data-attribute value, knob, partial name or argument, or React prop
 is breaking.
+
+## 0.2.0 (2026-11-06)
+
+No breaking changes. One visible change at the default density:
+
+- Inputs and selects are now 40px, the same as a button. They were about 44px, because their padding and line height added up to more than their minimum height. A multi-select and a textarea keep their padding.
+
+New:
+
+- `data-density="compact"`: 32px controls, 14px text and 4px corners, on the page or on one part of it.
+- `ui-scoped.css`: styles library markup only, for a page whose other content must stay untouched. Generated from the same sources as `ui.css`.
 
 ## 0.1.1 (2026-11-06)
```

`README.md`:

````diff
--- a/README.md
+++ b/README.md
@@ -8,5 +8,5 @@
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.1.1
+go get github.com/fuad-daoud/ai-ui@v0.2.0
 ```
 
@@ -33,5 +33,5 @@
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.1.1" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.2.0" }
 ```
 
````

`package.json`:

```diff
--- a/package.json
+++ b/package.json
@@ -1,5 +1,5 @@
 {
   "name": "ai-ui",
-  "version": "0.1.1",
+  "version": "0.2.0",
   "description": "Shared UI for Fuad's platforms: tokens, layered CSS, Go partials and React components.",
   "license": "MIT",
```

`package-lock.json`:

```diff
--- a/package-lock.json
+++ b/package-lock.json
@@ -1,5 +1,5 @@
 {
   "name": "ai-ui",
-  "version": "0.1.1",
+  "version": "0.2.0",
   "lockfileVersion": 3,
   "requires": true,
@@ -7,5 +7,5 @@
     "": {
       "name": "ai-ui",
-      "version": "0.1.1",
+      "version": "0.2.0",
       "license": "MIT",
       "devDependencies": {
```

- [ ] **Step 2: Compare the table**

Compare every row of the new `DESIGN.md` table with `e2e/density.spec.ts`. Report any difference, and don't change either side.

- [ ] **Step 3: Run the gate, then commit**

Expected: PASS, and `git diff --stat dist/` shows nothing.

```bash
git add DESIGN.md CHANGELOG.md README.md package.json package-lock.json
git commit -m "docs: density and ui-scoped.css in the contract; 0.2.0 (plan 03 of 2026-11-06-ai-ui-v0.2)"
```
