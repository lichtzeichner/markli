package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func render(input []byte) (map[string][]byte, error) {
	blocks := &scriptBlocks{}
	retval := make(map[string][]byte)

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
	if err != nil {
		return retval, err
	}

	for k, v := range blocks.renderer.Output {
		retval[k] = v.content
	}
	return retval, nil
}

func writeRendered(outDir string, output map[string][]byte) error {
	for filename, content := range output {
		path := filepath.Clean(filepath.Join(outDir, filename))
		dir := filepath.Dir(path)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(path, content, 0755); err != nil {
			return err
		}
	}
	return nil
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

	if err := writeRendered(outDir, rendered); err != nil {
		panic(err)
	}
}
