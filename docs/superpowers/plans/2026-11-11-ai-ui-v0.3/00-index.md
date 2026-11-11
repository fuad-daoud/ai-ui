# ai-ui v0.3: plan index

**Spec:** `docs/superpowers/specs/2026-11-11-v0.3-options-design.md`. Fuad approved the spec, the
plans, the build, the merge, the release and the site deploy on 2026-11-11.

**What ships:** options on the components that exist (spec §4, §5) and two changed looks (§6):
disabled buttons, and compact cards and dialogs. Nothing is renamed or removed. Tables are v0.3.1.

## Plans

| # | Plan | What it adds | `dist/ui.css` |
|---|---|---|---|
| 01 | `01-library.md` | CSS, Go partials, React, parity fixtures, the Go icon test, `e2e/options.spec.ts`, the gallery | 23,822 B (18,084 B in v0.2.1) |
| 02 | `02-docs-site-version.md` | `DESIGN.md`, `CHANGELOG.md`, 0.3.0, install pins, site demos, table region names | Unchanged |

They run as one relevo chain on zen, in order, from `feat/v0.3` with the plans committed:

```sh
relevo chain --server zen --name ai-ui-v03 --feature ai-ui-v0.3 --no-security \
  --plan docs/superpowers/plans/2026-11-11-ai-ui-v0.3/01-library.md \
  --plan docs/superpowers/plans/2026-11-11-ai-ui-v0.3/02-docs-site-version.md
```

`--no-security`: no new served path, no new input handling; the only new markup input is `Icon`,
which takes `template.HTML` from the app and escapes anything else (pinned by
`TestIconIsTheAppsMarkup`).

## What the planner already measured

The plans' patches and files were generated from a trial tree (`trial/v0.3`), and both plans were
replayed mechanically on a clean `feat/v0.3`: the result equals the trial tree, 0 diff lines.

| Check, on the trial tree | Result |
|---|---|
| Go, vet, gofmt, site Go | Pass |
| `tsc`, vitest | Pass, 61 tests |
| Browser (`--ignore-snapshots`) | 54 passed, axe on all 16 gallery pages |
| Site browser | 57 passed, axe on every page |

Each break below was made on purpose, on a copy, and the named tests failed:

| Broken on purpose | Tests that failed |
|---|---|
| Disabled button fades to 0.5 | disabled changes colour, not opacity (both densities) |
| Outline rules without `:where()` | disabled changes colour, not opacity (both) |
| Touch sizes removed | compact controls grow to finger sizes |
| Compact card keeps its shadow | cards and dialogs are flat only in compact |
| Small select loses its chevron room | small and large fields are as tall as the button (both) |
| Icon padding removed | the input icon sits inside the input (both) |

## How the planner verifies, before merging

1. **The chain's head tree equals `trial/v0.3`,** apart from `docs/superpowers/` and the baselines.
2. **The gate on the head,** Go locally and Node on contabo (`dev run` from a git checkout, so it
   keeps `dist/` and `cmd/build/`).
3. **Baselines on contabo:** `npx playwright test e2e/gallery.spec.ts --update-snapshots`, copy
   `e2e/__screenshots__/` back, look at one light and one dark compact page, then commit them on the
   chain's branch. All 16 change: new gallery content and the §6 looks.

## Release

1. Push and open a PR to `main`; CI must pass. Merge with a merge commit.
2. Tag `v0.3.0` on the merge commit and push the tag; CI checks it matches `package.json`.
3. `go list -m github.com/fuad-daoud/ai-ui@v0.3.0` resolves.
4. Redeploy the site (`srv.fish deploy ai-ui-site`, see the servers repo's
   `services/ai-ui-site/NOTES.md`); its footer and release notes show 0.3.0.
