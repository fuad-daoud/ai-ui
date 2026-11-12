# ai-ui v0.3.2: plan index

**Spec:** `docs/superpowers/specs/2026-11-12-v0.3.2-fixes-design.md`. Fuad approved the release,
its merge, tag and site deploy on 2026-11-12.

One plan, one relevo chain on `home` (home-01), from `feat/v0.3.2` with the plan committed:

```sh
relevo chain --server home --name ai-ui-v032 --feature ai-ui-v0.3.2 --no-security \
  --plan docs/superpowers/plans/2026-11-12-ai-ui-v0.3.2/01-fixes.md
```

The plan's patches came from `trial/v0.3.2`; replaying them on a clean `feat/v0.3.2` gives that
tree exactly. On the trial: Go and site Go pass; 61 vitest; 64 browser tests; 57 site tests; the 18
gallery screenshot tests pass against v0.3.1's baselines. Each new test fails with its fix removed.

## Before merging

1. The chain's head tree equals `trial/v0.3.2`, apart from `docs/superpowers/`.
2. The gate on the head, Go locally and Node on contabo, with snapshots.
3. PR to `main`, CI green, merge, tag `v0.3.2`, redeploy the site.
4. Career: pin v0.3.2, drop its two stopgap rules, captures unchanged.
