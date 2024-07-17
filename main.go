package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
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
		linkText := v.Text(source)
		fmt.Fprintf(w, "[%s](%s)", linkText, v.Destination)
	case *ast.Text:
		fmt.Fprintln(w)
		fmt.Fprint(w, string(v.Text(source)))
	case *ast.String:
		fmt.Fprint(w, string(v.Value))
	case *ast.AutoLink:
		fmt.Fprintf(w, "%s", v.URL(source))
	default:
		if n.Type() == ast.TypeBlock {
			for i := 0; i < n.Lines().Len(); i++ {
				line := n.Lines().At(i)
				fmt.Fprint(w, string(line.Value(source)))
			}
		}
	}
}

func main() {
	markdown := goldmark.New()
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
