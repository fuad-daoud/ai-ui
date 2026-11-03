package ui

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
)

//go:embed templates/*.tmpl
var partials embed.FS

// Option is one <option> of a select field.
type Option struct{ Value, Label string }

// Parse adds the ui.* partials and their helpers to t and returns t.
//
// App templates get dict (build a map for a partial's arguments) and href
// (ui.Href). The partials call only the ui_* helpers, so an app that later
// redefines dict or href does not break them.
func Parse(t *template.Template) (*template.Template, error) {
	if t == nil {
		return nil, errors.New("ui: Parse needs a template")
	}
	return t.Funcs(funcs).ParseFS(partials, "templates/*.tmpl")
}

var funcs = template.FuncMap{
	"dict":        dict,
	"href":        Href,
	"ui_str":      str,
	"ui_or":       or,
	"ui_variant":  variant,
	"ui_required": required,
}

func dict(kv ...any) (map[string]any, error) {
	if len(kv)%2 != 0 {
		return nil, fmt.Errorf("dict: odd number of arguments (%d)", len(kv))
	}
	m := make(map[string]any, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key %d is %T, want string", i/2, kv[i])
		}
		m[k] = kv[i+1]
	}
	return m, nil
}

// str renders a template value as text. A missing map key (nil) is "".
func str(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(x)
	}
}

func or(v any, fallback string) string {
	if s := str(v); s != "" {
		return s
	}
	return fallback
}

// variant returns v if it is one of allowed, else "". The caller then omits
// the attribute and the component's default applies (spec §6).
func variant(v any, allowed ...string) string {
	s := str(v)
	for _, a := range allowed {
		if s == a {
			return s
		}
	}
	return ""
}

func required(component, key string, v any) (string, error) {
	if s := str(v); s != "" {
		return s, nil
	}
	return "", fmt.Errorf("%s: %s is required", component, key)
}
