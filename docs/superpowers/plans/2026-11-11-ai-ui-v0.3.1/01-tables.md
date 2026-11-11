# ai-ui v0.3.1: table options

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Precondition (STOP if unmet):** `git log --oneline -3` shows `docs: plan ai-ui v0.3.1` at the top, and `git status --short` is empty.

**Goal:** tables take four options (no frame, a quiet head, a sticky head, rows stacked on a phone) and new knobs; row hover applies only where a pointer can hover; the docs, the version (0.3.1) and the site follow.

**Architecture:** CSS only in the library (`table.css`, one rule in `utilities.css` because `.ui-num` lives in the later utilities layer). Tables stay plain markup: no Go partial or React change.

**Tech Stack:** CSS layers, Playwright with axe, Markdown, the Go site generator.

**Spec:** `docs/superpowers/specs/2026-11-11-v0.3.1-tables-design.md`. Read it first.

## Global Constraints

- Go module `github.com/fuad-daoud/ai-ui`, `go 1.25`. No new Go or npm dependency.
- Modify **only** the files this plan names. Need another? **STOP and say so.**
- **Every patch and file in this plan was applied to a trial copy and passed the gate.** Apply each one exactly: write the patch to the named file under `/tmp` byte for byte, apply it with `git apply`, then delete it. If a step's result differs from what it says, find out why before changing anything, and report it. Don't change a test to make it pass.
- `dist/ui.css` and `dist/ui-scoped.css` are generated: after the CSS task run `go run ./cmd/build` and commit `dist/` with it. Never edit `dist/` by hand.
- **Never run `--update-snapshots` and never add, change or delete anything in `e2e/__screenshots__`.** The planner makes baselines on contabo. Always pass `--ignore-snapshots`.
- Run `npm ci` first if `node_modules/` is missing.
- Never run git except `git apply <patch>`, `git add <named files>`, `git commit`, `git status --short`, `git diff`, `git log`, `git show`. No amend, `add -A`, `add .`, stash, reset, checkout, branch, push, merge, rebase or tag. Never `git stash` in a relevo worktree.
- No scratch files in the repository; patches go in `/tmp` and are deleted after applying.

**Gate:** `test -z "$(gofmt -l .)" && go vet ./... && go test -count=1 ./... && (cd site && go vet ./... && go test -count=1 ./...) && npx tsc --noEmit && npx vitest run && npx playwright test --ignore-snapshots && npx playwright test -c playwright.site.config.ts`

**Report:** every file touched; every step's real result; every place the plan's wording did not fit the code and what you did. `git status --short` is empty after the last commit.

## Review Focus

1. **`.ui-num` in a stacked row:** the utilities layer comes after components, so the start alignment has to live in `utilities.css`. Pinned by `rows stack into labelled lines`.
2. **The header row stays in the markup** when stacked (hidden from sight only). Pinned by the same test.
3. **Hover only where it exists.** Pinned by `a tapped row keeps no tint` (touch emulation) and `a hovered row is tinted` (desktop).
4. **Nothing changes for a table without the new attributes,** apart from hover on touch screens: the default head still has its `--band` fill (checked in the quiet-head test).

---

### Task 1: Table CSS

**Files:** `css/components/table.css`, `css/utilities.css`

