# ai-ui v0.3, round 1: the options, in CSS, Go, React and the gallery

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do the steps in order.

**Precondition (STOP if unmet):** `git log --oneline -3` shows `docs: plan ai-ui v0.3` at the top, and `git status --short` is empty.

**Goal:** every option in the spec exists in `ui.css`, the Go partials and the React components, with parity fixtures, a Go test for icons, browser tests that measure each option, and gallery sections that show them.

**Architecture:** additive rules in the existing component files; options are `data-*` attributes and `--ui-*` knobs (spec §3). The partials and components gain arguments of the same names. Two looks change on purpose (spec §6): disabled buttons, and compact cards and dialogs.

**Tech Stack:** CSS layers, Go `html/template`, React 19 (TypeScript), vitest, Playwright with axe.

**Spec:** `docs/superpowers/specs/2026-11-11-v0.3-options-design.md`. Read it first.

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

- The 16 gallery baselines **will** differ from the new gallery: that is why every run passes `--ignore-snapshots`. axe still runs on every gallery page.

## Review Focus

1. **Disabled beats every variant and tone:** the outline rules use `:where()` so they weigh the same as a variant rule, and the disabled rules come after them. Pinned by `disabled changes colour, not opacity` (an outlined danger button, a ghost and a link).
2. **A select keeps its chevron room at `sm` and `lg`:** pinned by the `padding-right: 36px` check.
3. **Touch sizes only in compact:** every coarse rule sits under `[data-density="compact"]`. Pinned by the touch test and by `compact has no effect without the attribute`.
4. **An icon is the app's markup:** a `template.HTML` goes in raw, a string is escaped. Pinned by `TestIconIsTheAppsMarkup`.
5. **Parity:** Go and React render the same markup for every new argument, including unknown values, which render the default.

---

### Task 1: CSS

**Files:** `css/components/button.css`, `css/components/field.css`, `css/components/density.css`, `css/components/card.css`, `css/components/dialog.css`, `css/components/pill.css`, `css/components/disclosure.css`, `dist/ui.css`, `dist/ui-scoped.css`

- [ ] **1. Modify `css/components/button.css`.** Write this patch to `/tmp/v03-css-components-button.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-button.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/button.css b/css/components/button.css
index 15823ac..c8b5fda 100644
--- a/css/components/button.css
+++ b/css/components/button.css
@@ -48,7 +48,40 @@
   --ui-btn-bg-hover: var(--ink);
   --ui-btn-fg-hover: var(--paper);
 }
+/* An outlined primary or danger button. :where() keeps these at a variant's weight,
+   so the disabled rules below still win. Secondary and ghost don't change. */
+.ui-btn[data-tone="outline"]:where(:not([data-variant]), [data-variant="danger"]) {
+  --ui-btn-bg: var(--sheet);
+  --ui-btn-border: var(--rule);
+  --ui-btn-bg-hover: var(--band);
+}
+.ui-btn[data-tone="outline"]:where(:not([data-variant])) { --ui-btn-fg: var(--primary); --ui-btn-fg-hover: var(--primary); }
+.ui-btn[data-tone="outline"]:where([data-variant="danger"]) { --ui-btn-fg: var(--bad); --ui-btn-fg-hover: var(--bad); }
 .ui-btn[data-size="sm"] { --ui-btn-height: 2rem; --ui-btn-pad: var(--space-3); --ui-btn-font-size: var(--text-sm); }
 .ui-btn[data-size="lg"] { --ui-btn-height: 3rem; --ui-btn-pad: var(--space-5); --ui-btn-font-size: var(--text-lg); }
+.ui-btn > svg { flex: none; inline-size: 1em; block-size: 1em; }
+/* An icon-only button: as wide as it is tall. */
+.ui-btn[data-shape="square"] { inline-size: var(--ui-btn-height); padding: 0; }
+.ui-btn[data-width="full"] {
+  inline-size: 100%;
+  padding-block: var(--space-1);
+  white-space: normal;
+  text-wrap: balance;
+}
+/* Disabled changes colour, not opacity, so the label stays readable. */
 .ui-btn:disabled,
-.ui-btn[aria-disabled="true"] { opacity: 0.5; cursor: not-allowed; pointer-events: none; }
+.ui-btn[aria-disabled="true"] {
+  --ui-btn-bg: var(--band);
+  --ui-btn-fg: var(--ink-3);
+  --ui-btn-border: var(--rule);
+  --ui-btn-bg-hover: var(--band);
+  --ui-btn-fg-hover: var(--ink-3);
+  cursor: not-allowed;
+}
+.ui-btn[data-variant="ghost"]:is(:disabled, [aria-disabled="true"]) {
+  --ui-btn-bg: transparent;
+  --ui-btn-bg-hover: transparent;
+  --ui-btn-border: transparent;
+}
+/* A disabled link has no href, but must not take clicks either. */
+.ui-btn[aria-disabled="true"] { pointer-events: none; }
```

- [ ] **2. Modify `css/components/field.css`.** Write this patch to `/tmp/v03-css-components-field.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-field.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/field.css b/css/components/field.css
index 57a54a9..9ab3e66 100644
--- a/css/components/field.css
+++ b/css/components/field.css
@@ -11,8 +11,9 @@
   --ui-input-border: var(--rule-field);
   --ui-input-radius: var(--radius-md);
   --ui-input-height: 2.5rem;
+  --ui-input-width: 100%;
 
-  inline-size: 100%;
+  inline-size: var(--ui-input-width);
   min-block-size: var(--ui-input-height);
   padding: var(--space-2) var(--space-3);
   background-color: var(--ui-input-bg);
@@ -36,6 +37,31 @@
   background-size: 0.3rem 0.3rem, 0.3rem 0.3rem;
   background-repeat: no-repeat;
 }
+/* Sizes match the button of the same size. A select keeps its end padding for the chevron. */
+:is(.ui-input, .ui-select, .ui-textarea)[data-size="sm"] { --ui-input-height: 2rem; font-size: var(--text-sm); }
+:is(.ui-input, .ui-textarea)[data-size="sm"] { padding-inline: var(--space-2); }
+.ui-textarea[data-size="sm"] { padding-block: var(--space-1); }
+.ui-select[data-size="sm"] { padding-inline-start: var(--space-2); }
+:is(.ui-input, .ui-select, .ui-textarea)[data-size="lg"] { --ui-input-height: 3rem; font-size: 1rem; }
+:is(.ui-input, .ui-textarea)[data-size="lg"] { padding-inline: var(--space-4); }
+.ui-select[data-size="lg"] { padding-inline-start: var(--space-4); }
+
+/* A leading icon inside an input. */
+.ui-input-group { --ui-input-icon-inset: var(--space-3); position: relative; }
+.ui-input-group:has(> .ui-input[data-size="sm"]) { --ui-input-icon-inset: var(--space-2); }
+.ui-input-group:has(> .ui-input[data-size="lg"]) { --ui-input-icon-inset: var(--space-4); }
+.ui-input-group > :is(.ui-input-icon, svg) {
+  position: absolute;
+  inset-inline-start: var(--ui-input-icon-inset);
+  inset-block-start: 50%;
+  translate: 0 -50%;
+  inline-size: 1rem;
+  block-size: 1rem;
+  color: var(--ink-3);
+  pointer-events: none;
+}
+.ui-input-group > .ui-input { padding-inline-start: calc(var(--ui-input-icon-inset) + 1rem + var(--space-2)); }
+
 .ui-input::placeholder,
 .ui-textarea::placeholder { color: var(--ink-3); opacity: 1; }
 .ui-input:hover,
@@ -52,4 +78,4 @@
 .ui-textarea:disabled { background-color: var(--band); color: var(--ink-3); cursor: not-allowed; }
 
 .ui-form-row { display: flex; flex-wrap: wrap; align-items: end; gap: var(--space-4); }
-.ui-form-row > * { flex: 1 1 14rem; }
+.ui-form-row > * { flex: var(--ui-row-grow, 1) 1 var(--ui-row-basis, 14rem); }
```

- [ ] **3. Modify `css/components/density.css`.** Write this patch to `/tmp/v03-css-components-density.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-density.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/density.css b/css/components/density.css
index 52ef7b0..f0816c4 100644
--- a/css/components/density.css
+++ b/css/components/density.css
@@ -14,6 +14,20 @@
 [data-density="compact"] .ui-btn[data-size="sm"] { --ui-btn-height: 1.75rem; --ui-btn-pad: 0.625rem; }
 [data-density="compact"] .ui-btn[data-size="lg"] { --ui-btn-height: 2.75rem; --ui-btn-pad: 1.25rem; --ui-btn-font-size: 1rem; }
 [data-density="compact"] :is(.ui-input, .ui-select) { --ui-input-height: 2rem; }
+[data-density="compact"] :is(.ui-input, .ui-select, .ui-textarea)[data-size="sm"] { --ui-input-height: 1.75rem; }
+[data-density="compact"] :is(.ui-input, .ui-select, .ui-textarea)[data-size="lg"] { --ui-input-height: 2.75rem; }
 [data-density="compact"] .ui-field { --ui-field-gap: var(--space-1); }
 [data-density="compact"] .ui-pill { border-radius: var(--radius-sm); }
 [data-density="compact"] .ui-table { font-size: var(--text-md); }
+/* Compact cards and dialogs are flat: a hairline, no shadow. */
+[data-density="compact"] .ui-card { --ui-card-shadow: none; }
+[data-density="compact"] .ui-dialog { --ui-dialog-shadow: none; }
+/* Finger-sized controls on a touch screen. */
+@media (pointer: coarse) {
+  [data-density="compact"] .ui-btn { --ui-btn-height: 2.5rem; }
+  [data-density="compact"] .ui-btn[data-size="sm"] { --ui-btn-height: 2.125rem; }
+  [data-density="compact"] .ui-btn[data-size="lg"] { --ui-btn-height: 2.75rem; }
+  [data-density="compact"] :is(.ui-input, .ui-select) { --ui-input-height: 2.5rem; }
+  [data-density="compact"] :is(.ui-input, .ui-select, .ui-textarea)[data-size="sm"] { --ui-input-height: 2.125rem; }
+  [data-density="compact"] :is(.ui-input, .ui-select, .ui-textarea)[data-size="lg"] { --ui-input-height: 2.75rem; }
+}
```

