package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

func renderMarkdown(w io.Writer, n ast.Node, source []byte) {
	switch v := n.(type) {
	case *ast.Document:
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			renderMarkdown(w, c, source)
			if c.NextSibling() != nil {
				fmt.Fprintln(w)
			}
		}
	case *ast.Paragraph:
		fmt.Fprintln(w)
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			renderMarkdown(w, c, source)
		}
	case *ast.Link:
		var linkText strings.Builder
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			switch t := c.(type) {
			case *ast.Text:
				linkText.Write(t.Text(source))
			case *ast.String:
				linkText.WriteString(string(t.Value))
			}
		}
		fmt.Fprintf(w, "[%s](%s)", linkText.String(), v.Destination)
	case *ast.Text:
		fmt.Fprint(w, string(v.Text(source)))
	case *ast.String:
		fmt.Fprint(w, string(v.Value))
	case *ast.AutoLink:
		fmt.Fprintf(w, "[](%s)", v.URL(source))
	case *east.Table:
		renderTable(w, v, source)
	}
}

func renderTable(w io.Writer, table *east.Table, source []byte) {
	fmt.Fprint(w, "\n")
	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		fmt.Fprint(w, "|")
		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			fmt.Fprint(w, " ")
			for node := cell.FirstChild(); node != nil; node = node.NextSibling() {
				renderMarkdown(w, node, source)
			}
			fmt.Fprint(w, " |")
		}
		fmt.Fprintln(w)

		if row == table.FirstChild() {
			fmt.Fprint(w, "|")
			for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
				fmt.Fprint(w, "---|")
			}
			fmt.Fprintln(w)
		}
	}
}

func main() {
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	input, err := os.ReadFile("testdata/test.md")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	doc := markdown.Parser().Parse(text.NewReader(input))

	var buf bytes.Buffer
	renderMarkdown(&buf, doc, input)

	fmt.Println(buf.String())
}
