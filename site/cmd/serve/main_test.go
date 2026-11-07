package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestRefusesAnEmptySite(t *testing.T) {
	if _, err := handler(fstest.MapFS{".gitkeep": {}}); err == nil {
		t.Fatal("a binary built without the pages must not start")
	}
}

func TestServesPages(t *testing.T) {
	h, err := handler(fstest.MapFS{
		"index.html":          {Data: []byte("home")},
		"releases/index.html": {Data: []byte("releases")},
		"assets/site.css":     {Data: []byte("css")},
	})
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		path string
		code int
	}{
		{"/", http.StatusOK},
		{"/releases/", http.StatusOK},
		{"/releases", http.StatusMovedPermanently},
		{"/assets/site.css", http.StatusOK},
		{"/missing", http.StatusNotFound},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, c.path, nil))
		if rec.Code != c.code {
			t.Errorf("%s: got %d, want %d", c.path, rec.Code, c.code)
		}
		if c.code == http.StatusOK {
			if got := rec.Header().Get("Cache-Control"); got != "public, max-age=300" {
				t.Errorf("%s: Cache-Control %q", c.path, got)
			}
			if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("%s: X-Content-Type-Options %q", c.path, got)
			}
		}
	}
}
