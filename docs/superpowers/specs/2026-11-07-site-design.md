# ai-ui site: first version

Date: 2026-11-07
Status: draft, awaiting Fuad's review

## 1. Purpose

ai-ui is open source and has no front door. This builds the smallest real site that can be shown to
people for feedback: a home page that makes the library's point in one screen, its release notes, and
the file a visitor's coding agent gets sent to.

**The reader** is a developer whose coding agent writes their UI. The main action on the site is
copying a prompt for that agent. Installing by hand comes second.

**The claim** is "a small, stable contract your agent can read". The site doesn't say agents produce
better UI, because nothing here measures that. It can say what is true of the repository: the contract
is one file, renaming anything in it is a breaking change, and tests hold Go and React to the same
markup.

**Fuad's decisions (2026-11-06 and 2026-11-07):**

- Audience: other developers. Their agents are routed to the repo by a copyable prompt, with no
  separate agent docs.
- Positioning: a UI library built for AI agents to use.
- Build: a small Go generator in this repo producing static HTML.
- Add an `AGENTS.md`.
- Home page: draft 15, "The pile", from the drafts canvas. This replaces an earlier pick of draft 3.
- The home page shows components off. Browsing them belongs on a separate page.
- First version: home, release notes and `AGENTS.md`.
- Hosting: GitHub Pages.
- Site copy never states how many components there are, because the number will change.

## 2. Non-goals

- **No Components page and no guides.** Both wait for feedback on the home page. Until then their nav
  links go to `DESIGN.md` on GitHub.
- **No change to the library.** No token, class, partial or React component changes, no version bump
  and no tag.
- **No search, versioned docs, blog, analytics or feedback form.** Fuad collects feedback directly.
- **No custom domain.** It can be added later without changing the site, because every link is
  relative (§3).
- **No screenshot baselines for the site.** The home page animates, and its design is still being
  judged.
- **No physics engine.** The pile is hand-placed and animated with CSS (§4).
- **The test gallery stays as it is.** It remains a fixture for screenshots and axe.

## 3. Shape

A new Go module in `site/`, so people installing ai-ui don't inherit its Markdown dependency:

```
site/
  go.mod              module github.com/fuad-daoud/ai-ui/site; replace ai-ui => ../
  main.go             the generator
  main_test.go
  templates/          layout.html, home.html, releases.html
  assets/             site.css, site.js, the font and its licence
  e2e/                Playwright specs, with their own config (§8)
```

`go run . -out out` in `site/` writes the finished site to `site/out/`, which is gitignored. `-repo`
points at the repository root and defaults to `..`. `-serve 127.0.0.1:4318` serves the output under
`/ai-ui/`, matching GitHub Pages, for local preview and for the browser tests.

The generator does four things:

1. Parses its templates with `ui.Parse`, so the home page renders the real partials.
2. Copies `ui.css` out of `ui.Assets` and its own `assets/` into the output.
3. Renders `CHANGELOG.md` to HTML with goldmark.
4. Reads the current version from `package.json`.

**Output:**

```
out/index.html
out/releases/index.html
out/assets/ui.css, site.css, site.js, the font
```

Every internal link and asset path is relative. The site therefore works under `/ai-ui/` on GitHub
Pages, at a root path on a custom domain, and from the preview server, with no base-path setting.

## 4. Home page

Built from draft 15, top to bottom:

- **Nav:** the `ai-ui` wordmark, Components, Guides, Release notes, GitHub.
- **Headline:** "Here's the whole UI kit. Your agent can read it."
- **Lede:** "A small pile of real, working components, written down in one contract. Not much to
  learn, and not much to get wrong."
- **Actions:** a primary **Copy the prompt** button, a secondary **Drop it again** button and a
  **Browse the components** link.
- **The pile:** real components fall from above, bounce and settle in a heap on a floor line.
- **Prompt section** (new, see below).
- **Footer:** MIT licence, the current version linked to the release notes, GitHub.

**The pile.** Each item has a hand-chosen position and tilt, set as custom properties, and falls on a
staggered delay using CSS keyframes. The whole fall takes about two seconds and then the page is still.
**Drop it again** swaps between two identical animations, which restarts the fall without resetting
anything typed into the pile. Without JavaScript the first fall still plays, and that button is hidden.

Button, Field, Check and Pill come from the Go partials. Callout, Table, Disclosure and Card are plain
library markup, which is how the library ships them. Dialog is new compared with the draft: it appears
as a static panel using the dialog's classes, so the pile holds the whole library. Sample content is
generic invoicing data.

**The pile is real.** That's the point of the draft, so the components stay clickable, typeable and
reachable by keyboard, though the sample buttons do nothing and the dialog is static, and the page
doesn't promise more than that, and screen readers announce the pile as a labelled group. Items
lift slightly on hover and on focus. Hand placement has to leave every control usable, which a browser
test checks at both desktop and phone widths (§8).

**Adding a component takes a person.** It needs a spot in the pile at desktop and phone widths, unlike
a layout that flows. The site test fails until a new component has been placed, so the page can't fall
behind unnoticed. A physics engine would avoid this, and would let people grab and throw things, but
it adds a dependency and a lot of code to a page meant to be simple. It can come later if the pile
goes down well.

**Motion.** Under `prefers-reduced-motion` the components start in their final places and **Drop it
again** is hidden. No Pause control is needed, because the fall ends within five seconds and only
replays when someone asks.

**The prompt is visible,** which is new compared with the draft. A section below the pile shows the
full prompt text with its own copy button. People can see what they're pasting, and it still works if
the clipboard is refused or scripts are off.

**Links out, until those pages exist:**

