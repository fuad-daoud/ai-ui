# ai-ui v0.3.2: dialog title and quiet row headers

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Precondition (STOP if unmet):** `git log --oneline -3` shows `docs: plan ai-ui v0.3.2` at the top, and `git status --short` is empty.

**Goal:** a dialog title shrinks and wraps, so a long email address stays inside the dialog; a quiet table's row headers go quiet too. Docs, the version (0.3.2) and tests follow.

**Architecture:** two CSS rules (`dialog.css`, `table.css`), a long-title dialog and a row header in the library fixture, two browser tests. No Go partial or React change.

**Tech Stack:** CSS layers, Playwright, Markdown.

**Spec:** `docs/superpowers/specs/2026-11-12-v0.3.2-fixes-design.md`. Read it first.

## Global Constraints

- Go module `github.com/fuad-daoud/ai-ui`, `go 1.25`. No new Go or npm dependency.
- Modify **only** the files this plan names. Need another? **STOP and say so.**
- **Every patch in this plan was applied to a trial copy and passed the gate.** Apply each one exactly: write the patch to the named file under `/tmp` byte for byte, apply it with `git apply`, then delete it. If a step's result differs from what it says, find out why before changing anything, and report it. Don't change a test to make it pass.
- `dist/ui.css` and `dist/ui-scoped.css` are generated: after the CSS task run `go run ./cmd/build` and commit `dist/` with it. Never edit `dist/` by hand.
- **Never run `--update-snapshots` and never add, change or delete anything in `e2e/__screenshots__`.** Always pass `--ignore-snapshots`.
- Run `npm ci` first if `node_modules/` is missing.
- Never run git except `git apply <patch>`, `git add <named files>`, `git commit`, `git status --short`, `git diff`, `git log`, `git show`. No amend, `add -A`, `add .`, stash, reset, checkout, branch, push, merge, rebase or tag. Never `git stash` in a relevo worktree.
- No scratch files in the repository; patches go in `/tmp` and are deleted after applying.

**Gate:** `test -z "$(gofmt -l .)" && go vet ./... && go test -count=1 ./... && (cd site && go vet ./... && go test -count=1 ./...) && npx tsc --noEmit && npx vitest run && npx playwright test --ignore-snapshots && npx playwright test -c playwright.site.config.ts`

**Report:** every file touched; every step's real result; every place the plan's wording did not fit the code and what you did. `git status --short` is empty after the last commit.

## Review Focus

1. **The title's min-content.** The reset gives headings `overflow-wrap: break-word`, which does not lower a flex item's min-content; `min-inline-size: 0` plus `anywhere` is what lets it shrink. Pinned by `a long title wraps inside its dialog` (without the fix the title starts at x = -197).
2. **Hover still wins over the quiet row header.** `:where(tbody th)` keeps the rule at (0,2,0), under the hover rule's (0,2,2). Pinned by the last line of the row-header test.
3. **Case stays.** The quiet head's `--ui-table-head-case` is not applied to row headers; the fixture's caps table checks it.
4. **Nothing else moves.** Gallery baselines are unchanged; the planner ran them on the trial.

---

### Task 1: CSS

**Files:** `css/components/dialog.css`, `css/components/table.css`

- [ ] **1. Modify `css/components/dialog.css`.** Write this patch to `/tmp/plan-css-components-dialog.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-css-components-dialog.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/dialog.css b/css/components/dialog.css
index a9648ea..bf41d2f 100644
--- a/css/components/dialog.css
+++ b/css/components/dialog.css
@@ -39,7 +39,9 @@
   text-transform: uppercase;
 }
 .ui-dialog-eyebrow + .ui-dialog-head { padding-block-start: var(--space-1); }
-.ui-dialog-title { font-size: var(--text-lg); }
+/* The title shrinks and wraps, so a long word such as an email address can't push the × or the
+   dialog past the viewport. */
+.ui-dialog-title { min-inline-size: 0; font-size: var(--text-lg); overflow-wrap: anywhere; }
 .ui-dialog-x {
   display: grid;
   place-items: center;
```

