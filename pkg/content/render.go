package content

import (
    "os"

    "github.com/gomarkdown/markdown"
)

type Renderer interface {
    Render() ([]byte, error)
}

type RenderFunc func() ([]byte, error)

type MarkdownRenderer struct {
    File string
}

func (r *MarkdownRenderer) Render() ([]byte, error) {
	in, err := os.ReadFile(r.File)
	if err != nil {
		return nil, err
	}

	out := markdown.ToHTML(in, nil, nil)

	return out, nil
}