| Link | Goes to |
|---|---|
| Components, Browse the components | `DESIGN.md#components` on GitHub |
| Guides | `DESIGN.md` on GitHub |
| GitHub | the repository |

**Mode.** The whole site is dark, as in the draft, using the library's own dark tokens. It doesn't
follow the visitor's setting in this version, and there is no toggle, so everyone giving feedback sees
the same design. Palette is the default blue. `theme.js` isn't loaded, since nothing sets a choice.

**Type.** Outfit, as in the draft, self-hosted as a Latin-subset WOFF2 with its OFL licence beside it,
so visitors make no request to Google. It's set through `--font-display` and `--font-sans`, the way an
app would.

**JavaScript.** `site.js` handles the copy buttons and **Drop it again**, and nothing else. Without it,
the page reads correctly and the prompt can be selected by hand.

**Phone width.** The headline scales down and the actions wrap. The pile gets its own hand-chosen
positions at narrow widths, because scaling the desktop heap down would leave controls too small and
overlapping. The page itself never scrolls sideways.

## 5. Release notes page

`CHANGELOG.md` rendered as written, under the same nav and footer, in a readable single column. Each
version's heading gets an anchor. The page has no hand-kept copy of anything, so a release that updates
the changelog updates the site.

## 6. `AGENTS.md` and the prompt

**The prompt:**

> Use ai-ui for this project's UI. Read
> https://raw.githubusercontent.com/fuad-daoud/ai-ui/main/AGENTS.md first and follow it.

This differs from the drafts, which linked the GitHub page. An agent fetching the raw file gets plain
Markdown and no page around it. The prompt text is defined once in the generator.

**`AGENTS.md`**, at the repository root, stays short and sends agents to `DESIGN.md` for detail:

- What ai-ui is, in two sentences.
- How to install it for Go and for Next.js, pinned to the current version. It tells the agent to check
  which stack the project uses first.
- A link to the raw `DESIGN.md`, stating that it is the contract and should be read before writing UI.
- The working rules:
  - Use a library component before writing your own.
  - Use only names in the contract, and never invent a `ui-*` class.
  - Don't write a default as an attribute.
  - Theme through tokens and knobs in the app's own unlayered CSS, and name app-only tokens `--app-*`.
  - Never edit the library's files, and never use `!important` to beat it.
- For a project that already has CSS, a pointer to the compat section of `DESIGN.md`.

`DESIGN.md` is linked on `main`, so it can describe something newer than the pinned tag between
releases. That's acceptable for now, because docs and code normally land together at a release.

## 7. Wording for strangers

The public docs currently assume the reader knows Fuad's apps. These change, with no effect on
behaviour:

| File | Now | Becomes |
|---|---|---|
| `README.md` opening | "The shared UI for Fuad's platforms" | What the library is and who it's for, with links to the site and `AGENTS.md` |
| `README.md` | "not wired up in v0.1" | The same limitation, without the stale version |
| `package.json` description | "Shared UI for Fuad's platforms" | A neutral description, with the version unchanged |
| `doc.go` | "the shared UI library for Fuad's platforms" | A neutral description |
| `DESIGN.md` | "money moved this way" | "One app moved this way" |
| `CHANGELOG.md` 0.1.1 | "as verified moving money" | "as verified moving a real app" |

"Delete resume?" in the dialog example stays, because it's pinned by a parity fixture and reads fine.

## 8. Testing

**`go test ./...` in `site/`** builds the site into a temp directory and checks:

- Both pages exist, and every internal link and asset path resolves inside the output.
- Every version heading in `CHANGELOG.md` appears on the release notes page.
- The current version from `package.json` is in the footer.
- The prompt on the page contains the `AGENTS.md` URL, and `AGENTS.md` exists at the repository root.
- Every component in `css/components/`, apart from density, has its class somewhere on the home page.
- Nothing on either page references a host other than GitHub.

**`go test ./...` at the root** gains one test: `AGENTS.md` names the version in `package.json`, so a
release can't leave it pinned to an old one.

**Playwright**, with a separate config so the gallery's config doesn't change. It serves `site/out`
under `/ai-ui/` to prove the relative links hold:

- axe passes on both pages once the pile has settled.
- The page doesn't scroll sideways at 360px.
- No duplicate IDs.
- The pile settles within five seconds, and **Drop it again** replays it without clearing typed text.
- At 1280px and 360px, once settled, the centre of every control in the pile is the control itself,
  so no item covers another's button, input or checkbox, and none is cut off by the pile's edge.
- Tab reaches every control in the pile.
- Under reduced motion the components are already in place and **Drop it again** is hidden.
- Copy puts the exact prompt on the clipboard.
- Every asset loads.

**CI.** `ci.yml` gains a `site` job running `go vet` and `go test` in `site/`, since the root's
`go test ./...` skips a nested module. Existing steps don't change.

## 9. Publishing

A new `.github/workflows/pages.yml` runs on pushes to `main` and on demand. It builds with
`go run . -out out` in `site/` and deploys through `actions/upload-pages-artifact` and
`actions/deploy-pages`. Pages' source is set to GitHub Actions, once, on the repository.

The address is `https://fuad-daoud.github.io/ai-ui/`.

Pages has to be enabled before the first merge, or that first deploy fails. Enabling it serves nothing
until a deploy runs.

The work goes through a pull request. **Merging it makes the site and the rewritten docs public**, so
the merge is the point of no return for the first version, and it's Fuad's to make. The Pages setting
is also changed only with Fuad's go-ahead.

## 10. After the feedback round

Decided then, not now: the Components page, guides, a custom domain, and whether to revisit the home
page. The drafts canvas stays available for comparison.
