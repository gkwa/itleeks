package main

import (
   "testing"
   "testing/fstest"
)

func TestParseMarkdown(t *testing.T) {
   fsys := fstest.MapFS{
   	"test01.md": &fstest.MapFile{Data: []byte(`[name](https://name.com)`)},
   	"test02.md": &fstest.MapFile{Data: []byte(`Hello, world!`)},
   	"test03.md": &fstest.MapFile{Data: []byte(`# Header 1

This is a paragraph with a [link](https://example.com).

* List item 1
* List item 2

## Header 2

1. Numbered item 1
2. Numbered item 2

> This is a blockquote.

   This is a code block.
`)},
   }

   tests := []struct {
   	name     string
   	filename string
   	expected string
   }{
   	{
   		name:     "Simple link",
   		filename: "test01.md",
   		expected: "[name](https://name.com)",
   	},
   	{
   		name:     "Simple text",
   		filename: "test02.md",
   		expected: "Hello, world!",
   	},
   	{
   		name:     "Complex markdown",
   		filename: "test03.md",
   		expected: `# Header 1

This is a paragraph with a [link](https://example.com).

* List item 1
* List item 2

## Header 2

1. Numbered item 1
2. Numbered item 2

> This is a blockquote.

   This is a code block.`,
   	},
   }

   for _, tt := range tests {
   	t.Run(tt.name, func(t *testing.T) {
   		result, err := parseMarkdown(fsys, tt.filename)
   		if err != nil {
   			t.Fatalf("Error parsing markdown: %v", err)
   		}
   		if result != tt.expected {
   			t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
   		}
   	})
   }
}
