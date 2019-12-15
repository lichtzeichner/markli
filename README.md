# markli

markli is a simple commandline tool to support [literate programming](https://en.wikipedia.org/wiki/Literate_programming). It's main focus is better support for setup scripts that configure build machines or there alike. Basically you have much more documentation than actual code.

## Basic Usage

This utility has a simple commandline syntax:

    markli -i your-markdown.md -o output-folder

When called like this, all code-blocks containing `###FILE: ` within the first line will be converted into standalone files contained within `output-folder`.

See the examples subfolder for some use cases

## Acknowledgements

Thanks to [simonfxr](https://github.com/simonfxr) for sharing the idea!

## Boostrapping

Since markli uses itself as a literate programming environment we have a boostrapping issue. This could be overcome by commiting its own output code into the repository. As an alternative we provide a boostrapping script here:

```sh
# ### FILE: bootstrap_stage1.sh
#!/usr/bin/env bash
###BEGIN_BOOTSTRAP
d=$(dirname "${BASH_SOURCE[0]}")
f() { re='.*### FILE: (.*)' c=0 f=0
while IFS= read -r l; do case "$c$f$l" in
    ??'```'*) ((c=!c)); o=/dev/null f=1;;
    11*'### FILE:'*) [[ $l =~ $re ]] && o="${BASH_REMATCH[1]}" f=0 && :>"$o";;
    1*) f=0; echo "$l" >>"$o";;
  esac
done <"$d"/README.md; }; f
###END_BOOTSTRAP
```

For convenience we provide a simple shell script which extracts the code block
above and executes it, just run `bootstrap.sh` to execute all steps in the correct order.
You now should have all the extracted sources of a "stage 2" markli.

We still commit the extracted sources into repository, but just as a
convienience for users, e.g. go get will work as expected.

## Implementation

The code of markli follows

### Main code file

```go
// ### FILE: main.go
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func main() {

	input := flag.String("i", "", "Markdown file to process")
	flag.StringVar(input, "input", "", "Markdown file to process")

	output := flag.String("o", ".", "Output directory.")
	flag.StringVar(output, "out-dir", ".", "Output directory.")

	flag.Parse()

	if input == nil || *input == "" {
		println("No valid input file specified. Parameters:")
		flag.PrintDefaults()
		return
	}

	if output == nil || *output == "" {
		println("No valid output directory specified. Parameters:")
		flag.PrintDefaults()
		return
	}

	file, err := ioutil.ReadFile(*input)
	if err != nil {
		panic(err)

	}

	blocks := &scriptBlocks{}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			blocks,
		),
	)

	var buf bytes.Buffer
	md.Convert(file, &buf)

	outDir := *output

	for filename, content := range blocks.renderer.Output {
		path := filepath.Join(outDir, filename)
		ioutil.WriteFile(path, content, 0755)
	}
}
```

### Script Renderer code file

```go
// ### FILE: scriptrenderer.go
package main

import (
	"fmt"
	"path/filepath"
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
		path := ""
		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			value := line.Value(source)
			if i == 0 {
				filePragmaRE.Longest()
				if match := filePragmaRE.FindSubmatch(value); match != nil {
					p := filepath.ToSlash(string(match[1]))
					switch {
					case filepath.IsAbs(p):
						fmt.Printf("Warning: absolute paths are not allowed, ignoring path: %s\n", p)
					case filepath.Clean("/"+p) != "/"+p:
						fmt.Printf("Warning: using .. in paths is not allowed, ignoring path: %s\n", p)
					default:
						// accept this path
						path = p
					}
				}
			} else {
				r.Output[path] = append(r.Output[path], value...)
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
```