- [ ] **4. Modify `css/components/card.css`.** Write this patch to `/tmp/v03-css-components-card.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-card.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/card.css b/css/components/card.css
index f58ad4e..e3af8a8 100644
--- a/css/components/card.css
+++ b/css/components/card.css
@@ -1,9 +1,11 @@
 .ui-card {
   --ui-card-pad: var(--space-5);
   --ui-card-radius: var(--radius-lg);
+  --ui-card-shadow: 0 1px 2px var(--shadow);
+  --ui-card-border: var(--rule-soft);
   padding: var(--ui-card-pad);
-  border: 1px solid var(--rule-soft);
+  border: 1px solid var(--ui-card-border);
   border-radius: var(--ui-card-radius);
   background: var(--sheet);
-  box-shadow: 0 1px 2px var(--shadow);
+  box-shadow: var(--ui-card-shadow);
 }
```

- [ ] **5. Modify `css/components/dialog.css`.** Write this patch to `/tmp/v03-css-components-dialog.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-dialog.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/dialog.css b/css/components/dialog.css
index bdbcd04..a9648ea 100644
--- a/css/components/dialog.css
+++ b/css/components/dialog.css
@@ -3,11 +3,14 @@
   --ui-dialog-bg: var(--sheet);
   --ui-dialog-radius: var(--radius-lg);
   --ui-dialog-pad: var(--space-5);
+  --ui-dialog-shadow: 0 1px 2px var(--shadow), 0 16px 48px var(--shadow);
+  --ui-dialog-border: var(--rule-soft);
+  --ui-dialog-eyebrow-fg: var(--ink-3);
   inline-size: min(var(--ui-dialog-width), calc(100vw - 2 * var(--space-4)));
   max-block-size: calc(100dvh - 2 * var(--space-4));
   margin: auto; /* the reset zeroes margins; a modal is centred by margin: auto */
   padding: 0;
-  border: 1px solid var(--rule-soft);
+  border: 1px solid var(--ui-dialog-border);
   border-radius: var(--ui-dialog-radius);
   background: var(--ui-dialog-bg);
   color: var(--ink);
@@ -15,7 +18,7 @@
      often an end-aligned, unwrapped table cell. */
   text-align: start;
   white-space: normal;
-  box-shadow: 0 1px 2px var(--shadow), 0 16px 48px var(--shadow);
+  box-shadow: var(--ui-dialog-shadow);
 }
 .ui-dialog[open] { display: flex; flex-direction: column; }
 .ui-dialog::backdrop { background: var(--scrim); }
@@ -26,6 +29,16 @@
   gap: var(--space-3);
   padding: var(--ui-dialog-pad) var(--ui-dialog-pad) 0;
 }
+/* A small line above the head, such as "Can't be undone". */
+.ui-dialog-eyebrow {
+  padding: var(--ui-dialog-pad) var(--ui-dialog-pad) 0;
+  color: var(--ui-dialog-eyebrow-fg);
+  font-size: var(--text-xs);
+  font-weight: 600;
+  letter-spacing: 0.06em;
+  text-transform: uppercase;
+}
+.ui-dialog-eyebrow + .ui-dialog-head { padding-block-start: var(--space-1); }
 .ui-dialog-title { font-size: var(--text-lg); }
 .ui-dialog-x {
   display: grid;
```

- [ ] **6. Modify `css/components/pill.css`.** Write this patch to `/tmp/v03-css-components-pill.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-pill.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/pill.css b/css/components/pill.css
index 55780e5..6218979 100644
--- a/css/components/pill.css
+++ b/css/components/pill.css
@@ -19,3 +19,14 @@
 .ui-pill[data-variant="good"] { --ui-pill-bg: var(--good-bg); --ui-pill-fg: var(--good); --ui-pill-border: transparent; }
 .ui-pill[data-variant="warning"] { --ui-pill-bg: var(--warning-bg); --ui-pill-fg: var(--warning); --ui-pill-border: transparent; }
 .ui-pill[data-variant="bad"] { --ui-pill-bg: var(--bad-bg); --ui-pill-fg: var(--bad); --ui-pill-border: transparent; }
+/* After the variants: an outline pill drops its fill, and its border takes the text colour. */
+.ui-pill[data-tone="outline"] { --ui-pill-bg: transparent; --ui-pill-border: var(--ui-pill-fg); }
+.ui-pill[data-tone="outline"]:not([data-variant]) { --ui-pill-border: var(--rule-field); }
+.ui-pill[data-dot]::before {
+  content: "";
+  flex: none;
+  inline-size: 0.375rem;
+  block-size: 0.375rem;
+  border-radius: 50%;
+  background: currentColor;
+}
```

- [ ] **7. Modify `css/components/disclosure.css`.** Write this patch to `/tmp/v03-css-components-disclosure.css.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-css-components-disclosure.css.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/css/components/disclosure.css b/css/components/disclosure.css
index 7226ce8..c4c5852 100644
--- a/css/components/disclosure.css
+++ b/css/components/disclosure.css
@@ -1,4 +1,4 @@
-.ui-disclosure { border-block-end: 1px solid var(--rule-soft); }
+.ui-disclosure { --ui-disclosure-rule: var(--rule-soft); border-block-end: 1px solid var(--ui-disclosure-rule); }
 .ui-disclosure-summary {
   display: flex;
   align-items: center;
@@ -22,3 +22,18 @@
 }
 .ui-disclosure[open] > .ui-disclosure-summary::after { transform: rotate(-135deg); }
 .ui-disclosure-body { padding-block-end: var(--space-4); color: var(--ink-2); }
+/* The chevron before the text: it points to the inline end, and down when open. */
+.ui-disclosure[data-marker="start"] > .ui-disclosure-summary { justify-content: flex-start; }
+.ui-disclosure[data-marker="start"] > .ui-disclosure-summary::after { content: none; }
+.ui-disclosure[data-marker="start"] > .ui-disclosure-summary::before {
+  content: "";
+  flex: none;
+  inline-size: 0.5rem;
+  block-size: 0.5rem;
+  border-inline-end: 2px solid var(--ink-3);
+  border-block-end: 2px solid var(--ink-3);
+  transform: rotate(-45deg);
+  transition: transform var(--dur) var(--ease);
+}
+.ui-disclosure[data-marker="start"] > .ui-disclosure-summary:dir(rtl)::before { transform: rotate(135deg); }
+.ui-disclosure[data-marker="start"][open] > .ui-disclosure-summary::before { transform: rotate(45deg); }
```

- [ ] **8.** Run `go run ./cmd/build`. Expect `wrote dist/ui.css (23822 bytes)` and `wrote dist/ui-scoped.css (27256 bytes)`.

- [ ] **9.** Run `go test -count=1 .`. Expect `ok`.

- [ ] **10. Commit:** `git add css/components/button.css css/components/field.css css/components/density.css css/components/card.css css/components/dialog.css css/components/pill.css css/components/disclosure.css dist/ui.css dist/ui-scoped.css` then `git commit -m "feat(css): options on buttons, fields, cards, dialogs, pills and disclosures"`.

### Task 2: Go partials and their tests

**Files:** `templates/ui-button.tmpl`, `templates/ui-field.tmpl`, `templates/ui-pill.tmpl`, `templates_test.go`, `classes_test.go`, `testdata/parity/*`

- [ ] **11. Modify `templates/ui-button.tmpl`.** Write this patch to `/tmp/v03-templates-ui-button.tmpl.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-templates-ui-button.tmpl.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/templates/ui-button.tmpl b/templates/ui-button.tmpl
index bfb6bd0..c4eccb0 100644
--- a/templates/ui-button.tmpl
+++ b/templates/ui-button.tmpl
@@ -2,9 +2,13 @@
 {{- $label := ui_required "ui.button" "Label" (index . "Label") -}}
 {{- $variant := ui_variant (index . "Variant") "secondary" "ghost" "danger" -}}
 {{- $size := ui_variant (index . "Size") "sm" "lg" -}}
+{{- $tone := ui_variant (index . "Tone") "outline" -}}
+{{- $shape := ui_variant (index . "Shape") "square" -}}
+{{- $width := ui_variant (index . "Width") "full" -}}
+{{- if $shape}}{{$_ := ui_required "ui.button" "Icon" (index . "Icon")}}{{end -}}
 {{- if (index . "Href") -}}
-<a class="ui-btn"{{if (index . "Disabled")}} aria-disabled="true"{{else}} href="{{(index . "Href")}}"{{end}}{{with $variant}} data-variant="{{.}}"{{end}}{{with $size}} data-size="{{.}}"{{end}}>{{$label}}</a>
+<a class="ui-btn"{{if (index . "Disabled")}} aria-disabled="true"{{else}} href="{{(index . "Href")}}"{{end}}{{with $variant}} data-variant="{{.}}"{{end}}{{with $size}} data-size="{{.}}"{{end}}{{with $tone}} data-tone="{{.}}"{{end}}{{with $shape}} data-shape="{{.}}" aria-label="{{$label}}"{{end}}{{with $width}} data-width="{{.}}"{{end}}>{{with (index . "Icon")}}{{.}}{{end}}{{if not $shape}}{{$label}}{{end}}</a>
 {{- else -}}
-<button class="ui-btn" type="{{ui_or (ui_variant (index . "Type") "submit" "reset") "button"}}"{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with ui_str (index . "Value")}} value="{{.}}"{{end}}{{with ui_str (index . "Command")}} command="{{.}}"{{end}}{{with ui_str (index . "CommandFor")}} commandfor="{{.}}"{{end}}{{with $variant}} data-variant="{{.}}"{{end}}{{with $size}} data-size="{{.}}"{{end}}{{if (index . "Disabled")}} disabled{{end}}>{{$label}}</button>
+<button class="ui-btn" type="{{ui_or (ui_variant (index . "Type") "submit" "reset") "button"}}"{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with ui_str (index . "Value")}} value="{{.}}"{{end}}{{with ui_str (index . "Command")}} command="{{.}}"{{end}}{{with ui_str (index . "CommandFor")}} commandfor="{{.}}"{{end}}{{with $variant}} data-variant="{{.}}"{{end}}{{with $size}} data-size="{{.}}"{{end}}{{with $tone}} data-tone="{{.}}"{{end}}{{with $shape}} data-shape="{{.}}" aria-label="{{$label}}"{{end}}{{with $width}} data-width="{{.}}"{{end}}{{if (index . "Disabled")}} disabled{{end}}>{{with (index . "Icon")}}{{.}}{{end}}{{if not $shape}}{{$label}}{{end}}</button>
 {{- end -}}
 {{- end}}
```

