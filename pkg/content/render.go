package content

import (
	"os"

	md2term "github.com/MichaelMure/go-term-markdown"
	md2html "github.com/gomarkdown/markdown"
)

type Renderer interface {
	Render() ([]byte, error)
}

type RenderFunc func() ([]byte, error)

type HTMLRenderer struct {
	File string
}

func (r *HTMLRenderer) Render() ([]byte, error) {
	in, err := os.ReadFile(r.File)
	if err != nil {
		return nil, err
	}

	out := md2html.ToHTML(in, nil, nil)

	return out, nil
}

type TerminalRenderer struct {
	File    string
	Width   int
	LeftPad int
}

func (r *TerminalRenderer) Render() ([]byte, error) {
	in, err := os.ReadFile(r.File)
	if err != nil {
		return nil, err
	}

	width := r.Width
	if width == 0 {
		width = 80
	}
	pad := r.LeftPad
	if pad < 0 {
		pad = 0
	}

	out := md2term.Render(string(in), width, pad)

	return out, nil
}
