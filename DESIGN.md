# ai-ui design contract

This file is the contract between ai-ui and the apps that use it. A change
that renames or removes anything named here is a breaking change (see
CHANGELOG.md and spec §9).

## Override model

Every library rule lives in a cascade layer:

```css
@layer ui.reset, ui.base, ui.components, ui.utilities;
```

Any unlayered rule in an app beats every library rule, with no `!important`.
You can override at three levels:

1. **Theme**: set tokens on `:root` in your app CSS: `:root { --primary: oklch(52% 0.2 150); }`.
   Use `light-dark(<light>, <dark>)` to cover both modes in one line.
2. **Component**: set a knob: `.ui-btn { --ui-btn-radius: 0; }`. A knob set on the
   root class also overrides the knob values the variants set, so target the default
   variant alone with `.ui-btn:not([data-variant])`.
3. **Anything else**: write a plain rule. An unlayered app rule also beats the
   library's state rules (`:hover`, `[aria-invalid]`) for the same property, so
   restate the states you need.

Tokens that belong to one app are named `--app-*`. The library never defines or reads them.

## Adopting ai-ui in an app that has its own CSS

A new app needs none of this. An app with existing markup does, because `ui.reset` and
`ui.base` change every element its CSS says nothing about: list bullets, heading sizes, button
frames. A compat layer between `ui.base` and `ui.components` keeps those elements as they were
and lets the app move one component at a time. money moved this way, and the only pixels that
changed were the focus ring on its buttons, which it chose to take.

Load the compat sheet first, then `ui.css`, then the app's own CSS. The first line has to come
before `ui.css`, because a layer's place is fixed the first time it is named:

```css
@layer ui.reset, ui.base, ui.compat, ui.components, ui.utilities;
@layer ui.compat {
  :not(svg, svg *, [class^="ui-"], [class^="ui-"] *, [class*=" ui-"], [class*=" ui-"] *),
  :not(svg, svg *, [class^="ui-"], [class^="ui-"] *, [class*=" ui-"], [class*=" ui-"] *)::before,
  :not(svg, svg *, [class^="ui-"], [class^="ui-"] *, [class*=" ui-"], [class*=" ui-"] *)::after,
  :not(svg, svg *, [class^="ui-"], [class^="ui-"] *, [class*=" ui-"], [class*=" ui-"] *)::selection { all: revert; }
  svg:not([class^="ui-"] *, [class*=" ui-"] *) { display: revert; max-width: revert; }
}
```

Each exclusion is there for a reason:

- **`svg`:** `all` also wipes an SVG's presentation attributes, so icons and charts would lose
  their fills and strokes.
- **`ui-*` subtrees:** library markup needs the reset. This means a `ui-*` element skips the
  layer, so a heading with a `ui-*` class takes the library's heading rules, not the browser's.

Don't reach for the two obvious alternatives. Wrapping the app's old CSS in a layer below `ui`
fails because a later layer wins whatever the specificity, so the reset overrides the app.
Unlayered "fix" rules fail because they also beat the components.

**App markup inside a library container** skips the layer too, because it sits in a `ui-*`
subtree. Give those containers a rule of their own, listing the ones you use, above the `svg`
line. Keep `[class^="ui-"]` and `[class*=" ui-"]` in it, and add the subtrees of any library
component that goes inside the container:

```css
:is(.ui-callout, .ui-dialog) :not(svg, svg *, [class^="ui-"], [class*=" ui-"], .ui-btn *, .ui-pill *) { all: revert; }
```

**Moving a component.** Unlayered app rules beat the library on every property they state,
including its state rules, so renaming a class keeps most of an app's rule. To keep the app's
look, override only what the library sets and the app doesn't:

- **Pin those overrides with a test.** Removing one can leave every other test green, and only a
  browser notices. Many read like dead code, because they only cancel a library rule.
- **Keep a variant's class in its selector:** `.ui-callout[data-variant="bad"]`, never a bare
  `[data-variant="bad"]`. A bare rule is global and unlayered, so it also repaints every other
  component with that variant. Screenshots can't catch this until something else uses the value.
