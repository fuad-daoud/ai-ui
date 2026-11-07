# ai-ui

A small UI library written to be used by coding agents: one token vocabulary, layered CSS that any
app can override, `html/template` partials for Go apps, and React components for Next.js apps. Both
render the same markup. [DESIGN.md](DESIGN.md) is the contract, and [AGENTS.md](AGENTS.md) is where
an agent starts. The site is at https://fuad-daoud.github.io/ai-ui/.

## Go apps

```bash
go get github.com/fuad-daoud/ai-ui@v0.2.0
```

```go
import ui "github.com/fuad-daoud/ai-ui"

mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
tmpl := template.Must(ui.Parse(template.New("app")))
tmpl = template.Must(tmpl.ParseFS(pages, "*.html"))
```

```html
<head>
  <script src="{{href "theme.js"}}"></script>
  <link rel="stylesheet" href="{{href "ui.css"}}">
  <link rel="stylesheet" href="/static/app.css">  <!-- unlayered: always wins -->
</head>
…
{{template "ui.field" (dict "Label" "Email" "Name" "email" "Type" "email" "Error" .Errs.Email)}}
{{template "ui.button" (dict "Label" "Sign in" "Type" "submit")}}
```

## Next.js apps

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

`theme.js` has to run before first paint. Serving it from a Next app isn't wired up yet, so that
step is yours.

## Developing

- `go test ./...`: contrast, build freshness, Go parity, handler, gallery.
- `dev run npx vitest run`: theme.js, TS canon, React parity, Tailwind mapping.
- `dev run npx tsc --noEmit`: typecheck.
- `dev run npx playwright test`: gallery screenshots and axe (contabo only).
- Site: `cd site && go test ./... && go run . -out out -serve 127.0.0.1:4318`, then open http://127.0.0.1:4318/ai-ui/.
- `dev run npm run e2e:site`: the site in a browser.
- After editing `css/`: `go run ./cmd/build`, then commit `dist/ui.css` too.
- After changing a partial on purpose: `go test . -run TestParityGo -update`, then read the diff.
