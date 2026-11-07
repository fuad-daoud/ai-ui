package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/html"
)

var (
	builtOnce sync.Once
	builtDir  string
	builtErr  error
)

// built builds the real site once, from the real repository.
func built(t *testing.T) string {
	t.Helper()
	builtOnce.Do(func() {
		builtDir, builtErr = os.MkdirTemp("", "ai-ui-site")
		if builtErr == nil {
			builtErr = build(builtDir, "..")
		}
	})
	if builtErr != nil {
		t.Fatal(builtErr)
	}
	return builtDir
}

func parse(t *testing.T, path string) *html.Node {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	n, err := html.Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func walk(n *html.Node, visit func(*html.Node)) {
	visit(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, visit)
	}
}

func attr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func text(n *html.Node) string {
	var b strings.Builder
	walk(n, func(c *html.Node) {
		if c.Type == html.TextNode {
			b.WriteString(c.Data)
			b.WriteString(" ")
		}
	})
	return b.String()
}

var pages = []string{
	"index.html",
	"releases/index.html",
	"guides/index.html",
	"guides/getting-started/index.html",
	"guides/theming/index.html",
	"guides/density/index.html",
	"guides/existing-apps/index.html",
	"components/index.html",
}

func TestPagesAndAssetsExist(t *testing.T) {
	out := built(t)
	for _, p := range append(pages, "assets/ui.css", "assets/site.css", "assets/site.js", "assets/outfit-latin.woff2", "assets/OFL.txt") {
		if _, err := os.Stat(filepath.Join(out, p)); err != nil {
			t.Error(err)
		}
	}
}

func hasID(t *testing.T, path, id string) bool {
	t.Helper()
	found := false
	walk(parse(t, path), func(n *html.Node) {
		if n.Type == html.ElementNode && attr(n, "id") == id {
			found = true
		}
	})
	return found
}

// The nav's Components and Guides go to site pages, and each page marks the
// part of the site it belongs to, and nothing else.
func TestNavPointsAtSitePagesAndMarksTheCurrentOne(t *testing.T) {
	out := built(t)
	want := map[string]string{
		"index.html":                      "",
		"releases/index.html":             "Release notes",
		"components/index.html":           "Components",
		"guides/index.html":               "Guides",
		"guides/theming/index.html":       "Guides",
		"guides/existing-apps/index.html": "Guides",
	}
	for p, current := range want {
		var marked []string
		walk(parse(t, filepath.Join(out, p)), func(n *html.Node) {
			if n.Type != html.ElementNode || n.Data != "nav" || attr(n, "class") != "site-nav" {
				return
			}
			walk(n, func(a *html.Node) {
				if a.Type != html.ElementNode || a.Data != "a" {
					return
				}
				name := strings.TrimSpace(text(a))
				if (name == "Components" || name == "Guides") && strings.Contains(attr(a, "href"), "://") {
					t.Errorf("%s: nav %s leaves the site: %s", p, name, attr(a, "href"))
				}
				if attr(a, "aria-current") == "page" {
					marked = append(marked, name)
				}
			})
		})
		if got := strings.Join(marked, ","); got != current {
			t.Errorf("%s: nav marks %q, want %q", p, got, current)
		}
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	out := built(t)
	for _, p := range pages {
		seen := map[string]bool{}
		walk(parse(t, filepath.Join(out, p)), func(n *html.Node) {
			if id := attr(n, "id"); n.Type == html.ElementNode && id != "" {
				if seen[id] {
					t.Errorf("%s: id %q is used twice", p, id)
				}
				seen[id] = true
			}
		})
	}
}

// Every internal path must resolve in the output, and anything else must be GitHub.
func TestLinksResolveAndStayOnGitHub(t *testing.T) {
	out := built(t)
	for _, p := range pages {
		walk(parse(t, filepath.Join(out, p)), func(n *html.Node) {
			if n.Type != html.ElementNode {
				return
			}
			ref := attr(n, "href")
			if ref == "" {
				ref = attr(n, "src")
			}
			switch {
			case ref == "", strings.HasPrefix(ref, "data:"):
			case strings.HasPrefix(ref, "#"):
				if !hasID(t, filepath.Join(out, p), ref[1:]) {
					t.Errorf("%s: %q has no target on the page", p, ref)
				}
			case strings.HasPrefix(ref, "/"):
				t.Errorf("%s: %q isn't relative", p, ref)
			case strings.Contains(ref, "://"):
				if !strings.HasPrefix(ref, "https://github.com/fuad-daoud/ai-ui") && !strings.HasPrefix(ref, "https://raw.githubusercontent.com/fuad-daoud/ai-ui/") {
					t.Errorf("%s: %q leaves GitHub", p, ref)
				}
			default:
				file, frag, _ := strings.Cut(ref, "#")
				target := filepath.Join(out, filepath.Dir(p), file)
				if info, err := os.Stat(target); err == nil && info.IsDir() {
					target = filepath.Join(target, "index.html")
				}
				if _, err := os.Stat(target); err != nil {
					t.Errorf("%s: %q doesn't resolve", p, ref)
				} else if frag != "" && !hasID(t, target, frag) {
					t.Errorf("%s: %q has no #%s", p, ref, frag)
				}
			}
		})
	}
}

func TestCSSFontResolves(t *testing.T) {
	out := built(t)
	css, err := os.ReadFile(filepath.Join(out, "assets/site.css"))
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range regexp.MustCompile(`url\("?([^")]+)"?\)`).FindAllStringSubmatch(string(css), -1) {
		if strings.Contains(m[1], "://") {
			t.Errorf("site.css reaches out to %s", m[1])
		} else if _, err := os.Stat(filepath.Join(out, "assets", m[1])); err != nil {
			t.Error(err)
		}
	}
}

func TestReleaseNotesShowEveryVersion(t *testing.T) {
	out := built(t)
	changelog, err := os.ReadFile("../CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	headings := map[string]bool{} // heading text, for those with an anchor
	walk(parse(t, filepath.Join(out, "releases/index.html")), func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h2" && attr(n, "id") != "" {
			headings[strings.Fields(text(n))[0]] = true
		}
	})
	versions := regexp.MustCompile(`(?m)^## (\S+)`).FindAllStringSubmatch(string(changelog), -1)
	if len(versions) == 0 {
		t.Fatal("no versions found in CHANGELOG.md")
	}
	for _, v := range versions {
		if !headings[v[1]] {
			t.Errorf("release notes have no anchored heading for %s", v[1])
		}
	}
}