- [ ] **2. Modify `css/components/table.css`.** Write this patch to `/tmp/plan-css-components-table.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-css-components-table.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/table.css b/css/components/table.css
index 7b93a81..b666e80 100644
--- a/css/components/table.css
+++ b/css/components/table.css
@@ -55,6 +55,13 @@
   --ui-table-head-fg: var(--ink-3);
   --ui-table-head-size: var(--text-xs);
 }
+/* Its row headers go quiet too: no fill, the head's colour and size, but not its case. :where()
+   keeps the row hover tint winning. */
+.ui-table[data-head="quiet"] :where(tbody th) {
+  background: transparent;
+  color: var(--ui-table-head-fg);
+  font-size: var(--ui-table-head-size);
+}
 /* The header row stays at the top of a wrap that scrolls (bound it with --ui-table-max-height). */
 .ui-table[data-sticky="head"] thead th {
   position: sticky;
```

- [ ] **3.** Run `go run ./cmd/build`. Expect `wrote dist/ui.css (27495 bytes)` and `wrote dist/ui-scoped.css (30929 bytes)`.

- [ ] **4.** Run `go test -count=1 .`. Expect `ok`.

- [ ] **5. Commit:** `git add css/components/dialog.css css/components/table.css dist/ui.css dist/ui-scoped.css` then `git commit -m "fix(css): a dialog title shrinks and wraps; a quiet table's row headers go quiet"`.

### Task 2: Tests

**Files:** `gallery/fixture_library.html`, `e2e/options.spec.ts`, `e2e/tables.spec.ts`

- [ ] **6. Modify `gallery/fixture_library.html`.** Write this patch to `/tmp/plan-gallery-fixture_library.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-gallery-fixture_library.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/gallery/fixture_library.html b/gallery/fixture_library.html
index c417757..d2b813e 100644
--- a/gallery/fixture_library.html
+++ b/gallery/fixture_library.html
@@ -56,7 +56,7 @@
   <div class="ui-callout" id="callout" data-variant="bad"><p class="ui-callout-title">Title</p><p>Body with <a id="callout-link" href="#library">a link</a> and <strong>bold</strong>.</p></div>
   <div class="ui-table-wrap" id="wrap"><table class="ui-table"><caption>Caption</caption><thead><tr><th scope="col">A</th><th scope="col" class="ui-num">B</th></tr></thead><tbody><tr><td id="td">x</td><td class="ui-num">1</td></tr></tbody></table></div>
   <div class="ui-table-wrap" id="plain-wrap" data-frame="none"><table class="ui-table" data-head="quiet"><thead><tr><th scope="col" id="plain-th">Name</th><th scope="col" class="ui-num" id="plain-th-num">Count</th></tr></thead><tbody><tr id="plain-row"><td id="plain-first">a</td><td class="ui-num" id="plain-last">1</td></tr></tbody></table></div>
-  <div class="ui-table-wrap"><table class="ui-table" data-head="quiet" style="--ui-table-head-case: uppercase; --ui-table-head-rule: var(--rule)"><thead><tr><th scope="col" id="caps-th">Email</th></tr></thead><tbody><tr><td>x</td></tr></tbody></table></div>
+  <div class="ui-table-wrap"><table class="ui-table" data-head="quiet" style="--ui-table-head-case: uppercase; --ui-table-head-rule: var(--rule)"><thead><tr><th scope="col" id="caps-th">Email</th><th scope="col">Plan</th></tr></thead><tbody><tr><th scope="row" id="caps-rowhead">Total</th><td id="caps-td">x</td></tr></tbody></table></div>
   <div class="ui-table-wrap" id="sticky-wrap" style="--ui-table-max-height: 8rem"><table class="ui-table" data-sticky="head"><thead><tr><th scope="col" id="sticky-th">Day</th></tr></thead><tbody><tr><td>Day 1</td></tr><tr><td>Day 2</td></tr><tr><td>Day 3</td></tr><tr><td>Day 4</td></tr><tr><td>Day 5</td></tr><tr><td>Day 6</td></tr><tr><td>Day 7</td></tr><tr><td>Day 8</td></tr><tr><td>Day 9</td></tr><tr><td>Day 10</td></tr><tr><td>Day 11</td></tr><tr><td>Day 12</td></tr></tbody></table></div>
   <div class="ui-table-wrap" id="stack-wrap" data-stack="phone"><table class="ui-table" id="stack"><thead id="stack-head"><tr><th scope="col">Email</th><th scope="col">Role</th><th scope="col" class="ui-num">Logins</th></tr></thead><tbody><tr id="stack-row"><td id="stack-name">ada@example.com</td><td id="stack-role" data-label="Role">owner</td><td id="stack-num" data-label="Logins" class="ui-num">12</td></tr><tr><th scope="row" id="stack-rowhead">grace@example.com</th><td data-label="Role">member</td><td data-label="Logins" class="ui-num">3</td></tr></tbody></table></div>
   <details class="ui-disclosure" open><summary class="ui-disclosure-summary">More</summary><div class="ui-disclosure-body"><p>Body</p><ul><li>Item</li></ul></div></details>
@@ -73,6 +73,10 @@
     <div class="ui-dialog-head"><h2 class="ui-dialog-title" id="dialog-eyebrow-title">Delete?</h2></div>
     <div class="ui-dialog-body"><p>Body</p></div>
   </dialog>
+  <dialog class="ui-dialog" id="dialog-long" aria-labelledby="dialog-long-title">
+    <div class="ui-dialog-head"><h2 class="ui-dialog-title" id="dialog-long-title">Remove grace.hopper.long.address@subdomain.example.com?</h2><button class="ui-dialog-x" id="dialog-long-x" aria-label="Close">x</button></div>
+    <div class="ui-dialog-body"><p>Body</p></div>
+  </dialog>
 </main>
 </body>
 </html>
```

