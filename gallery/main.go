// Command gallery serves every ai-ui v0.1 component in each palette and mode,
// for screenshots and accessibility checks (spec §8.2):
//
//	go run ./gallery -addr 127.0.0.1:4317
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"html/template"
	"log"
	"net/http"
	"slices"

	ui "github.com/fuad-daoud/ai-ui"
)

//go:embed page.html
var pageSrc string

//go:embed fixture_library.html
var fixtureLibrarySrc string

//go:embed fixture_document.html
var fixtureDocumentSrc string

var (
	palettes   = []string{"blue", "mono"}
	modes      = []string{"light", "dark"}
	currencies = []ui.Option{{Value: "usd", Label: "US dollar"}, {Value: "jod", Label: "Jordanian dinar"}, {Value: "eur", Label: "Euro"}}
)

type link struct {
	Href, Label string
	Current     bool
}

// fixtureSheets are the sheets a fixture page can load. "none" loads no library
// CSS, for comparing a document with and without it.
var fixtureSheets = []string{"none", "ui.css", "ui-scoped.css"}

type fixture struct{ Sheet, Density string }

type page struct {
	Palette, Mode string
	Dialog        bool
	Compact       bool
	Links         []link
	Currencies    []ui.Option
	Icons         icons
}

// icons are the app's own markup for the partials' Icon argument.
type icons struct{ Download, Theme, Search, Avatar template.HTML }

var galleryIcons = icons{
	Download: `<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M8 2v8m-3.5-3.5L8 10l3.5-3.5M3 13h10"/></svg>`,
	Theme:    `<svg viewBox="0 0 16 16" aria-hidden="true"><circle cx="8" cy="8" r="6" fill="none" stroke="currentColor" stroke-width="1.6"/><path d="M8 2a6 6 0 0 1 0 12z" fill="currentColor"/></svg>`,
	Search:   `<svg class="ui-input-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" aria-hidden="true"><circle cx="7" cy="7" r="4.5"/><path d="m10.5 10.5 3 3"/></svg>`,
	Avatar:   `<span class="ui-avatar" aria-hidden="true">G</span>`,
}

func main() {
	addr := flag.String("addr", "127.0.0.1:4317", "listen address")
	flag.Parse()
	h, err := newHandler()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("gallery on http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, h))
}

func newHandler() (http.Handler, error) {
	set, err := ui.Parse(template.New("gallery"))
	if err != nil {
		return nil, err
	}
	for name, src := range map[string]string{
		"page":             pageSrc,
		"fixture-library":  fixtureLibrarySrc,
		"fixture-document": fixtureDocumentSrc,
	} {
		if _, err := set.New(name).Parse(src); err != nil {
			return nil, err
		}
	}
	mux := http.NewServeMux()
	mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/p/blue/light", http.StatusFound)
	})
	serve := func(dialog, compact bool) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			p, m := r.PathValue("palette"), r.PathValue("mode")
			if !slices.Contains(palettes, p) || !slices.Contains(modes, m) {
				http.NotFound(w, r)
				return
			}
			var b bytes.Buffer
			data := page{Palette: p, Mode: m, Dialog: dialog, Compact: compact, Links: links(p, m), Currencies: currencies, Icons: galleryIcons}
			if err := set.ExecuteTemplate(&b, "page", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = b.WriteTo(w)
		}
	}
	mux.HandleFunc("GET /p/{palette}/{mode}", serve(false, false))
	mux.HandleFunc("GET /p/{palette}/{mode}/dialog", serve(true, false))
	mux.HandleFunc("GET /p/{palette}/{mode}/compact", serve(false, true))
	mux.HandleFunc("GET /p/{palette}/{mode}/compact/dialog", serve(true, true))
	serveFixture := func(name string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sheet := r.PathValue("sheet")
			if !slices.Contains(fixtureSheets, sheet) {
				http.NotFound(w, r)
				return
			}
			var b bytes.Buffer
			data := fixture{Sheet: sheet, Density: r.URL.Query().Get("density")}
			if err := set.ExecuteTemplate(&b, name, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = b.WriteTo(w)
		}
	}
	mux.HandleFunc("GET /fixture/library/{sheet}", serveFixture("fixture-library"))
	mux.HandleFunc("GET /fixture/document/{sheet}", serveFixture("fixture-document"))
	return mux, nil
}

func links(palette, mode string) []link {
	var out []link
	for _, p := range palettes {
		for _, m := range modes {
			out = append(out, link{Href: "/p/" + p + "/" + m, Label: p + " " + m, Current: p == palette && m == mode})
		}
	}
	return out
}