- [ ] **12. Modify `templates/ui-field.tmpl`.** Write this patch to `/tmp/v03-templates-ui-field.tmpl.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-templates-ui-field.tmpl.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/templates/ui-field.tmpl b/templates/ui-field.tmpl
index 2128007..da202b1 100644
--- a/templates/ui-field.tmpl
+++ b/templates/ui-field.tmpl
@@ -3,16 +3,24 @@
 {{- $label := ui_required "ui.field" "Label" (index . "Label") -}}
 {{- $desc := ui_describedby $id (index . "Hint") (index . "Error") -}}
 {{- $control := ui_or (ui_variant (index . "Control") "textarea" "select") "input" -}}
+{{- $size := ui_variant (index . "Size") "sm" "lg" -}}
 <div class="ui-field">
 <label class="ui-field-label" for="{{$id}}">{{$label}}</label>
 {{- if eq $control "textarea"}}
-<textarea class="ui-textarea" id="{{$id}}"{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with ui_str (index . "Rows")}} rows="{{.}}"{{end}}{{with ui_str (index . "Placeholder")}} placeholder="{{.}}"{{end}}{{with $desc}} aria-describedby="{{.}}"{{end}}{{if (index . "Error")}} aria-invalid="true"{{end}}{{if (index . "Required")}} required{{end}}{{if (index . "Disabled")}} disabled{{end}}>{{ui_str (index . "Value")}}</textarea>
+<textarea class="ui-textarea" id="{{$id}}"{{with $size}} data-size="{{.}}"{{end}}{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with ui_str (index . "Rows")}} rows="{{.}}"{{end}}{{with ui_str (index . "Placeholder")}} placeholder="{{.}}"{{end}}{{with $desc}} aria-describedby="{{.}}"{{end}}{{if (index . "Error")}} aria-invalid="true"{{end}}{{if (index . "Required")}} required{{end}}{{if (index . "Disabled")}} disabled{{end}}>{{ui_str (index . "Value")}}</textarea>
 {{- else if eq $control "select"}}
-<select class="ui-select" id="{{$id}}"{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with $desc}} aria-describedby="{{.}}"{{end}}{{if (index . "Error")}} aria-invalid="true"{{end}}{{if (index . "Required")}} required{{end}}{{if (index . "Disabled")}} disabled{{end}}>
+<select class="ui-select" id="{{$id}}"{{with $size}} data-size="{{.}}"{{end}}{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with $desc}} aria-describedby="{{.}}"{{end}}{{if (index . "Error")}} aria-invalid="true"{{end}}{{if (index . "Required")}} required{{end}}{{if (index . "Disabled")}} disabled{{end}}>
 {{- range (index . "Options")}}<option value="{{ui_str .Value}}"{{if ui_selected (index $ "Value") .Value}} selected{{end}}>{{ui_str .Label}}</option>{{end -}}
 </select>
 {{- else}}
-<input class="ui-input" id="{{$id}}" type="{{ui_or (ui_variant (index . "Type") "email" "password" "number" "search" "tel" "url" "date") "text"}}"{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with ui_str (index . "Value")}} value="{{.}}"{{end}}{{with ui_str (index . "Placeholder")}} placeholder="{{.}}"{{end}}{{with ui_str (index . "Autocomplete")}} autocomplete="{{.}}"{{end}}{{with $desc}} aria-describedby="{{.}}"{{end}}{{if (index . "Error")}} aria-invalid="true"{{end}}{{if (index . "Required")}} required{{end}}{{if (index . "Disabled")}} disabled{{end}}>
+{{- with (index $ "Icon")}}
+<div class="ui-input-group">
+{{.}}
+{{- end}}
+<input class="ui-input" id="{{$id}}"{{with $size}} data-size="{{.}}"{{end}} type="{{ui_or (ui_variant (index . "Type") "email" "password" "number" "search" "tel" "url" "date") "text"}}"{{with ui_str (index . "Name")}} name="{{.}}"{{end}}{{with ui_str (index . "Value")}} value="{{.}}"{{end}}{{with ui_str (index . "Placeholder")}} placeholder="{{.}}"{{end}}{{with ui_str (index . "Autocomplete")}} autocomplete="{{.}}"{{end}}{{with $desc}} aria-describedby="{{.}}"{{end}}{{if (index . "Error")}} aria-invalid="true"{{end}}{{if (index . "Required")}} required{{end}}{{if (index . "Disabled")}} disabled{{end}}>
+{{- if (index . "Icon")}}
+</div>
+{{- end}}
 {{- end}}
 {{- with ui_str (index . "Hint")}}
 <p class="ui-field-hint" id="{{$id}}-hint">{{.}}</p>
```

- [ ] **13. Modify `templates/ui-pill.tmpl`.** Write this patch to `/tmp/v03-templates-ui-pill.tmpl.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-templates-ui-pill.tmpl.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/templates/ui-pill.tmpl b/templates/ui-pill.tmpl
index dd0ae71..19e6a18 100644
--- a/templates/ui-pill.tmpl
+++ b/templates/ui-pill.tmpl
@@ -1,3 +1,3 @@
 {{define "ui.pill" -}}
-<span class="ui-pill"{{with ui_variant (index . "Variant") "info" "good" "warning" "bad"}} data-variant="{{.}}"{{end}}>{{ui_required "ui.pill" "Label" (index . "Label")}}</span>
+<span class="ui-pill"{{with ui_variant (index . "Variant") "info" "good" "warning" "bad"}} data-variant="{{.}}"{{end}}{{with ui_variant (index . "Tone") "outline"}} data-tone="{{.}}"{{end}}{{if (index . "Dot")}} data-dot{{end}}>{{ui_required "ui.pill" "Label" (index . "Label")}}</span>
 {{- end}}
```

- [ ] **14. Modify `templates_test.go`.** Write this patch to `/tmp/v03-templates_test.go.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-templates_test.go.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/templates_test.go b/templates_test.go
index 0b44c55..4ee1112 100644
--- a/templates_test.go
+++ b/templates_test.go
@@ -69,3 +69,25 @@ func TestHrefFunc(t *testing.T) {
 		t.Errorf("err = %v, want unknown asset", err)
 	}
 }
+
+// An icon is markup the app wrote, so a template.HTML value goes in as it is,
+// and a plain string is escaped like any other argument.
+func TestIconIsTheAppsMarkup(t *testing.T) {
+	svg := template.HTML(`<svg class="ui-input-icon" aria-hidden="true"></svg>`)
+	for _, c := range []struct {
+		name, src, want string
+		data            any
+	}{
+		{"button", `{{template "ui.button" (dict "Label" "Theme" "Shape" "square" "Icon" .)}}`, `data-shape="square" aria-label="Theme">` + string(svg) + `</button>`, svg},
+		{"field", `{{template "ui.field" (dict "Label" "Search" "Name" "q" "Icon" .)}}`, `<div class="ui-input-group">` + "\n" + string(svg) + "\n" + `<input class="ui-input" id="q"`, svg},
+		{"escaped", `{{template "ui.button" (dict "Label" "Go" "Icon" .)}}`, `>&lt;b&gt;Go</button>`, "<b>"},
+	} {
+		got, err := render(t, c.src, c.data)
+		if err != nil {
+			t.Fatalf("%s: %v", c.name, err)
+		}
+		if !strings.Contains(got, c.want) {
+			t.Errorf("%s: got %s\nwant it to contain %s", c.name, got, c.want)
+		}
+	}
+}
```

- [ ] **15. Modify `classes_test.go`.** Write this patch to `/tmp/v03-classes_test.go.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-classes_test.go.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/classes_test.go b/classes_test.go
index 9dda907..fd1d66f 100644
--- a/classes_test.go
+++ b/classes_test.go
@@ -7,15 +7,15 @@ import (
 	"testing"
 )
 
-// documentedClasses is every class DESIGN.md documents for v0.1. Renaming
+// documentedClasses is every class DESIGN.md documents. Renaming
 // or removing one is a breaking change (spec §9).
 var documentedClasses = []string{
 	"ui-btn",
 	"ui-field", "ui-field-label", "ui-field-hint", "ui-field-error",
-	"ui-input", "ui-select", "ui-textarea", "ui-form-row",
+	"ui-input", "ui-select", "ui-textarea", "ui-form-row", "ui-input-group", "ui-input-icon",
 	"ui-check", "ui-check-input",
 	"ui-pill",
-	"ui-dialog", "ui-dialog-head", "ui-dialog-title", "ui-dialog-x", "ui-dialog-body", "ui-dialog-actions",
+	"ui-dialog", "ui-dialog-eyebrow", "ui-dialog-head", "ui-dialog-title", "ui-dialog-x", "ui-dialog-body", "ui-dialog-actions",
 	"ui-callout", "ui-callout-title",
 	"ui-table-wrap", "ui-table",
 	"ui-disclosure", "ui-disclosure-summary", "ui-disclosure-body",
```

