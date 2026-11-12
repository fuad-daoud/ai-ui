# Changelog

Breaking changes come first in each entry. Renaming or removing a token,
class, data-attribute value, knob, partial name or argument, or React prop
is breaking.

## 0.5.0 (2026-11-12)

No breaking changes.

New, all plain markup in Go and React:

- Page parts: `.ui-page-head` (with `.ui-page-title`, `.ui-page-sub`, `.ui-page-aside`), `.ui-eyebrow`, `.ui-section-head` (with `.ui-section-index`, `.ui-section-title`, `.ui-section-count`, `.ui-section-aside`), stat tiles (`.ui-stats`, `.ui-stat`, `.ui-stat-value`, `.ui-stat-note`, `data-variant`, `data-frame="none"`), and the empty state (`.ui-empty`, `.ui-empty-title`, `data-frame="none"`). Knobs `--ui-page-title-size`, `--ui-page-head-rule`, `--ui-eyebrow-fg`, `--ui-eyebrow-case`, `--ui-section-rule`, `--ui-section-title-size`, `--ui-section-index-fg`, `--ui-stats-min`, `--ui-stat-value-size`, `--ui-stat-label-case`.
- Chips (`.ui-chips`, `.ui-chip`, current by `aria-current` or `aria-pressed`), the count badge (`.ui-badge`), the segmented control (`.ui-segmented`, `.ui-segment`) and the avatar (`.ui-avatar`, which fills a square button). Knobs `--ui-chip-height`, `--ui-badge-bg`, `--ui-badge-fg`, `--ui-segment-height`, `--ui-avatar-size`.
- Meter (`.ui-meter` on `<progress>`): a fill, `data-variant`, and a sliding segment when indeterminate. Knobs `--ui-meter-height`, `--ui-meter-track`, `--ui-meter-fill`.
- Table: the current row (`aria-current`, or `aria-selected` in a grid) keeps a `--ui-table-selected` tint.

## 0.4.0 (2026-11-12)

No breaking changes. One look can change:

- `data-palette` and `data-theme` now apply on any element, not only `<html>`. An element that already carried either attribute takes that palette or mode. `data-theme="system"` follows the OS again.

New:

- Menu (`.ui-menu`): a popover anchored to the button that opens it (`popovertarget`), with items (`.ui-menu-item`, `data-variant="danger"`, `role="menuitemradio"` with a check), a label, a separator, `data-align="start"`, and `data-sheet="phone"` to become a bottom sheet below 800px with a head (`.ui-menu-head`, `.ui-menu-title`, `.ui-menu-x`). Knobs `--ui-menu-width`, `--ui-menu-bg`, `--ui-menu-border`, `--ui-menu-radius`, `--ui-menu-shadow`, `--ui-menu-item-height`, `--ui-menu-hover`. React `Menu`, `MenuItem`, `MenuLink`, `MenuLabel`, `MenuSeparator`.
- Sheet (`.ui-sheet`): a modal `<dialog>` docked to the bottom, or with `data-side="end"` to the side as a drawer; `closedby="any"` closes it on a backdrop tap. It slides in unless reduced motion is asked for, and is flat in compact. Knobs `--ui-sheet-bg`, `--ui-sheet-border`, `--ui-sheet-radius`, `--ui-sheet-shadow`, `--ui-sheet-pad`, `--ui-sheet-max-height`, `--ui-sheet-width`. React `Sheet`.
- Theme picker (`.ui-swatches`, `.ui-swatch`, `.ui-swatch-face`): swatch buttons carrying `data-ui-set-*`, faces that preview any palette and mode, a split face for `system`. React `Swatches`, `Swatch`, `SwatchFace`.
- A light-only palette: `color-scheme: light` in its block keeps it light in either mode.
- `theme.js`: `aiui:themechange` on `document` after every change and on an OS flip while following it; `aria-pressed` / `aria-checked` on every `data-ui-set-*` element; `window.aiui.theme.sync(root)`.
- `ui.js` (`{{href "ui.js"}}`, `import "ai-ui/ui.js"`): a menu's keys and focus, swatch arrows, a placement fallback where CSS anchor positioning is missing, and `closedby="any"` where the browser lacks it.
- `DESIGN.md`: moving an app's stored theme choice to `theme.js`, and why `ui.css` belongs in `<head>`.

## 0.3.2 (2026-11-12)

No breaking changes. Two looks change:

- A dialog's title shrinks and wraps anywhere it must, so a long word such as an email address no longer pushes the × or the dialog past the viewport.
- In a table with `data-head="quiet"`, row headers in the body (`<th scope="row">`) go quiet too: no fill, in the head's colour and size. They keep their case, and their row still tints on hover.

## 0.3.1 (2026-11-11)

No breaking changes. One look changes:

- A table row is tinted on hover only where the pointer can hover, so a tapped row on a phone no longer keeps its tint. The tint is the new `--ui-table-hover` knob.

New:

- Tables: `data-frame="none"` on the wrap for no frame, with flush outer columns; `data-head="quiet"` for a header on no fill in small grey text; `data-sticky="head"` for a header that stays at the top of a scrolling wrap; `data-stack="phone"` on the wrap to stack each row into labelled lines below 800px, labels from each cell's `data-label`.
- Table knobs: `--ui-table-head-bg`, `--ui-table-head-fg`, `--ui-table-head-size`, `--ui-table-head-case`, `--ui-table-head-rule`, `--ui-table-hover`, `--ui-table-sticky-bg`, `--ui-table-label-width`, and `--ui-table-max-height` on the wrap.

## 0.3.0 (2026-11-11)

No breaking changes. Two looks change:

- A disabled button is drawn in `--band` with `--ink-3` text and a `--rule` border, at full opacity, at every density. It used to fade to half opacity. A disabled ghost button keeps no fill, and a disabled `<button>` now shows the not-allowed cursor.
- In compact density, cards and dialogs have no shadow. Their border stays.

New:

- Button: `data-tone="outline"` for an outlined primary or danger button, `data-shape="square"` for an icon-only button, `data-width="full"` for a full-width button whose label wraps, and an icon before the label. Go `Tone`, `Shape`, `Width`, `Icon`; React `tone`, `shape`, `width`, `icon`.
- Inputs, selects and textareas: `data-size="sm"` and `"lg"`, as tall as the button of that size. Go `Size`, React `size`.
- A leading icon inside an input: `.ui-input-group` and `.ui-input-icon`. Go `Icon`, React `icon`.
- Knobs `--ui-input-width`, and `--ui-row-grow` and `--ui-row-basis` on a child of `.ui-form-row`.
- In compact density, on a touch screen, controls grow to 40px (34px small, 44px large).
- Card: knobs `--ui-card-shadow` and `--ui-card-border`.
- Dialog: knobs `--ui-dialog-shadow`, `--ui-dialog-border` and `--ui-dialog-eyebrow-fg`; `.ui-dialog-eyebrow` above the title; the × is optional. React `eyebrow` and `closeButton`.
- Pill: `data-dot` for a status dot and `data-tone="outline"`. Go `Dot`, `Tone`; React `dot`, `tone`.
- Disclosure: `data-marker="start"` puts the chevron before the summary, and the knob `--ui-disclosure-rule`.
- `DESIGN.md`: every option above, and a dialog shown without script.

## 0.2.1 (2026-11-09)

No breaking changes.

- A dialog's text starts at the start and wraps, wherever its markup sits. It used to inherit `text-align` and `white-space`, so a dialog in an end-aligned table cell came out right-aligned.
- `AGENTS.md`: npm 12 refuses git dependencies, so the Next.js steps allow them in `.npmrc`; Go's version floor; the app's own `ParseFS` pattern; `theme.js` only for a page with a theme choice; what is a component, what is plain markup, and that page layout is the app's.
- `DESIGN.md` documents what the partials and components already took: the button's `Type` (default `button`), `Name` and `Value`; the check's `Type`, `ID`, `Required` and `Disabled`; `.ui-select` and `.ui-textarea`; the options shape; controlled fields and checks; `Dialog`'s `closeLabel`; `.ui-dialog-actions` and a confirm per row; roles for callouts.
- `DESIGN.md`: where to set a knob, library sizes as minimums, hover colours, a renamed selector that wakes a rule, shared token names, and a page with a design of its own on `ui-scoped.css`.

## 0.2.0 (2026-11-06)

No breaking changes. One visible change at the default density:

- Inputs and selects are now 40px, the same as a button. They were about 44px, because their padding and line height added up to more than their minimum height. A multi-select and a textarea keep their padding.

New:

- `data-density="compact"`: 32px controls, 14px text and 4px corners, on the page or on one part of it.
- `ui-scoped.css`: styles library markup only, for a page whose other content must stay untouched. Generated from the same sources as `ui.css`.

## 0.1.1 (2026-11-06)

Docs only. `ui.css`, the templates and the React components do not change.

- `DESIGN.md`: adopting ai-ui in an app that has its own CSS. The compat layer, app markup inside a library container, and how to move a component and keep the app's look, as verified moving a real app.

## 0.1.0 (2026-11-04)

First release.

- Tokens with the `blue` (default) and `mono` palettes, light and dark through `light-dark()`.
- Layered CSS: reset, base, utilities; button, field (input, select, textarea), check, pill, dialog, callout, table, disclosure, card.
- Go: `ui.Handler`, `ui.Href`, `ui.Parse`, and the partials `ui.button`, `ui.field`, `ui.check`, `ui.pill`.
- React: `Button`, `ButtonLink`, `Field`, `Check`, `Pill`, `Dialog`; Tailwind v4 layer order in `layers.css` and token mapping in `theme.css`.
- `theme.js`: pre-paint theme and palette, persisted, with `data-ui-set-*` buttons.