- [ ] **1. Modify `css/components/table.css`.** Write this patch to `/tmp/v031-css-components-table.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-css-components-table.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/table.css b/css/components/table.css
index d1921cb..7b93a81 100644
--- a/css/components/table.css
+++ b/css/components/table.css
@@ -1,5 +1,7 @@
 .ui-table-wrap {
-  overflow-x: auto;
+  --ui-table-max-height: none;
+  overflow: auto;
+  max-block-size: var(--ui-table-max-height);
   border: 1px solid var(--rule-soft);
   border-radius: var(--radius-md);
   background: var(--sheet);
@@ -7,6 +9,14 @@
 .ui-table {
   --ui-table-pad-y: var(--space-2);
   --ui-table-pad-x: var(--space-3);
+  --ui-table-head-bg: var(--band);
+  --ui-table-head-fg: var(--ink-2);
+  --ui-table-head-size: 1em;
+  --ui-table-head-case: none;
+  --ui-table-head-rule: var(--rule-soft);
+  --ui-table-hover: color-mix(in oklch, var(--band) 60%, transparent);
+  --ui-table-sticky-bg: var(--sheet);
+  --ui-table-label-width: 8rem;
   inline-size: 100%;
   border-collapse: collapse;
   font-size: var(--text-sm);
@@ -19,8 +29,18 @@
   vertical-align: top;
 }
 .ui-table th { background: var(--band); color: var(--ink-2); font-weight: 600; white-space: nowrap; }
+.ui-table thead th {
+  background: var(--ui-table-head-bg);
+  color: var(--ui-table-head-fg);
+  font-size: var(--ui-table-head-size);
+  text-transform: var(--ui-table-head-case);
+  border-block-end-color: var(--ui-table-head-rule);
+}
 .ui-table tbody tr:last-child > * { border-block-end: 0; }
-.ui-table tbody tr:hover > * { background: color-mix(in oklch, var(--band) 60%, transparent); }
+/* Only where a pointer can hover, so a tapped row on a phone doesn't keep its tint. */
+@media (hover: hover) {
+  .ui-table tbody tr:hover > * { background: var(--ui-table-hover); }
+}
 .ui-table caption {
   caption-side: top;
   padding: var(--space-2) var(--ui-table-pad-x);
@@ -28,3 +48,59 @@
   font-size: var(--text-sm);
   text-align: start;
 }
+
+/* A quiet header row: no fill, small grey text. */
+.ui-table[data-head="quiet"] {
+  --ui-table-head-bg: transparent;
+  --ui-table-head-fg: var(--ink-3);
+  --ui-table-head-size: var(--text-xs);
+}
+/* The header row stays at the top of a wrap that scrolls (bound it with --ui-table-max-height). */
+.ui-table[data-sticky="head"] thead th {
+  position: sticky;
+  inset-block-start: 0;
+  z-index: 1;
+  background: var(--ui-table-sticky-bg);
+}
+
+/* No frame: the table sits on the page or in a card, its outer columns flush with the content. */
+.ui-table-wrap[data-frame="none"] { border: 0; border-radius: 0; background: transparent; }
+.ui-table-wrap[data-frame="none"] .ui-table :is(th, td):first-child { padding-inline-start: 0; }
+.ui-table-wrap[data-frame="none"] .ui-table :is(th, td):last-child { padding-inline-end: 0; }
+
+/* On a phone, each row is a block and each cell a labelled line. 800px is the apps' phone boundary. */
+@media (max-width: 799.98px) {
+  .ui-table-wrap[data-stack="phone"] { overflow: visible; max-block-size: none; }
+  .ui-table-wrap[data-stack="phone"] :is(.ui-table, .ui-table > tbody, .ui-table > tbody > tr) { display: block; }
+  .ui-table-wrap[data-stack="phone"] .ui-table > thead {
+    position: absolute;
+    inline-size: 1px;
+    block-size: 1px;
+    overflow: hidden;
+    clip-path: inset(50%);
+    white-space: nowrap;
+  }
+  .ui-table-wrap[data-stack="phone"] .ui-table > tbody > tr {
+    padding-block: var(--space-2);
+    border-block-end: 1px solid var(--rule-soft);
+  }
+  .ui-table-wrap[data-stack="phone"] .ui-table > tbody > tr:last-child { border-block-end: 0; }
+  .ui-table-wrap[data-stack="phone"] .ui-table > tbody > tr > * {
+    display: block;
+    padding-block: var(--space-1);
+    border-block-end: 0;
+    text-align: start;
+  }
+  .ui-table-wrap[data-stack="phone"] .ui-table > tbody > tr > [data-label] {
+    display: grid;
+    grid-template-columns: var(--ui-table-label-width) minmax(0, 1fr);
+    gap: var(--space-3);
+  }
+  .ui-table-wrap[data-stack="phone"] .ui-table > tbody > tr > [data-label]::before {
+    content: attr(data-label);
+    color: var(--ink-3);
+    font-size: var(--text-xs);
+    font-weight: 600;
+  }
+  .ui-table-wrap[data-stack="phone"][data-frame="none"] .ui-table > tbody > tr > * { padding-inline: 0; }
+}
```

