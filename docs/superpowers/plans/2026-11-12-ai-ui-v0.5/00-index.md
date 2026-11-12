# ai-ui v0.5: plan index

**Spec:** `docs/superpowers/specs/2026-11-12-v0.5-page-parts-design.md`. Fuad approved the
release, its merge, tag and site deploy on 2026-11-12.

One plan, one relevo chain on `home-02`, from `feat/v0.5` with the plan committed:

```sh
relevo chain --server home-02 --name ai-ui-v05 --feature ai-ui-v0.5 --no-security \
  --plan docs/superpowers/plans/2026-11-12-ai-ui-v0.5/01-v0.5.md
```

## What the planner measured

The plan's patches and new files came from `trial/v0.5`, and replaying them on a clean `feat/v0.5`
gives that tree exactly (apart from the 16 gallery baselines, which the planner commits after the
chain). On the trial: Go and site Go pass; 70 vitest; 99 browser tests with the new baselines, 13
of them in `e2e/parts.spec.ts`; 57 site tests. Money builds against the trial with only its render
fixture's checksum moving (+28,790 bytes, the bigger `ui.css`).

## Before merging

1. The chain's head tree equals `trial/v0.5`, apart from `docs/superpowers/` and the baselines.
2. The planner commits the 16 baselines from the trial, then runs the gate on the head with
   snapshots.
3. PR to `main`, CI green, merge, tag `v0.5.0`, redeploy the site.