- [ ] **7. Modify `e2e/options.spec.ts`.** Write this patch to `/tmp/plan-e2e-options.spec.ts.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-e2e-options.spec.ts.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/e2e/options.spec.ts b/e2e/options.spec.ts
index bcd3ca5..a594b8d 100644
--- a/e2e/options.spec.ts
+++ b/e2e/options.spec.ts
@@ -155,3 +155,21 @@ test.describe("on a touch screen", () => {
     expect((await rect(page, "#input-sm")).height).toBe(32);
   });
 });
+
+test.describe("on a phone", () => {
+  test.use({ viewport: { width: 360, height: 740 } });
+
+  test("a long title wraps inside its dialog and keeps the × in view", async ({ page }) => {
+    await page.goto("/fixture/library/ui.css");
+    await showModal(page, "#dialog-long");
+    const dialog = await rect(page, "#dialog-long");
+    const title = await rect(page, "#dialog-long-title");
+    const x = await rect(page, "#dialog-long-x");
+    expect(dialog.right).toBeLessThanOrEqual(360);
+    // Without the fix the title keeps the email's width (about 500px) and the dialog scrolls sideways.
+    expect(title.left).toBeGreaterThanOrEqual(dialog.left);
+    expect(title.right).toBeLessThanOrEqual(x.left);
+    expect(x.right).toBeLessThanOrEqual(dialog.right);
+    expect(await page.locator("#dialog-long").evaluate((el) => el.scrollWidth <= el.clientWidth)).toBe(true);
+  });
+});
```

