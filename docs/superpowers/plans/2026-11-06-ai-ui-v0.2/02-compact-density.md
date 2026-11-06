# ai-ui v0.2, round 2: compact density and the field fix

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Precondition (STOP if unmet):** `git log --oneline` shows the commit `feat: ui-scoped.css styles library markup and nothing else`, and `git status --short` is empty.

**Goal:** `data-density="compact"` gives 32px controls, 14px text and 4px corners. Inputs and selects match buttons at both densities.

**Architecture:** One new source file, `css/components/density.css`, adds rules under the attribute and changes no existing rule. It redefines four tokens and sets the knobs those tokens don't reach. Separately, the field fix removes block padding from single-line controls in `field.css`.

**Tech Stack:** CSS layers, Go `html/template` (gallery), Playwright with axe.

**Spec:** `docs/superpowers/specs/2026-11-06-v0.2-compact-and-scoped-design.md`, §3, §5 and §7. Read it first.

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

- The gallery's compact pages have no baselines yet, which is why every run passes `--ignore-snapshots`. The axe checks still run on them.

## Review Focus

1. **`rem` shrinking:** `font-size` on `:root` would turn every `2rem` control into 28px. Pinned by the `html` font-size check.
2. **A `<select multiple>`:** it lists rows, so the field fix must not take its padding. Pinned by the `#multi` check.
3. **A compact dialog:** the top layer must not lose compact. Pinned by the `#dialog` check.
4. **Compact in another palette or mode:** colours must still pass. Pinned by the gallery's compact pages, which run axe in blue, mono, light and dark.
5. **A page with no attribute:** it must render exactly as before, apart from the field fix. Pinned by `compact has no effect without the attribute`, and by the baselines the planner regenerates.

---

### Task 1: The field fix

**Files:**
- Modify: `css/components/field.css`
- Create: `e2e/density.spec.ts`

**Interfaces:**
- Consumes: `/fixture/library/{sheet}`, and `snapshot` and `differences` from `e2e/support.ts` (round 1).

- [ ] **Step 1: Write the failing tests**

`e2e/density.spec.ts`:

