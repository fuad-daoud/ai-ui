# Using ai-ui

You are a coding agent, and this project's UI should use ai-ui: a small UI library with tokens,
layered CSS, Go `html/template` partials and React components. Go and React render the same markup.

## Before you write any UI

Read the contract. It names every token, class, knob, partial argument and React prop:
https://raw.githubusercontent.com/fuad-daoud/ai-ui/main/DESIGN.md

## Install

Check which stack this project uses, then follow one of these.

### Go with `html/template`

```bash
go get github.com/fuad-daoud/ai-ui@v0.2.0
```

```go
import ui "github.com/fuad-daoud/ai-ui"

mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
tmpl := template.Must(ui.Parse(template.New("app")))
```

```html
<script src="{{href "theme.js"}}"></script>
<link rel="stylesheet" href="{{href "ui.css"}}">
```

### Next.js

```json
"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.2.0" }
```

```ts
// next.config.ts
export default { transpilePackages: ["ai-ui"] };
```

```css
/* app/globals.css: this order matters */
@import "ai-ui/layers.css";
@import "tailwindcss";
@import "ai-ui/theme.css";
@import "ai-ui/ui.css";
```

```tsx
import { Button, Field } from "ai-ui/react";
```

`theme.js` has to run before first paint. Serving it from a Next app isn't wired up yet, so add it
yourself if the project needs a stored theme choice.

## Rules

- Use a library component before writing your own.
- Use only names in the contract. Never invent a `ui-*` class, and don't give your own classes that
  prefix.
- Don't write a default as an attribute. An unknown value renders the default, so a typo fails
  quietly: check names against the contract.
- Theme through tokens and knobs in the app's own unlayered CSS. Name tokens that belong to the app
  `--app-*`.
- Never edit the library's files, and never use `!important` to beat it. An unlayered rule already
  wins.

## A project that already has CSS

`ui.css` changes every element its rules touch. Before adding it to an existing app, read "Adopting
ai-ui in an app that has its own CSS" in the contract, and consider `ui-scoped.css` for a page whose
other content must stay as it is.