- [ ] **8. Modify `e2e/tables.spec.ts`.** Write this patch to `/tmp/plan-e2e-tables.spec.ts.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-e2e-tables.spec.ts.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/e2e/tables.spec.ts b/e2e/tables.spec.ts
index f72ee48..7852bf1 100644
--- a/e2e/tables.spec.ts
+++ b/e2e/tables.spec.ts
@@ -45,6 +45,17 @@ test.describe("at desktop width", () => {
     expect(await css(page, "#wrap thead th >> nth=0", "background-color")).toBe(await token(page, "#library", "--band"));
   });
 
+  test("a quiet head's row headers go quiet too, keep their case and still tint on hover", async ({ page }) => {
+    expect(await css(page, "#caps-rowhead", "background-color")).toBe("rgba(0, 0, 0, 0)");
+    expect(await css(page, "#caps-rowhead", "color")).toBe(await token(page, "#library", "--ink-3"));
+    expect(await css(page, "#caps-rowhead", "font-size")).toBe("12px");
+    expect(await css(page, "#caps-rowhead", "text-transform")).toBe("none");
+    // A row header in a default table keeps its fill.
+    expect(await css(page, "#stack-rowhead", "background-color")).toBe(await token(page, "#library", "--band"));
+    await page.locator("#caps-td").hover();
+    expect(await css(page, "#caps-rowhead", "background-color")).not.toBe("rgba(0, 0, 0, 0)");
+  });
+
   test("a sticky head stays at the top of a bounded wrap as it scrolls", async ({ page }) => {
     const wrap = await box(page, "#sticky-wrap");
     expect(wrap.height).toBeLessThanOrEqual(128);
```

- [ ] **9.** Run `npx playwright test --ignore-snapshots e2e/options.spec.ts e2e/tables.spec.ts`. Expect every test to pass.

- [ ] **10. Commit:** `git add gallery/fixture_library.html e2e/options.spec.ts e2e/tables.spec.ts` then `git commit -m "test: a long dialog title and a quiet row header"`.

### Task 3: Docs and version

**Files:** `DESIGN.md`, `CHANGELOG.md`, `AGENTS.md`, `README.md`, `package.json`, `package-lock.json`

- [ ] **11. Modify `DESIGN.md`.** Write this patch to `/tmp/plan-DESIGN.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-DESIGN.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/DESIGN.md b/DESIGN.md
index e69a07e..b7f26ee 100644
--- a/DESIGN.md
+++ b/DESIGN.md
@@ -342,7 +342,8 @@ In Go, write this markup (it is pinned by `testdata/parity/dialog-basic.html`):
 </dialog>
 ```
 
-`.ui-dialog-actions` is the row of buttons at the end, and it can be the `<form>` itself.
+`.ui-dialog-actions` is the row of buttons at the end, and it can be the `<form>` itself. The title
+wraps anywhere it must, so a long email address in it stays inside the dialog.
 
 **An eyebrow** is a `<p class="ui-dialog-eyebrow">` as the dialog's first child, before the head:
 small uppercase text in `--ink-3`, such as "Can't be undone". Set `--ui-dialog-eyebrow-fg: var(--bad)`
@@ -403,7 +404,7 @@ card flat.
 | Option | Markup | What it does |
 |---|---|---|
 | No frame | `data-frame="none"` on `.ui-table-wrap` | No border, corners or fill; the first and last columns sit flush with the content around them. |
-| Quiet head | `data-head="quiet"` on `.ui-table` | Header cells on no fill, in small `--ink-3` text. |
+| Quiet head | `data-head="quiet"` on `.ui-table` | Header cells on no fill, in small `--ink-3` text. Row headers in the body (`<th scope="row">`) go quiet too, but keep their case. |
 | Sticky head | `data-sticky="head"` on `.ui-table` | Header cells stay at the top of the wrap as it scrolls. Bound the wrap with `--ui-table-max-height`. |
 | Stacked on a phone | `data-stack="phone"` on `.ui-table-wrap` | Below 800px of viewport, each row is a block and each cell a line with its label. |
 
````