- [ ] **16. Create `testdata/parity/button-icon.html`** with exactly this content:

```html
<button class="ui-btn" data-variant="secondary" type="button">↓Download</button>
```

- [ ] **17. Create `testdata/parity/button-icon.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Download", "Variant": "secondary", "Icon": "↓" } },
  "react": { "component": "Button", "props": { "variant": "secondary", "icon": "↓", "children": "Download" } }
}
```

- [ ] **18. Create `testdata/parity/button-link-square.html`** with exactly this content:

```html
<a aria-label="Download PDF" class="ui-btn" data-shape="square" href="/resume.pdf">↓</a>
```

- [ ] **19. Create `testdata/parity/button-link-square.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Download PDF", "Href": "/resume.pdf", "Shape": "square", "Icon": "↓" } },
  "react": { "component": "ButtonLink", "props": { "href": "/resume.pdf", "shape": "square", "icon": "↓", "children": "Download PDF" } }
}
```

- [ ] **20. Create `testdata/parity/button-square-no-icon.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Theme", "Shape": "square" } },
  "error": "ui.button: Icon is required"
}
```

- [ ] **21. Create `testdata/parity/button-square-no-text.json`** with exactly this content:

```json
{
  "react": { "component": "Button", "props": { "shape": "square", "icon": "◐" } },
  "error": "ui.button: a square button's label must be text"
}
```

- [ ] **22. Create `testdata/parity/button-square.html`** with exactly this content:

```html
<button aria-label="Theme" class="ui-btn" data-shape="square" data-variant="secondary" type="button">◐</button>
```

- [ ] **23. Create `testdata/parity/button-square.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Theme", "Variant": "secondary", "Shape": "square", "Icon": "◐" } },
  "react": { "component": "Button", "props": { "variant": "secondary", "shape": "square", "icon": "◐", "children": "Theme" } }
}
```

- [ ] **24. Create `testdata/parity/button-tone-outline.html`** with exactly this content:

```html
<button class="ui-btn" data-tone="outline" data-variant="danger" type="button">Delete</button>
```

- [ ] **25. Create `testdata/parity/button-tone-outline.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Delete", "Variant": "danger", "Tone": "outline" } },
  "react": { "component": "Button", "props": { "variant": "danger", "tone": "outline", "children": "Delete" } }
}
```

- [ ] **26. Create `testdata/parity/button-unknown-options.html`** with exactly this content:

```html
<button class="ui-btn" type="button">Save</button>
```

- [ ] **27. Create `testdata/parity/button-unknown-options.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Save", "Tone": "bold", "Shape": "round", "Width": "half" } },
  "react": { "component": "Button", "props": { "tone": "bold", "shape": "round", "width": "half", "children": "Save" } }
}
```

- [ ] **28. Create `testdata/parity/button-width-full.html`** with exactly this content:

```html
<button class="ui-btn" data-width="full" type="button">Continue to the next step</button>
```

- [ ] **29. Create `testdata/parity/button-width-full.json`** with exactly this content:

```json
{
  "go": { "template": "ui.button", "args": { "Label": "Continue to the next step", "Width": "full" } },
  "react": { "component": "Button", "props": { "width": "full", "children": "Continue to the next step" } }
}
```

- [ ] **30. Create `testdata/parity/dialog-eyebrow.html`** with exactly this content:

```html
<dialog class="ui-dialog" id="del" aria-labelledby="del-title">
  <p class="ui-dialog-eyebrow">Can't be undone</p>
  <header class="ui-dialog-head">
    <h2 class="ui-dialog-title" id="del-title">Delete grace@example.com?</h2>
  </header>
  <div class="ui-dialog-body">This disables the account.</div>
</dialog>
```

- [ ] **31. Create `testdata/parity/dialog-eyebrow.json`** with exactly this content:

```json
{
  "react": { "component": "Dialog", "props": { "id": "del", "title": "Delete grace@example.com?", "eyebrow": "Can't be undone", "closeButton": false, "children": "This disables the account." } }
}
```

- [ ] **32. Create `testdata/parity/field-icon.html`** with exactly this content:

```html
<div class="ui-field"><label class="ui-field-label" for="q">Search</label><div class="ui-input-group">⌕<input aria-describedby="q-hint" class="ui-input" id="q" name="q" type="search"></div><p class="ui-field-hint" id="q-hint">Merchant or note.</p></div>
```

- [ ] **33. Create `testdata/parity/field-icon.json`** with exactly this content:

```json
{
  "go": { "template": "ui.field", "args": { "Label": "Search", "Name": "q", "Type": "search", "Icon": "⌕", "Hint": "Merchant or note." } },
  "react": { "component": "Field", "props": { "label": "Search", "name": "q", "type": "search", "icon": "⌕", "hint": "Merchant or note." } }
}
```

- [ ] **34. Create `testdata/parity/field-select-lg.html`** with exactly this content:

```html
<div class="ui-field"><label class="ui-field-label" for="plan">Plan</label><select class="ui-select" data-size="lg" id="plan" name="plan"><option selected value="m">Monthly</option><option value="y">Yearly</option></select></div>
```

- [ ] **35. Create `testdata/parity/field-select-lg.json`** with exactly this content:

```json
{
  "go": { "template": "ui.field", "args": { "Label": "Plan", "Name": "plan", "Control": "select", "Size": "lg", "Value": "m", "Options": [{ "Value": "m", "Label": "Monthly" }, { "Value": "y", "Label": "Yearly" }] } },
  "react": { "component": "Field", "props": { "label": "Plan", "name": "plan", "control": "select", "size": "lg", "defaultValue": "m", "options": [{ "value": "m", "label": "Monthly" }, { "value": "y", "label": "Yearly" }] } }
}
```

- [ ] **36. Create `testdata/parity/field-size-sm.html`** with exactly this content:

```html
<div class="ui-field"><label class="ui-field-label" for="port">Port</label><input class="ui-input" data-size="sm" id="port" name="port" type="number"></div>
```

- [ ] **37. Create `testdata/parity/field-size-sm.json`** with exactly this content:

```json
{
  "go": { "template": "ui.field", "args": { "Label": "Port", "Name": "port", "Type": "number", "Size": "sm" } },
  "react": { "component": "Field", "props": { "label": "Port", "name": "port", "type": "number", "size": "sm" } }
}
```

- [ ] **38. Create `testdata/parity/field-textarea-sm.html`** with exactly this content:

```html
<div class="ui-field"><label class="ui-field-label" for="note">Note</label><textarea class="ui-textarea" data-size="sm" id="note" name="note"></textarea></div>
```

- [ ] **39. Create `testdata/parity/field-textarea-sm.json`** with exactly this content:

```json
{
  "go": { "template": "ui.field", "args": { "Label": "Note", "Name": "note", "Control": "textarea", "Size": "sm" } },
  "react": { "component": "Field", "props": { "label": "Note", "name": "note", "control": "textarea", "size": "sm" } }
}
```

- [ ] **40. Create `testdata/parity/field-unknown-size.html`** with exactly this content:

```html
<div class="ui-field"><label class="ui-field-label" for="n">Name</label><input class="ui-input" id="n" name="n" type="text"></div>
```

- [ ] **41. Create `testdata/parity/field-unknown-size.json`** with exactly this content:

```json
{
  "go": { "template": "ui.field", "args": { "Label": "Name", "Name": "n", "Size": "xl" } },
  "react": { "component": "Field", "props": { "label": "Name", "name": "n", "size": "xl" } }
}
```

- [ ] **42. Create `testdata/parity/pill-dot.html`** with exactly this content:

```html
<span class="ui-pill" data-dot data-variant="bad">1 violation</span>
```

- [ ] **43. Create `testdata/parity/pill-dot.json`** with exactly this content:

```json
{
  "go": { "template": "ui.pill", "args": { "Label": "1 violation", "Variant": "bad", "Dot": true } },
  "react": { "component": "Pill", "props": { "variant": "bad", "dot": true, "children": "1 violation" } }
}
```

- [ ] **44. Create `testdata/parity/pill-outline.html`** with exactly this content:

```html
<span class="ui-pill" data-tone="outline">meta</span>
```

- [ ] **45. Create `testdata/parity/pill-outline.json`** with exactly this content:

```json
{
  "go": { "template": "ui.pill", "args": { "Label": "meta", "Tone": "outline" } },
  "react": { "component": "Pill", "props": { "tone": "outline", "children": "meta" } }
}
```

- [ ] **46. Create `testdata/parity/pill-unknown-tone.html`** with exactly this content:

```html
<span class="ui-pill">Draft</span>
```

- [ ] **47. Create `testdata/parity/pill-unknown-tone.json`** with exactly this content:

```json
{
  "go": { "template": "ui.pill", "args": { "Label": "Draft", "Tone": "dashed", "Dot": false } },
  "react": { "component": "Pill", "props": { "tone": "dashed", "dot": false, "children": "Draft" } }
}
```

- [ ] **48.** Run `go test -count=1 -run 'TestParity|TestIcon|TestDocumented' .`. Expect `ok`. The new Go goldens are given above; never run `-update`.

- [ ] **49.** Run `test -z "$(gofmt -l .)" && go vet ./...`. Expect no output.