func TestFooterNamesCurrentVersion(t *testing.T) {
	out := built(t)
	version, err := readVersion("..")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pages {
		found := false
		walk(parse(t, filepath.Join(out, p)), func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "footer" && strings.Contains(text(n), "v"+version) {
				found = true
			}
		})
		if !found {
			t.Errorf("%s: footer doesn't name v%s", p, version)
		}
	}
}


func TestPromptPointsAtAgentsFile(t *testing.T) {
	if !strings.Contains(prompt, agentsURL) {
		t.Errorf("prompt doesn't contain %s", agentsURL)
	}
	if _, err := os.Stat("../AGENTS.md"); err != nil {
		t.Error(err)
	}
	if !strings.Contains(text(parse(t, filepath.Join(built(t), "index.html"))), prompt) {
		t.Error("home page doesn't show the prompt")
	}
}

// componentClasses is the first .ui-* class of each file in css/components/, except density.
func componentClasses(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("../css/components/*.css")
	if err != nil || len(files) == 0 {
		t.Fatalf("no component CSS found: %v", err)
	}
	first := regexp.MustCompile(`(?m)^\.(ui-[a-z]+)`)
	var classes []string
	for _, f := range files {
		if filepath.Base(f) == "density.css" {
			continue
		}
		css, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		m := first.FindSubmatch(css)
		if m == nil {
			t.Errorf("%s: no component class found", f)
			continue
		}
		classes = append(classes, string(m[1]))
	}
	return classes
}

// classesIn collects the classes used under every element that has class root.
func classesIn(n *html.Node, root string) map[string]bool {
	found := map[string]bool{}
	walk(n, func(c *html.Node) {
		if c.Type != html.ElementNode || !strings.Contains(" "+attr(c, "class")+" ", " "+root+" ") {
			return
		}
		walk(c, func(d *html.Node) {
			for _, cl := range strings.Fields(attr(d, "class")) {
				found[cl] = true
			}
		})
	})
	return found
}

// A new component has to be placed in the pile before this passes.
func TestHomeShowsEveryComponent(t *testing.T) {
	found := classesIn(parse(t, filepath.Join(built(t), "index.html")), "pile")
	for _, c := range componentClasses(t) {
		if !found[c] {
			t.Errorf("home page has no .%s", c)
		}
	}
}

// A new component needs a demo on the Components page before this passes.
func TestComponentsPageDemosEveryComponent(t *testing.T) {
	found := classesIn(parse(t, filepath.Join(built(t), "components/index.html")), "demo")
	for _, c := range componentClasses(t) {
		if !found[c] {
			t.Errorf("Components page has no demo using .%s", c)
		}
	}
}

