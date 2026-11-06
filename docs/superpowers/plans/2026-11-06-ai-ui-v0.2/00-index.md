# ai-ui v0.2: plan index

**Spec:** `docs/superpowers/specs/2026-11-06-v0.2-compact-and-scoped-design.md`. Fuad approved it on
2026-11-06, including the field fix in its §5.

**What ships:** `data-density="compact"`, `ui-scoped.css`, and inputs and selects that match buttons.
Nothing is renamed or removed.

## Plans

| # | Plan | What it adds | `dist/ui.css` | `dist/ui-scoped.css` |
|---|---|---|---|---|
| 01 | `01-scoped-stylesheet.md` | The generator, the published file, fixture pages and browser proof | Unchanged, 16,455 B | New, 19,889 B |
| 02 | `02-compact-density.md` | Compact density, the field fix and gallery compact pages | 17,911 B | 21,345 B |
| 03 | `03-docs-and-version.md` | `DESIGN.md`, `CHANGELOG.md`, 0.2.0 | Unchanged | Unchanged |

They run as one relevo chain on zen, in order:

```sh
relevo chain --server zen --name ai-ui-v02 --feature ai-ui-v0.2 --no-security \
  --plan docs/superpowers/plans/2026-11-06-ai-ui-v0.2/01-scoped-stylesheet.md \
  --plan docs/superpowers/plans/2026-11-06-ai-ui-v0.2/02-compact-density.md \
  --plan docs/superpowers/plans/2026-11-06-ai-ui-v0.2/03-docs-and-version.md
```

Run it from `feat/v0.2`, with the plans committed. `--no-security` because the only new served path
is one more static file through the handler that already exists. A chain on a server holds its own
copy of the plans, so correcting one means closing the chain and starting a new one.

## What the planner already measured

The plans' code was generated from trial trees, not retyped, and each round's tree was run on its own:

| Check | Result |
|---|---|
| Round 1 alone | Go, vet, gofmt, `tsc`, 45 vitest, 16 browser tests pass. `ui.css` byte-identical to `v0.1.1`. |
| Round 2 | Go, vet, gofmt, `tsc`, 45 vitest, 30 browser tests pass, including axe on 8 new compact runs: 4 pages and 4 dialogs. |
| Field fix at the default size | Measured: a 40px input and select beside a 40px button. |

Each break below was made on purpose, and the named tests failed:

| Broken on purpose | Tests that failed |
|---|---|
| `font-size` on `:root` | Compact sizes on the root, fields match buttons on the root |
| Field fix removed | Compact sizes (both), fields match buttons (both) |
| Multi-select loses its padding | Fields match buttons (both) |
| Compact table rule removed | Compact sizes (both) |
| Scope matches everything | `ui-scoped.css` leaves it untouched (light, dark) |
| Scoped file loses the base | Library markup matches (both), focus ring |
| Scope misses subtrees | Library markup matches (both) |

## How the planner verifies, before merging

A `check: none` chain runs no gate of its own, and the builder's report isn't proof.

1. **The diff against the plans.** The plans' code is exact, so any other difference needs a reason.
2. **The gate on the chain's head,** from a `git archive` copy under `$HOME`: Go locally, Node on
   contabo. `dev run` doesn't mirror `dist/` or any directory named `build`, so carry `dist/` over as
   a tarball and unpack it first. A vitest count of 41 means the tarball is missing.
3. **Baselines on contabo,** because they depend on fonts, with the same copy:
   `npx playwright test --update-snapshots`. Bring back `e2e/__screenshots__/` with `scp`.
   - The 8 existing shots may change only where a field sits. Look at one pair, and check that the
     difference starts at the first field and that the page above it is identical.
   - The 4 compact shots and 4 compact dialog shots are new. Look at one of each.
   - Re-run without `--ignore-snapshots`: all 30 pass.
4. **Commit the baselines** on the chain's branch. Check for stray runs with `pgrep -af playwright`
   on contabo.

## Release, once Fuad has approved this plan

1. Push the branch and open a PR to `main`. CI must pass.
2. Merge it, with a merge commit.
3. Tag `v0.2.0` on the merge commit and push the tag. CI checks that it matches `package.json`.
4. Check that `go list -m github.com/fuad-daoud/ai-ui@v0.2.0` resolves.

## After this

Career's stage 1 (foundation) and stage 2 (sign-in pages) get their own plans, pinned to `v0.2.0`. Money
stays on `v0.1.0` unless Fuad wants it moved.
