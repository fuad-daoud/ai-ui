package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	mdtext "github.com/yuin/goldmark/text"
)

// rawURL serves a repository file's bytes, which is what an image needs.
const rawURL = "https://raw.githubusercontent.com/fuad-daoud/ai-ui/main/"

// emptyHeader is a header cell with nothing in it, aligned or not.
var emptyHeader = regexp.MustCompile(`<th([^>]*)></th>`)

// markdown renders every Markdown source on the site: tables and heading IDs
// on, and raw HTML dropped (goldmark's default), so a source file can't put
// script on the page.
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.Table),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

// section is one heading of a Markdown file and what follows it.
type section struct {
	ID, Title string        // the heading's id and plain text
	Heading   template.HTML // the heading element, as an h2
	Body      template.HTML
	IDs       []string // every heading id in the section, its own first
}

// findSection returns the section under the top-level heading whose plain
// text is heading. It runs to the next heading at the same level or higher,
// or, with toAnyHeading, to the next heading of any level. Headings move so
// the section's own heading is an h2, and links point at GitHub where the
// site doesn't hold their target. A link that climbs out of the repository is
// an error.
func findSection(src []byte, file, heading string, toAnyHeading bool) (section, error) {
	doc := markdown.Parser().Parse(mdtext.NewReader(src))
	var h *ast.Heading
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if c, ok := n.(*ast.Heading); ok && plain(c, src) == heading {
			h = c
			break
		}
	}
	if h == nil {
		return section{}, fmt.Errorf("%s has no heading %q", file, heading)
	}
	body := ast.NewDocument()
	for n := h.NextSibling(); n != nil; {
		next := n.NextSibling()
		if c, ok := n.(*ast.Heading); ok && (toAnyHeading || c.Level <= h.Level) {
			break
		}
		body.AppendChild(body, n) // moves n out of doc
		n = next
	}
	s := section{Title: plain(h, src)}
	head := ast.NewDocument()
	head.AppendChild(head, h)
	var err error
	if s.IDs, err = fitSection(head, body, file, 2-h.Level); err != nil {
		return section{}, err
	}
	s.ID = s.IDs[0]
	if s.Heading, err = renderMD(head, src); err != nil {
		return section{}, fmt.Errorf("%s: %w", file, err)
	}
	if s.Body, err = renderMD(body, src); err != nil {
		return section{}, fmt.Errorf("%s: %w", file, err)
	}
	return s, nil
}

// fitSection moves every heading in head and body by shift, no deeper than h6,
// and points their links at GitHub where the site doesn't hold the target. It
// returns the heading IDs, the section's own first.
func fitSection(head, body *ast.Document, file string, shift int) ([]string, error) {
	var ids []string
	var links []ast.Node
	for _, d := range []*ast.Document{head, body} {
		_ = ast.Walk(d, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}
			switch c := n.(type) {
			case *ast.Heading:
				c.Level = min(c.Level+shift, 6)
				id, _ := c.AttributeString("id")
				b, _ := id.([]byte)
				ids = append(ids, string(b))
			case *ast.Link, *ast.Image:
				links = append(links, n)
			}
			return ast.WalkContinue, nil
		})
	}
	own := map[string]bool{}
	for _, id := range ids {
		own[id] = true
	}
	for _, n := range links {
		var err error
		switch c := n.(type) {
		case *ast.Link:
			c.Destination, err = repoLink(c.Destination, file, false, own)
		case *ast.Image:
			c.Destination, err = repoLink(c.Destination, file, true, own)
		}
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// repoLink resolves a link written in a repository file at the root, as GitHub
// would: a path goes to the file on GitHub (an image to its raw bytes), and a
// fragment the section doesn't hold goes to that place in the whole file. A
// link with a scheme or a host is kept.
func repoLink(dest []byte, file string, image bool, own map[string]bool) ([]byte, error) {
	u, err := url.Parse(string(dest))
	if err != nil {
		return nil, fmt.Errorf("%s: link %q: %w", file, dest, err)
	}
	if u.Scheme != "" || u.Host != "" {
		return dest, nil
	}
	if u.Path == "" {
		if u.Fragment == "" || own[u.Fragment] {
			return dest, nil
		}
		return []byte(repoURL + "/blob/main/" + file + "#" + u.Fragment), nil
	}
	p := path.Clean(strings.TrimPrefix(u.Path, "/"))
	if p == ".." || strings.HasPrefix(p, "../") {
		return nil, fmt.Errorf("%s: link %q leaves the repository", file, dest)
	}
	out := repoURL + "/blob/main/" + p
	if image {
		out = rawURL + p
	}
	if u.Fragment != "" {
		out += "#" + u.Fragment
	}
	return []byte(out), nil
}

// plain is a node's text with the Markdown taken out: "Button: `.ui-btn`"
// reads "Button: .ui-btn".
func plain(n ast.Node, src []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := c.(type) {
			case *ast.Text:
				b.Write(t.Segment.Value(src))
			case *ast.String:
				b.Write(t.Value)
			}
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

// renderMD renders doc and frames each table as a library table, in a named
// region that takes focus so a table wider than a phone can scroll. An empty
// header cell, the corner of a table like DESIGN.md's, names no column, so it
// becomes a plain cell. Raw HTML is dropped, so the renderer is the only
// source of these tags.
func renderMD(doc ast.Node, src []byte) (template.HTML, error) {
	var b bytes.Buffer
	if err := markdown.Renderer().Render(&b, src, doc); err != nil {
		return "", err
	}
	out := strings.ReplaceAll(b.String(), "<table>", `<div class="ui-table-wrap" tabindex="0" role="region" aria-label="Table"><table class="ui-table">`)
	out = strings.ReplaceAll(out, "</table>", "</table></div>")
	out = emptyHeader.ReplaceAllString(out, "<td$1></td>")
	return template.HTML(out), nil
}