- [ ] **50. Commit:** `git add templates/ui-button.tmpl templates/ui-field.tmpl templates/ui-pill.tmpl templates_test.go classes_test.go testdata/parity/button-icon.html testdata/parity/button-icon.json testdata/parity/button-link-square.html testdata/parity/button-link-square.json testdata/parity/button-square-no-icon.json testdata/parity/button-square-no-text.json testdata/parity/button-square.html testdata/parity/button-square.json testdata/parity/button-tone-outline.html testdata/parity/button-tone-outline.json testdata/parity/button-unknown-options.html testdata/parity/button-unknown-options.json testdata/parity/button-width-full.html testdata/parity/button-width-full.json testdata/parity/dialog-eyebrow.html testdata/parity/dialog-eyebrow.json testdata/parity/field-icon.html testdata/parity/field-icon.json testdata/parity/field-select-lg.html testdata/parity/field-select-lg.json testdata/parity/field-size-sm.html testdata/parity/field-size-sm.json testdata/parity/field-textarea-sm.html testdata/parity/field-textarea-sm.json testdata/parity/field-unknown-size.html testdata/parity/field-unknown-size.json testdata/parity/pill-dot.html testdata/parity/pill-dot.json testdata/parity/pill-outline.html testdata/parity/pill-outline.json testdata/parity/pill-unknown-tone.html testdata/parity/pill-unknown-tone.json` then `git commit -m "feat(go): Tone, Shape, Width and Icon on buttons; Size and Icon on fields; Tone and Dot on pills"`.

### Task 3: React components

**Files:** `react/button.tsx`, `react/field.tsx`, `react/pill.tsx`, `react/dialog.tsx`, `react/index.ts`

- [ ] **51. Modify `react/button.tsx`.** Write this patch to `/tmp/v03-react-button.tsx.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-react-button.tsx.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/react/button.tsx b/react/button.tsx
index 3576c7f..91b70d2 100644
--- a/react/button.tsx
+++ b/react/button.tsx
@@ -1,4 +1,4 @@
-import type { ComponentPropsWithRef } from "react";
+import type { ComponentPropsWithRef, ReactNode } from "react";
 import { cx, invoker, pick } from "./util";
 
 export type ButtonVariant = "primary" | "secondary" | "ghost" | "danger";
@@ -7,43 +7,68 @@ export type ButtonSize = "sm" | "md" | "lg";
 const VARIANTS = ["secondary", "ghost", "danger"] as const;
 const SIZES = ["sm", "lg"] as const;
 const TYPES = ["submit", "reset"] as const;
+const TONES = ["outline"] as const;
+const SHAPES = ["square"] as const;
+const WIDTHS = ["full"] as const;
 
-export type ButtonProps = Omit<ComponentPropsWithRef<"button">, "type"> & {
+export type ButtonTone = "fill" | "outline";
+export type ButtonShape = "auto" | "square";
+export type ButtonWidth = "auto" | "full";
+
+type Options = {
   variant?: ButtonVariant;
   size?: ButtonSize;
+  tone?: ButtonTone;
+  shape?: ButtonShape;
+  width?: ButtonWidth;
+  icon?: ReactNode;
+};
+
+// options returns the data attributes and the content shared by Button and
+// ButtonLink. A square button shows only its icon, and its text label becomes
+// the accessible name, as in the Go partial.
+function options({ variant, size, tone, shape, width, icon }: Options, children: ReactNode) {
+  const square = pick(shape, SHAPES);
+  if (square && typeof children !== "string") throw new Error("ui.button: a square button's label must be text");
+  return {
+    attrs: {
+      "data-variant": pick(variant, VARIANTS),
+      "data-size": pick(size, SIZES),
+      "data-tone": pick(tone, TONES),
+      "data-shape": square,
+      "aria-label": square ? (children as string) : undefined,
+      "data-width": pick(width, WIDTHS),
+    },
+    content: square ? icon : <>{icon}{children}</>,
+  };
+}
+
+export type ButtonProps = Omit<ComponentPropsWithRef<"button">, "type" | "children"> & Options & {
+  children?: ReactNode;
   type?: "button" | "submit" | "reset";
   command?: string;
   commandFor?: string;
 };
 
-export function Button({ variant, size, type, command, commandFor, className, ...rest }: ButtonProps) {
+export function Button({ variant, size, tone, shape, width, icon, type, command, commandFor, className, children, ...rest }: ButtonProps) {
+  const o = options({ variant, size, tone, shape, width, icon }, children);
   return (
-    <button
-      {...rest}
-      {...invoker(command, commandFor)}
-      className={cx("ui-btn", className)}
-      type={pick(type, TYPES) ?? "button"}
-      data-variant={pick(variant, VARIANTS)}
-      data-size={pick(size, SIZES)}
-    />
+    <button {...rest} {...invoker(command, commandFor)} className={cx("ui-btn", className)} type={pick(type, TYPES) ?? "button"} {...o.attrs}>
+      {o.content}
+    </button>
   );
 }
 
-export type ButtonLinkProps = ComponentPropsWithRef<"a"> & {
-  variant?: ButtonVariant;
-  size?: ButtonSize;
+export type ButtonLinkProps = Omit<ComponentPropsWithRef<"a">, "children"> & Options & {
+  children?: ReactNode;
   disabled?: boolean;
 };
 
-export function ButtonLink({ variant, size, disabled, href, className, ...rest }: ButtonLinkProps) {
+export function ButtonLink({ variant, size, tone, shape, width, icon, disabled, href, className, children, ...rest }: ButtonLinkProps) {
+  const o = options({ variant, size, tone, shape, width, icon }, children);
   return (
-    <a
-      {...rest}
-      className={cx("ui-btn", className)}
-      href={disabled ? undefined : href}
-      aria-disabled={disabled ? "true" : undefined}
-      data-variant={pick(variant, VARIANTS)}
-      data-size={pick(size, SIZES)}
-    />
+    <a {...rest} className={cx("ui-btn", className)} href={disabled ? undefined : href} aria-disabled={disabled ? "true" : undefined} {...o.attrs}>
+      {o.content}
+    </a>
   );
 }
```

- [ ] **52. Modify `react/field.tsx`.** Write this patch to `/tmp/v03-react-field.tsx.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-react-field.tsx.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/react/field.tsx b/react/field.tsx
index 6e0fafd..8b047fc 100644
--- a/react/field.tsx
+++ b/react/field.tsx
@@ -5,8 +5,12 @@ const INPUT_TYPES = ["email", "password", "number", "search", "tel", "url", "dat
 export type InputType = "text" | (typeof INPUT_TYPES)[number];
 export type FieldOption = { value: string; label: string };
 
+export type FieldSize = "sm" | "md" | "lg";
+const SIZES = ["sm", "lg"] as const;
+
 type FieldBase = {
   label: ReactNode;
+  size?: FieldSize;
   name?: string;
   id?: string;
   hint?: ReactNode;
@@ -15,16 +19,16 @@ type FieldBase = {
 };
 type Own = keyof FieldBase | "control";
 
-export type InputFieldProps = FieldBase & { control?: "input"; type?: InputType } & Omit<ComponentPropsWithRef<"input">, Own | "type">;
+export type InputFieldProps = FieldBase & { control?: "input"; type?: InputType; icon?: ReactNode } & Omit<ComponentPropsWithRef<"input">, Own | "type">;
 export type TextareaFieldProps = FieldBase & { control: "textarea" } & Omit<ComponentPropsWithRef<"textarea">, Own>;
 export type SelectFieldProps = FieldBase & { control: "select"; options: readonly FieldOption[] } & Omit<ComponentPropsWithRef<"select">, Own>;
 export type FieldProps = InputFieldProps | TextareaFieldProps | SelectFieldProps;
 
-type Wiring = { id: string; name?: string; "aria-describedby"?: string; "aria-invalid"?: ControlAria["aria-invalid"] };
+type Wiring = { id: string; name?: string; "data-size"?: string; "aria-describedby"?: string; "aria-invalid"?: ControlAria["aria-invalid"] };
 type ControlAria = Pick<ComponentPropsWithRef<"input">, "aria-describedby" | "aria-invalid">;
 
 export function Field(props: FieldProps) {
-  const { label, name, id: idProp, hint, error, className, ...controlProps } = props;
+  const { label, name, id: idProp, hint, error, className, size, ...controlProps } = props;
   const { "aria-describedby": callerDescribedBy, "aria-invalid": callerInvalid } = controlProps as ControlAria;
   const id = idProp || name;
   if (!id) throw new Error("ui.field: needs ID or Name");
@@ -36,6 +40,7 @@ export function Field(props: FieldProps) {
   const wiring: Wiring = {
     id,
     name,
+    "data-size": pick(size, SIZES),
     "aria-describedby": describedBy || undefined,
     "aria-invalid": hasError ? "true" : callerInvalid,
   };
@@ -51,6 +56,7 @@ export function Field(props: FieldProps) {
 
 type ControlProps = {
   control?: string;
+  icon?: ReactNode;
   type?: string;
   options?: readonly FieldOption[];
   defaultValue?: string | number | readonly string[];
@@ -61,7 +67,7 @@ type ControlProps = {
 const blank = (v: ControlProps["defaultValue"]) => (v === "" ? undefined : v);
 
 function control(props: object, wiring: Wiring): ReactNode {
-  const { control: kind, type, options, defaultValue, ...rest } = props as ControlProps;
+  const { control: kind, type, options, defaultValue, icon, ...rest } = props as ControlProps;
   if (kind === "textarea") {
     return <textarea {...(rest as ComponentPropsWithRef<"textarea">)} {...wiring} className="ui-textarea" defaultValue={blank(defaultValue)} />;
   }
@@ -74,5 +80,12 @@ function control(props: object, wiring: Wiring): ReactNode {
       </select>
     );
   }
-  return <input {...(rest as ComponentPropsWithRef<"input">)} {...wiring} className="ui-input" type={pick(type, INPUT_TYPES) ?? "text"} defaultValue={blank(defaultValue)} />;
+  const input = <input {...(rest as ComponentPropsWithRef<"input">)} {...wiring} className="ui-input" type={pick(type, INPUT_TYPES) ?? "text"} defaultValue={blank(defaultValue)} />;
+  if (!present(icon)) return input;
+  return (
+    <div className="ui-input-group">
+      {icon}
+      {input}
+    </div>
+  );
 }
```

