// Package canon turns an HTML fragment into one canonical string, so markup
// from the Go partials and the React components can be compared (spec §8.1).
// test/canon.ts mirrors these rules; testdata/canon/cases.json pins both.
// Change them together.
package canon

import (
	"sort"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var void = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true, "hr": true, "img": true,
	"input": true, "link": true, "meta": true, "source": true, "track": true, "wbr": true,
}

var (
	textEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	attrEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
)

// HTML parses s as a fragment in <body> context and serialises it canonically.
func HTML(s string) (string, error) {
	ctx := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	nodes, err := html.ParseFragment(strings.NewReader(s), ctx)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, n := range nodes {
		write(&b, n)
	}
	return b.String(), nil
}

// words collapses runs of HTML whitespace (space, tab, LF, CR, FF) and trims.
func words(s string) string {
	return strings.Join(strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f'
	}), " ")
}

func write(b *strings.Builder, n *html.Node) {
	switch n.Type {
	case html.TextNode:
		if t := words(n.Data); t != "" {
			b.WriteString(textEscaper.Replace(t))
		}
	case html.ElementNode:
		b.WriteString("<" + n.Data)
		attrs := append([]html.Attribute(nil), n.Attr...)
		sort.Slice(attrs, func(i, j int) bool { return attrs[i].Key < attrs[j].Key })
		for _, a := range attrs {
			v := a.Val
			if a.Key == "class" {
				v = words(v)
			}
			b.WriteString(" " + a.Key)
			if v != "" {
				b.WriteString(`="` + attrEscaper.Replace(v) + `"`)
			}
		}
		b.WriteString(">")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			write(b, c)
		}
		if !void[n.Data] {
			b.WriteString("</" + n.Data + ">")
		}
	}
}
