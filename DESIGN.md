# ai-ui design contract

This file is the contract between ai-ui and the apps that use it. A change
that renames or removes anything named here is a breaking change (see
CHANGELOG.md).

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
   variant alone with `.ui-btn:not([data-variant])`. Set a knob on a selector that matches the
   component itself: the library sets each knob's default on the component, so a value set on an
   ancestor or on `:root` never reaches it. `.toolbar .ui-btn { --ui-btn-height: 2rem; }` works;
   `.toolbar { --ui-btn-height: 2rem; }` doesn't.
3. **Anything else**: write a plain rule. An unlayered app rule also beats the
   library's state rules (`:hover`, `[aria-invalid]`) for the same property, so
   restate the states you need.

Tokens that belong to one app are named `--app-*`. The library never defines or reads them.

## Adopting ai-ui in an app that has its own CSS

A new app needs none of this. An app with existing markup does, because `ui.reset` and
`ui.base` change every element its CSS says nothing about: list bullets, heading sizes, button
frames. A compat layer between `ui.base` and `ui.components` keeps those elements as they were
and lets the app move one component at a time. One app moved this way, and the only pixels that
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

**Sizes, colours and names.** What two apps learned keeping their own look:

- **A library size is a minimum.** A button's height is a `min-height` of `--ui-btn-height`, and
  an input's is a `min-block-size` of `--ui-input-height`, so an app's `height: 34px` on a 40px
  button loses. Name the size through the knob as well: `#theme-btn { --ui-btn-height: 34px; height: 34px; }`.
- **Hover changes colour.** A filled button turns `--primary-hover` under the pointer. If the app's
  buttons don't change, set `--ui-btn-bg-hover` and `--ui-btn-fg-hover` to their resting colours.
- **A renamed selector can wake a rule.** Retargeting `.input` to `.ui-input` changes which of the
  app's own rules win their specificity ties. A rule that used to lose to another app rule on the old
  class can start to apply, or the other way round. Before renaming, check that each retargeted
  rule was already the winner where it applies, and look at screenshots in every theme afterwards.
- **Token names are shared.** The library reads its own token names, and an unlayered app rule
  sets them. An app token with a library name and another meaning repaints the library: one app's
  `--rule` was the library's `--rule-soft`. Rename those first. When the app's tokens have other
  names, give the library's the app's values in one block:
  `:root { --ink-3: var(--app-muted); --rule: var(--app-line); --rule-soft: var(--app-hairline); }`.

**If you serve `ui.css` inline** under a hash-based CSP from Go's `html/template`, pass it as a
`template.CSS` value. Literal text in a template has its CSS comments stripped, so the bytes
served no longer match a hash taken over the source.

Delete the compat sheet once no legacy markup is left.

## Density

Controls are 40px at 16px by default. Working tools can ask for compact, on the page or on one part of it:

```html
<html data-density="compact">          <!-- the whole page -->
<div data-density="compact">…</div>    <!-- one part of a page -->
```

| | Default | Compact |
|---|---|---|
| Button | 40px, 16px text | 32px, 14px text, 12px side padding |
| Button `data-size="sm"` | 32px, 14px | 28px, 13px, 10px side padding |
| Button `data-size="lg"` | 48px, 18px | 44px, 16px, 20px side padding |
| Input, select | 40px, 16px | 32px, 14px |
| Input, select `data-size="sm"` | 32px, 14px | 28px, 13px |
| Input, select `data-size="lg"` | 48px, 16px | 44px, 16px |
| Textarea | 80px, 16px | 80px, 14px |
| Field label, hint, error | 14px | 13px |
| Field gap | 8px | 4px |
| Controls, callout, table wrap | 6px corners | 4px corners |
| Card, dialog | 10px corners, a soft shadow | 8px corners, no shadow |
| Pill | fully round | 4px corners |
| Table text | 14px | 14px |
| Body text (attribute on the page) | 16px | 14px |