```ts
import { expect, test, type Page } from "@playwright/test";
import { differences, snapshot } from "./support";

type Box = { height: number; fontSize: string; radius: string };

const box = (page: Page, sel: string): Promise<Box> =>
  page.locator(sel).evaluate((el) => {
    const c = getComputedStyle(el);
    return { height: el.getBoundingClientRect().height, fontSize: c.fontSize, radius: c.borderTopLeftRadius };
  });

const css = (page: Page, sel: string, prop: string): Promise<string> =>
  page.locator(sel).evaluate((el, prop) => getComputedStyle(el).getPropertyValue(prop), prop);

for (const density of ["root", "container"]) {
  test(`compact sizes with the attribute on the ${density}`, async ({ page }) => {
    await page.goto(`/fixture/library/ui.css?density=${density}`);
    expect(await box(page, "#btn")).toEqual({ height: 32, fontSize: "14px", radius: "4px" });
    expect(await box(page, "#btn-sm")).toEqual({ height: 28, fontSize: "13px", radius: "4px" });
    expect(await box(page, "#btn-lg")).toEqual({ height: 44, fontSize: "16px", radius: "4px" });
    expect(await box(page, "#input")).toEqual({ height: 32, fontSize: "14px", radius: "4px" });
    expect(await box(page, "#select")).toEqual({ height: 32, fontSize: "14px", radius: "4px" });
    expect(await box(page, "#area")).toEqual({ height: 80, fontSize: "14px", radius: "4px" });
    for (const sel of ["#label", "#hint", "#error"]) expect(await css(page, sel, "font-size")).toBe("13px");
    expect(await css(page, "#field", "row-gap")).toBe("4px");
    expect(await css(page, "#pill", "border-top-left-radius")).toBe("4px");
    expect(await css(page, "#callout", "border-top-left-radius")).toBe("4px");
    expect(await css(page, "#wrap", "border-top-left-radius")).toBe("4px");
    expect(await css(page, "#card", "border-top-left-radius")).toBe("8px");
    expect(await css(page, "#td", "font-size")).toBe("14px");
    // A dialog stays in its place in the tree, so it inherits compact in the top layer too.
    await page.locator("#dialog").evaluate((d: HTMLDialogElement) => d.showModal());
    expect(await css(page, "#dialog", "border-top-left-radius")).toBe("8px");
    // font-size on :root would shrink rem and every rem-sized control with it.
    expect(await css(page, "html", "font-size")).toBe("16px");
    expect(await css(page, "body", "font-size")).toBe(density === "root" ? "14px" : "16px");
    expect(await css(page, "#library", "font-size")).toBe("14px");
  });
}

test("default sizes", async ({ page }) => {
  await page.goto("/fixture/library/ui.css");
  expect(await box(page, "#btn")).toEqual({ height: 40, fontSize: "16px", radius: "6px" });
  expect(await box(page, "#btn-sm")).toEqual({ height: 32, fontSize: "14px", radius: "6px" });
  expect(await box(page, "#btn-lg")).toEqual({ height: 48, fontSize: "18px", radius: "6px" });
  expect(await css(page, "#card", "border-top-left-radius")).toBe("10px");
  expect(await css(page, "#pill", "border-top-left-radius")).toBe("999px");
});

for (const query of ["", "?density=root"]) {
  test(`fields match buttons (${query || "default"})`, async ({ page }) => {
    await page.goto(`/fixture/library/ui.css${query}`);
    const button = (await box(page, "#btn")).height;
    expect((await box(page, "#input")).height).toBe(button);
    expect((await box(page, "#select")).height).toBe(button);
    // A multi-select lists rows, so it keeps its padding.
    expect(await css(page, "#multi", "padding-top")).toBe("8px");
    expect(await css(page, "#area", "padding-top")).toBe("8px");
  });
}

test("compact has no effect without the attribute", async ({ page }) => {
  await page.goto("/p/blue/light");
  const before = await snapshot(page, "body");
  const removed = await page.evaluate(() => {
    let n = 0;
    const walk = (list: CSSRuleList, remove: (i: number) => void) => {
      for (let i = list.length - 1; i >= 0; i--) {
        const rule = list[i];
        if (rule instanceof CSSStyleRule && rule.selectorText.includes("data-density")) {
          remove(i);
          n++;
        } else if ("cssRules" in rule) {
          const group = rule as CSSGroupingRule;
          walk(group.cssRules, (j) => group.deleteRule(j));
        }
      }
    };
    for (const sheet of Array.from(document.styleSheets)) walk(sheet.cssRules, (i) => sheet.deleteRule(i));
    return n;
  });
  expect(removed, "the page had compact rules to remove").toBeGreaterThan(0);
  expect(differences(before, await snapshot(page, "body"))).toEqual([]);
});
```

- [ ] **Step 2: Run them**

Run: `npx playwright test e2e/density.spec.ts --ignore-snapshots`
Expected: FAIL, every test except `default sizes`. `fields match buttons (default)` reports the input at about 43.6px against a 40px button.

- [ ] **Step 3: Fix the fields**

`css/components/field.css`:

```diff
--- a/css/components/field.css
+++ b/css/components/field.css
@@ -22,4 +22,8 @@
   transition: border-color var(--dur) var(--ease);
 }
+/* A single-line control is centred by its height, so block padding would only make it taller
+   than a button beside it. A multi-select lists rows, and keeps its padding. */
+.ui-input,
+.ui-select:not([multiple]) { padding-block: 0; }
 .ui-textarea { min-block-size: calc(var(--ui-input-height) * 2); resize: vertical; line-height: 1.5; }
 .ui-select {
```

Run: `go run ./cmd/build && npx playwright test e2e/density.spec.ts -g "fields match buttons \(default\)" --ignore-snapshots`
Expected: PASS.

### Task 2: Compact density

**Files:**
- Create: `css/components/density.css`
- Modify: `classes_test.go`, `gallery/page.html`, `gallery/main.go`, `gallery/main_test.go`, `e2e/gallery.spec.ts`, `dist/ui.css`, `dist/ui-scoped.css`

**Interfaces:**
- Produces: `data-density="compact"`, and gallery routes `/p/{palette}/{mode}/compact` and `/p/{palette}/{mode}/compact/dialog`.

- [ ] **Step 1: Write the failing Go tests**

