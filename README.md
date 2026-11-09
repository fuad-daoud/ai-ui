# ai-ui

A small UI library written to be used by coding agents: one token vocabulary, layered CSS that any
app can override, `html/template` partials for Go apps, and React components for Next.js apps. Both
render the same markup. [DESIGN.md](DESIGN.md) is the contract, and [AGENTS.md](AGENTS.md) is where
an agent starts. The site is at https://ui.fuad-daoud.com/.

## Go apps

```bash
go get github.com/fuad-daoud/ai-ui@v0.2.0
```

It needs Go 1.25.

```go
import ui "github.com/fuad-daoud/ai-ui"

mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
tmpl := template.Must(ui.Parse(template.New("app")))
tmpl = template.Must(tmpl.ParseFS(pages, "*.html")) // your own pattern, after ui.Parse
```

`theme.js` is only for a page that offers a light/dark or palette choice.

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

This needs React 19. npm 12 and later refuse a dependency fetched from git unless the project
allows it, so allow it in a `.npmrc` (which `npm ci` and CI read too), then install:

```bash
echo "allow-git=root" >> .npmrc
npm install github:fuad-daoud/ai-ui#v0.2.0
```

```json
"dependencies": { "ai-ui": "github:fuad-daoud/ai-ui#v0.2.0" }
```

```ts
// next.config.ts
export default { transpilePackages: ["ai-ui"] };
```

With Tailwind v4:

```css
/* app/globals.css: this order matters */
@import "ai-ui/layers.css";
@import "tailwindcss";
@import "ai-ui/theme.css";
@import "ai-ui/ui.css";
```

Without Tailwind, import only `ai-ui/ui.css`. `layers.css` and `theme.css` exist only for Tailwind v4.

```tsx
import { Button, Field } from "ai-ui/react";
```

`theme.js` has to run before first paint. Serving it from a Next app isn't wired up yet, so that
step is yours.

## Developing

- `go test ./...`: contrast, build freshness, Go parity, handler, gallery.
- `npx vitest run`: theme.js, TS canon, React parity, Tailwind mapping.
- `npx tsc --noEmit`: typecheck.
- `npx playwright test`: gallery screenshots and axe. The baselines were made on Linux with `/usr/bin/chromium`, so screenshots may differ elsewhere.
- Site: `cd site && go test ./... && go run . -out out -serve 127.0.0.1:4318`, then open http://127.0.0.1:4318/ai-ui/.
- `npm run e2e:site`: the site in a browser.
- The site runs on contabo as `ai-ui-site`, built with `site/cmd/serve`. Its deploy steps are in the servers repo.
- After editing `css/`: `go run ./cmd/build`, then commit `dist/ui.css` too.
- After changing a partial on purpose: `go test . -run TestParityGo -update`, then read the diff.