On a touch screen (`pointer: coarse`), compact controls grow to finger sizes: buttons, inputs and
selects are 40px, 34px at `sm` and 44px at `lg`. The default density is already 40px there.

Headings and the spacing tokens don't change. `compact` is the only value, and a default region
nested inside a compact one isn't supported. Inside the attribute, `--text-md`, `--text-sm`,
`--radius-md` and `--radius-lg` take compact values, so an app rule that reads them follows along.
An unlayered app rule still beats all of this.

## Choosing a stylesheet

Two files are published, and a page loads one of them:

- **`ui.css`** is the normal choice. Its reset and base style the whole page.
- **`ui-scoped.css`** styles library markup only, meaning an element with a `ui-*` class and
  everything inside it. Use it on a page whose other content must stay exactly as it is, such as a
  document drawn under your nav. Don't give your own classes the `ui-` prefix, because they would
  count as library markup.

```html
<link rel="stylesheet" href="{{href "ui-scoped.css"}}">
```

Library markup looks the same under either file. With `ui-scoped.css`, the page gives up:

- **Page-level styling.** `body` gets no background, font or colour, and links outside library
  markup get no colour.
- **An inherited font.** Library markup takes its font from around it, so the surrounding app sets
  `font-family`.
- **Page-level compact.** With no `body` rule, `<html data-density="compact">` doesn't set body
  text, so put the attribute on the element that holds the library markup.

**A page with a design of its own,** a marketing home for one, can take `ui-scoped.css` and keep
its look: load it before the page's stylesheet, put library classes on the page's buttons and
fields, and let the page's unlayered rules restyle them. Map the library's token names to the
page's in one block (see "Sizes, colours and names" above), and name sizes through the knobs.

Tokens still apply to `:root`, including `color-scheme: light dark`. A page that sets no
`color-scheme` of its own follows the system's dark mode for its default colours, so set one if
the page must not.

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

- The mode follows the OS. `<html data-theme="light|dark">` forces one; `"system"` follows the OS
  again.
- `<html data-palette="mono">` selects the neutral palette. `blue` is the default and needs no attribute.
- Both apply on any element, not only `<html>`: a region, or a swatch face, can show another
  palette or mode. Custom properties inherit, so everything inside takes the palette's tokens; the
  element paints with them where its CSS uses them (`background: var(--paper); color: var(--ink)`).
- An app adds a palette in its own CSS: `[data-palette="whiteboard"] { --paper: …; }`. An app that
  changes blue writes `:root, [data-palette="blue"] { … }`, so a blue swatch face shows the change
  inside another palette.
- **A light-only palette** has no dark values. Put `color-scheme: light` in its block:
  `[data-palette="rose"] { color-scheme: light; --paper: …; }`. The app's rule is unlayered, so it
  beats `[data-theme="dark"]` and the palette stays light whatever the mode. Offer it once in a
  picker.
- `theme.js` (load it blocking in `<head>`) applies and persists the visitor's choice.
  Use `<button data-ui-set-theme="dark">`, `data-ui-set-theme="system"` or
  `<button data-ui-set-palette="mono">`; an empty value clears the choice. From script,
  use `window.aiui.theme.get()` and `.set({ theme, palette })`. Until the visitor
  chooses, the server's own `data-theme` / `data-palette` stay untouched.