`classes_test.go`:

```diff
--- a/classes_test.go
+++ b/classes_test.go
@@ -4,4 +4,5 @@
 	"os"
 	"regexp"
+	"strings"
 	"testing"
 )
@@ -34,2 +35,14 @@
 	}
 }
+
+// The compact density is part of the contract (DESIGN.md), and the default rules
+// must not depend on it.
+func TestCompactDensityIsDefined(t *testing.T) {
+	css, err := os.ReadFile("dist/ui.css")
+	if err != nil {
+		t.Fatal(err)
+	}
+	if !strings.Contains(string(css), `[data-density="compact"]`) {
+		t.Error(`[data-density="compact"] is documented but not defined in dist/ui.css`)
+	}
+}
```

`gallery/main_test.go`:

```diff
--- a/gallery/main_test.go
+++ b/gallery/main_test.go
@@ -25,5 +25,5 @@
 	for _, p := range palettes {
 		for _, m := range modes {
-			for _, suffix := range []string{"", "/dialog"} {
+			for _, suffix := range []string{"", "/dialog", "/compact", "/compact/dialog"} {
 				path := "/p/" + p + "/" + m + suffix
 				rec := fetch(t, h, path)
@@ -43,6 +43,9 @@
 					t.Errorf("%s: a missing key leaked into the markup", path)
 				}
-				if opens := strings.Contains(body, "showModal()"); opens != (suffix == "/dialog") {
+				if opens := strings.Contains(body, "showModal()"); opens != strings.HasSuffix(suffix, "/dialog") {
 					t.Errorf("%s: showModal present = %v", path, opens)
+				}
+				if compact := strings.Contains(body, `data-density="compact"`); compact != strings.HasPrefix(suffix, "/compact") {
+					t.Errorf("%s: data-density present = %v", path, compact)
 				}
 			}
@@ -56,5 +59,5 @@
 		t.Fatal(err)
 	}
-	for _, path := range []string{"/p/rose/light", "/p/blue/sepia", "/p/blue", "/fixture/library/other.css", "/fixture/document/theme.js"} {
+	for _, path := range []string{"/p/rose/light", "/p/blue/sepia", "/p/blue", "/p/blue/light/compact/x", "/fixture/library/other.css", "/fixture/document/theme.js"} {
 		if rec := fetch(t, h, path); rec.Code != http.StatusNotFound {
 			t.Errorf("%s: status %d, want 404", path, rec.Code)
```

Run: `go test -count=1 . ./gallery`
Expected: FAIL. `TestCompactDensityIsDefined` finds no rule, and the `/compact` pages return 404.

- [ ] **Step 2: Add the density rules**

`css/components/density.css`:

```css
/* Compact density (DESIGN.md). Every rule sits under the attribute, so a page
   without it renders as before. Tokens are redefined here, in a later layer than
   their definition, and an app's unlayered override still beats all of this. */
[data-density="compact"] {
  --text-md: 0.875rem;
  --text-sm: 0.8125rem;
  --radius-md: 4px;
  --radius-lg: 8px;
}
/* Not on :root: font-size there would shrink rem, and every rem-sized control
   with it. On the root, body picks up --text-md through its own rule. */
[data-density="compact"]:not(:root) { font-size: var(--text-md); }
[data-density="compact"] .ui-btn { --ui-btn-height: 2rem; --ui-btn-pad: var(--space-3); }
[data-density="compact"] .ui-btn[data-size="sm"] { --ui-btn-height: 1.75rem; --ui-btn-pad: 0.625rem; }
[data-density="compact"] .ui-btn[data-size="lg"] { --ui-btn-height: 2.75rem; --ui-btn-pad: 1.25rem; --ui-btn-font-size: 1rem; }
[data-density="compact"] :is(.ui-input, .ui-select) { --ui-input-height: 2rem; }
[data-density="compact"] .ui-field { --ui-field-gap: var(--space-1); }
[data-density="compact"] .ui-pill { border-radius: var(--radius-sm); }
[data-density="compact"] .ui-table { font-size: var(--text-md); }
```

- [ ] **Step 3: Add the gallery pages**

`gallery/page.html`:

