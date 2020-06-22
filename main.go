package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type logger struct {
	verbosity int
	outstream io.Writer
}

var NoInputError = errors.New("No inputs specified")
var StdinMustBeOnlyArgumentError = errors.New("If stdin is specified, it must be the only input argument")

func (l *logger) printf(verbosity int, format string, a ...interface{}) {
	if l.outstream != nil && verbosity <= l.verbosity {
		fmt.Fprintf(l.outstream, format, a...)
	}
}

func (l *logger) verbosef(format string, a ...interface{}) {
	l.printf(1, format, a...)
}

func (l *logger) verbose2f(format string, a ...interface{}) {
	l.printf(2, format, a...)
}

func (l *logger) verbose3f(format string, a ...interface{}) {
	l.printf(3, format, a...)
}

var log *logger = &logger{
	verbosity: 0,
	outstream: nil,
}

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
		log.verbosef("Writing output: %s\n", path)

		if err := os.MkdirAll(dir, 0755); err == nil {
			if err := ioutil.WriteFile(path, content, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

func exec(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	log.outstream = stdout

	var inputFiles []string
	var outDir string
	var rendered map[string][]byte

	appname := filepath.Base(args[0])

	flags := flag.NewFlagSet(appname, flag.ContinueOnError)
	flags.StringArrayVarP(&inputFiles, "input", "i", []string{}, "Markdown file to process, can be given multiple times. Use - to read from stdin.")
	flags.StringVarP(&outDir, "out-dir", "o", ".", "Output directory.")
	flags.CountVarP(&log.verbosity, "verbose", "v", "Control verbosity, shorthand can be given multiple times")
	flags.SetOutput(stderr)

	if err := flags.Parse(args); err != nil {
		return err
	}

	inputCnt := len(inputFiles)

	if inputCnt == 0 {
		fmt.Fprintf(stderr, "Usage of %s:\n", appname)
		flags.PrintDefaults()
		return NoInputError
	}

	for _, file := range inputFiles {
		if strings.TrimSpace(file) == "-" && inputCnt > 1 {
			return StdinMustBeOnlyArgumentError
		}
	}

	inputs := make([][]byte, 0, inputCnt)

	trimmedFiles := make([]string, inputCnt)
	for idx, file := range inputFiles {
		trimmedFiles[idx] = strings.TrimSpace(file)
	}

	for _, file := range trimmedFiles {
		log.verbose2f("Processing file %s\n", file)
		var input []byte
		var err error

		if file == "-" {
			input, err = ioutil.ReadAll(stdin)
		} else {
			input, err = ioutil.ReadFile(file)
		}
		if err != nil {
			return err
		}
		inputs = append(inputs, input)
	}

	rendered, err := render(inputs)
	if err != nil {
		return err
	}

	if err := writeRendered(strings.TrimSpace(outDir), rendered); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := exec(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