- **State `font-family` beside a `font` shorthand.** The app's shorthand beats `.ui-mono`, so
  write `.ui-input.ui-mono { font-family: var(--font-mono); }` if you use both.
- **Compare computed styles as well as pixels,** with one control focused or hovered at a time.
  A screenshot with nothing focused can't show a focus ring.

**If you serve `ui.css` inline** under a hash-based CSP from Go's `html/template`, pass it as a
`template.CSS` value. Literal text in a template has its CSS comments stripped, so the bytes
served no longer match a hash taken over the source.

Delete the compat sheet once no legacy markup is left.

## Tokens

| Group | Tokens |
|---|---|
| Surfaces | `--paper` page, `--band` secondary strip, `--sheet` card/panel |
| Text | `--ink`, `--ink-2`, `--ink-3` |
| Lines | `--rule`, `--rule-soft`, `--rule-field` (input borders) |
| Action | `--primary`, `--on-primary`, `--primary-hover`, `--link`, `--focus`, `--accent` |
| Status | `--good`, `--good-bg`, `--warning`, `--warning-bg`, `--bad`, `--bad-bg` |
| Effects | `--mark-bg`, `--shadow`, `--scrim` |
| Type | `--font-sans`, `--font-display` (defaults to `--font-sans`), `--font-mono`, `--text-xs` … `--text-3xl` |
| Shape | `--space-1` … `--space-8`, `--radius-sm`, `--radius-md`, `--radius-lg` |
| Motion | `--dur` (0ms under reduced motion), `--ease` |

The contrast pairs `go test` enforces in every palette and mode: `ink`, `ink-2` and `ink-3`
on `paper`, `band` and `sheet`; `on-primary` on `primary`, `primary-hover` and `bad`;
`link` on `paper`, `band` and `sheet`; each status colour on its `-bg`, all at 4.5:1 or more;
and `rule-field` on `sheet` at 3:1 or more.

## Modes and palettes

- The mode follows the OS. `<html data-theme="light|dark">` forces one.
- `<html data-palette="mono">` selects the neutral palette. `blue` is the default and needs no attribute.
- An app adds a palette in its own CSS: `:root[data-palette="whiteboard"] { --paper: …; }`.
- `theme.js` (load it blocking in `<head>`) applies and persists the visitor's choice.
  Use `<button data-ui-set-theme="dark">`, `data-ui-set-theme="system"` or
  `<button data-ui-set-palette="mono">`; an empty value clears the choice. From script,
  use `window.aiui.theme.get()` and `.set({ theme, palette })`. Until the visitor
  chooses, the server's own `data-theme` / `data-palette` stay untouched.

## Components

Variants and sizes are `data-variant` / `data-size`. A default is never written as an
attribute, and an unknown value renders the default. State comes from native HTML and
ARIA only.

### Button: `.ui-btn`

Knobs: `--ui-btn-bg`, `--ui-btn-fg`, `--ui-btn-border`, `--ui-btn-bg-hover`, `--ui-btn-fg-hover`, `--ui-btn-radius`, `--ui-btn-height`, `--ui-btn-pad`, `--ui-btn-font-size`.

| | Go | React |
|---|---|---|
| Button | `{{template "ui.button" (dict "Label" "Save" "Variant" "secondary" "Size" "sm" "Type" "submit")}}` | `<Button variant="secondary" size="sm" type="submit">Save</Button>` |
| Link | `{{template "ui.button" (dict "Label" "Docs" "Href" "/docs")}}` | `<ButtonLink href="/docs">Docs</ButtonLink>` |
| Dialog opener | `(dict "Label" "Open" "Command" "show-modal" "CommandFor" "confirm")` | `<Button command="show-modal" commandFor="confirm">Open</Button>` |

Variants: `primary` (default), `secondary`, `ghost`, `danger`. Sizes: `sm`, `md` (default), `lg`.
`Disabled` on a link drops `href` and sets `aria-disabled="true"`. Labels are text.

```html
<button class="ui-btn" type="button" data-variant="secondary" data-size="sm">Save</button>
```

### Field: `.ui-field`

Knobs: `--ui-field-gap`; on controls, `--ui-input-bg`, `--ui-input-fg`, `--ui-input-border`, `--ui-input-radius`, `--ui-input-height`.

