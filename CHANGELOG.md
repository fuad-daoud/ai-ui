# Changelog

Breaking changes come first in each entry. Renaming or removing a token,
class, data-attribute value, knob, partial name or argument, or React prop
is breaking.

## 0.1.0 (unreleased)

First release.

- Tokens with the `blue` (default) and `mono` palettes, light and dark through `light-dark()`.
- Layered CSS: reset, base, utilities; button, field (input, select, textarea), check, pill, dialog, callout, table, disclosure, card.
- Go: `ui.Handler`, `ui.Href`, `ui.Parse`, and the partials `ui.button`, `ui.field`, `ui.check`, `ui.pill`.
- React: `Button`, `ButtonLink`, `Field`, `Check`, `Pill`, `Dialog`; Tailwind v4 layer order in `layers.css` and token mapping in `theme.css`.
- `theme.js`: pre-paint theme and palette, persisted, with `data-ui-set-*` buttons.