- [ ] **53. Modify `react/pill.tsx`.** Write this patch to `/tmp/v03-react-pill.tsx.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-react-pill.tsx.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/react/pill.tsx b/react/pill.tsx
index f60d601..3eb79b3 100644
--- a/react/pill.tsx
+++ b/react/pill.tsx
@@ -2,11 +2,21 @@ import type { ComponentPropsWithRef } from "react";
 import { cx, pick } from "./util";
 
 export type PillVariant = "neutral" | "info" | "good" | "warning" | "bad";
+export type PillTone = "fill" | "outline";
 
 const VARIANTS = ["info", "good", "warning", "bad"] as const;
+const TONES = ["outline"] as const;
 
-export type PillProps = ComponentPropsWithRef<"span"> & { variant?: PillVariant };
+export type PillProps = ComponentPropsWithRef<"span"> & { variant?: PillVariant; tone?: PillTone; dot?: boolean };
 
-export function Pill({ variant, className, ...rest }: PillProps) {
-  return <span {...rest} className={cx("ui-pill", className)} data-variant={pick(variant, VARIANTS)} />;
+export function Pill({ variant, tone, dot, className, ...rest }: PillProps) {
+  return (
+    <span
+      {...rest}
+      className={cx("ui-pill", className)}
+      data-variant={pick(variant, VARIANTS)}
+      data-tone={pick(tone, TONES)}
+      data-dot={dot ? "" : undefined}
+    />
+  );
 }
```

- [ ] **54. Modify `react/dialog.tsx`.** Write this patch to `/tmp/v03-react-dialog.tsx.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-react-dialog.tsx.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/react/dialog.tsx b/react/dialog.tsx
index 97753b0..093f2cf 100644
--- a/react/dialog.tsx
+++ b/react/dialog.tsx
@@ -1,22 +1,26 @@
 import type { ComponentPropsWithRef, ReactNode } from "react";
-import { cx, invoker } from "./util";
+import { cx, invoker, present } from "./util";
 
 export type DialogProps = Omit<ComponentPropsWithRef<"dialog">, "title" | "id"> & {
   id: string;
   title: ReactNode;
   closeLabel?: string;
+  closeButton?: boolean;
+  eyebrow?: ReactNode;
 };
 
 // Dialog is a modal <dialog>. Open it from any button with
 // command="show-modal" commandfor={id}; the × button closes it with
-// command="close". It needs no JavaScript (spec §4.3).
-export function Dialog({ id, title, closeLabel = "Close", className, children, ...rest }: DialogProps) {
+// command="close". It needs no JavaScript (spec §4.3). Without the ×
+// (closeButton={false}), give it a button that closes it.
+export function Dialog({ id, title, closeLabel = "Close", closeButton = true, eyebrow, className, children, ...rest }: DialogProps) {
   const titleId = `${id}-title`;
   return (
     <dialog {...rest} id={id} className={cx("ui-dialog", className)} aria-labelledby={titleId}>
+      {present(eyebrow) && <p className="ui-dialog-eyebrow">{eyebrow}</p>}
       <header className="ui-dialog-head">
         <h2 className="ui-dialog-title" id={titleId}>{title}</h2>
-        <button className="ui-dialog-x" type="button" aria-label={closeLabel} {...invoker("close", id)}>×</button>
+        {closeButton && <button className="ui-dialog-x" type="button" aria-label={closeLabel} {...invoker("close", id)}>×</button>}
       </header>
       <div className="ui-dialog-body">{children}</div>
     </dialog>
```

- [ ] **55. Modify `react/index.ts`.** Write this patch to `/tmp/v03-react-index.ts.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-react-index.ts.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/react/index.ts b/react/index.ts
index 759dd66..483fa0d 100644
--- a/react/index.ts
+++ b/react/index.ts
@@ -1,10 +1,10 @@
 export { Button, ButtonLink } from "./button";
-export type { ButtonLinkProps, ButtonProps, ButtonSize, ButtonVariant } from "./button";
+export type { ButtonLinkProps, ButtonProps, ButtonShape, ButtonSize, ButtonTone, ButtonVariant, ButtonWidth } from "./button";
 export { Field } from "./field";
-export type { FieldOption, FieldProps, InputType } from "./field";
+export type { FieldOption, FieldProps, FieldSize, InputType } from "./field";
 export { Check } from "./check";
 export type { CheckProps } from "./check";
 export { Pill } from "./pill";
-export type { PillProps, PillVariant } from "./pill";
+export type { PillProps, PillTone, PillVariant } from "./pill";
 export { Dialog } from "./dialog";
 export type { DialogProps } from "./dialog";
```

- [ ] **56.** Run `npx tsc --noEmit && npx vitest run`. Expect every test to pass, the new parity fixtures among them.

- [ ] **57. Commit:** `git add react/button.tsx react/field.tsx react/pill.tsx react/dialog.tsx react/index.ts` then `git commit -m "feat(react): the same options on Button, ButtonLink, Field, Pill and Dialog"`.

### Task 4: Gallery, fixture and browser tests

**Files:** `gallery/main.go`, `gallery/page.html`, `gallery/fixture_library.html`, `e2e/options.spec.ts`