- `theme.js` marks every element carrying `data-ui-set-theme` or `data-ui-set-palette`:
  `aria-pressed="true|false"`, or `aria-checked` when its role is `menuitemradio`,
  `menuitemcheckbox` or `radio`. It matches when each attribute it carries is the current value
  (an empty theme means `system`; an empty palette means the page's default). Markup added later is
  marked by `window.aiui.theme.sync(root)`.
- After every change, and when the OS flips light/dark while the mode follows it, `theme.js`
  fires `aiui:themechange` on `document`, its `detail` the new `get()`. Relabel a theme button
  there.
- **Moving an app's stored choice.** Run this before `theme.js`, mapping the app's old values:

  ```html
  <script>try{var o=localStorage.getItem("app.theme");if(o!==null&&localStorage.getItem("ai-ui:theme")===null){
  var m={"blue-dark":{theme:"dark",palette:"blue"},"mono":{theme:"light",palette:"mono"}}[o]||{};
  localStorage.setItem("ai-ui:theme",JSON.stringify(m));}localStorage.removeItem("app.theme");}catch(e){}</script>
  ```
- Load `ui.css` in `<head>` too. `theme.js` sets attributes only, so a stylesheet loaded later
  in `<body>` paints the page once in the wrong ground.

## Scripts

| Script | Load | What it does |
|---|---|---|
| `theme.js` | blocking, first in `<head>` | The visitor's mode and palette (above). |
| `ui.js` | `defer` | The menu's keys and placement fallback, the swatch row's arrows, and `closedby="any"` where the browser lacks it. It listens on `document`, so markup added later needs nothing. |

Go: `<script src="{{href "ui.js"}}" defer></script>`. React: `import "ai-ui/ui.js"` once, in a
client component. Without `ui.js`, menus and sheets still open, close on Escape and an outside
click, and return focus; only the keys are missing.

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

| Option | Markup | Go | React |
|---|---|---|---|
| Outlined | `data-tone="outline"` | `Tone "outline"` | `tone="outline"` |
| Icon-only | `data-shape="square"` | `Shape "square"` | `shape="square"` |
| Full width | `data-width="full"` | `Width "full"` | `width="full"` |
| Icon | an `<svg>` before the label | `Icon` | `icon` |

- **Outlined** draws a primary or danger button on the sheet with a `--rule` border, its text in
  `--primary` or `--bad`. Use it for a destructive action in a row, and keep the filled danger
  button for a dialog's confirm. Secondary and ghost are unchanged by it.
- **Icon-only** is as wide as it is tall. It needs both `Icon` and `Label`: the label becomes its
  `aria-label` and the icon its only content. In Go, a square button without `Icon` is an error; in
  React, a square button's children must be text.
- **Full width** fills its container, and a long label wraps.
- **Icon** is the app's own markup. In Go, pass a `template.HTML` value, from the app's own
  templates or code, never from user input; a plain string is escaped and shows as text. An `<svg>`
  in a button is sized to the text.

A disabled button is drawn in `--band` with `--ink-3` text and a `--rule` border, at full opacity,
so its label stays readable. A disabled ghost button keeps no fill.

Go args: `Label` (required), `Variant`, `Size`, `Tone`, `Shape`, `Width`, `Icon`, `Type` (`submit`,
`reset`; default `button`), `Name`, `Value`, `Command`, `CommandFor`, `Disabled`, and `Href`, which
makes it a link. **A button that
submits its form needs `"Type" "submit"`:** the default is `button`, which submits nothing.
React: `Button` passes every other `<button>` prop through, and its `type` also defaults to
`button`. `ButtonLink` passes every `<a>` prop through.
`Disabled` on a link (`disabled` on `ButtonLink`) drops `href` and sets `aria-disabled="true"`.
Labels are text.

```html
<button class="ui-btn" type="button" data-variant="secondary" data-size="sm">Save</button>
```

### Field: `.ui-field`

Knobs: `--ui-field-gap`; on controls, `--ui-input-bg`, `--ui-input-fg`, `--ui-input-border`, `--ui-input-radius`, `--ui-input-height`, `--ui-input-width` (default 100%).

Go args: `Label` (required), `ID` or `Name` (at least one), `Control` (`textarea`, `select`; default input), `Type` (`email`, `password`, `number`, `search`, `tel`, `url`, `date`; default `text`), `Size` (`sm`, `lg`), `Icon` (input only), `Value`, `Placeholder`, `Autocomplete`, `Rows`, `Options` (`[]ui.Option`), `Hint`, `Error`, `Required`, `Disabled`.
React: `<Field label name id control type size icon options defaultValue hint error … />`. React's
`size` is the library's size, `sm` or `lg`, in place of the native numeric `size` attribute. Every other prop, including `ref`, goes to the control,
so a controlled field takes `value` and `onChange`. `options` is `{ value, label }[]` in React and
`[]ui.Option{{Value: "m", Label: "Monthly"}}` in Go; the option whose value equals `Value` is selected.

The controls are `.ui-input`, `.ui-select` and `.ui-textarea`. A textarea is twice the input
height at least, and `Rows` (`rows`) makes it taller.

```html
<div class="ui-field">
  <label class="ui-field-label" for="email">Email</label>
  <input class="ui-input" id="email" name="email" type="email" aria-describedby="email-hint email-error" aria-invalid="true" required>
  <p class="ui-field-hint" id="email-hint">We never share it.</p>
  <p class="ui-field-error" id="email-error">That address bounced.</p>
</div>
```

**Sizes.** `data-size="sm"` or `"lg"` on a control (`Size`, `size`) makes it as tall as the
button of that size, at either density.

**A leading icon.** `Icon` (`icon`) puts an icon inside an input, before its text:

```html
<div class="ui-input-group">
  <svg class="ui-input-icon" aria-hidden="true">…</svg>
  <input class="ui-input" type="search" aria-label="Search merchants">
</div>
```

The icon is 1rem in `--ink-3` and takes no clicks. Knob on the group: `--ui-input-icon-inset`.

Lay fields side by side with `<div class="ui-form-row">`. Each child grows from 14rem; set
`--ui-row-grow` and `--ui-row-basis` on one child to change that, for a narrow field beside a wide
one: `<div style="--ui-row-grow: 0; --ui-row-basis: 6rem">…port…</div>`.

### Check: `.ui-check`

Knobs: `--ui-check-size`, `--ui-check-accent`.
Go: `{{template "ui.check" (dict "Label" "Monthly" "Name" "period" "Value" "m" "Type" "radio" "Checked" true)}}`. React: `<Check label="Monthly" name="period" value="m" type="radio" defaultChecked />`.

Go args: `Label` (required), `Type` (`radio`; default `checkbox`), `ID`, `Name`, `Value`, `Checked`,
`Required`, `Disabled`. React: every other prop goes to the `<input>`, so a controlled check takes
`checked` and `onChange`.

```html
<label class="ui-check"><input class="ui-check-input" type="radio" name="period" value="m" checked><span class="ui-check-label">Monthly</span></label>
```

### Pill: `.ui-pill`

Knobs: `--ui-pill-bg`, `--ui-pill-fg`, `--ui-pill-border`. Variants: neutral (default), `info`, `good`, `warning`, `bad`.
Go: `{{template "ui.pill" (dict "Label" "Paid" "Variant" "good")}}`. React: `<Pill variant="good">Paid</Pill>`.

- **Dot:** `data-dot` (Go `Dot true`, React `dot`) draws a small dot in the pill's text colour
  before the label, for a status.
- **Outlined:** `data-tone="outline"` (`Tone "outline"`, `tone="outline"`) drops the fill; the
  border takes the text colour, or `--rule-field` on a neutral pill. `info` is already outlined.

Go args: `Label` (required), `Variant`, `Tone`, `Dot`.

### Dialog: `.ui-dialog`

Knobs: `--ui-dialog-width`, `--ui-dialog-bg`, `--ui-dialog-radius`, `--ui-dialog-pad`, `--ui-dialog-shadow`, `--ui-dialog-border`, `--ui-dialog-eyebrow-fg`. It needs no JavaScript:
open it with a button carrying `command="show-modal" commandfor="<id>"`. React: `<Dialog id="confirm" title="Delete resume?">…</Dialog>`;
`closeLabel` (default `Close`) is the × button's accessible name, `closeButton={false}` leaves the ×
out, and `eyebrow` adds a small line above the title.
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

`.ui-dialog-actions` is the row of buttons at the end, and it can be the `<form>` itself. The title
wraps anywhere it must, so a long email address in it stays inside the dialog.

**An eyebrow** is a `<p class="ui-dialog-eyebrow">` as the dialog's first child, before the head:
small uppercase text in `--ink-3`, such as "Can't be undone". Set `--ui-dialog-eyebrow-fg: var(--bad)`
on the dialog for red.

**The × is optional.** Leave its button out of the head. The dialog must still close, so give it a
Cancel button with `command="close" commandfor="<id>"`.

**Shown without script.** A `<dialog open>` lays out as a dialog in the page's flow, with no
backdrop, so a confirm can be served open when script is off. Placing it is the app's. A button
with `command="close" commandfor="<id>"` closes the dialog without script.

**A confirm per row.** In Go, give each row its own dialog with its own id
(`printf "remove-%d" $i`), and put the action's form in it, so it works without script:

```html
<form class="ui-dialog-actions" method="post" action="/remove">
  <input type="hidden" name="email" value="{{.Email}}">
  {{template "ui.button" (dict "Label" "Cancel" "Variant" "secondary" "Command" "close" "CommandFor" (printf "remove-%d" $i))}}
  {{template "ui.button" (dict "Label" "Remove" "Variant" "danger" "Type" "submit")}}
</form>
```

In React, keep one dialog and the row it is for in state: each row's opener sets it, and the
confirm button acts on it and closes the dialog with `command="close"`.

### Menu: `.ui-menu`

A popover anchored to the button that opens it. The browser gives light dismiss, Escape, focus
return, `aria-expanded` and the anchor; `ui.js` adds the keys. Go writes this markup (pinned by
`testdata/parity/menu-basic.html`); React: `<Menu id="acct" label="Account">` with `MenuItem`,
`MenuLink`, `MenuLabel` and `MenuSeparator`, and `<Button popoverTarget="acct">` to open it.

```html
<button class="ui-btn" type="button" popovertarget="acct" aria-haspopup="menu">Account</button>
<div class="ui-menu" id="acct" popover role="menu" aria-label="Account">
  <p class="ui-menu-label" role="none">grace@example.com</p>
  <a class="ui-menu-item" role="menuitem" href="/account/">Account</a>
  <hr class="ui-menu-sep">
  <form method="post" action="/logout" role="none"><button class="ui-menu-item" role="menuitem" type="submit">Sign out</button></form>
  <button class="ui-menu-item" role="menuitem" type="button" data-variant="danger">Delete</button>
</div>
```

| Part | What it is |
|---|---|
| `.ui-menu-item` | A row, `<button>` or `<a>`, `role="menuitem"`. `data-variant="danger"` paints it `--bad`. `role="menuitemradio"` with `aria-checked` draws a check when checked. A leading `svg` or `.ui-swatch-face` is 1em. |
| `.ui-menu-label` | A caption or identity row: small `--ink-3`, one line with an ellipsis, `role="none"`. |
| `.ui-menu-sep` | A hairline: `<hr>`. |
| `.ui-menu-head`, `.ui-menu-title`, `.ui-menu-x` | A head shown only when the menu is a sheet on a phone (below). |

| Menu option | Attribute | What it does |
|---|---|---|
| Start-aligned | `data-align="start"` | Lines up the menu's start edge with the opener's. By default the end edges line up. Either flips to fit. |
| A sheet on a phone | `data-sheet="phone"` | Below 800px the menu is a bottom sheet over a `--scrim` backdrop, rows 44px, and its head shows. |

A phone-sheet menu's head, first in the menu:

```html
<div class="ui-menu-head" role="none"><span class="ui-menu-title">Theme</span><button class="ui-menu-x" type="button" role="menuitem" popovertarget="theme" popovertargetaction="hide" aria-label="Close">×</button></div>
```

With `ui.js`: opening focuses the checked item, else the first; Down/Up move and wrap; Home/End go
to the ends; Tab closes the menu and moves on from the opener; activating a `menuitem` closes it,
and a `menuitemradio` or swatch keeps it open, so a theme can be tried in place.

Knobs: `--ui-menu-width` (`16rem`), `--ui-menu-bg`, `--ui-menu-border`, `--ui-menu-radius`,
`--ui-menu-shadow`, `--ui-menu-item-height` (`2.25rem`; 34px in compact; 44px on a coarse pointer),
`--ui-menu-hover`. A menu keeps its shadow in compact: it floats over content.

### Sheet: `.ui-sheet`

A modal `<dialog>` docked to the bottom edge, or with `data-side="end"` to the inline end as a
drawer. Open it like a dialog, with `command="show-modal" commandfor="<id>"`; the page behind is
inert, Escape closes it and focus returns. `closedby="any"` closes it on a backdrop tap as well.
Go writes this markup (pinned by `testdata/parity/sheet-basic.html`); React: `<Sheet id="facts"
title="Facts" dismissible>…</Sheet>`, with `side="end"`, `closeLabel` and `closeButton={false}`.

```html
<dialog class="ui-sheet" id="facts" aria-labelledby="facts-title" closedby="any">
  <header class="ui-sheet-head">
    <h2 class="ui-sheet-title" id="facts-title">Facts</h2>
    <button class="ui-sheet-x" type="button" command="close" commandfor="facts" aria-label="Close">×</button>
  </header>
  <div class="ui-sheet-body">…</div>
</dialog>
```

At the bottom it is full width and at most `--ui-sheet-max-height` (`70dvh`) tall; at the end it is
full height and `--ui-sheet-width` (`28rem`) wide, never wider than the viewport. The body scrolls.
A modal sheet slides in from its edge unless reduced motion is asked for; one drawn in place
(`open`, not modal) doesn't move.

Knobs: `--ui-sheet-bg`, `--ui-sheet-border`, `--ui-sheet-radius`, `--ui-sheet-shadow`,
`--ui-sheet-pad`, `--ui-sheet-max-height`, `--ui-sheet-width`. Flat in compact.

### Theme picker: `.ui-swatches`

A row of swatches, each a button carrying a choice for `theme.js`, which marks the current one
pressed (a 2px `--ink` ring). Go writes this markup (pinned by `testdata/parity/swatches.html`);
React: `<Swatches label="Theme"><Swatch theme="dark" palette="mono" label="Black and white dark" /></Swatches>`.

```html
<div class="ui-swatches" role="group" aria-label="Theme">
  <button class="ui-swatch" type="button" data-ui-set-theme="system" data-ui-set-palette="blue" aria-label="System"><span class="ui-swatch-face" data-theme="system" data-palette="blue" aria-hidden="true"></span></button>
  <button class="ui-swatch" type="button" data-ui-set-theme="dark" data-ui-set-palette="mono" aria-label="Black and white dark"><span class="ui-swatch-face" data-theme="dark" data-palette="mono" aria-hidden="true"></span></button>
</div>
```

- A face (`.ui-swatch-face`) previews its `data-palette` and `data-theme` from the tokens: the
  palette's `--paper` with a `--primary` corner. `data-theme="system"` splits it, light over dark.
- Swatches share the row equally, 2rem tall (2.5rem on a coarse pointer). With `ui.js`,
  Left/Right/Home/End move between them.
- In a menu, a face can lead a `menuitemradio` row instead (React `<MenuItem checked icon={<SwatchFace … />}>`),
  with `data-ui-set-*` on the row.

### Page parts: `.ui-page-head`

Plain markup in Go and React: a page's head, eyebrows, section heads, stat tiles and empty states.
The page's width, gutters and the space between sections stay the app's.

```html
<header class="ui-page-head">
  <div><span class="ui-eyebrow">Settings</span><h1 class="ui-page-title">Account</h1><p class="ui-page-sub">Your sign-in and your theme.</p></div>
  <div class="ui-page-aside">Basis: 12 months …</div>
</header>
<div class="ui-section-head"><span class="ui-section-index">01</span><h2 class="ui-section-title">Cash</h2><span class="ui-section-count">4 accounts</span><div class="ui-section-aside">…</div></div>
<dl class="ui-stats">
  <div class="ui-stat"><dt>Spent</dt><dd class="ui-stat-value">1,234.56</dd><dd class="ui-stat-note" data-variant="bad">+412.86 vs usual</dd></div>
  <div class="ui-stat" data-variant="good"><dt>Left</dt><dd class="ui-stat-value">1,865.44</dd></div>
</dl>
<div class="ui-empty"><p class="ui-empty-title">No archived resumes</p><p>Archive one from its menu and it waits here.</p></div>
```

| Page part | What it is |
|---|---|
| `.ui-page-head` | The page's title row: its first child (eyebrow, `.ui-page-title`, `.ui-page-sub`) and an `.ui-page-aside` at the end, on the title's last line, with a `--rule` under both. The aside wraps under the title when there's no room. |
| `.ui-eyebrow` | A small uppercase `--ink-3` label above a title. Knobs `--ui-eyebrow-fg`, `--ui-eyebrow-case`. |
| `.ui-section-head` | A section's title row: `.ui-section-index` ("01", `--primary`), `.ui-section-title`, `.ui-section-count`, and `.ui-section-aside` pushed to the end. |
| `.ui-stats` | A `<dl>` of equal tiles with hairlines between them. Each `.ui-stat` holds a `<dt>` label, a `.ui-stat-value` (tabular) and an optional `.ui-stat-note`. `data-variant="good|warning|bad|muted"` on a tile colours its value; on a note, the note. `data-frame="none"` drops the frame, leaving hairlines between tiles. |
| `.ui-empty` | Nothing here yet: a dashed frame, centred, with `.ui-empty-title`, a sentence and maybe a button. `data-frame="none"` makes it flush and start-aligned. |

Knobs: `--ui-page-title-size` (28px on a phone to 38px wide), `--ui-page-head-rule`;
`--ui-section-rule`, `--ui-section-title-size`, `--ui-section-index-fg`; `--ui-stats-min` (a tile's
narrowest, 9rem), `--ui-stat-value-size`, `--ui-stat-label-case`.

### Chips: `.ui-chips`

Plain markup in Go and React: a row of chips, the segmented control, the count badge and the avatar.

```html
<nav class="ui-chips" aria-label="Scope">
  <a class="ui-chip" href="/claims/" aria-current="page">All resumes <span class="ui-badge">23</span></a>
  <a class="ui-chip" href="/claims/?r=acme">acme-lead</a>
  <button class="ui-chip" type="button" aria-pressed="false">Needs review</button>
</nav>
<div class="ui-segmented" role="group" aria-label="Month or year">
  <button class="ui-segment" type="button" aria-pressed="true">Month</button>
  <button class="ui-segment" type="button" aria-pressed="false" aria-disabled="true">Year</button>
</div>
```

| Control | What it is |
|---|---|
| `.ui-chip` | A filter or scope, a link or a button, 28px (40px on a coarse pointer). The current one, `aria-current` on a link or `aria-pressed="true"` on a toggle, gets a `--primary` border on `--mark-bg`. A leading `svg` is 1em. Knob `--ui-chip-height`. |
| `.ui-badge` | A count beside a label, on a chip, a tab or a link. Knobs `--ui-badge-bg`, `--ui-badge-fg`. |
| `.ui-segmented`, `.ui-segment` | One choice of a few: buttons with `aria-pressed`, or links with `aria-current`. The chosen one is filled `--primary`; `aria-disabled="true"` greys one and keeps it focusable. Knob `--ui-segment-height`. |
| `.ui-avatar` | A person's initial in a `--mark-bg` circle. Inside a square button (Go `Shape "square"` with `Icon` set to the avatar's markup) it fills the button with no circle. Knob `--ui-avatar-size`. |