- [ ] **2. Modify `css/utilities.css`.** Write this patch to `/tmp/v031-css-utilities.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-css-utilities.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/utilities.css b/css/utilities.css
index 9b8ddc0..61c2d94 100644
--- a/css/utilities.css
+++ b/css/utilities.css
@@ -1,6 +1,10 @@
 .ui-muted { color: var(--ink-3); }
 .ui-mono { font-family: var(--font-mono); }
 .ui-num { font-variant-numeric: tabular-nums; text-align: end; }
+/* In a table stacked on a phone, a number sits beside its label, so it starts at the start. */
+@media (max-width: 799.98px) {
+  .ui-table-wrap[data-stack="phone"] .ui-table > tbody > tr > .ui-num { text-align: start; }
+}
 .ui-vh {
   position: absolute;
   inline-size: 1px;
```

- [ ] **3.** Run `go run ./cmd/build`. Expect `wrote dist/ui.css (27033 bytes)` and `wrote dist/ui-scoped.css (30467 bytes)`.

- [ ] **4.** Run `go test -count=1 .`. Expect `ok`.

- [ ] **5. Commit:** `git add css/components/table.css css/utilities.css dist/ui.css dist/ui-scoped.css` then `git commit -m "feat(css): table options: no frame, quiet head, sticky head, stacked on a phone"`.

### Task 2: Fixture, gallery and browser tests

**Files:** `gallery/fixture_library.html`, `gallery/page.html`, `e2e/tables.spec.ts`

- [ ] **6. Modify `gallery/fixture_library.html`.** Write this patch to `/tmp/v031-gallery-fixture_library.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-gallery-fixture_library.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/gallery/fixture_library.html b/gallery/fixture_library.html
index 9eb5b25..c417757 100644
--- a/gallery/fixture_library.html
+++ b/gallery/fixture_library.html
@@ -55,6 +55,10 @@
   <div class="ui-cluster"><span class="ui-pill" id="pill">Neutral</span><span class="ui-pill" data-variant="good">Good</span><span class="ui-pill" data-variant="bad">Bad</span><span class="ui-pill" id="pill-dot" data-variant="bad" data-dot>1 violation</span><span class="ui-pill" id="pill-outline" data-tone="outline">meta</span><span class="ui-pill" id="pill-outline-warn" data-variant="warning" data-tone="outline">Price up</span></div>
   <div class="ui-callout" id="callout" data-variant="bad"><p class="ui-callout-title">Title</p><p>Body with <a id="callout-link" href="#library">a link</a> and <strong>bold</strong>.</p></div>
   <div class="ui-table-wrap" id="wrap"><table class="ui-table"><caption>Caption</caption><thead><tr><th scope="col">A</th><th scope="col" class="ui-num">B</th></tr></thead><tbody><tr><td id="td">x</td><td class="ui-num">1</td></tr></tbody></table></div>
+  <div class="ui-table-wrap" id="plain-wrap" data-frame="none"><table class="ui-table" data-head="quiet"><thead><tr><th scope="col" id="plain-th">Name</th><th scope="col" class="ui-num" id="plain-th-num">Count</th></tr></thead><tbody><tr id="plain-row"><td id="plain-first">a</td><td class="ui-num" id="plain-last">1</td></tr></tbody></table></div>
+  <div class="ui-table-wrap"><table class="ui-table" data-head="quiet" style="--ui-table-head-case: uppercase; --ui-table-head-rule: var(--rule)"><thead><tr><th scope="col" id="caps-th">Email</th></tr></thead><tbody><tr><td>x</td></tr></tbody></table></div>
+  <div class="ui-table-wrap" id="sticky-wrap" style="--ui-table-max-height: 8rem"><table class="ui-table" data-sticky="head"><thead><tr><th scope="col" id="sticky-th">Day</th></tr></thead><tbody><tr><td>Day 1</td></tr><tr><td>Day 2</td></tr><tr><td>Day 3</td></tr><tr><td>Day 4</td></tr><tr><td>Day 5</td></tr><tr><td>Day 6</td></tr><tr><td>Day 7</td></tr><tr><td>Day 8</td></tr><tr><td>Day 9</td></tr><tr><td>Day 10</td></tr><tr><td>Day 11</td></tr><tr><td>Day 12</td></tr></tbody></table></div>
+  <div class="ui-table-wrap" id="stack-wrap" data-stack="phone"><table class="ui-table" id="stack"><thead id="stack-head"><tr><th scope="col">Email</th><th scope="col">Role</th><th scope="col" class="ui-num">Logins</th></tr></thead><tbody><tr id="stack-row"><td id="stack-name">ada@example.com</td><td id="stack-role" data-label="Role">owner</td><td id="stack-num" data-label="Logins" class="ui-num">12</td></tr><tr><th scope="row" id="stack-rowhead">grace@example.com</th><td data-label="Role">member</td><td data-label="Logins" class="ui-num">3</td></tr></tbody></table></div>
   <details class="ui-disclosure" open><summary class="ui-disclosure-summary">More</summary><div class="ui-disclosure-body"><p>Body</p><ul><li>Item</li></ul></div></details>
   <details class="ui-disclosure" id="disc-start" data-marker="start"><summary class="ui-disclosure-summary" id="disc-start-summary">Folder</summary><div class="ui-disclosure-body"><p>File</p></div></details>
   <details class="ui-disclosure" id="disc-norule" style="--ui-disclosure-rule: transparent"><summary class="ui-disclosure-summary">No rule</summary><div class="ui-disclosure-body"><p>Body</p></div></details>
```

