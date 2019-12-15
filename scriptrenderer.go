package main

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type scriptRenderer struct {
	Output map[string][]byte
	html.Config
}

var filePragmaRE = regexp.MustCompile(`###\s*FILE:\s*(.*)\s*$`)

func newScriptRenderer(opts ...html.Option) *scriptRenderer {
	r := &scriptRenderer{
		Config: html.NewConfig(),
		Output: make(map[string][]byte),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

func (r *scriptRenderer) renderNoop(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *scriptRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		filepath := ""
		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			value := line.Value(source)
			if i == 0 {
				filePragmaRE.Longest()
				if match := filePragmaRE.FindSubmatch(value); match != nil {
					filepath = string(match[1])
				}
			} else {
				r.Output[filepath] = append(r.Output[filepath], value...)
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *scriptRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Things we care for
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)

	// Everything else get's ignored
	reg.Register(ast.KindAutoLink, r.renderNoop)
	reg.Register(ast.KindBlockquote, r.renderNoop)
	reg.Register(ast.KindDocument, r.renderNoop)
	reg.Register(ast.KindEmphasis, r.renderNoop)
	reg.Register(ast.KindHTMLBlock, r.renderNoop)
	reg.Register(ast.KindHeading, r.renderNoop)
	reg.Register(ast.KindImage, r.renderNoop)
	reg.Register(ast.KindLink, r.renderNoop)
	reg.Register(ast.KindList, r.renderNoop)
	reg.Register(ast.KindListItem, r.renderNoop)
	reg.Register(ast.KindParagraph, r.renderNoop)
	reg.Register(ast.KindRawHTML, r.renderNoop)
	reg.Register(ast.KindString, r.renderNoop)
	reg.Register(ast.KindText, r.renderNoop)
	reg.Register(ast.KindTextBlock, r.renderNoop)
	reg.Register(ast.KindThematicBreak, r.renderNoop)
}

type scriptBlocks struct {
	renderer *scriptRenderer
}

func (e *scriptBlocks) Extend(m goldmark.Markdown) {
	if e.renderer != nil {
		panic("scriptBlocks can only be used once")
	}
	e.renderer = newScriptRenderer()

	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(e.renderer, 500),
	))
}
