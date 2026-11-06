package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/fuad-daoud/ai-ui/internal/canon"
)

func fetch(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestGalleryPages(t *testing.T) {
	h, err := newHandler()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range palettes {
		for _, m := range modes {
			for _, suffix := range []string{"", "/dialog", "/compact", "/compact/dialog"} {
				path := "/p/" + p + "/" + m + suffix
				rec := fetch(t, h, path)
				if rec.Code != http.StatusOK {
					t.Fatalf("%s: status %d: %s", path, rec.Code, rec.Body)
				}
				body := rec.Body.String()
				for _, want := range []string{`data-theme="` + m + `"`, `class="ui-btn"`, `class="ui-field"`, `class="ui-check"`, `class="ui-pill"`, `href="/_ui/ui.css?v=`} {
					if !strings.Contains(body, want) {
						t.Errorf("%s: missing %s", path, want)
					}
				}
				if p == "mono" && !strings.Contains(body, `data-palette="mono"`) {
					t.Errorf("%s: missing data-palette", path)
				}
				if strings.Contains(body, "no value") {
					t.Errorf("%s: a missing key leaked into the markup", path)
				}
				if opens := strings.Contains(body, "showModal()"); opens != strings.HasSuffix(suffix, "/dialog") {
					t.Errorf("%s: showModal present = %v", path, opens)
				}
				if compact := strings.Contains(body, `data-density="compact"`); compact != strings.HasPrefix(suffix, "/compact") {
					t.Errorf("%s: data-density present = %v", path, compact)
				}
			}
		}
	}
}

func TestGalleryRoutes(t *testing.T) {
	h, err := newHandler()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/p/rose/light", "/p/blue/sepia", "/p/blue", "/p/blue/light/compact/x", "/fixture/library/other.css", "/fixture/document/theme.js"} {
		if rec := fetch(t, h, path); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", path, rec.Code)
		}
	}
	if rec := fetch(t, h, "/"); rec.Code != http.StatusFound || rec.Header().Get("Location") != "/p/blue/light" {
		t.Errorf("/: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
	if rec := fetch(t, h, "/_ui/ui.css"); rec.Code != http.StatusOK {
		t.Errorf("/_ui/ui.css: status %d", rec.Code)
	}
}

func TestGalleryDialogMatchesGolden(t *testing.T) {
	h, err := newHandler()
	if err != nil {
		t.Fatal(err)
	}
	body := fetch(t, h, "/p/blue/light").Body.String()
	i, j := strings.Index(body, "<dialog"), strings.Index(body, "</dialog>")
	if i < 0 || j < i {
		t.Fatal("no <dialog> in the gallery")
	}
	got, err := canon.HTML(body[i : j+len("</dialog>")])
	if err != nil {
		t.Fatal(err)
	}
	golden, err := os.ReadFile("../testdata/parity/dialog-basic.html")
	if err != nil {
		t.Fatal(err)
	}
	want, err := canon.HTML(string(golden))
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("gallery dialog differs from the documented markup\n got: %s\nwant: %s", got, want)
	}
}

func TestFixturePages(t *testing.T) {
	h, err := newHandler()
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{
		"/fixture/library/ui.css":                   `/_ui/ui.css?v=`,
		"/fixture/library/ui-scoped.css":            `/_ui/ui-scoped.css?v=`,
		"/fixture/library/ui.css?density=root":      `<html lang="en" data-density="compact">`,
		"/fixture/library/ui.css?density=container": `class="ui-stack" data-density="compact"`,
		"/fixture/document/ui-scoped.css":           `/_ui/ui-scoped.css?v=`,
	} {
		rec := fetch(t, h, path)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d: %s", path, rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("%s: missing %s", path, want)
		}
	}
	if body := fetch(t, h, "/fixture/document/none").Body.String(); strings.Contains(body, "/_ui/") {
		t.Error("/fixture/document/none links a library sheet")
	}
}
