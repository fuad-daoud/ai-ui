# ai-ui v0.3.1: plan index

**Spec:** `docs/superpowers/specs/2026-11-11-v0.3.1-tables-design.md`. Fuad said "go ahead" on
2026-11-11; merge, tag and site deploy wait for his word.

**What ships:** table options (no frame, quiet head, sticky head, stacked on a phone), table knobs,
and row hover only where a pointer can hover. CSS, docs, gallery, tests and the site; no Go partial
or React change. `dist/ui.css` 27,033 B (23,822 B in v0.3.0).

One plan, one relevo chain on `home` (home-01), from `feat/v0.3.1` with the plan committed:

```sh
relevo chain --server home --name ai-ui-v031 --feature ai-ui-v0.3.1 --no-security \
  --plan docs/superpowers/plans/2026-11-11-ai-ui-v0.3.1/01-tables.md
```

## What the planner already measured

The plan's patches came from a trial tree (`trial/v0.3.1`), and replaying the plan on a clean
`feat/v0.3.1` gives that tree exactly. On the trial: Go and site Go pass; 61 vitest; 62 browser
tests; 57 site tests; axe on every gallery page and on the gallery at 390px with its tables stacked.

| Broken on purpose | Test that failed |
|---|---|
| Hover on touch screens | a tapped row keeps no tint |
| Sticky head not sticky | a sticky head stays at the top of a bounded wrap |
| Unframed table keeps its edge padding | an unframed table has no frame, and its outer columns sit flush |
| Header row shown when stacked | rows stack into labelled lines |
| Numbers end-aligned when stacked | rows stack into labelled lines |
| Quiet head keeps its fill | a quiet head has no fill and small grey text |

## Before merging

1. The chain's head tree equals `trial/v0.3.1`, apart from `docs/superpowers/` and the baselines.
2. The gate on the head, Go locally and Node on contabo, against the 16 new baselines (they change:
   the gallery gains a table).
3. PR to `main`, CI green. Merge, tag `v0.3.1` and redeploy the site at Fuad's word.
