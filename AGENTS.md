# Using ai-ui

You are a coding agent, and this project's UI should use ai-ui: a small UI library with tokens,
layered CSS, Go `html/template` partials and React components. Go and React render the same markup.

## Before you write any UI

Read the contract. It's the library's public API: tokens, classes, knobs, partial arguments and
React props:
https://raw.githubusercontent.com/fuad-daoud/ai-ui/main/DESIGN.md

What the library gives you:

- **Components** with a Go partial and a React component: button, field (input, select, textarea),
  check (checkbox, radio), pill.
- **Components with a React component whose markup Go writes from the contract:** dialog, menu,
  sheet (a bottom sheet or a side drawer), and the theme picker (swatches).
- **Plain markup** you write yourself in both stacks, with the classes the contract shows: callout,
  table, card, disclosure, page parts (page head, eyebrow, section head, stat tiles, empty state),
  chips (with the segmented control, badge and avatar), meter, and the utilities (`ui-stack`,
  `ui-cluster`, `ui-vh`, …).
- **Nothing for page layout.** A page's width, gutters and section spacing are the app's own CSS.

Before deciding the library can't do something, check its partials in `templates/` and its
components in `react/`.

## Install

Check which stack this project uses, then follow one of these.

### Go with `html/template`

```bash
go get github.com/fuad-daoud/ai-ui@v0.5.0
```

It needs Go 1.25, and `go get` raises the `go` line in an older `go.mod` to match.

```go
import ui "github.com/fuad-daoud/ai-ui"

mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
tmpl := template.Must(ui.Parse(template.New("app")))
tmpl = template.Must(tmpl.ParseFS(pages, "*.html"))
```

Parse the app's own templates after `ui.Parse`, because that call is what defines `href` and `dict`.
Use the app's own pattern in `ParseFS`: `"templates/*.html"` if that's what it embeds.

```html
<link rel="stylesheet" href="{{href "ui.css"}}">
```

Add `<script src="{{href "theme.js"}}"></script>` before it only if the page offers a light/dark
or palette choice. Without one, the page follows the system's mode and needs no script. Add
`<script src="{{href "ui.js"}}" defer></script>` if the page has a menu, a theme picker or a sheet:
it gives the menu its arrow keys.

### Next.js

This needs React 19.

npm 12 and later refuse a dependency fetched from git unless the project allows it. Allow it
for the project's own dependencies first, in a `.npmrc` beside `package.json`, so `npm ci` and CI
read it too, then install:

```bash
echo "allow-git=root" >> .npmrc
npm install github:fuad-daoud/ai-ui#v0.5.0
```

That leaves this in `package.json`:

```json
"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.5.0" }
```

```ts
// next.config.ts
export default { transpilePackages: ["ai-ui"] };
```

If the project uses Tailwind v4:

```css
/* app/globals.css: this order matters */
@import "ai-ui/layers.css";
@import "tailwindcss";
@import "ai-ui/theme.css";
@import "ai-ui/ui.css";
```

If it doesn't use Tailwind, import only `ai-ui/ui.css`, and don't add Tailwind for this. `layers.css`
and `theme.css` exist only for Tailwind v4. Tailwind v3 isn't covered, so ask before going further.

```tsx
import { Button, Field } from "ai-ui/react";
```

A page with a `Menu`, `Swatches` or `Sheet` also imports `"ai-ui/ui.js"` once, in a client
component.

`theme.js` has to run before first paint. Serving it from a Next app isn't wired up yet, so add it
yourself if the project needs a stored theme choice.

## Rules

- Use the library's components and markup before writing your own.
- Use names from the contract or the library's source. Never invent a `ui-*` class, and don't give
  your own classes that prefix.
- Don't write a default as an attribute. An unknown value renders the default, so a typo fails
  quietly: check names against the contract.
- Theme through tokens and knobs in the app's own unlayered CSS. Name tokens that belong to the app
  `--app-*`.
- Never edit the library's files, and never use `!important` to beat it. An unlayered rule already
  wins.

## A project that already has CSS

`ui.css` changes every element its rules touch. Before adding it to an existing app, read "Adopting
ai-ui in an app that has its own CSS" in the contract, and consider `ui-scoped.css` for a page whose
other content must stay as it is. Tailwind's own reset on its own doesn't count: a Tailwind v4 app
with no styles of its own follows the steps above and needs no compat layer.
