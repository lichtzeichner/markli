package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// filepath.isAbs checks only for the platform the program is running on
// this function checks for *ANY* kind of absolute path
func isAbs(path string) bool {
	if strings.HasPrefix(path, "/") {
		return true
	}
	if driveLetterRE.Match([]byte(path)) {
		return true
	}
	return false
}

// filepath.ToSlash() and .Clean() have platform-dependent behavior
// this is not helpful in this case
func hasDirUp(path string) bool {
	return strings.Contains(path, "..")
}

func normalizePath(path string) string {
	return strings.ReplaceAll(strings.TrimSpace(path), "\\", "/")
}

func detectLineEnding(line []byte) string {
	l := len(line)
	sl := l - 2
	if sl < 0 {
		sl = 0
	}
	ending := line[sl:]
	if bytes.Compare(ending, []byte{'\r', '\n'}) == 0 {
		return "CRLF"
	} else if bytes.Compare(ending, []byte{'\r'}) == 0 {
		return "CR"
	}
	return "LF"
}

type script struct {
	content    []byte
	lineEnding string
}

func (s *script) append(value []byte) {

	if s.lineEnding != detectLineEnding(value) {
		cp := make([]byte, len(value))
		copy(cp, value)
		cp = bytes.TrimRight(cp, "\r\n")
		if s.lineEnding == "CRLF" {
			cp = append(cp, []byte{'\r', '\n'}...)
		} else if s.lineEnding == "CR" {
			cp = append(cp, '\r')
		} else {
			cp = append(cp, '\n')
		}
		s.content = append(s.content, cp...)
	} else {
		s.content = append(s.content, value...)
	}
}

func (s *script) initLineEnding(lineEnding string) {
	if s.lineEnding == "" {
		s.lineEnding = lineEnding
	}
}

type scriptRenderer struct {
	Output map[string]script
	html.Config
}

var filePragmaRE = regexp.MustCompile(`###\s*FILE(-CR|-LF|-CRLF)?:\s*(.*)\s*$`)
var driveLetterRE = regexp.MustCompile(`^[a-zA-Z]:[\/]`)

func newScriptRenderer(rendered map[string]script, opts ...html.Option) *scriptRenderer {
	r := &scriptRenderer{
		Config: html.NewConfig(),
		Output: rendered,
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
		path := ""
		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			value := line.Value(source)
			ending := detectLineEnding(value)
			if i == 0 {
				filePragmaRE.Longest()
				if match := filePragmaRE.FindSubmatch(value); match != nil {
					desiredEnding := match[1]
					if len(desiredEnding) > 0 {
						ending = string(desiredEnding[1:]) // Cut the - from -CRLF
					}
					p := normalizePath(string(match[2]))
					switch {
					case p == "":
						fmt.Fprintln(os.Stderr, "Warning: ingoring empty path")
					case isAbs(p):
						fmt.Fprintf(os.Stderr, "Warning: absolute paths are not allowed, ignoring path: %s\n", p)
					case hasDirUp(p):
						fmt.Fprintf(os.Stderr, "Warning: using .. in paths is not allowed, ignoring path: %s\n", p)
					default:
						// accept this path
						path = p
					}
				}
				if path == "" {
					return ast.WalkContinue, nil
				}
				sc := r.Output[path]
				sc.initLineEnding(ending)
				r.Output[path] = sc
			} else {
				sc := r.Output[path]
				sc.append(value)
				r.Output[path] = sc
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
	rendered map[string]script
	renderer *scriptRenderer
}

func newScriptBlocks(rendered map[string]script) scriptBlocks {
	if rendered == nil {
		panic("output struct must be initialized")
	}
	e := scriptBlocks{}
	e.rendered = rendered
	return e
}

func (e *scriptBlocks) Extend(m goldmark.Markdown) {
	if e.renderer != nil {
		panic("scriptBlocks can only be used once")
	}
	e.renderer = newScriptRenderer(e.rendered)

	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(e.renderer, 500),
	))
}
