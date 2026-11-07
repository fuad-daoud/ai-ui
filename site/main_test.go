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

var pages = []string{"index.html", "releases/index.html"}

func TestPagesAndAssetsExist(t *testing.T) {
	out := built(t)
	for _, p := range append(pages, "assets/ui.css", "assets/site.css", "assets/site.js", "assets/outfit-latin.woff2", "assets/OFL.txt") {
		if _, err := os.Stat(filepath.Join(out, p)); err != nil {
			t.Error(err)
		}
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
			case ref == "", strings.HasPrefix(ref, "#"), strings.HasPrefix(ref, "data:"):
			case strings.HasPrefix(ref, "/"):
				t.Errorf("%s: %q isn't relative", p, ref)
			case strings.Contains(ref, "://"):
				if !strings.HasPrefix(ref, "https://github.com/fuad-daoud/ai-ui") && !strings.HasPrefix(ref, "https://raw.githubusercontent.com/fuad-daoud/ai-ui/") {
					t.Errorf("%s: %q leaves GitHub", p, ref)
				}
			default:
				target := filepath.Join(out, filepath.Dir(p), ref)
				if info, err := os.Stat(target); err == nil && info.IsDir() {
					target = filepath.Join(target, "index.html")
				}
				if _, err := os.Stat(target); err != nil {
					t.Errorf("%s: %q doesn't resolve", p, ref)
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

// The site links into these headings on GitHub, and a rename would break them quietly.
func TestLinkedHeadingsExist(t *testing.T) {
	for file, heading := range map[string]string{
		"../DESIGN.md":    "## Components",
		"../README.md":    "## Go apps",
		"../CHANGELOG.md": "# Changelog",
	} {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(heading) + `$`).Match(raw) {
			t.Errorf("%s has no %q", file, heading)
		}
	}
	readme, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(readme), "## Next.js apps") {
		t.Error(`README.md has no "## Next.js apps"`)
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

// A new component has to be placed in the pile before this passes.
func TestHomeShowsEveryComponent(t *testing.T) {
	out := built(t)
	classes := map[string]bool{}
	walk(parse(t, filepath.Join(out, "index.html")), func(n *html.Node) {
		for _, c := range strings.Fields(attr(n, "class")) {
			classes[c] = true
		}
	})
	files, err := filepath.Glob("../css/components/*.css")
	if err != nil || len(files) == 0 {
		t.Fatalf("no component CSS found: %v", err)
	}
	first := regexp.MustCompile(`(?m)^\.(ui-[a-z]+)`)
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
		if !classes[string(m[1])] {
			t.Errorf("home page has no .%s", m[1])
		}
	}
}

func TestHomeDoesNotCountComponents(t *testing.T) {
	page := text(parse(t, filepath.Join(built(t), "index.html")))
	if m := regexp.MustCompile(`(?i)\b(\d+|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve)\s+components\b`).FindString(page); m != "" {
		t.Errorf("home page counts components: %q", m)
	}
}

// fakeRepo holds only what build reads from the repository.
func fakeRepo(t *testing.T, pkg, changelog string) string {
	t.Helper()
	dir := t.TempDir()
	if pkg != "" {
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if changelog != "" {
		if err := os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte(changelog), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
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