- [ ] **7. Modify `gallery/page.html`.** Write this patch to `/tmp/v031-gallery-page.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-gallery-page.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/gallery/page.html b/gallery/page.html
index ec83674..19de7a0 100644
--- a/gallery/page.html
+++ b/gallery/page.html
@@ -123,6 +123,16 @@
         </tbody>
       </table>
     </div>
+    <p class="gallery-row ui-muted">Unframed, with a quiet head; on a phone its rows stack.</p>
+    <div class="ui-table-wrap" data-frame="none" data-stack="phone">
+      <table class="ui-table" data-head="quiet">
+        <thead><tr><th scope="col">Account</th><th scope="col">Role</th><th scope="col">State</th><th scope="col" class="ui-num">Sign-ins</th></tr></thead>
+        <tbody>
+          <tr><td>ada@example.com</td><td data-label="Role">owner</td><td data-label="State">{{template "ui.pill" (dict "Label" "active" "Variant" "good" "Dot" true)}}</td><td data-label="Sign-ins" class="ui-num">41</td></tr>
+          <tr><td>grace@example.com</td><td data-label="Role">member</td><td data-label="State">{{template "ui.pill" (dict "Label" "disabled" "Variant" "bad" "Dot" true)}}</td><td data-label="Sign-ins" class="ui-num">3</td></tr>
+        </tbody>
+      </table>
+    </div>
   </section>
 
   <section aria-labelledby="s-disclosure">
```

- [ ] **8. Create `e2e/tables.spec.ts`** with exactly this content:

