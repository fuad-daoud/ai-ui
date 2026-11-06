# Changelog

Breaking changes come first in each entry. Renaming or removing a token,
class, data-attribute value, knob, partial name or argument, or React prop
is breaking.

## 0.2.0 (2026-11-06)

No breaking changes. One visible change at the default density:

- Inputs and selects are now 40px, the same as a button. They were about 44px, because their padding and line height added up to more than their minimum height. A multi-select and a textarea keep their padding.

New:

- `data-density="compact"`: 32px controls, 14px text and 4px corners, on the page or on one part of it.
- `ui-scoped.css`: styles library markup only, for a page whose other content must stay untouched. Generated from the same sources as `ui.css`.

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