### Meter: `.ui-meter`

A native `<progress>`, so it is a progressbar to assistive tech. Give it a name.

```html
<progress class="ui-meter" value="40" max="100" aria-label="Import">40%</progress>
<progress class="ui-meter" aria-label="Waiting for the server"></progress>
```

With no `value` it is indeterminate: a segment slides along it, unless reduced motion is asked for.
`data-variant="good|warning|bad"` colours the fill. Knobs: `--ui-meter-height` (0.5rem),
`--ui-meter-track` (`--band`), `--ui-meter-fill` (`--primary`).

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

Knobs: `--ui-callout-bg`, `--ui-callout-fg`, `--ui-callout-accent`; `--ui-table-pad-y`, `--ui-table-pad-x` and the table knobs below; `--ui-card-pad`, `--ui-card-radius`, `--ui-card-shadow`, `--ui-card-border`; `--ui-disclosure-rule`.

A card has a soft shadow by default and none in compact density. `--ui-card-shadow: none` makes any
card flat.

**Tables** take four options, all plain markup:

| Option | Markup | What it does |
|---|---|---|
| No frame | `data-frame="none"` on `.ui-table-wrap` | No border, corners or fill; the first and last columns sit flush with the content around them. |
| Quiet head | `data-head="quiet"` on `.ui-table` | Header cells on no fill, in small `--ink-3` text. Row headers in the body (`<th scope="row">`) go quiet too, but keep their case. |
| Sticky head | `data-sticky="head"` on `.ui-table` | Header cells stay at the top of the wrap as it scrolls. Bound the wrap with `--ui-table-max-height`. |
| Stacked on a phone | `data-stack="phone"` on `.ui-table-wrap` | Below 800px of viewport, each row is a block and each cell a line with its label. |

