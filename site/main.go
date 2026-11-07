// Command site builds the ai-ui website into a directory of static files:
//
//	go run . -out out
//	go run . -out out -serve 127.0.0.1:4318
//
// Every internal link is relative, so the output works under any path.
package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ui "github.com/fuad-daoud/ai-ui"
)

//go:embed templates/*.html
var templates embed.FS

//go:embed assets
var assets embed.FS

const (
	repoURL     = "https://github.com/fuad-daoud/ai-ui"
	contractURL = repoURL + "/blob/main/DESIGN.md"
	readmeURL   = repoURL + "/blob/main/README.md"
	agentsURL   = "https://raw.githubusercontent.com/fuad-daoud/ai-ui/main/AGENTS.md"
	prompt      = "Use ai-ui for this project's UI. Read " + agentsURL + " first and follow it."
)

// part is one heading's section of a Markdown file in the repository.
type part struct{ File, Heading string }

// guide is a page made of parts, in order.
type guide struct {
	Slug, Title, Blurb string
	Parts              []part
}

var guides = []guide{
	{"getting-started", "Getting started", "Install it in a Go app or a Next.js app.",
		[]part{{"README.md", "Go apps"}, {"README.md", "Next.js apps"}}},
	{"theming", "Theming", "Tokens, knobs, light and dark, palettes, and why your CSS always wins.",
		[]part{{"DESIGN.md", "Override model"}, {"DESIGN.md", "Tokens"}, {"DESIGN.md", "Modes and palettes"}}},
	{"density", "Density", "A compact size for working tools, for a whole page or one part of it.",
		[]part{{"DESIGN.md", "Density"}}},
	{"existing-apps", "Adding to an existing app", "Bring it into an app that already has its own CSS.",
		[]part{{"DESIGN.md", "Adopting ai-ui in an app that has its own CSS"}, {"DESIGN.md", "Choosing a stylesheet"}}},
}

// page is what every template renders from.
type page struct {
	Title, Description string
	Root               string // relative path back to the site root: "./", "../" or "../../"
	Section            string // the nav link this page sits under: "components", "guides", "releases", or "" for home
	Version            string
	Prompt             string
	RepoURL            string
	ContractURL        string
	ReadmeURL          string
	Changelog          template.HTML
	Currencies         []ui.Option
	Guides             []guide   // every guide, for the index and the side list
	Guide              guide     // the guide this page shows
	Sections           []section // that guide's parts, rendered
}

func main() {
	out := flag.String("out", "out", "directory to write the site to")
	repo := flag.String("repo", "..", "the repository root")
	serve := flag.String("serve", "", "after building, serve the site under /ai-ui/ on this address")
	flag.Parse()
	if err := build(*out, *repo); err != nil {
		log.Fatal(err)
	}
	if *serve == "" {
		return
	}
	mux := http.NewServeMux()
	mux.Handle("/ai-ui/", http.StripPrefix("/ai-ui/", http.FileServer(http.Dir(*out))))
	mux.Handle("/{$}", http.RedirectHandler("/ai-ui/", http.StatusFound))
	log.Printf("site on http://%s/ai-ui/", *serve)
	log.Fatal(http.ListenAndServe(*serve, mux))
}