// Every ### under ## Components in the contract must be on the Components page.
func TestEveryContractComponentIsListed(t *testing.T) {
	raw, err := os.ReadFile("../DESIGN.md")
	if err != nil {
		t.Fatal(err)
	}
	_, rest, ok := strings.Cut(string(raw), "\n## Components\n")
	if !ok {
		t.Fatal(`DESIGN.md has no "## Components"`)
	}
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	listed := map[string]bool{}
	for _, c := range components {
		listed[c.Heading] = true
	}
	for _, m := range regexp.MustCompile(`(?m)^### (.+)$`).FindAllStringSubmatch(rest, -1) {
		if h := strings.ReplaceAll(m[1], "`", ""); !listed[h] {
			t.Errorf("DESIGN.md's %q isn't on the Components page", h)
		}
	}
}

// The contract shows Go template calls as code; the site must not run them.
func TestContractCodeShowsAsText(t *testing.T) {
	page := text(parse(t, filepath.Join(built(t), "components/index.html")))
	if !strings.Contains(page, `{{template "ui.button" (dict "Label" "Save"`) {
		t.Error("the Components page doesn't show the Go usage from DESIGN.md as text")
	}
}

func TestNoPageCountsComponents(t *testing.T) {
	count := regexp.MustCompile(`(?i)\b(\d+|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve)\s+components\b`)
	for _, p := range pages {
		if m := count.FindString(text(parse(t, filepath.Join(built(t), p)))); m != "" {
			t.Errorf("%s counts components: %q", p, m)
		}
	}
}

// fakeRepo holds only what build reads from the repository. DESIGN.md and
// README.md are the real ones unless a test changes them.
func fakeRepo(t *testing.T, pkg, changelog string) string {
	t.Helper()
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if pkg != "" {
		write("package.json", pkg)
	}
	if changelog != "" {
		write("CHANGELOG.md", changelog)
	}
	for _, name := range []string{"DESIGN.md", "README.md"} {
		raw, err := os.ReadFile(filepath.Join("..", name))
		if err != nil {
			t.Fatal(err)
		}
		write(name, string(raw))
	}
	return dir
}

// A guide is built from headings in DESIGN.md and README.md, so renaming one
// must fail the build rather than publish a guide with a hole in it.
func TestBrokenGuideSourcesFailAndWriteNothing(t *testing.T) {
	for name, change := range map[string]func(dir string) error{
		"no DESIGN.md": func(dir string) error { return os.Remove(filepath.Join(dir, "DESIGN.md")) },
		"heading renamed": func(dir string) error {
			p := filepath.Join(dir, "DESIGN.md")
			raw, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			return os.WriteFile(p, []byte(strings.Replace(string(raw), "\n## Density\n", "\n## Sizes\n", 1)), 0o644)
		},
	} {
		t.Run(name, func(t *testing.T) {
			repo := fakeRepo(t, `{"version":"1.0.0"}`, "# Changelog\n\n## 1.0.0\n")
			if err := change(repo); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(t.TempDir(), "out")
			if err := build(out, repo); err == nil {
				t.Fatal("build succeeded")
			}
			if _, err := os.Stat(out); !os.IsNotExist(err) {
				t.Errorf("build left %s behind", out)
			}
		})
	}
}

func TestGuidesShowTheirSections(t *testing.T) {
	out := built(t)
	for _, g := range guides {
		page := parse(t, filepath.Join(out, "guides", g.Slug, "index.html"))
		ids := map[string]bool{}
		walk(page, func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "h2" {
				ids[strings.TrimSpace(text(n))] = true
			}
		})
		for _, p := range g.Parts {
			if !ids[p.Heading] {
				t.Errorf("%s: no h2 %q", g.Slug, p.Heading)
			}
		}
	}
}

// CI publishes whatever build writes, so a broken source must leave nothing behind.
func TestBrokenSourcesFailAndWriteNothing(t *testing.T) {
	cases := map[string][2]string{
		"no package.json":  {"", "## 1.0.0\n"},
		"bad package.json": {"{", "## 1.0.0\n"},
		"no version":       {`{"name":"ai-ui"}`, "## 1.0.0\n"},
		"no changelog":     {`{"version":"1.0.0"}`, ""},
		"empty changelog":  {`{"version":"1.0.0"}`, "  \n"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "out")
			if err := build(out, fakeRepo(t, c[0], c[1])); err == nil {
				t.Fatal("build succeeded")
			}
			if _, err := os.Stat(out); !os.IsNotExist(err) {
				t.Errorf("build left %s behind", out)
			}
		})
	}
}

func TestChangelogRawHTMLIsDropped(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")
	repo := fakeRepo(t, `{"version":"1.0.0"}`, "# Changelog\n\n## 1.0.0\n\n<script>alert(1)</script>\n\n<img src=x onerror=alert(2)>\n")
	if err := build(out, repo); err != nil {
		t.Fatal(err)
	}
	page, err := os.ReadFile(filepath.Join(out, "releases/index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(page), "alert(") {
		t.Errorf("raw HTML reached the page:\n%s", page)
	}
}