- [ ] **12. Modify `CHANGELOG.md`.** Write this patch to `/tmp/plan-CHANGELOG.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-CHANGELOG.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/CHANGELOG.md b/CHANGELOG.md
index ad5f9be..9c11a1b 100644
--- a/CHANGELOG.md
+++ b/CHANGELOG.md
@@ -4,6 +4,13 @@ Breaking changes come first in each entry. Renaming or removing a token,
 class, data-attribute value, knob, partial name or argument, or React prop
 is breaking.
 
+## 0.3.2 (2026-11-12)
+
+No breaking changes. Two looks change:
+
+- A dialog's title shrinks and wraps anywhere it must, so a long word such as an email address no longer pushes the × or the dialog past the viewport.
+- In a table with `data-head="quiet"`, row headers in the body (`<th scope="row">`) go quiet too: no fill, in the head's colour and size. They keep their case, and their row still tints on hover.
+
 ## 0.3.1 (2026-11-11)
 
 No breaking changes. One look changes:
```

- [ ] **13. Modify `AGENTS.md`.** Write this patch to `/tmp/plan-AGENTS.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-AGENTS.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/AGENTS.md b/AGENTS.md
index 45ea1be..6ca7e65 100644
--- a/AGENTS.md
+++ b/AGENTS.md
@@ -27,7 +27,7 @@ Check which stack this project uses, then follow one of these.
 ### Go with `html/template`
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.3.1
+go get github.com/fuad-daoud/ai-ui@v0.3.2
 ```
 
 It needs Go 1.25, and `go get` raises the `go` line in an older `go.mod` to match.
@@ -60,13 +60,13 @@ read it too, then install:
 
 ```bash
 echo "allow-git=root" >> .npmrc
-npm install github:fuad-daoud/ai-ui#v0.3.1
+npm install github:fuad-daoud/ai-ui#v0.3.2
 ```
 
 That leaves this in `package.json`:
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.1" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.2" }
 ```
 
 ```ts
````

- [ ] **14. Modify `README.md`.** Write this patch to `/tmp/plan-README.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-README.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/README.md b/README.md
index b916342..91ba372 100644
--- a/README.md
+++ b/README.md
@@ -8,7 +8,7 @@ an agent starts. The site is at https://ui.fuad-daoud.com/.
 ## Go apps
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.3.1
+go get github.com/fuad-daoud/ai-ui@v0.3.2
 ```
 
 It needs Go 1.25.
@@ -41,11 +41,11 @@ allows it, so allow it in a `.npmrc` (which `npm ci` and CI read too), then inst
 
 ```bash
 echo "allow-git=root" >> .npmrc
-npm install github:fuad-daoud/ai-ui#v0.3.1
+npm install github:fuad-daoud/ai-ui#v0.3.2
 ```
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.1" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.2" }
 ```
 
 ```ts
````

- [ ] **15. Modify `package.json`.** Write this patch to `/tmp/plan-package.json.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-package.json.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/package.json b/package.json
index 5da215f..1f1dda0 100644
--- a/package.json
+++ b/package.json
@@ -1,6 +1,6 @@
 {
   "name": "ai-ui",
-  "version": "0.3.1",
+  "version": "0.3.2",
   "description": "A small UI library written to be used by coding agents: tokens, layered CSS, Go partials and React components.",
   "license": "MIT",
   "type": "module",
```

- [ ] **16. Modify `package-lock.json`.** Write this patch to `/tmp/plan-package-lock.json.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/plan-package-lock.json.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/package-lock.json b/package-lock.json
index 0fe9292..2208f54 100644
--- a/package-lock.json
+++ b/package-lock.json
@@ -1,12 +1,12 @@
 {
   "name": "ai-ui",
-  "version": "0.3.1",
+  "version": "0.3.2",
   "lockfileVersion": 3,
   "requires": true,
   "packages": {
     "": {
       "name": "ai-ui",
-      "version": "0.3.1",
+      "version": "0.3.2",
       "license": "MIT",
       "devDependencies": {
         "@axe-core/playwright": "^4.14.0",
```

- [ ] **17.** Run the gate above. Expect it all to pass: 61 vitest tests, 64 browser tests, 57 site tests.

- [ ] **18. Commit:** `git add DESIGN.md CHANGELOG.md AGENTS.md README.md package.json package-lock.json` then `git commit -m "docs: v0.3.2 in DESIGN.md and the changelog; 0.3.2"`.

