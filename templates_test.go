package ui_test

import (
	"html/template"
	"strings"
	"testing"

	ui "github.com/fuad-daoud/ai-ui"
)

func TestParseNil(t *testing.T) {
	if _, err := ui.Parse(nil); err == nil {
		t.Fatal("Parse(nil) succeeded")
	}
}

func TestPartialsSurviveAppFuncOverride(t *testing.T) {
	tmpl, err := ui.Parse(template.New("app"))
	if err != nil {
		t.Fatal(err)
	}
	tmpl.Funcs(template.FuncMap{
		"dict": func() string { return "app dict" },
		"href": func() string { return "app href" },
	})
	var b strings.Builder
	if err := tmpl.ExecuteTemplate(&b, "ui.button", map[string]any{"Label": "Save"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), ">Save</button>") {
		t.Fatalf("got %s", b.String())
	}
}

func render(t *testing.T, src string, data any) (string, error) {
	t.Helper()
	tmpl, err := ui.Parse(template.New("app"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpl.New("page").Parse(src); err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	err = tmpl.ExecuteTemplate(&b, "page", data)
	return b.String(), err
}

func TestDictErrors(t *testing.T) {
	for src, want := range map[string]string{
		`{{template "ui.button" (dict "Label")}}`:  "odd number of arguments",
		`{{template "ui.button" (dict 1 "Save")}}`: "key 0 is int",
	} {
		if _, err := render(t, src, nil); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err = %v, want %q", src, err, want)
		}
	}
}

func TestHrefFunc(t *testing.T) {
	got, err := render(t, `<link rel="stylesheet" href="{{href "ui.css"}}">`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := `href="` + ui.Href("ui.css") + `"`; !strings.Contains(got, want) {
		t.Errorf("got %s, want it to contain %s", got, want)
	}
	if _, err := render(t, `{{href "nope.css"}}`, nil); err == nil || !strings.Contains(err.Error(), `unknown asset "nope.css"`) {
		t.Errorf("err = %v, want unknown asset", err)
	}
}

// An icon is markup the app wrote, so a template.HTML value goes in as it is,
// and a plain string is escaped like any other argument.
func TestIconIsTheAppsMarkup(t *testing.T) {
	svg := template.HTML(`<svg class="ui-input-icon" aria-hidden="true"></svg>`)
	for _, c := range []struct {
		name, src, want string
		data            any
	}{
		{"button", `{{template "ui.button" (dict "Label" "Theme" "Shape" "square" "Icon" .)}}`, `data-shape="square" aria-label="Theme">` + string(svg) + `</button>`, svg},
		{"field", `{{template "ui.field" (dict "Label" "Search" "Name" "q" "Icon" .)}}`, `<div class="ui-input-group">` + "\n" + string(svg) + "\n" + `<input class="ui-input" id="q"`, svg},
		{"escaped", `{{template "ui.button" (dict "Label" "Go" "Icon" .)}}`, `>&lt;b&gt;Go</button>`, "<b>"},
	} {
		got, err := render(t, c.src, c.data)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if !strings.Contains(got, c.want) {
			t.Errorf("%s: got %s\nwant it to contain %s", c.name, got, c.want)
		}
	}
}
