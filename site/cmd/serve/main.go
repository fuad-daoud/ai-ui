// Command serve runs the ai-ui site on contabo, serving the pages built into
// it. Build the pages into www first, then the binary, from site/:
//
//	go run . -out cmd/serve/www
//	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/ai-ui-site ./cmd/serve
//
// The deploy steps are in the servers repo, services/ai-ui-site/NOTES.md.
package main

import (
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:www
var www embed.FS

func main() {
	addr := flag.String("addr", "127.0.0.1:8083", "listen address")
	flag.Parse()
	site, err := fs.Sub(www, "www")
	if err != nil {
		log.Fatal(err)
	}
	h, err := handler(site)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ai-ui site on http://%s/", *addr)
	log.Fatal(http.ListenAndServe(*addr, h))
}

// handler serves site. It refuses a site without a home page, so a binary
// built before the generator ran fails its health check and is rolled back.
func handler(site fs.FS) (http.Handler, error) {
	if _, err := fs.Stat(site, "index.html"); err != nil {
		return nil, errors.New("no pages built into this binary: run the generator into cmd/serve/www first")
	}
	files := http.FileServer(http.FS(site))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Cache-Control", "public, max-age=300")
		h.Set("X-Content-Type-Options", "nosniff")
		files.ServeHTTP(w, r)
	}), nil
}
