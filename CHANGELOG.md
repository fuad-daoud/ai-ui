# Changelog

Breaking changes come first in each entry. Renaming or removing a token,
class, data-attribute value, knob, partial name or argument, or React prop
is breaking.

## 0.1.1 (2026-11-06)

Docs only. `ui.css`, the templates and the React components do not change.

- `DESIGN.md`: adopting ai-ui in an app that has its own CSS. The compat layer, app markup inside a library container, and how to move a component and keep the app's look, as verified moving money.

## 0.1.0 (2026-11-04)

First release.

- Tokens with the `blue` (default) and `mono` palettes, light and dark through `light-dark()`.
- Layered CSS: reset, base, utilities; button, field (input, select, textarea), check, pill, dialog, callout, table, disclosure, card.
- Go: `ui.Handler`, `ui.Href`, `ui.Parse`, and the partials `ui.button`, `ui.field`, `ui.check`, `ui.pill`.
- React: `Button`, `ButtonLink`, `Field`, `Check`, `Pill`, `Dialog`; Tailwind v4 layer order in `layers.css` and token mapping in `theme.css`.
- `theme.js`: pre-paint theme and palette, persisted, with `data-ui-set-*` buttons.
