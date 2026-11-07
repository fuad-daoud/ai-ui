# ai-ui site: Components and Guides pages

Date: 2026-11-07
Status: draft, awaiting Fuad's review
Extends: [2026-11-07-site-design.md](2026-11-07-site-design.md)

## 1. Purpose

The nav's Components and Guides links go to `DESIGN.md` on GitHub. Fuad wants them to be pages of the
site. The first spec held both back until the home page had feedback (its §2), and this reverses that.

**Fuad's decisions (2026-11-07):**

- Guides are built from sections of `DESIGN.md` and `README.md`, chosen by heading at build time. No
  guide text is written twice.
- Components is one page with a section per component: a live demo made with the real partials, then
  that component's section of `DESIGN.md` as written.

Everything in the first spec still holds unless this one says otherwise. That includes the rule that
site copy never states how many components there are.

## 2. Non-goals

- **No hand-written guide text** beyond each guide's title and the one line that describes it on the
  guides index.
- **No change to the library** apart from the README edit in §6. No token, class, partial or React
  component changes, no version bump and no tag.
- **No page per component, no search, no code tabs and no syntax highlighting.** Code blocks are
  plain, as GitHub shows them.
- **No theme switcher.** The site stays dark.

## 3. Pages

```
out/components/index.html
out/guides/index.html
out/guides/getting-started/index.html
out/guides/theming/index.html
out/guides/density/index.html
out/guides/existing-apps/index.html
```

| Page | Title | Built from |
|---|---|---|
| Getting started | Getting started | README: `## Go apps`, `## Next.js apps` |
| Theming | Theming | DESIGN: `## Override model`, `## Tokens`, `## Modes and palettes` |
| Density | Density | DESIGN: `## Density` |
| Existing apps | Adding to an existing app | DESIGN: `## Adopting ai-ui in an app that has its own CSS`, `## Choosing a stylesheet` |

**The guides index** lists the four guides, each a link with one line under it:

- Getting started: Install it in a Go app or a Next.js app.
- Theming: Tokens, knobs, light and dark, palettes, and why your CSS always wins.
- Density: A compact size for working tools, for a whole page or one part of it.
- Adding to an existing app: Bring it into an app that already has its own CSS.

**A guide page** has the guide's title as its `h1`, then its sections in the order above, each opening
with its own heading as an `h2`. A side list links all four guides, and the current one carries
`aria-current="page"`.

**The Components page** has the `h1` "Components" and the paragraph under `DESIGN.md`'s
`## Components` heading. Then comes one section per `###` heading under it, in the contract's order:

| Side list | `DESIGN.md` heading | Demo |
|---|---|---|
| Button | `### Button: .ui-btn` | every variant, every size, disabled, a link button |
| Field | `### Field: .ui-field` | input with a hint, input with an error, select, textarea, disabled, two fields in a `ui-form-row` |
| Check | `### Check: .ui-check` | a checkbox, a checked checkbox, a radio pair, a disabled box |
| Pill | `### Pill: .ui-pill` | every variant |
| Dialog | `### Dialog: .ui-dialog` | a button that opens a real modal with `command="show-modal"`; Cancel and × close it |
| Containers | `### Containers (plain markup in Go and React)` | callouts in every variant, a table, a disclosure pair sharing a `name`, a card |
| Utilities | `### Utilities` | muted, mono and number text, a stack and a cluster |

Each section is the heading as an `h2`, the demo, then the rest of that section of `DESIGN.md`. Demo
IDs start with `demo-` so they can't collide with heading IDs. The side list uses the short names
above.

## 4. Layout

- **Nav:** Components goes to `components/` and Guides to `guides/`, both relative to `.Root`. The
  nav link for the current part of the site gets `aria-current="page"` and a visible underline.
  Release notes gets the same treatment. GitHub stays the last link.
- **Home page:** "Browse the components" goes to `components/`. The Go and Next.js links under the
  prompt go to `guides/getting-started/#go-apps` and `#nextjs-apps`.
- **Docs layout** for the Components page and the guide pages:
  - At 60rem and wider, a side list about 12rem wide sits beside the text and is sticky.
  - Narrower, the side list becomes a wrapping row above the text.
  - The text column uses the existing `.prose` measure.
  - Demos sit in a framed panel at the text column's width.
- **Rendered Markdown:**
  - Code blocks wrap their lines rather than scroll, so they never push the page sideways and no
    scroll region needs focus.
  - Tables render as `ui-table` inside a `ui-table-wrap` with `tabindex="0"`, so a table too wide for
    a phone scrolls inside its frame and can be reached from the keyboard.

## 5. How it's built

`section(src, heading)` parses one Markdown file with goldmark and finds the heading whose plain text
equals `heading`; inline code counts as its text. The section runs from that heading up to the next
heading at the same level or higher. It returns the heading and the body separately, rendered to HTML,
with every heading in it moved so the section's own heading becomes an `h2`.

- **The Components page's opening paragraph** comes from the same lookup, which stops at the first
  heading of any level, so it takes only the text before the first `###`.
- **A heading the site names that's missing** fails the build, and nothing is written. This is the
  rule broken sources already follow.
- **Rendering:** goldmark with tables and automatic heading IDs, and raw HTML dropped, as for the
  changelog. The renderer is the only source of `<table>` tags, so wrapping them is a string
  replacement on its output.
- **Relative links** in a section are rewritten to the file on GitHub at `main`. Links that start with
  `#` are kept.
- **The component list** is a Go slice of short name, `DESIGN.md` heading and demo template. A demo is
  a template in `templates/components.html` named `demo <short name>`. It's rendered to HTML by the
  generator, so a missing demo fails the build.
- **`page`** gains `Section`, which says which nav link is current, plus the fields the new templates
  read: the guide list, the current guide and its sections, and the component sections.

## 6. README

`README.md`'s Next.js section only covers Tailwind v4. It gains what `AGENTS.md` already says:

- it needs React 19
- without Tailwind, import only `ai-ui/ui.css`, and `layers.css` and `theme.css` exist only for
  Tailwind v4

This keeps Getting started complete for a human.

## 7. Testing

**Go (`site/main_test.go`):**

- Every page in §3 is built. The existing link, footer and no-count tests cover every page, not just
  the first two.
- The nav's Components and Guides links resolve to site pages. Each page marks its own nav link and
  no other with `aria-current="page"`.
- Every `###` under `## Components` in the real `DESIGN.md` is in the component list. Adding a
  component to the contract therefore fails until it has a demo.
- The Components page holds the first `.ui-*` class of every file in `css/components/`, except
  density, inside its demos.
- No page has a duplicate `id`.
- `section`:
  - finds a heading that has inline code
  - stops at the next heading at the same level or higher, and keeps deeper ones
  - moves headings to `h2`
  - rewrites relative links and keeps `#` links
  - drops raw HTML
  - errors on a missing heading
- A fake repository whose `DESIGN.md` lacks a heading the guides use fails the build and writes nothing.

**Browser (`site/e2e/site.spec.ts`):**

- axe passes on the Components page, the guides index and every guide.
- None of those pages scrolls sideways at 320 and 360 px.
- The nav goes to Components and to Guides, and from the guides index into a guide.
- The dialog demo opens, Cancel closes it, and focus returns to the opener.
- Tab reaches the side list and every demo control on the Components page.
- No duplicate IDs on any page.

## 8. Shipping

One PR. After it merges, the site is rebuilt and deployed with the two commands in the servers repo's
`services/ai-ui-site/NOTES.md`.