Table knobs, on `.ui-table`: `--ui-table-head-bg`, `--ui-table-head-fg`, `--ui-table-head-size`,
`--ui-table-head-case` (`text-transform`), `--ui-table-head-rule` (the line under the header),
`--ui-table-hover` (`transparent` turns the row tint off), `--ui-table-sticky-bg` (default
`--sheet`; set it to whatever is behind the table), `--ui-table-label-width` (default 8rem),
`--ui-table-selected` (default `--mark-bg`). On `.ui-table-wrap`: `--ui-table-max-height`. A row is
tinted on hover only where the pointer can hover. The current row, `<tr aria-current="true">` (or
`aria-selected="true"` in a grid), keeps a `--ui-table-selected` tint, under the pointer too.

A stacked table labels each cell from its `data-label`. A cell without one, such as the row's name,
spans the row, and so does a row header (`<th scope="row">`). The header row is hidden from sight
but stays in the markup.

```html
<div class="ui-table-wrap" data-stack="phone">
  <table class="ui-table">
    <thead><tr><th scope="col">Email</th><th scope="col">Role</th><th scope="col" class="ui-num">Logins</th></tr></thead>
    <tbody>
      <tr><td>ada@example.com</td><td data-label="Role">owner</td><td data-label="Logins" class="ui-num">12</td></tr>
    </tbody>
  </table>
</div>
```

Some browsers stop treating a table as a table once its rows are blocks, so a screen reader may lose
the column names on a phone. If that matters, add the ARIA table roles to the markup: `role="table"`
on the table, `rowgroup`, `row`, `columnheader` and `cell`.

`data-marker="start"` on a disclosure puts its chevron before the summary text, pointing to the
end when closed and down when open, for a list of folders or groups. `--ui-disclosure-rule:
transparent` drops the line under it.

A callout that reports what just happened takes `role="status"`, and one that reports an error
takes `role="alert"`, so a screen reader announces it. An error about one field belongs in that
field's `Error` (`error`) instead.

### Utilities

`.ui-muted` (ink-3), `.ui-mono`, `.ui-num` (tabular, end-aligned), `.ui-vh` (visually hidden),
`.ui-stack` (column, gap `--ui-stack-gap`), `.ui-cluster` (wrapping row, gap `--ui-cluster-gap`).
