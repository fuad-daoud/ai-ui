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

var (
	palettes   = []string{"blue", "mono"}
	modes      = []string{"light", "dark"}
	currencies = []ui.Option{{Value: "usd", Label: "US dollar"}, {Value: "jod", Label: "Jordanian dinar"}, {Value: "eur", Label: "Euro"}}
)

type link struct {
	Href, Label string
	Current     bool
}

type page struct {
	Palette, Mode string
	Dialog        bool
	Links         []link
	Currencies    []ui.Option
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
	if _, err := set.New("page").Parse(pageSrc); err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/p/blue/light", http.StatusFound)
	})
	serve := func(dialog bool) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			p, m := r.PathValue("palette"), r.PathValue("mode")
			if !slices.Contains(palettes, p) || !slices.Contains(modes, m) {
				http.NotFound(w, r)
				return
			}
			var b bytes.Buffer
			data := page{Palette: p, Mode: m, Dialog: dialog, Links: links(p, m), Currencies: currencies}
			if err := set.ExecuteTemplate(&b, "page", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = b.WriteTo(w)
		}
	}
	mux.HandleFunc("GET /p/{palette}/{mode}", serve(false))
	mux.HandleFunc("GET /p/{palette}/{mode}/dialog", serve(true))
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
