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

func render(input []byte) (map[string][]byte, error) {
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
	err := md.Convert(input, &buf)

	return blocks.renderer.Output, err
}

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

	outDir := *output

	rendered, err := render(file)
	if err != nil {
		panic(err)
	}

	for filename, content := range rendered {
		path := filepath.Join(outDir, filename)
		ioutil.WriteFile(path, content, 0755)
	}
}
