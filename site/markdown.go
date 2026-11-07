package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	mdtext "github.com/yuin/goldmark/text"
)

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
}

// findSection returns the section under the top-level heading whose plain
// text is heading. It runs to the next heading at the same level or higher,
// or, with toAnyHeading, to the next heading of any level. Headings move so
// the section's own heading is an h2, and relative links point at the file
// on GitHub.
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
	shift := 2 - h.Level
	h.Level = 2
	_ = ast.Walk(body, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch c := n.(type) {
		case *ast.Heading:
			c.Level += shift
		case *ast.Link:
			if d := string(c.Destination); !strings.Contains(d, ":") && !strings.HasPrefix(d, "#") {
				c.Destination = []byte(repoURL + "/blob/main/" + d)
			}
		}
		return ast.WalkContinue, nil
	})
	s := section{Title: plain(h, src)}
	if id, ok := h.AttributeString("id"); ok {
		if b, ok := id.([]byte); ok {
			s.ID = string(b)
		}
	}
	head := ast.NewDocument()
	head.AppendChild(head, h)
	var err error
	if s.Heading, err = renderMD(head, src); err != nil {
		return section{}, fmt.Errorf("%s: %w", file, err)
	}
	if s.Body, err = renderMD(body, src); err != nil {
		return section{}, fmt.Errorf("%s: %w", file, err)
	}
	return s, nil
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

// renderMD renders doc and frames each table as a library table. Raw HTML is
// dropped, so the renderer is the only source of a <table> tag.
func renderMD(doc ast.Node, src []byte) (template.HTML, error) {
	var b bytes.Buffer
	if err := markdown.Renderer().Render(&b, src, doc); err != nil {
		return "", err
	}
	out := strings.ReplaceAll(b.String(), "<table>", `<div class="ui-table-wrap" tabindex="0"><table class="ui-table">`)
	out = strings.ReplaceAll(out, "</table>", "</table></div>")
	return template.HTML(out), nil
}