Go args: `Label` (required), `ID` or `Name` (at least one), `Control` (`textarea`, `select`; default input), `Type` (`email`, `password`, `number`, `search`, `tel`, `url`, `date`; default `text`), `Value`, `Placeholder`, `Autocomplete`, `Rows`, `Options` (`[]ui.Option`), `Hint`, `Error`, `Required`, `Disabled`.
React: `<Field label name id control type options defaultValue hint error … />`. Every other prop, including `ref`, goes to the control.

```html
<div class="ui-field">
  <label class="ui-field-label" for="email">Email</label>
  <input class="ui-input" id="email" name="email" type="email" aria-describedby="email-hint email-error" aria-invalid="true" required>
  <p class="ui-field-hint" id="email-hint">We never share it.</p>
  <p class="ui-field-error" id="email-error">That address bounced.</p>
</div>
```

Lay fields side by side with `<div class="ui-form-row">`.

### Check: `.ui-check`

Knobs: `--ui-check-size`, `--ui-check-accent`.
Go: `{{template "ui.check" (dict "Label" "Monthly" "Name" "period" "Value" "m" "Type" "radio" "Checked" true)}}`. React: `<Check label="Monthly" name="period" value="m" type="radio" defaultChecked />`.

```html
<label class="ui-check"><input class="ui-check-input" type="radio" name="period" value="m" checked><span class="ui-check-label">Monthly</span></label>
```

### Pill: `.ui-pill`

Knobs: `--ui-pill-bg`, `--ui-pill-fg`, `--ui-pill-border`. Variants: neutral (default), `info`, `good`, `warning`, `bad`.
Go: `{{template "ui.pill" (dict "Label" "Paid" "Variant" "good")}}`. React: `<Pill variant="good">Paid</Pill>`.

### Dialog: `.ui-dialog`

Knobs: `--ui-dialog-width`, `--ui-dialog-bg`, `--ui-dialog-radius`, `--ui-dialog-pad`. It needs no JavaScript:
open it with a button carrying `command="show-modal" commandfor="<id>"`. React: `<Dialog id="confirm" title="Delete resume?">…</Dialog>`.
In Go, write this markup (it is pinned by `testdata/parity/dialog-basic.html`):

```html
<dialog class="ui-dialog" id="confirm" aria-labelledby="confirm-title">
  <header class="ui-dialog-head">
    <h2 class="ui-dialog-title" id="confirm-title">Delete resume?</h2>
    <button class="ui-dialog-x" type="button" command="close" commandfor="confirm" aria-label="Close">×</button>
  </header>
  <div class="ui-dialog-body">
    This cannot be undone.
    <div class="ui-dialog-actions">…buttons…</div>
  </div>
</dialog>
```

### Containers (plain markup in Go and React)

```html
<div class="ui-callout" data-variant="bad">          <!-- info (default), good, warning, bad -->
  <p class="ui-callout-title">Build failed</p>
  <p>The gate found 3 unsupported claims.</p>
</div>

<div class="ui-table-wrap">
  <table class="ui-table">
    <caption>October spending</caption>
    <thead><tr><th scope="col">Category</th><th scope="col" class="ui-num">Amount</th></tr></thead>
    <tbody><tr><td>Groceries</td><td class="ui-num">212.40</td></tr></tbody>
  </table>
</div>

<details class="ui-disclosure" name="faq">
  <summary class="ui-disclosure-summary">Can I export it?</summary>
  <div class="ui-disclosure-body"><p>Yes, as CSV or JSON.</p></div>
</details>

<section class="ui-card">…</section>
```

Knobs: `--ui-callout-bg`, `--ui-callout-fg`, `--ui-callout-accent`; `--ui-table-pad-y`, `--ui-table-pad-x`; `--ui-card-pad`, `--ui-card-radius`.

### Utilities

`.ui-muted` (ink-3), `.ui-mono`, `.ui-num` (tabular, end-aligned), `.ui-vh` (visually hidden),
`.ui-stack` (column, gap `--ui-stack-gap`), `.ui-cluster` (wrapping row, gap `--ui-cluster-gap`).
