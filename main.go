package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

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
	case *ast.Heading:
		fmt.Fprintf(w, "%s ", strings.Repeat("#", v.Level))
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			renderMarkdown(w, c, source)
		}
		fmt.Fprintln(w)
	case *ast.Paragraph:
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			renderMarkdown(w, c, source)
		}
		fmt.Fprintln(w)
	case *ast.Link:
		linkText := v.Text(source)
		fmt.Fprintf(w, "[%s](%s)", linkText, v.Destination)
	case *ast.Text:
		fmt.Fprint(w, string(v.Text(source)))
	case *ast.String:
		fmt.Fprint(w, string(v.Value))
	case *ast.AutoLink:
		fmt.Fprintf(w, "%s", v.URL(source))
	case *ast.List:
		for i, c := 0, v.FirstChild(); c != nil; i, c = i+1, c.NextSibling() {
			if v.IsOrdered() {
				fmt.Fprintf(w, "%d. ", i+1)
			} else {
				fmt.Fprint(w, "* ")
			}
			renderMarkdown(w, c, source)
			fmt.Fprintln(w)
		}
	case *ast.ListItem:
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			renderMarkdown(w, c, source)
		}
	case *ast.Blockquote:
		fmt.Fprint(w, "> ")
		for c := v.FirstChild(); c != nil; c = c.NextSibling() {
			renderMarkdown(w, c, source)
		}
		fmt.Fprintln(w)
	case *ast.FencedCodeBlock:
		fmt.Fprintln(w, "```")
		for i := 0; i < v.Lines().Len(); i++ {
			line := v.Lines().At(i)
			fmt.Fprintf(w, "%s", line.Value(source))
		}
		fmt.Fprintln(w, "```")
	default:
		if n.Type() == ast.TypeBlock {
			for i := 0; i < n.Lines().Len(); i++ {
				line := n.Lines().At(i)
				fmt.Fprint(w, string(line.Value(source)))
			}
		}
	}
}

func parseMarkdown(fsys fs.FS, filename string) (string, error) {
	markdown := goldmark.New()
	input, err := fs.ReadFile(fsys, filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	doc := markdown.Parser().Parse(text.NewReader(input))
	var buf bytes.Buffer
	renderMarkdown(&buf, doc, input)
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	result, err := parseMarkdown(os.DirFS("."), "testdata/test.md")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
}