```ts
import AxeBuilder from "@axe-core/playwright";
import { expect, test, type Page } from "@playwright/test";

// The v0.3.1 table options, measured on the library fixture.

const css = (page: Page, sel: string, prop: string, pseudo: string | null = null): Promise<string> =>
  page.locator(sel).evaluate((el, [prop, pseudo]) => getComputedStyle(el, pseudo).getPropertyValue(prop), [prop, pseudo] as const);

const token = (page: Page, sel: string, name: string): Promise<string> =>
  page.locator(sel).evaluate((el, name) => {
    const probe = document.createElement("span");
    probe.style.color = `var(${name})`;
    el.appendChild(probe);
    const c = getComputedStyle(probe).color;
    probe.remove();
    return c;
  }, name);

const box = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => {
    const r = el.getBoundingClientRect();
    return { top: r.top, left: r.left, width: r.width, height: r.height };
  });

test.describe("at desktop width", () => {
  test.beforeEach(async ({ page }) => page.goto("/fixture/library/ui.css"));

  test("an unframed table has no frame, and its outer columns sit flush", async ({ page }) => {
    expect(await css(page, "#plain-wrap", "border-top-width")).toBe("0px");
    expect(await css(page, "#plain-wrap", "border-top-left-radius")).toBe("0px");
    expect(await css(page, "#plain-wrap", "background-color")).toBe("rgba(0, 0, 0, 0)");
    expect(await css(page, "#plain-first", "padding-left")).toBe("0px");
    expect(await css(page, "#plain-last", "padding-right")).toBe("0px");
    expect(await css(page, "#plain-th", "padding-left")).toBe("0px");
  });

  test("a quiet head has no fill and small grey text, and its knobs take effect", async ({ page }) => {
    expect(await css(page, "#plain-th", "background-color")).toBe("rgba(0, 0, 0, 0)");
    expect(await css(page, "#plain-th", "color")).toBe(await token(page, "#library", "--ink-3"));
    expect(await css(page, "#plain-th", "font-size")).toBe("12px");
    expect(await css(page, "#plain-th", "text-transform")).toBe("none");
    expect(await css(page, "#caps-th", "text-transform")).toBe("uppercase");
    expect(await css(page, "#caps-th", "border-bottom-color")).toBe(await token(page, "#library", "--rule"));
    // A default head keeps its fill.
    expect(await css(page, "#wrap thead th >> nth=0", "background-color")).toBe(await token(page, "#library", "--band"));
  });

  test("a sticky head stays at the top of a bounded wrap as it scrolls", async ({ page }) => {
    const wrap = await box(page, "#sticky-wrap");
    expect(wrap.height).toBeLessThanOrEqual(128);
    await page.locator("#sticky-wrap").evaluate((el) => el.scrollTo(0, 120));
    expect(await page.locator("#sticky-wrap").evaluate((el) => el.scrollTop)).toBeGreaterThan(0);
    const th = await box(page, "#sticky-th");
    expect(Math.abs(th.top - (wrap.top + 1))).toBeLessThan(1);
    expect(await css(page, "#sticky-th", "background-color")).toBe(await token(page, "#library", "--sheet"));
  });

  test("a stacked table is still a table at desktop width", async ({ page }) => {
    expect(await css(page, "#stack", "display")).toBe("table");
    expect(await css(page, "#stack-row", "display")).toBe("table-row");
    expect(await css(page, "#stack-role", "content", "::before")).toBe("none");
  });

  test("a hovered row is tinted where the pointer can hover", async ({ page }) => {
    await page.locator("#td").hover();
    expect(await css(page, "#td", "background-color")).not.toBe("rgba(0, 0, 0, 0)");
  });
});

test.describe("on a phone", () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test("rows stack into labelled lines", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    expect(await css(page, "#stack", "display")).toBe("block");
    expect(await css(page, "#stack-row", "display")).toBe("block");
    // A labelled cell shows its label; the name cell has none and spans the row.
    expect(await css(page, "#stack-role", "content", "::before")).toBe('"Role"');
    expect(await page.locator("#stack-role").evaluate((el) => getComputedStyle(el, "::before").getPropertyValue("width"))).toBe("128px");
    expect(await css(page, "#stack-name", "content", "::before")).toBe("none");
    const row = await box(page, "#stack-row");
    expect((await box(page, "#stack-name")).width).toBe(row.width);
    expect(await css(page, "#stack-rowhead", "display")).toBe("block");
    // The header row is hidden from sight but stays in the markup.
    expect((await box(page, "#stack-head")).width).toBeLessThanOrEqual(1);
    expect(await page.locator("#stack-head th").count()).toBe(3);
    // A number sits beside its label.
    expect(await css(page, "#stack-num", "text-align")).toBe("start");
    // Nothing is wider than the screen.
    expect(await css(page, "#stack-wrap", "overflow-x")).toBe("visible");
    expect(row.width).toBeLessThanOrEqual(390);
  });

  test("axe passes on the gallery with its tables stacked", async ({ page }) => {
    await page.goto("/p/blue/light/compact");
    const { violations } = await new AxeBuilder({ page }).analyze();
    expect(violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(" ")).join(", ")}`)).toEqual([]);
  });
});

