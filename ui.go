package ui

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

//go:embed dist/ui.css dist/ui-scoped.css js/theme.js
var raw embed.FS

// Prefix is the path Handler expects to be mounted under:
//
//	mux.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
const Prefix = "/_ui/"

// published maps each public asset name to its file in the repository.
var published = map[string]string{
	"ui.css":        "dist/ui.css",
	"ui-scoped.css": "dist/ui-scoped.css",
	"theme.js":      "js/theme.js",
}

var contentTypes = map[string]string{
	".css": "text/css; charset=utf-8",
	".js":  "text/javascript; charset=utf-8",
}

type asset struct {
	body  []byte
	hash  string
	ctype string
}

var assets = loadAssets()

func loadAssets() map[string]asset {
	m := make(map[string]asset, len(published))
	for name, file := range published {
		b, err := raw.ReadFile(file)
		if err != nil {
			panic("ui: embedded asset missing: " + file)
		}
		sum := sha256.Sum256(b)
		m[name] = asset{body: b, hash: hex.EncodeToString(sum[:])[:12], ctype: contentTypes[path.Ext(name)]}
	}
	return m
}

// Assets holds the published files, opened by their public names
// ("ui.css", "ui-scoped.css", "theme.js").
var Assets fs.FS = publicFS{}

type publicFS struct{}

func (publicFS) Open(name string) (fs.File, error) {
	file, ok := published[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return raw.Open(file)
}

// Href returns the URL of a published asset with its content hash, such as
// "/_ui/ui.css?v=3f2a9c0b1d4e". It panics on an unknown name, so a typo
// fails at the first render.
func Href(name string) string {
	a, ok := assets[name]
	if !ok {
		panic(fmt.Sprintf("ui: unknown asset %q (published: ui.css, ui-scoped.css, theme.js)", name))
	}
	return Prefix + name + "?v=" + a.hash
}

// Handler serves the published assets. A request whose ?v= matches the
// asset's hash is cached for a year as immutable. Any other request gets
// Cache-Control: no-cache, so a link written without Href is never pinned
// to a stale copy.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/")
		a, ok := assets[name]
		if !ok {
			http.NotFound(w, r)
			return
		}
		h := w.Header()
		h.Set("Content-Type", a.ctype)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("ETag", `"`+a.hash+`"`)
		if r.URL.Query().Get("v") == a.hash {
			h.Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			h.Set("Cache-Control", "no-cache")
		}
		http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(a.body))
	})
}
