package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func render(inputs [][]byte) (map[string][]byte, error) {

	output := make(map[string]script)
	blocks := newScriptBlocks(output)
	retval := make(map[string][]byte)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&blocks,
		),
	)

	var buf bytes.Buffer
	for _, input := range inputs {
		err := md.Convert(input, &buf)
		if err != nil {
			return retval, err
		}
	}

	for k, v := range output {
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

	var inputFiles []string
	var outDir string
	var rendered map[string][]byte
	var inputs [][]byte

	flag.StringArrayVarP(&inputFiles, "input", "i", []string{}, "Markdown file to process, use - for stdin")
	flag.StringVarP(&outDir, "out-dir", "o", ".", "Output directory.")
	flag.Parse()

	if len(inputFiles) == 0 {
		fmt.Fprint(os.Stderr, "No inputs specified\n")
		flag.Usage()
		os.Exit(1)
	}

	for _, file := range inputFiles {
		input, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		inputs = append(inputs, input)
	}

	rendered, err := render(inputs)
	if err != nil {
		panic(err)
	}

	if err := writeRendered(outDir, rendered); err != nil {
		panic(err)
	}
}