test.describe("on a touch screen", () => {
  test.use({ hasTouch: true, isMobile: true, viewport: { width: 1100, height: 900 } });

  test("a tapped row keeps no tint", async ({ page }) => {
    await page.goto("/fixture/library/ui.css");
    expect(await page.evaluate(() => matchMedia("(hover: hover)").matches), "the page sees no hover").toBe(false);
    await page.locator("#td").tap();
    expect(await css(page, "#td", "background-color")).toBe("rgba(0, 0, 0, 0)");
  });
});
```

- [ ] **9.** Run `go test -count=1 ./gallery`. Expect `ok`.

- [ ] **10.** Run `npx playwright test e2e/tables.spec.ts --ignore-snapshots`. Expect 8 passed.

- [ ] **11.** Run `npx playwright test --ignore-snapshots`. Expect 62 passed.

- [ ] **12. Commit:** `git add gallery/fixture_library.html gallery/page.html e2e/tables.spec.ts` then `git commit -m "test: measure the table options; show a stacked table in the gallery"`.

### Task 3: Docs, version and site

**Files:** `DESIGN.md`, `CHANGELOG.md`, `AGENTS.md`, `README.md`, `package.json`, `package-lock.json`, `site/templates/components.html`

- [ ] **13. Modify `DESIGN.md`.** Write this patch to `/tmp/v031-DESIGN.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-DESIGN.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/DESIGN.md b/DESIGN.md
index 789f093..e69a07e 100644
--- a/DESIGN.md
+++ b/DESIGN.md
@@ -393,11 +393,45 @@ confirm button acts on it and closes the dialog with `command="close"`.
 <section class="ui-card">…</section>
 ```
 
-Knobs: `--ui-callout-bg`, `--ui-callout-fg`, `--ui-callout-accent`; `--ui-table-pad-y`, `--ui-table-pad-x`; `--ui-card-pad`, `--ui-card-radius`, `--ui-card-shadow`, `--ui-card-border`; `--ui-disclosure-rule`.
+Knobs: `--ui-callout-bg`, `--ui-callout-fg`, `--ui-callout-accent`; `--ui-table-pad-y`, `--ui-table-pad-x` and the table knobs below; `--ui-card-pad`, `--ui-card-radius`, `--ui-card-shadow`, `--ui-card-border`; `--ui-disclosure-rule`.
 
 A card has a soft shadow by default and none in compact density. `--ui-card-shadow: none` makes any
 card flat.
 
+**Tables** take four options, all plain markup:
+
+| Option | Markup | What it does |
+|---|---|---|
+| No frame | `data-frame="none"` on `.ui-table-wrap` | No border, corners or fill; the first and last columns sit flush with the content around them. |
+| Quiet head | `data-head="quiet"` on `.ui-table` | Header cells on no fill, in small `--ink-3` text. |
+| Sticky head | `data-sticky="head"` on `.ui-table` | Header cells stay at the top of the wrap as it scrolls. Bound the wrap with `--ui-table-max-height`. |
+| Stacked on a phone | `data-stack="phone"` on `.ui-table-wrap` | Below 800px of viewport, each row is a block and each cell a line with its label. |
+
+Table knobs, on `.ui-table`: `--ui-table-head-bg`, `--ui-table-head-fg`, `--ui-table-head-size`,
+`--ui-table-head-case` (`text-transform`), `--ui-table-head-rule` (the line under the header),
+`--ui-table-hover` (`transparent` turns the row tint off), `--ui-table-sticky-bg` (default
+`--sheet`; set it to whatever is behind the table), `--ui-table-label-width` (default 8rem). On
+`.ui-table-wrap`: `--ui-table-max-height`. A row is tinted on hover only where the pointer can hover.
+
+A stacked table labels each cell from its `data-label`. A cell without one, such as the row's name,
+spans the row, and so does a row header (`<th scope="row">`). The header row is hidden from sight
+but stays in the markup.
+
+```html
+<div class="ui-table-wrap" data-stack="phone">
+  <table class="ui-table">
+    <thead><tr><th scope="col">Email</th><th scope="col">Role</th><th scope="col" class="ui-num">Logins</th></tr></thead>
+    <tbody>
+      <tr><td>ada@example.com</td><td data-label="Role">owner</td><td data-label="Logins" class="ui-num">12</td></tr>
+    </tbody>
+  </table>
+</div>
+```
+
+Some browsers stop treating a table as a table once its rows are blocks, so a screen reader may lose
+the column names on a phone. If that matters, add the ARIA table roles to the markup: `role="table"`
+on the table, `rowgroup`, `row`, `columnheader` and `cell`.
+
 `data-marker="start"` on a disclosure puts its chevron before the summary text, pointing to the
 end when closed and down when open, for a list of folders or groups. `--ui-disclosure-rule:
 transparent` drops the line under it.
````

- [ ] **14. Modify `CHANGELOG.md`.** Write this patch to `/tmp/v031-CHANGELOG.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-CHANGELOG.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/CHANGELOG.md b/CHANGELOG.md
index e2cf3c1..ad5f9be 100644
--- a/CHANGELOG.md
+++ b/CHANGELOG.md
@@ -4,6 +4,17 @@ Breaking changes come first in each entry. Renaming or removing a token,
 class, data-attribute value, knob, partial name or argument, or React prop
 is breaking.
 