// build renders everything before it writes anything, so a broken source
// leaves no output behind for CI to publish.
func build(out, repo string) error {
	version, err := readVersion(repo)
	if err != nil {
		return err
	}
	changelog, err := renderChangelog(repo)
	if err != nil {
		return err
	}
	sources, err := readSources(repo)
	if err != nil {
		return err
	}
	set, err := ui.Parse(template.New("site"))
	if err != nil {
		return err
	}
	if set, err = set.ParseFS(templates, "templates/*.html"); err != nil {
		return err
	}
	base := page{
		Version:     version,
		Prompt:      prompt,
		RepoURL:     repoURL,
		ContractURL: contractURL,
		ReadmeURL:   readmeURL,
		Guides:      guides,
	}

	files := map[string][]byte{}

	home := base
	home.Title = "ai-ui: the UI kit your agent can read"
	home.Description = "A small pile of real, working UI components, written down in one contract your coding agent can read."
	home.Root = "./"
	home.Currencies = []ui.Option{{Value: "usd", Label: "US dollar"}, {Value: "eur", Label: "Euro"}}
	if files["index.html"], err = render(set, "home.html", home); err != nil {
		return err
	}

	releases := base
	releases.Title = "Release notes: ai-ui"
	releases.Description = "What changed in each release of ai-ui."
	releases.Root = "../"
	releases.Section = "releases"
	releases.Changelog = changelog
	if files["releases/index.html"], err = render(set, "releases.html", releases); err != nil {
		return err
	}

	if err := buildGuides(set, base, sources, files); err != nil {
		return err
	}

	if files["assets/ui.css"], err = fs.ReadFile(ui.Assets, "ui.css"); err != nil {
		return err
	}
	err = fs.WalkDir(assets, "assets", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		files[p], err = fs.ReadFile(assets, p)
		return err
	})
	if err != nil {
		return err
	}

	for name, body := range files {
		path := filepath.Join(out, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, body, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// readSources reads the Markdown the guides and the Components page are built from.
func readSources(repo string) (map[string][]byte, error) {
	src := map[string][]byte{}
	for _, name := range []string{"DESIGN.md", "README.md"} {
		raw, err := os.ReadFile(filepath.Join(repo, name))
		if err != nil {
			return nil, err
		}
		src[name] = raw
	}
	return src, nil
}

func buildGuides(set *template.Template, base page, src, files map[string][]byte) error {
	index := base
	index.Title = "Guides: ai-ui"
	index.Description = "Install ai-ui, theme it, size it, and bring it into an app that has its own CSS."
	index.Root = "../"
	index.Section = "guides"
	var err error
	if files["guides/index.html"], err = render(set, "guides.html", index); err != nil {
		return err
	}
	for _, g := range guides {
		p := index
		p.Title = g.Title + ": ai-ui"
		p.Description = g.Blurb
		p.Root = "../../"
		p.Guide = g
		p.Sections = nil
		for _, pt := range g.Parts {
			s, err := findSection(src[pt.File], pt.File, pt.Heading, false)
			if err != nil {
				return err
			}
			p.Sections = append(p.Sections, s)
		}
		if files["guides/"+g.Slug+"/index.html"], err = render(set, "guide.html", p); err != nil {
			return err
		}
	}
	return nil
}

func render(set *template.Template, name string, data page) ([]byte, error) {
	var b bytes.Buffer
	if err := set.ExecuteTemplate(&b, name, data); err != nil {
		return nil, fmt.Errorf("render %s: %w", name, err)
	}
	return b.Bytes(), nil
}

func readVersion(repo string) (string, error) {
	raw, err := os.ReadFile(filepath.Join(repo, "package.json"))
	if err != nil {
		return "", err
	}
	var pkg struct{ Version string }
	if err := json.Unmarshal(raw, &pkg); err != nil {
		return "", fmt.Errorf("package.json: %w", err)
	}
	if pkg.Version == "" {
		return "", errors.New("package.json has no version")
	}
	return pkg.Version, nil
}

// renderChangelog leaves goldmark's default in place: raw HTML in the Markdown
// is dropped, not passed through.
func renderChangelog(repo string) (template.HTML, error) {
	raw, err := os.ReadFile(filepath.Join(repo, "CHANGELOG.md"))
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(raw)) == "" {
		return "", errors.New("CHANGELOG.md is empty")
	}
	var b bytes.Buffer
	if err := markdown.Convert(raw, &b); err != nil {
		return "", fmt.Errorf("CHANGELOG.md: %w", err)
	}
	return template.HTML(b.String()), nil
}
