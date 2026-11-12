package ui_test

import (
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	ui "github.com/fuad-daoud/ai-ui"
)

func mux() http.Handler {
	m := http.NewServeMux()
	m.Handle(ui.Prefix, http.StripPrefix(ui.Prefix, ui.Handler()))
	return m
}

func get(t *testing.T, h http.Handler, target string, header map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHrefCarriesContentHash(t *testing.T) {
	for _, name := range []string{"ui.css", "theme.js", "ui.js"} {
		if got := ui.Href(name); !regexp.MustCompile(`^/_ui/` + regexp.QuoteMeta(name) + `\?v=[0-9a-f]{12}$`).MatchString(got) {
			t.Errorf("Href(%q) = %q", name, got)
		}
	}
}

func TestHrefPanicsOnUnknownName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || !strings.Contains(r.(string), `"nope.css"`) {
			t.Fatalf("recover() = %v, want a panic naming nope.css", r)
		}
	}()
	ui.Href("nope.css")
}

func TestHandlerServesVersionedCSSAsImmutable(t *testing.T) {
	rec := get(t, mux(), ui.Href("ui.css"), nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/css; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
		t.Errorf("Cache-Control = %q", cc)
	}
	want, err := os.ReadFile("dist/ui.css")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Body.String() != string(want) {
		t.Error("body differs from dist/ui.css")
	}
}

func TestHandlerServesTheScripts(t *testing.T) {
	for _, name := range []string{"theme.js", "ui.js"} {
		rec := get(t, mux(), ui.Href(name), nil)
		if ct := rec.Header().Get("Content-Type"); rec.Code != http.StatusOK || ct != "text/javascript; charset=utf-8" {
			t.Fatalf("%s: status %d, Content-Type %q", name, rec.Code, ct)
		}
	}
}

func TestHandlerStaleOrMissingVersionIsNoCache(t *testing.T) {
	for _, target := range []string{"/_ui/ui.css", "/_ui/ui.css?v=000000000000"} {
		rec := get(t, mux(), target, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d", target, rec.Code)
		}
		if cc := rec.Header().Get("Cache-Control"); cc != "no-cache" {
			t.Errorf("%s: Cache-Control = %q, want no-cache", target, cc)
		}
	}
}

func TestHandlerRevalidatesWithETag(t *testing.T) {
	first := get(t, mux(), ui.Href("ui.css"), nil)
	etag := first.Header().Get("ETag")
	if etag == "" {
		t.Fatal("no ETag")
	}
	if rec := get(t, mux(), ui.Href("ui.css"), map[string]string{"If-None-Match": etag}); rec.Code != http.StatusNotModified {
		t.Errorf("status %d, want 304", rec.Code)
	}
}

func TestHandlerRejectsUnknownAndNonGET(t *testing.T) {
	if rec := get(t, ui.Handler(), "/nope.css", nil); rec.Code != http.StatusNotFound {
		t.Errorf("unknown asset: status %d", rec.Code)
	}
	if rec := get(t, ui.Handler(), "/dist/ui.css", nil); rec.Code != http.StatusNotFound {
		t.Errorf("source path: status %d", rec.Code)
	}
	rec := httptest.NewRecorder()
	ui.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ui.css", nil))
	if rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("Allow") != "GET, HEAD" {
		t.Errorf("POST: status %d, Allow %q", rec.Code, rec.Header().Get("Allow"))
	}
}

func TestAssetsOpensByPublicName(t *testing.T) {
	got, err := fs.ReadFile(ui.Assets, "ui.css")
	if err != nil {
		t.Fatal(err)
	}
	want, _ := os.ReadFile("dist/ui.css")
	if string(got) != string(want) {
		t.Error("Assets ui.css differs from dist/ui.css")
	}
	if _, err := ui.Assets.Open("dist/ui.css"); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Open(dist/ui.css) err = %v, want fs.ErrNotExist", err)
	}
}

func TestScopedStylesheetIsPublished(t *testing.T) {
	rec := get(t, mux(), ui.Href("ui-scoped.css"), nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/css; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
		t.Errorf("Cache-Control = %q", cc)
	}
	want, err := os.ReadFile("dist/ui-scoped.css")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Body.String() != string(want) {
		t.Error("body differs from dist/ui-scoped.css")
	}
	got, err := fs.ReadFile(ui.Assets, "ui-scoped.css")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Error("Assets ui-scoped.css differs from dist/ui-scoped.css")
	}
}