```diff
--- a/gallery/page.html
+++ b/gallery/page.html
@@ -1,4 +1,4 @@
 <!doctype html>
-<html lang="en" data-theme="{{.Mode}}"{{if ne .Palette "blue"}} data-palette="{{.Palette}}"{{end}}>
+<html lang="en" data-theme="{{.Mode}}"{{if ne .Palette "blue"}} data-palette="{{.Palette}}"{{end}}{{if .Compact}} data-density="compact"{{end}}>
 <head>
 <meta charset="utf-8">
```

`gallery/main.go`:

```diff
--- a/gallery/main.go
+++ b/gallery/main.go
@@ -46,4 +46,5 @@
 	Palette, Mode string
 	Dialog        bool
+	Compact       bool
 	Links         []link
 	Currencies    []ui.Option
@@ -80,5 +81,5 @@
 		http.Redirect(w, r, "/p/blue/light", http.StatusFound)
 	})
-	serve := func(dialog bool) http.HandlerFunc {
+	serve := func(dialog, compact bool) http.HandlerFunc {
 		return func(w http.ResponseWriter, r *http.Request) {
 			p, m := r.PathValue("palette"), r.PathValue("mode")
@@ -88,5 +89,5 @@
 			}
 			var b bytes.Buffer
-			data := page{Palette: p, Mode: m, Dialog: dialog, Links: links(p, m), Currencies: currencies}
+			data := page{Palette: p, Mode: m, Dialog: dialog, Compact: compact, Links: links(p, m), Currencies: currencies}
 			if err := set.ExecuteTemplate(&b, "page", data); err != nil {
 				http.Error(w, err.Error(), http.StatusInternalServerError)
@@ -97,6 +98,8 @@
 		}
 	}
-	mux.HandleFunc("GET /p/{palette}/{mode}", serve(false))
-	mux.HandleFunc("GET /p/{palette}/{mode}/dialog", serve(true))
+	mux.HandleFunc("GET /p/{palette}/{mode}", serve(false, false))
+	mux.HandleFunc("GET /p/{palette}/{mode}/dialog", serve(true, false))
+	mux.HandleFunc("GET /p/{palette}/{mode}/compact", serve(false, true))
+	mux.HandleFunc("GET /p/{palette}/{mode}/compact/dialog", serve(true, true))
 	serveFixture := func(name string) http.HandlerFunc {
 		return func(w http.ResponseWriter, r *http.Request) {
```

`e2e/gallery.spec.ts`:

```diff
--- a/e2e/gallery.spec.ts
+++ b/e2e/gallery.spec.ts
@@ -3,5 +3,10 @@
 
 const paths = ["blue", "mono"].flatMap((p) =>
-  ["light", "dark"].flatMap((m) => [`/p/${p}/${m}`, `/p/${p}/${m}/dialog`]),
+  ["light", "dark"].flatMap((m) => [
+    `/p/${p}/${m}`,
+    `/p/${p}/${m}/dialog`,
+    `/p/${p}/${m}/compact`,
+    `/p/${p}/${m}/compact/dialog`,
+  ]),
 );
 
```

- [ ] **Step 4: Build and run**

Run: `go run ./cmd/build && go test -count=1 ./... && npx playwright test --ignore-snapshots`
Expected: PASS, 30 browser tests. `dist/ui.css` is 17,911 bytes and `dist/ui-scoped.css` is 21,345 bytes.

- [ ] **Step 5: Prove the `rem` guard can fail**

In `css/components/density.css`, temporarily change `[data-density="compact"]:not(:root)` to `[data-density="compact"]`, then run `go run ./cmd/build && npx playwright test e2e/density.spec.ts --ignore-snapshots`.
Expected: `compact sizes with the attribute on the root` and `fields match buttons (?density=root)` fail. Paste that in your report.
Restore it and run `go run ./cmd/build`.

- [ ] **Step 6: Run the gate, then commit**

```bash
git add css/components/field.css css/components/density.css classes_test.go gallery/page.html gallery/main.go gallery/main_test.go e2e/gallery.spec.ts e2e/density.spec.ts dist/ui.css dist/ui-scoped.css
git commit -m "feat: compact density, and fields match buttons (plan 02 of 2026-11-06-ai-ui-v0.2)"
```
