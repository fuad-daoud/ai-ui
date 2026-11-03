# ai-ui — shared UI library design

Date: 2026-11-03
Status: draft, awaiting review

## 1. Purpose

One UI library for Fuad's own platforms, so that they:

- **stop rebuilding the same UI** — buttons, forms, tables, dialogs, nav are written once;
- **look like one family** — siblings from the same maker, with no silent drift;
- **start fast** — a new platform gets a finished look on day one.

All three goals carry equal weight. The platforms should look very similar,
and every app must still be able to override what it needs.

### In scope (consumers)

| App | Stack | Repo |
|---|---|---|
| money | Go, `html/template`, CSS embedded from `static/css` | github.com/fuad-daoud/money |
| career / career-platform | Go, `html/template`, CSS mostly inline in each page | github.com/fuad-daoud/career |
| relevo-site | Go, static HTML + `style.css` | github.com/fuad-daoud/relevo-site |
| unassisted (`forgot-how-code`) | Next.js, React 19, Tailwind v4 | github.com/fuad-daoud/unassisted |

### Out of scope

- spaceapi (it uses the separate `@dexpace` registry), uniqueperfumesjo, movies.
- Charts and data-visualisation colours (money's category colours stay app tokens).
- An icon set.
- The resume PDF / print template in career.

## 2. Decisions

| # | Decision | Chosen | Rejected |
|---|---|---|---|
| D1 | Architecture | **CSS core + thin Go and React adapters** that render identical markup | CSS-only (markup still rewritten per app); Web Components (shadow DOM breaks cascade overrides; JS-rendered markup causes FOUC on server-rendered pages) |
| D2 | Override model | **Layered CSS**: the library lives in cascade layers; apps theme through tokens and adjust components through `--ui-*` knobs; any unlayered app rule wins | tokens-only (too rigid); copy-and-own/eject (updates stop flowing) |
| D3 | Distribution | **Git tags, no registry**: Go via `go get`, npm via a `github:` dependency | GitHub Packages (publish workflow + install tokens); copy script (drift) |
| D4 | Visibility | **Public** repo `github.com/fuad-daoud/ai-ui` — no GOPRIVATE, no credentials where apps build | private |
| D5 | Behaviour | **Native HTML first** (`<dialog>`, `popover`, `<details>`, invoker commands); JS only where the platform has no answer | a JS component runtime |
| D6 | Token source | **Hand-written `css/tokens.css`**, documented in `DESIGN.md`; no token build pipeline | JSON → CSS generator (no consumer needs one yet) |
| D7 | Browsers | Evergreen, **Baseline 2025 and newer** (`light-dark()`, `@layer`, `:has()`, `popover`, `<dialog>`, invoker commands); CSS anchor positioning is used with a fallback (§4.3) | older-browser support |

## 3. Tokens and theming

### 3.1 Vocabulary

The vocabulary comes from the one career and money already share. It is
cleaned up and uses oklch throughout.

| Group | Tokens |
|---|---|
| Surfaces | `--paper` (page), `--band` (secondary strip), `--sheet` (card / panel) |
| Text | `--ink`, `--ink-2`, `--ink-3` |
| Lines | `--rule`, `--rule-soft`, `--rule-field` (input borders) |
| Action | `--primary`, `--on-primary`, `--primary-hover`, `--link`, `--focus`, `--accent` |
| Status | `--good`, `--good-bg`, `--warning`, `--warning-bg`, `--bad`, `--bad-bg` |
| Effects | `--mark-bg`, `--shadow`, `--scrim` |
| Type | `--font-sans`, `--font-mono`, `--font-display`; scale `--text-xs`, `--text-sm`, `--text-md`, `--text-lg`, `--text-xl`, `--text-2xl`, `--text-3xl` |
| Shape | `--space-1` … `--space-8`, `--radius-sm`, `--radius-md`, `--radius-lg` |
| Motion | `--dur`, `--ease` |

Font defaults: `--font-sans` and `--font-display` default to the system UI
stack, and `--font-mono` to the system monospace stack. The optional bundled
IBM Plex files (§5) come with an `@font-face` sheet, `plex.css`, which an app
links and then sets `--font-sans: "IBM Plex Sans", …`. Each app sets its own
fonts: money keeps Plus Jakarta Sans and relevo-site keeps IBM Plex.

### 3.2 Naming rules

1. **Library tokens are unprefixed** (`--ink`). An app themes itself by setting them.
2. **Component knobs are `--ui-<component>-<property>`** (`--ui-btn-radius`). Each one defaults to a token.
3. **Tokens that belong to one app are `--app-*`** (money's `--app-desk`, `--app-marker`). The library never defines or reads `--app-*`.

### 3.3 Light and dark

Each colour token is declared once with `light-dark()`:

```css
:root { color-scheme: light dark; --ink: light-dark(oklch(21% .03 262), oklch(96% .01 262)); }
:root[data-theme="light"] { color-scheme: light; }
:root[data-theme="dark"]  { color-scheme: dark; }
```

By default the mode follows the OS, and `data-theme` on `<html>` forces one.
This replaces the duplicated `blue-dark-*` / `mono-dark-*` lists in career and
money. An app override is one line that covers both modes.

### 3.4 Palettes

- Built in: **`blue`** (the default, active when no `data-palette` is set) and **`mono`**. Their values come from career's current DESIGN.md.
- Selected with `<html data-palette="mono">`.
- An app adds a palette by declaring tokens under its own selector, for example `:root[data-palette="whiteboard"] { … }` in unassisted.

### 3.5 `theme.js`

`theme.js` is the only JavaScript an app needs for theming, about 1 KB. It does two things:

- **Pre-paint**: an inline-able snippet that runs in `<head>`, reads the stored choice and sets `data-theme` / `data-palette` before first paint, so the page never flashes the wrong theme.
- **Persistence**: it wires the theme switcher (§4.4) to `localStorage` (wrapped in try/catch) and keeps `data-theme` / `data-palette` in sync.

### 3.6 Tailwind mapping

`react/theme.css` maps the tokens into Tailwind v4:

```css
@layer theme, base, ui, components, utilities;
@theme inline {
  --color-paper: var(--paper); --color-sheet: var(--sheet);
  --color-ink: var(--ink); --color-ink-2: var(--ink-2); /* … every colour token … */
  --font-sans: var(--font-sans); --font-mono: var(--font-mono);
}
```

Utilities such as `bg-sheet`, `text-ink-2` and `border-rule` then resolve to
the same variables the components use.

## 4. Component CSS and the override contract

### 4.1 Layers

`dist/ui.css` puts every rule inside these layers:

```css
@layer ui.reset, ui.base, ui.components, ui.utilities;
```

**The contract:** any unlayered app rule beats every library rule, with no
`!important`. In a Tailwind app the layer order from §3.6 puts `ui` below
`components` and `utilities`, so `class="ui-btn rounded-none"` works.

### 4.2 Markup conventions

- **One root class per component:** `.ui-btn`, `.ui-field`, `.ui-callout`, and so on. Parts are `.ui-<component>-<part>` (`.ui-field-hint`).
- **Variants and sizes are data attributes** with an enumerated value set: `data-variant="ghost"`, `data-size="sm"`.
- **State comes only from native HTML and ARIA:** `:disabled`, `[aria-invalid="true"]`, `[aria-current]`, `[aria-busy="true"]`, `:open`, `:checked`. There are no state classes.
- **Knobs:** each component reads its own `--ui-*` properties, and those default to tokens:

```css
@layer ui.components {
  .ui-btn {
    --ui-btn-bg: var(--primary);
    --ui-btn-fg: var(--on-primary);
    --ui-btn-radius: var(--radius-md);
    background: var(--ui-btn-bg); color: var(--ui-btn-fg); border-radius: var(--ui-btn-radius);
  }
}
```

An app changes one component everywhere with `.ui-btn { --ui-btn-radius: 0; }`.

### 4.3 Behaviour mechanisms

| Need | Mechanism | JS |
|---|---|---|
| Modal dialog | `<dialog>` opened by `command="show-modal"` / `commandfor` | none |
| Drawer / side sheet | `<dialog>` with `data-variant="drawer"` | none |
| Menu / popover | `popover` attribute; CSS anchor positioning, falling back to a fixed placement below the trigger where anchors are unsupported | none |
| Disclosure / accordion | `<details name="…">` | none |
| Toast | `popover="manual"`, shown and hidden by `ui.js` | `ui.js` |
| Tabs (ARIA tablist, arrow keys) | markup works as plain anchored sections without JS; `ui.js` upgrades it | `ui.js` |
| Copy button | `ui.js` | `ui.js` |

`ui.js` is optional. Every component it enhances stays usable without it.

### 4.4 v1 component set

| Group | Components |
|---|---|
| Actions | button (`primary` default, `secondary`, `ghost`, `danger`; sizes `sm`, `md` default, `lg`), link |
| Forms | field (label + control + hint + error, wired with ids and `aria-describedby`), input, select, textarea, checkbox, radio, switch, form-row |
| Display | card, callout (`info` default, `good`, `warning`, `bad`), pill, chip, badge, table (with `.ui-num` numeric columns), meter, empty state, skeleton, spinner |
| Overlays | dialog, drawer, menu, disclosure, toast |
| Navigation | tabs, segmented control, page shell (header, nav, footer, container, section heading), theme switcher |
| Base | typography and prose defaults; utilities `.ui-muted`, `.ui-mono`, `.ui-num`, `.ui-vh` (visually hidden), `.ui-stack`, `.ui-cluster` |

## 5. Repository layout

```
ai-ui/
  go.mod                  module github.com/fuad-daoud/ai-ui
  ui.go                   package ui — the Go adapter
  css/
    tokens.css            §3
    reset.css base.css utilities.css
    components/*.css      one file per component
  dist/ui.css             concatenation of css/ in layer order — generated, committed
  js/theme.js js/ui.js
  fonts/                  IBM Plex Sans + IBM Plex Mono, latin subset, woff2, + plex.css (optional)
  templates/*.tmpl        Go partials, {{define "ui.<name>"}}
  package.json            name "ai-ui"; exports ./ui.css ./theme.css ./theme.js ./ui.js ./react
  react/
    index.ts *.tsx theme.css
  testdata/parity/        fixtures (*.json) + golden HTML (*.html)
  gallery/                main.go: renders every component × palette × mode
  cmd/build/              generates dist/ui.css from css/
  DESIGN.md               the token and component contract for humans and agents
  CHANGELOG.md
```

Plain CSS needs no build for consumers. `cmd/build` exists only to
concatenate the files into one request, and CI checks that the output is fresh.

## 6. Go adapter (`package ui`)

```go
// Assets holds the published files under their public names (ui.css, theme.js, …).
var Assets fs.FS

// Handler serves Assets; mount it at Prefix. A request whose ?v= matches the
// asset's hash gets Cache-Control: public, max-age=31536000, immutable; any other
// request gets no-cache, so a link written without Href is never pinned stale.
func Handler() http.Handler
const Prefix = "/_ui/"

// Href returns Prefix+name+"?v="+<first 12 hex chars of the sha256 of that asset>.
// It panics on an unknown name, so a typo fails at the first render in tests.
func Href(name string) string

// Parse adds the ui.* partials and the helper funcs (dict, href) to t.
// The partials call only ui_* helpers, so an app that later redefines
// dict or href does not break them.
func Parse(t *template.Template) (*template.Template, error)

// Option is one <option> of a select field.
type Option struct{ Value, Label string }
```

Usage in an app:

```go
mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
tmpl := template.Must(ui.Parse(template.New("app").Funcs(appFuncs)))
tmpl = template.Must(tmpl.ParseFS(pages, "*.html"))
```

```html
<link rel="stylesheet" href="{{href "ui.css"}}">
{{template "ui.field" (dict "Label" "Email" "Name" "email" "Type" "email" "Error" .Errs.Email)}}
```

**Which components get a partial.** Only leaf components, whose markup is
more than a root class and which have no child content: button (and its
link form), field (input, select, textarea), check (checkbox, radio), pill;
later menu, tabs, segmented control, toast and theme switcher. Anything that
wraps child content stays plain markup with classes: card, callout, table,
disclosure, the page shell, the empty state, and also dialog and drawer.
`html/template` cannot pass a child block to a partial, and open/close
partial pairs would make the adapter heavier than the CSS it wraps. A
container's exact markup is documented in `DESIGN.md` and pinned by a golden
in `testdata/parity/` (§8.1).

**Argument errors.** `dict` returns an error on an odd number of arguments or
a non-string key. A partial given an unknown `Variant` or `Size` value renders
the default and adds nothing.

## 7. React adapter

- **Install:** `"ai-ui": "github:fuad-daoud/ai-ui#v0.1.0"`, plus `transpilePackages: ["ai-ui"]` in the Next config. React 19 is a peer dependency. There is no build step: the TSX ships as source.
- **Styles:** `import "ai-ui/ui.css"` and `import "ai-ui/theme.css"` in the root stylesheet or layout.
- **Components:** `ai-ui/react` exports one component per Go partial (`Button`, `ButtonLink`, `Field`, `Check`, `Pill`, later `Menu`, `Tabs`, `SegmentedControl`, `Toast`, `ThemeSwitcher`), plus components for the containers whose markup has wiring to get right (`Dialog`, later `Drawer`).
- **Rules for every component:**
  - It renders exactly the markup of its Go partial or its documented golden (enforced by §8.1).
  - It forwards `className` (appended after the library class) to the root element. It spreads the rest of the props, including `ref`, onto the element that carries the component's behaviour: the root for most components, the form control for `Field` and `Check`.
  - It contains no styling logic, only class and attribute mapping.
  - Only components that need effects (`Toast`, `Tabs`, `ThemeSwitcher`) are `'use client'`.
- Containers stay plain JSX with classes, the same as in Go.

## 8. Testing

### 8.1 Markup parity

Each case in `testdata/parity/` is a fixture `<case>.json` with a `go` part (`{"template": "ui.field", "args": {…}}`), a `react` part (`{"component": "Field", "props": {…}}`), or both, and a golden `<case>.html`. A case may instead name an `error` that both sides must raise.

- A Go test renders each `go` part through its partial.
- A Node test (vitest) renders each `react` part with `renderToStaticMarkup`.
- Both sides put the output and the golden into one canonical form and compare them. The canonical form: parsed as a `<body>` fragment, attributes sorted, empty-valued attributes printed bare, `class` whitespace normalised, text trimmed with HTML-whitespace runs collapsed, whitespace-only text and comments dropped, and one escaping. The two canonicalisers (Go and TypeScript) are pinned by shared vectors in `testdata/canon/`.
- Goldens for Go cases are regenerated only with an explicit `-update` flag. Goldens for container markup (react-only cases) are written by hand, and the gallery's Go-rendered dialog is checked against the same golden.

### 8.2 Gallery, visual regression and accessibility

- `go run ./gallery` serves one page per palette × mode with every component in every variant, size and state.
- Playwright screenshots each page, which gives visual regression, and runs axe on each, which fails on any violation.
- Like all Node and e2e work in this setup, this runs on contabo through `dev run`. It is not part of CI (§8.4): screenshots are only comparable on the machine that made the baselines.

### 8.3 Contrast

A Go test parses `css/tokens.css`, resolves each `light-dark()` pair, and
checks WCAG contrast at 4.5:1 or better in both modes and every built-in palette for:

- `ink`, `ink-2` and `ink-3` on `paper`, `band` and `sheet`;
- `on-primary` on `primary`, on `primary-hover` and on `bad` (the danger button);
- `link` on `paper`, `band` and `sheet`;
- `good`, `warning` and `bad` on their `-bg`.

`rule-field` must reach 3:1 or better on `sheet` (non-text contrast).

### 8.4 CI

GitHub Actions runs on every push:

- `go test ./...`
- the vitest parity suite
- `go run ./cmd/build` followed by `git diff --exit-code dist/` (the freshness check)

## 9. Versioning

- One semver tag (`vX.Y.Z`) covers both the Go module and the npm package. `package.json` `version` matches the tag, and CI checks this on tag pushes.
- **Breaking change:** renaming or removing a token, class, data-attribute value, knob, partial name or partial argument, or a React prop. A token keeps its name for its whole life. Only its value changes.
- During `0.x`, a breaking change bumps the minor version. From `1.0`, it bumps the major.
- `CHANGELOG.md` gets an entry for every tag, with breaking changes listed first.

## 10. Rollout

| Step | Work | Notes |
|---|---|---|
| 1 | **ai-ui v0.1.0** | tokens (blue + mono, light + dark), reset/base/utilities, button, field + controls, callout, pill, table, dialog, disclosure; Go adapter; React adapter for the same set; parity, contrast and gallery tests; CI |
| 2 | **money** | closest tokens and CSS already separate; `--app-*` and category colours stay in a thin app sheet |
| 3 | **career** | most UI and the biggest payoff; the CSS now inline in each page moves to library classes plus one app sheet |
| 4 | **relevo-site** | mostly tokens and base; a `relevo` palette as app CSS; display typography stays app CSS |
| 5 | **unassisted** | a `whiteboard` palette, the Tailwind mapping, the React components |

The rest of §4.4 (drawer, menu, toast, tabs, segmented control, page shell,
theme switcher, meter, empty state, skeleton, spinner, chip, badge, switch)
ships in later 0.x releases, pulled forward when a migration needs one.

**How each migration is done:**

- Before and after screenshots of the app's key pages, diffed and reviewed.
- A difference that comes from the shared look is accepted.
- Anything else is fixed with an app override, or fed back into the library if more than one app needs it.
- Each migration is its own branch and PR in the app's repo, pinned to a tagged ai-ui version.

## 11. Success criteria

- money, career, relevo-site and unassisted each import ai-ui at a tagged version, and none of them defines its own button, field, callout, pill, table or dialog styles.
- Each app's own CSS is limited to `--app-*` tokens, palette or token overrides, `--ui-*` knob overrides, and layout specific to that app.
- The parity, contrast and axe suites pass in CI on every tag.
- A new Go app can show a themed page with a form and a dialog with no CSS of its own: mount `ui.Handler()`, call `ui.Parse`, link `ui.css`.
