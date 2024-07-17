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

> This is a blockquote.`)},
		"test04.md": &fstest.MapFile{Data: []byte("```\nThis is a code block.\n```")},
		"test05.md": &fstest.MapFile{Data: []byte(`https://google.com`)},
		"test06.md": &fstest.MapFile{Data: []byte(`[link](https://example.com)`)},
		"test07.md": &fstest.MapFile{Data: []byte(`This is a [link](https://example.com) in some text.`)},
		"test08.md": &fstest.MapFile{Data: []byte(`[Link 1](https://example1.com) and [Link 2](https://example2.com)`)},
		"test09.md": &fstest.MapFile{Data: []byte(`Here's a [link with title](https://example.com "Example Title").`)},
		"test10.md": &fstest.MapFile{Data: []byte(`[Link 1](https://example1.com) Some text [Link 2](https://example2.com "Title 2")`)},
		"test11.md": &fstest.MapFile{Data: []byte(`Text before [Link 1](https://example1.com "Title 1") text between [Link 2](https://example2.com) text after`)},
		"test12.md": &fstest.MapFile{Data: []byte(`hello

https://d.com/d.txt

[a](https://b.com/c.txt)

| a | b |
|---|---|


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

> This is a blockquote.`,
		},
		{
			name:     "Code block only",
			filename: "test04.md",
			expected: "```\nThis is a code block.\n```",
		},
		{
			name:     "Links",
			filename: "test05.md",
			expected: "https://google.com",
		},
		{
			name:     "Complex markdown",
			filename: "test06.md",
			expected: `[link](https://example.com)`,
		},
		{
			name:     "Link in text",
			filename: "test07.md",
			expected: "This is a [link](https://example.com) in some text.",
		},
		{
			name:     "Multiple links",
			filename: "test08.md",
			expected: "[Link 1](https://example1.com) and [Link 2](https://example2.com)",
		},
		{
			name:     "Link with title",
			filename: "test09.md",
			expected: "Here's a [link with title](https://example.com).",
		},
		{
			name:     "Multiple links with and without title",
			filename: "test10.md",
			expected: "[Link 1](https://example1.com) Some text [Link 2](https://example2.com)",
		},
		{
			name:     "Links with text before, between, and after",
			filename: "test11.md",
			expected: "Text before [Link 1](https://example1.com) text between [Link 2](https://example2.com) text after",
		},

		{
			name:     "xxx",
			filename: "test12.md",
			expected: `hello

https://d.com/d.txt

[a](https://b.com/c.txt)

| a | b |
|---|---|
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseMarkdown(fsys, tt.filename)
			if err != nil {
				t.Fatalf("Error parsing markdown: %v", err)
			}
			if result != tt.expected {
				t.Errorf("\nExpected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
