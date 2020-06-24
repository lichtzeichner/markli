package main

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type lineEndingStyle int8

const (
	lineEndingUnknown lineEndingStyle = iota
	lineEndingCR
	lineEndingLF
	lineEndingCRLF
)

func (style lineEndingStyle) String() string {
	switch style {
	case lineEndingCR:
		return "CR"
	case lineEndingLF:
		return "LF"
	case lineEndingCRLF:
		return "CRLF"
	default:
		return "UNKNOWN"
	}
}

func parseLineEndingStyle(style string) lineEndingStyle {
	switch style {
	case "CR":
		return lineEndingCR
	case "LF":
		return lineEndingLF
	case "CRLF":
		return lineEndingCRLF
	default:
		return lineEndingUnknown
	}
}

// IsAbs() returns false for /etc/passwd on windows
// But if you do mkdir(/etc/passwd) you end up with C:/etc/passwd,
// this is *not* what we want
// therefore use filepath and check for / additionally.
func isAbs(path string) bool {
	return filepath.IsAbs(path) || strings.HasPrefix(path, "/")
}

// filepath.ToSlash() and .Clean() have platform-dependent behavior
// this is not helpful in this case.
func hasDirUp(path string) bool {
	for _, element := range strings.Split(path, "/") {
		if element == ".." {
			return true
		}
	}
	return false
}

func normalizePath(path string) string {
	return filepath.ToSlash(strings.TrimSpace(path))
}

func detectLineEnding(line []byte) lineEndingStyle {
	switch {
	case len(line) > 0 && line[len(line)-1] == '\r':
		return lineEndingCR
	case len(line) > 1 && line[len(line)-2] == '\r':
		return lineEndingCRLF
	default:
		return lineEndingLF
	}
}

type script struct {
	content    []byte
	lineEnding lineEndingStyle
}

func (s *script) append(value []byte) {
	if s.lineEnding != detectLineEnding(value) {
		cp := make([]byte, len(value))
		copy(cp, value)
		cp = bytes.TrimRight(cp, "\r\n")

		switch s.lineEnding {
		case lineEndingCRLF:
			cp = append(cp, []byte{'\r', '\n'}...)
		case lineEndingCR:
			cp = append(cp, '\r')
		default:
			cp = append(cp, '\n')
		}
		s.content = append(s.content, cp...)
	} else {
		s.content = append(s.content, value...)
	}
}

func (s *script) initLineEnding(lineEnding lineEndingStyle) {
	if s.lineEnding == lineEndingUnknown {
		s.lineEnding = lineEnding
	}
}

type scriptRenderer struct {
	Output map[string]script
}

var filePragmaRE = regexp.MustCompile(`###\s*FILE(-CR|-LF|-CRLF)?:(.*)\s*$`)

func parsePragma(input []byte) (string, lineEndingStyle) {
	ending := lineEndingUnknown
	if match := filePragmaRE.FindSubmatch(input); match != nil {
		desiredEnding := match[1]
		if len(desiredEnding) > 0 {
			// Cut the - from -CRLF
			ending = parseLineEndingStyle(string(desiredEnding[1:]))
		}
		p := normalizePath(string(match[2]))
		return p, ending
	}
	return "", ending
}

func newScriptRenderer(rendered map[string]script) *scriptRenderer {
	return &scriptRenderer{Output: rendered}
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
				if p, desiredEnding := parsePragma(value); p != "" {
					if desiredEnding != lineEndingUnknown {
						ending = desiredEnding
					}
					switch {
					case isAbs(p):
						log.verbosef("Warning: absolute paths are not allowed, ignoring path: %s\n", p)
					case hasDirUp(p):
						log.verbosef("Warning: using .. in paths is not allowed, ignoring path: %s\n", p)
					default:
						// accept this path
						path = p
					}
				}
				if path == "" {
					log.verbosef("Warning: ignoring empty path\n")
					return ast.WalkContinue, nil
				}
				log.verbose3f("Adding script '%s' with line ending '%s'\n", path, ending.String())
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
