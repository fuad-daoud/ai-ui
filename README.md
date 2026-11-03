# ai-ui

The shared UI for Fuad's platforms: one token vocabulary, layered CSS that any app
can override, `html/template` partials for the Go apps, and React components for
the Next.js apps. Both render the same markup. See [DESIGN.md](DESIGN.md) for the contract.

## Go apps

```bash
go get github.com/fuad-daoud/ai-ui@v0.1.0
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
"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.1.0" }
```

```ts
// next.config.ts
export default { transpilePackages: ["ai-ui"] };
```

```css
/* app/globals.css: this order matters */
@import "ai-ui/theme.css";
@import "tailwindcss";
@import "ai-ui/ui.css";
```

```tsx
import { Button, Field } from "ai-ui/react";
```

`theme.js` has to run before first paint. Serving it from a Next app is part of the
unassisted migration and is not wired up in v0.1.

## Developing

- `go test ./...`: contrast, build freshness, Go parity, handler, gallery.
- `dev run npx vitest run`: theme.js, TS canon, React parity, Tailwind mapping.
- `dev run npx tsc --noEmit`: typecheck.
- `dev run npx playwright test`: gallery screenshots and axe (contabo only).
- After editing `css/`: `go run ./cmd/build`, then commit `dist/ui.css` too.
- After changing a partial on purpose: `go test . -run TestParityGo -update`, then read the diff.