+## 0.3.1 (2026-11-11)
+
+No breaking changes. One look changes:
+
+- A table row is tinted on hover only where the pointer can hover, so a tapped row on a phone no longer keeps its tint. The tint is the new `--ui-table-hover` knob.
+
+New:
+
+- Tables: `data-frame="none"` on the wrap for no frame, with flush outer columns; `data-head="quiet"` for a header on no fill in small grey text; `data-sticky="head"` for a header that stays at the top of a scrolling wrap; `data-stack="phone"` on the wrap to stack each row into labelled lines below 800px, labels from each cell's `data-label`.
+- Table knobs: `--ui-table-head-bg`, `--ui-table-head-fg`, `--ui-table-head-size`, `--ui-table-head-case`, `--ui-table-head-rule`, `--ui-table-hover`, `--ui-table-sticky-bg`, `--ui-table-label-width`, and `--ui-table-max-height` on the wrap.
+
 ## 0.3.0 (2026-11-11)
 
 No breaking changes. Two looks change:
```

- [ ] **15. Modify `AGENTS.md`.** Write this patch to `/tmp/v031-AGENTS.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-AGENTS.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/AGENTS.md b/AGENTS.md
index 98d8393..45ea1be 100644
--- a/AGENTS.md
+++ b/AGENTS.md
@@ -27,7 +27,7 @@ Check which stack this project uses, then follow one of these.
 ### Go with `html/template`
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.3.0
+go get github.com/fuad-daoud/ai-ui@v0.3.1
 ```
 
 It needs Go 1.25, and `go get` raises the `go` line in an older `go.mod` to match.
@@ -60,13 +60,13 @@ read it too, then install:
 
 ```bash
 echo "allow-git=root" >> .npmrc
-npm install github:fuad-daoud/ai-ui#v0.3.0
+npm install github:fuad-daoud/ai-ui#v0.3.1
 ```
 
 That leaves this in `package.json`:
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.0" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.1" }
 ```
 
 ```ts
````

- [ ] **16. Modify `README.md`.** Write this patch to `/tmp/v031-README.md.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-README.md.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