- [ ] **58. Modify `gallery/main.go`.** Write this patch to `/tmp/v03-gallery-main.go.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-gallery-main.go.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/gallery/main.go b/gallery/main.go
index a42973c..39d0d7b 100644
--- a/gallery/main.go
+++ b/gallery/main.go
@@ -48,6 +48,16 @@ type page struct {
 	Compact       bool
 	Links         []link
 	Currencies    []ui.Option
+	Icons         icons
+}
+
+// icons are the app's own markup for the partials' Icon argument.
+type icons struct{ Download, Theme, Search template.HTML }
+
+var galleryIcons = icons{
+	Download: `<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M8 2v8m-3.5-3.5L8 10l3.5-3.5M3 13h10"/></svg>`,
+	Theme:    `<svg viewBox="0 0 16 16" aria-hidden="true"><circle cx="8" cy="8" r="6" fill="none" stroke="currentColor" stroke-width="1.6"/><path d="M8 2a6 6 0 0 1 0 12z" fill="currentColor"/></svg>`,
+	Search:   `<svg class="ui-input-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" aria-hidden="true"><circle cx="7" cy="7" r="4.5"/><path d="m10.5 10.5 3 3"/></svg>`,
 }
 
 func main() {
@@ -88,7 +98,7 @@ func newHandler() (http.Handler, error) {
 				return
 			}
 			var b bytes.Buffer
-			data := page{Palette: p, Mode: m, Dialog: dialog, Compact: compact, Links: links(p, m), Currencies: currencies}
+			data := page{Palette: p, Mode: m, Dialog: dialog, Compact: compact, Links: links(p, m), Currencies: currencies, Icons: galleryIcons}
 			if err := set.ExecuteTemplate(&b, "page", data); err != nil {
 				http.Error(w, err.Error(), http.StatusInternalServerError)
 				return
```

- [ ] **59. Modify `gallery/page.html`.** Write this patch to `/tmp/v03-gallery-page.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-gallery-page.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/gallery/page.html b/gallery/page.html
index 4a2f6a7..ec83674 100644
--- a/gallery/page.html
+++ b/gallery/page.html
@@ -11,6 +11,9 @@
   .gallery > section > h2 { margin-block-end: var(--space-4); font-size: var(--text-xl); }
   .gallery-nav a[aria-current="page"] { font-weight: 700; }
   .gallery-form { max-width: 34rem; }
+  .gallery-narrow { max-width: 16rem; }
+  .gallery-row { margin-block-start: var(--space-3); }
+  .gallery-inline { position: static; margin: var(--space-3) 0 0; }
 </style>
 </head>
 <body>
@@ -38,6 +41,18 @@
         {{template "ui.button" (dict "Label" "Medium")}}
         {{template "ui.button" (dict "Label" "Large" "Size" "lg")}}
       </div>
+      <div class="ui-cluster">
+        {{template "ui.button" (dict "Label" "Outline" "Tone" "outline")}}
+        {{template "ui.button" (dict "Label" "Remove" "Variant" "danger" "Tone" "outline")}}
+        {{template "ui.button" (dict "Label" "Off" "Variant" "danger" "Tone" "outline" "Disabled" true)}}
+        {{template "ui.button" (dict "Label" "Ghost off" "Variant" "ghost" "Disabled" true)}}
+        {{template "ui.button" (dict "Label" "Download" "Variant" "secondary" "Icon" .Icons.Download)}}
+        {{template "ui.button" (dict "Label" "Theme" "Variant" "secondary" "Shape" "square" "Icon" .Icons.Theme)}}
+        {{template "ui.button" (dict "Label" "Theme" "Shape" "square" "Size" "sm" "Icon" .Icons.Theme)}}
+      </div>
+      <div class="gallery-narrow">
+        {{template "ui.button" (dict "Label" "Use this category for every payment to this merchant" "Variant" "secondary" "Width" "full")}}
+      </div>
     </div>
   </section>
 
@@ -52,6 +67,12 @@
       </div>
       {{template "ui.field" (dict "Label" "Notes" "Name" "notes" "Control" "textarea" "Rows" 3 "Hint" "Optional.")}}
       {{template "ui.field" (dict "Label" "Account ID" "Name" "account" "Value" "acc_2041" "Disabled" true)}}
+      {{template "ui.field" (dict "Label" "Search" "Name" "q" "Type" "search" "Placeholder" "Merchant or note" "Icon" .Icons.Search)}}
+      <div class="ui-form-row">
+        {{template "ui.field" (dict "Label" "Host" "Name" "host" "Size" "sm" "Value" "imap.example.com")}}
+        <div style="--ui-row-grow: 0; --ui-row-basis: 6rem">{{template "ui.field" (dict "Label" "Port" "Name" "port" "Size" "sm" "Type" "number" "Value" "993")}}</div>
+      </div>
+      {{template "ui.field" (dict "Label" "Work email" "Name" "work" "Type" "email" "Size" "lg" "Placeholder" "you@company.com")}}
       <fieldset class="ui-stack">
         <legend class="ui-field-label">Billing period</legend>
         {{template "ui.check" (dict "Label" "Monthly" "Name" "period" "Value" "m" "Type" "radio" "Checked" true)}}
@@ -71,6 +92,12 @@
       {{template "ui.pill" (dict "Label" "Due soon" "Variant" "warning")}}
       {{template "ui.pill" (dict "Label" "Failed" "Variant" "bad")}}
     </div>
+    <div class="ui-cluster gallery-row">
+      {{template "ui.pill" (dict "Label" "1 violation" "Variant" "bad" "Dot" true)}}
+      {{template "ui.pill" (dict "Label" "Weak" "Variant" "warning" "Dot" true)}}
+      {{template "ui.pill" (dict "Label" "meta" "Tone" "outline")}}
+      {{template "ui.pill" (dict "Label" "Price up" "Variant" "warning" "Tone" "outline")}}
+    </div>
   </section>
 
   <section aria-labelledby="s-callouts">
@@ -103,6 +130,8 @@
     <div>
       <details class="ui-disclosure" name="faq" open><summary class="ui-disclosure-summary">Where is my data stored?</summary><div class="ui-disclosure-body"><p>On your own server, in one SQLite file.</p></div></details>
       <details class="ui-disclosure" name="faq"><summary class="ui-disclosure-summary">Can I export it?</summary><div class="ui-disclosure-body"><p>Yes, as CSV or JSON.</p></div></details>
+      <details class="ui-disclosure" data-marker="start" open><summary class="ui-disclosure-summary">experience</summary><div class="ui-disclosure-body"><p>invented_lead.md</p></div></details>
+      <details class="ui-disclosure" data-marker="start"><summary class="ui-disclosure-summary">education</summary><div class="ui-disclosure-body"><p>invented_degree.md</p></div></details>
     </div>
   </section>
 
@@ -125,6 +154,20 @@
       </header>
       <div class="ui-dialog-body">This cannot be undone.</div>
     </dialog>
+    <p class="gallery-row ui-muted">Shown without script: an open, non-modal dialog with an eyebrow and no ×.</p>
+    <dialog class="ui-dialog gallery-inline" id="remove" aria-labelledby="remove-title" open>
+      <p class="ui-dialog-eyebrow">Can't be undone</p>
+      <header class="ui-dialog-head">
+        <h2 class="ui-dialog-title" id="remove-title">Remove grace@example.com?</h2>
+      </header>
+      <div class="ui-dialog-body">
+        <p>Their sessions end now.</p>
+        <form class="ui-dialog-actions" method="dialog">
+          {{template "ui.button" (dict "Label" "Cancel" "Variant" "secondary" "Type" "submit")}}
+          {{template "ui.button" (dict "Label" "Remove" "Variant" "danger" "Type" "submit")}}
+        </form>
+      </div>
+    </dialog>
   </section>
 </main>
 {{if .Dialog}}<script>document.getElementById("confirm").showModal();</script>{{end}}
```

- [ ] **60. Modify `gallery/fixture_library.html`.** Write this patch to `/tmp/v03-gallery-fixture_library.html.patch`, byte for byte, then run `git apply --whitespace=nowarn /tmp/v03-gallery-fixture_library.html.patch` and delete the patch file. It must apply cleanly; if it doesn't, STOP.

```diff
diff --git a/gallery/fixture_library.html b/gallery/fixture_library.html
index 23131fd..9eb5b25 100644
--- a/gallery/fixture_library.html
+++ b/gallery/fixture_library.html
@@ -17,8 +17,31 @@
     <button class="ui-btn" id="btn-sm" data-size="sm">Small</button>
     <button class="ui-btn" id="btn-lg" data-size="lg">Large</button>
     <a class="ui-btn" href="#library">Link</a>
-    <button class="ui-btn" disabled>Off</button>
+    <button class="ui-btn" id="btn-off" disabled>Off</button>
   </div>
+  <div class="ui-cluster">
+    <button class="ui-btn" id="btn-outline" data-variant="danger" data-tone="outline">Remove</button>
+    <button class="ui-btn" id="btn-outline-primary" data-tone="outline">Outline</button>
+    <button class="ui-btn" id="btn-outline-off" data-variant="danger" data-tone="outline" disabled>Remove</button>
+    <button class="ui-btn" id="ghost-off" data-variant="ghost" disabled>Ghost</button>
+    <a class="ui-btn" id="link-off" aria-disabled="true">Link</a>
+    <button class="ui-btn" id="btn-square" data-variant="secondary" data-shape="square" aria-label="Theme"><svg viewBox="0 0 16 16" aria-hidden="true"><circle cx="8" cy="8" r="6" fill="currentColor"/></svg></button>
+    <button class="ui-btn" id="btn-square-sm" data-size="sm" data-shape="square" aria-label="Theme"><svg viewBox="0 0 16 16" aria-hidden="true"><circle cx="8" cy="8" r="6" fill="currentColor"/></svg></button>
+  </div>
+  <div id="full-box" style="inline-size: 10rem"><button class="ui-btn" id="btn-full" data-width="full">Use this category for every payment</button></div>
+  <div class="ui-cluster">
+    <input class="ui-input" id="input-sm" data-size="sm" aria-label="Small">
+    <select class="ui-select" id="select-sm" data-size="sm" aria-label="Small select"><option>One</option></select>
+    <input class="ui-input" id="input-lg" data-size="lg" aria-label="Large">
+    <select class="ui-select" id="select-lg" data-size="lg" aria-label="Large select"><option>One</option></select>
+    <textarea class="ui-textarea" id="area-sm" data-size="sm" aria-label="Small area"></textarea>
+  </div>
+  <div class="ui-input-group" id="group"><svg class="ui-input-icon" id="icon" viewBox="0 0 16 16" aria-hidden="true"><circle cx="7" cy="7" r="5" fill="none" stroke="currentColor"/></svg><input class="ui-input" id="input-icon" type="search" aria-label="Search"></div>
+  <div class="ui-form-row" id="row">
+    <div class="ui-field"><label class="ui-field-label" for="row-wide">Host</label><input class="ui-input" id="row-wide"></div>
+    <div class="ui-field" id="row-narrow" style="--ui-row-grow: 0; --ui-row-basis: 5rem"><label class="ui-field-label" for="row-port">Port</label><input class="ui-input" id="row-port"></div>
+  </div>
+  <input class="ui-input" id="input-width" style="--ui-input-width: 8rem" aria-label="Narrow">
   <div class="ui-field" id="field">
     <label class="ui-field-label" id="label" for="input">Label</label>
     <input class="ui-input" id="input" value="text">
@@ -29,16 +52,23 @@
   <div class="ui-field"><label class="ui-field-label" for="multi">Several</label><select class="ui-select" id="multi" multiple size="2"><option>One</option><option>Two</option></select></div>
   <div class="ui-field"><label class="ui-field-label" for="area">Area</label><textarea class="ui-textarea" id="area">x</textarea></div>
   <label class="ui-check"><input class="ui-check-input" type="checkbox" checked> Check</label>
-  <div class="ui-cluster"><span class="ui-pill" id="pill">Neutral</span><span class="ui-pill" data-variant="good">Good</span><span class="ui-pill" data-variant="bad">Bad</span></div>
+  <div class="ui-cluster"><span class="ui-pill" id="pill">Neutral</span><span class="ui-pill" data-variant="good">Good</span><span class="ui-pill" data-variant="bad">Bad</span><span class="ui-pill" id="pill-dot" data-variant="bad" data-dot>1 violation</span><span class="ui-pill" id="pill-outline" data-tone="outline">meta</span><span class="ui-pill" id="pill-outline-warn" data-variant="warning" data-tone="outline">Price up</span></div>
   <div class="ui-callout" id="callout" data-variant="bad"><p class="ui-callout-title">Title</p><p>Body with <a id="callout-link" href="#library">a link</a> and <strong>bold</strong>.</p></div>
   <div class="ui-table-wrap" id="wrap"><table class="ui-table"><caption>Caption</caption><thead><tr><th scope="col">A</th><th scope="col" class="ui-num">B</th></tr></thead><tbody><tr><td id="td">x</td><td class="ui-num">1</td></tr></tbody></table></div>
   <details class="ui-disclosure" open><summary class="ui-disclosure-summary">More</summary><div class="ui-disclosure-body"><p>Body</p><ul><li>Item</li></ul></div></details>
+  <details class="ui-disclosure" id="disc-start" data-marker="start"><summary class="ui-disclosure-summary" id="disc-start-summary">Folder</summary><div class="ui-disclosure-body"><p>File</p></div></details>
+  <details class="ui-disclosure" id="disc-norule" style="--ui-disclosure-rule: transparent"><summary class="ui-disclosure-summary">No rule</summary><div class="ui-disclosure-body"><p>Body</p></div></details>
   <section class="ui-card" id="card"><h2>Card head</h2><p>Text</p></section>
   <dialog class="ui-dialog" id="dialog" aria-labelledby="dialog-title">
     <div class="ui-dialog-head"><h2 class="ui-dialog-title" id="dialog-title">Title</h2><button class="ui-dialog-x" aria-label="Close">x</button></div>
     <div class="ui-dialog-body"><p>Body</p></div>
     <div class="ui-dialog-actions"><button class="ui-btn">OK</button></div>
   </dialog>
+  <dialog class="ui-dialog" id="dialog-eyebrow" aria-labelledby="dialog-eyebrow-title">
+    <p class="ui-dialog-eyebrow" id="eyebrow">Can't be undone</p>
+    <div class="ui-dialog-head"><h2 class="ui-dialog-title" id="dialog-eyebrow-title">Delete?</h2></div>
+    <div class="ui-dialog-body"><p>Body</p></div>
+  </dialog>
 </main>
 </body>
 </html>
```

- [ ] **61. Create `e2e/options.spec.ts`** with exactly this content:

```ts
import { expect, test, type Page } from "@playwright/test";

// The v0.3 options, measured on the library fixture at both densities.

const rect = (page: Page, sel: string) =>
  page.locator(sel).evaluate((el) => {
    const r = el.getBoundingClientRect();
    return { left: r.left, right: r.right, top: r.top, bottom: r.bottom, width: r.width, height: r.height };
  });

const css = (page: Page, sel: string, prop: string, pseudo: string | null = null): Promise<string> =>
  page.locator(sel).evaluate((el, [prop, pseudo]) => getComputedStyle(el, pseudo).getPropertyValue(prop), [prop, pseudo] as const);

// token resolves a colour token the way the browser paints it, inside sel.
const token = (page: Page, sel: string, name: string): Promise<string> =>
  page.locator(sel).evaluate((el, name) => {
    const probe = document.createElement("span");
    probe.style.color = `var(${name})`;
    el.appendChild(probe);
    const c = getComputedStyle(probe).color;
    probe.remove();
    return c;
  }, name);

const showModal = (page: Page, sel: string) => page.locator(sel).evaluate((d: HTMLDialogElement) => d.showModal());

for (const [name, query] of [["default", ""], ["compact", "?density=root"]] as const) {
  test.describe(name, () => {
    test.beforeEach(async ({ page }) => page.goto(`/fixture/library/ui.css${query}`));

    test("small and large fields are as tall as the button of their size", async ({ page }) => {
      const sm = (await rect(page, "#btn-sm")).height;
      const lg = (await rect(page, "#btn-lg")).height;
      expect(sm).toBe(name === "compact" ? 28 : 32);
      expect(lg).toBe(name === "compact" ? 44 : 48);
      expect((await rect(page, "#input-sm")).height).toBe(sm);
      expect((await rect(page, "#select-sm")).height).toBe(sm);
      expect((await rect(page, "#input-lg")).height).toBe(lg);
      expect((await rect(page, "#select-lg")).height).toBe(lg);
      expect((await rect(page, "#area-sm")).height).toBe(sm * 2);
      // A select keeps the room its chevron needs.
      expect(await css(page, "#select-sm", "padding-right")).toBe("36px");
    });

    test("a square button is as wide as it is tall", async ({ page }) => {
      for (const [sel, ref] of [["#btn-square", "#btn"], ["#btn-square-sm", "#btn-sm"]]) {
        const r = await rect(page, sel);
        expect(r.width).toBe(r.height);
        expect(r.height).toBe((await rect(page, ref)).height);
      }
    });

    test("a full-width button fills its box and wraps its label", async ({ page }) => {
      const box = await rect(page, "#full-box");
      const b = await rect(page, "#btn-full");
      expect(b.width).toBe(box.width);
      expect(b.height).toBeGreaterThan((await rect(page, "#btn")).height);
    });

    test("an outline button keeps the variant's colour as text on the sheet", async ({ page }) => {
      expect(await css(page, "#btn-outline", "color")).toBe(await token(page, "#library", "--bad"));
      expect(await css(page, "#btn-outline", "background-color")).toBe(await token(page, "#library", "--sheet"));
      expect(await css(page, "#btn-outline", "border-top-color")).toBe(await token(page, "#library", "--rule"));
      expect(await css(page, "#btn-outline-primary", "color")).toBe(await token(page, "#library", "--primary"));
    });

    test("disabled changes colour, not opacity", async ({ page }) => {
      const band = await token(page, "#library", "--band");
      const ink3 = await token(page, "#library", "--ink-3");
      for (const sel of ["#btn-off", "#btn-outline-off", "#link-off"]) {
        expect(await css(page, sel, "opacity"), sel).toBe("1");
        expect(await css(page, sel, "background-color"), sel).toBe(band);
        expect(await css(page, sel, "color"), sel).toBe(ink3);
      }
      expect(await css(page, "#ghost-off", "background-color")).toBe("rgba(0, 0, 0, 0)");
      expect(await css(page, "#ghost-off", "color")).toBe(ink3);
      expect(await css(page, "#btn-off", "cursor")).toBe("not-allowed");
      expect(await css(page, "#btn-off", "pointer-events")).toBe("auto");
      expect(await css(page, "#link-off", "pointer-events")).toBe("none");
    });

    test("the input icon sits inside the input, before its text", async ({ page }) => {
      const input = await rect(page, "#input-icon");
      const icon = await rect(page, "#icon");
      expect(icon.width).toBe(16);
      expect(icon.left).toBeGreaterThan(input.left);
      expect(icon.top).toBeGreaterThan(input.top);
      expect(icon.bottom).toBeLessThan(input.bottom);
      const start = parseFloat(await css(page, "#input-icon", "padding-left"));
      expect(input.left + start).toBeGreaterThan(icon.right);
    });

    test("width knobs set a field's width and a row child's basis", async ({ page }) => {
      expect((await rect(page, "#input-width")).width).toBe(128);
      expect((await rect(page, "#row-narrow")).width).toBe(80);
      expect((await rect(page, "#row-wide")).width).toBeGreaterThan(200);
    });

    test("cards and dialogs are flat only in compact", async ({ page }) => {
      const flat = name === "compact";
      const card = await css(page, "#card", "box-shadow");
      expect(card === "none").toBe(flat);
      await showModal(page, "#dialog");
      const dialog = await css(page, "#dialog", "box-shadow");
      expect(dialog === "none").toBe(flat);
      expect(await css(page, "#dialog", "border-top-width")).toBe("1px");
    });

    test("the eyebrow sits above the title, and a dialog lays out without its ×", async ({ page }) => {
      await showModal(page, "#dialog-eyebrow");
      const eyebrow = await rect(page, "#eyebrow");
      const title = await rect(page, "#dialog-eyebrow-title");
      expect(eyebrow.bottom).toBeLessThanOrEqual(title.top);
      // The eyebrow's text starts where the title's does.
      const start = parseFloat(await css(page, "#eyebrow", "padding-left"));
      expect(Math.abs(eyebrow.left + start - title.left)).toBeLessThan(1);
      expect(await css(page, "#eyebrow", "text-transform")).toBe("uppercase");
      expect(title.width).toBeGreaterThan(0);
    });

    test("a dotted pill draws its dot in its text colour; an outline pill has no fill", async ({ page }) => {
      expect(await css(page, "#pill-dot", "width", "::before")).toBe("6px");
      expect(await css(page, "#pill-dot", "background-color", "::before")).toBe(await css(page, "#pill-dot", "color"));
      expect(await css(page, "#pill-outline", "background-color")).toBe("rgba(0, 0, 0, 0)");
      expect(await css(page, "#pill-outline", "border-top-color")).toBe(await token(page, "#library", "--rule-field"));
      expect(await css(page, "#pill-outline-warn", "border-top-color")).toBe(await token(page, "#library", "--warning"));
    });

    test("a start-marker disclosure draws its chevron first and turns it when open", async ({ page }) => {
      expect(await css(page, "#disc-start-summary", "content", "::after")).toBe("none");
      expect(await css(page, "#disc-start-summary", "width", "::before")).toBe("8px");
      const closed = await css(page, "#disc-start-summary", "transform", "::before");
      await page.locator("#disc-start-summary").click();
      await expect.poll(() => css(page, "#disc-start-summary", "transform", "::before")).not.toBe(closed);
      expect(await css(page, "#disc-norule", "border-bottom-color")).toBe("rgba(0, 0, 0, 0)");
    });
  });
}

test.describe("on a touch screen", () => {
  test.use({ hasTouch: true, isMobile: true, viewport: { width: 1100, height: 900 } });

  test("compact controls grow to finger sizes; default ones don't change", async ({ page }) => {
    await page.goto("/fixture/library/ui.css?density=root");
    expect(await page.evaluate(() => matchMedia("(pointer: coarse)").matches), "the page sees a coarse pointer").toBe(true);
    expect((await rect(page, "#btn")).height).toBe(40);
    expect((await rect(page, "#btn-sm")).height).toBe(34);
    expect((await rect(page, "#btn-lg")).height).toBe(44);
    expect((await rect(page, "#input")).height).toBe(40);
    expect((await rect(page, "#input-sm")).height).toBe(34);
    expect((await rect(page, "#select-lg")).height).toBe(44);
    await page.goto("/fixture/library/ui.css");
    expect((await rect(page, "#btn")).height).toBe(40);
    expect((await rect(page, "#btn-sm")).height).toBe(32);
    expect((await rect(page, "#input-sm")).height).toBe(32);
  });
});
```

- [ ] **62.** Run `go test -count=1 ./gallery`. Expect `ok`: the gallery's confirm dialog still matches `testdata/parity/dialog-basic.html`.

- [ ] **63.** Run `npx playwright test e2e/options.spec.ts e2e/density.spec.ts e2e/scoped.spec.ts --ignore-snapshots`. Expect 36 passed.

- [ ] **64.** Run `npx playwright test e2e/gallery.spec.ts --ignore-snapshots`. Expect 18 passed: axe passes on all 16 gallery pages.

- [ ] **65. Commit:** `git add gallery/main.go gallery/page.html gallery/fixture_library.html e2e/options.spec.ts` then `git commit -m "test: measure every v0.3 option; show them in the gallery"`.

### Task 5: The gate

- [ ] **66.** Run the gate above. Expect it all to pass: Go, the site, `tsc`, 61 vitest tests, 54 browser tests and 57 site tests.