````diff
diff --git a/README.md b/README.md
index 67e5332..b916342 100644
--- a/README.md
+++ b/README.md
@@ -8,7 +8,7 @@ an agent starts. The site is at https://ui.fuad-daoud.com/.
 ## Go apps
 
 ```bash
-go get github.com/fuad-daoud/ai-ui@v0.3.0
+go get github.com/fuad-daoud/ai-ui@v0.3.1
 ```
 
 It needs Go 1.25.
@@ -41,11 +41,11 @@ allows it, so allow it in a `.npmrc` (which `npm ci` and CI read too), then inst
 
 ```bash
 echo "allow-git=root" >> .npmrc
-npm install github:fuad-daoud/ai-ui#v0.3.0
+npm install github:fuad-daoud/ai-ui#v0.3.1
 ```
 
 ```json
-"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.0" }
+"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.3.1" }
 ```
 
 ```ts
````

- [ ] **17. Modify `package.json`.** Write this patch to `/tmp/v031-package.json.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-package.json.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/package.json b/package.json
index 87649d9..5da215f 100644
--- a/package.json
+++ b/package.json
@@ -1,6 +1,6 @@
 {
   "name": "ai-ui",
-  "version": "0.3.0",
+  "version": "0.3.1",
   "description": "A small UI library written to be used by coding agents: tokens, layered CSS, Go partials and React components.",
   "license": "MIT",
   "type": "module",
```

- [ ] **18. Modify `package-lock.json`.** Write this patch to `/tmp/v031-package-lock.json.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-package-lock.json.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/package-lock.json b/package-lock.json
index fb14ee6..0fe9292 100644
--- a/package-lock.json
+++ b/package-lock.json
@@ -1,12 +1,12 @@
 {
   "name": "ai-ui",
-  "version": "0.3.0",
+  "version": "0.3.1",
   "lockfileVersion": 3,
   "requires": true,
   "packages": {
     "": {
       "name": "ai-ui",
-      "version": "0.3.0",
+      "version": "0.3.1",
       "license": "MIT",
       "devDependencies": {
         "@axe-core/playwright": "^4.14.0",
```

- [ ] **19. Modify `site/templates/components.html`.** Write this patch to `/tmp/v031-site-templates-components.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v031-site-templates-components.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/site/templates/components.html b/site/templates/components.html
index cb14327..432e855 100644
--- a/site/templates/components.html
+++ b/site/templates/components.html
@@ -131,6 +131,15 @@
       </tbody>
     </table>
   </div>
+  <div class="ui-table-wrap" data-frame="none" data-stack="phone">
+    <table class="ui-table" data-head="quiet">
+      <thead><tr><th scope="col">Invoice</th><th scope="col">Customer</th><th scope="col">State</th><th scope="col" class="ui-num">Amount</th></tr></thead>
+      <tbody>
+        <tr><td>1042</td><td data-label="Customer">Acme</td><td data-label="State">{{template "ui.pill" (dict "Label" "Paid" "Variant" "good" "Dot" true)}}</td><td data-label="Amount" class="ui-num">1,040.00</td></tr>
+        <tr><td>1043</td><td data-label="Customer">Northwind</td><td data-label="State">{{template "ui.pill" (dict "Label" "Due soon" "Variant" "warning" "Dot" true)}}</td><td data-label="Amount" class="ui-num">320.00</td></tr>
+      </tbody>
+    </table>
+  </div>
   <div>
     <details class="ui-disclosure" name="demo-faq">
       <summary class="ui-disclosure-summary">Can I export it?</summary>
```

- [ ] **20.** Run `go test -count=1 . && (cd site && go test -count=1 ./...)`. Expect `ok`.

- [ ] **21.** Run `npx playwright test -c playwright.site.config.ts`. Expect 57 passed.

- [ ] **22. Commit:** `git add DESIGN.md CHANGELOG.md AGENTS.md README.md package.json package-lock.json site/templates/components.html` then `git commit -m "docs: table options in DESIGN.md; 0.3.1"`.

### Task 4: The gate

- [ ] **23.** Run the gate above. Expect it all to pass: 61 vitest tests, 62 browser tests, 57 site tests.
