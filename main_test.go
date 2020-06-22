// Tests for various functions of markli using fixed inputs and outputs

package main

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gotest.tools/assert"
)

func TestMarkdownCrFileCrLf(t *testing.T) {
	// goldmark can't cope with \r as line-ending for markdown files
	// therefore this is not expected to have any output.
	input := "Foo\r```sh\r### FILE-CRLF: foo.txt\rshould have lf\rline ending\r```\r"

	output, err := render([][]byte{[]byte(input)})

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 0)
}

func TestMarkdownCrLfFileLf(t *testing.T) {
	// Test the FILE-pragma in a platform independent way
	input := "Foo\r\n```sh\n### FILE-LF: foo.txt\r\nshould have lf\r\nline ending\r\n```\r\n"
	expected := "should have lf\nline ending\n"

	lineEndingTestHelper(t, input, "foo.txt", expected)
}

func TestMarkdownCrLfFileCr(t *testing.T) {
	// Test the FILE-pragma in a platform independent way
	input := "Foo\r\n```sh\n### FILE-CR: foo.txt\r\nshould have lf\r\nline ending\r\n```\r\n"
	expected := "should have lf\rline ending\r"

	lineEndingTestHelper(t, input, "foo.txt", expected)
}

func TestMarkdownLfFileCrLf(t *testing.T) {
	// Test the FILE-pragma in a platform independent way
	input := "Foo\n```sh\n### FILE-CRLF: foo.txt\nshould have lf\nline ending\n```\n"
	expected := "should have lf\r\nline ending\r\n"

	lineEndingTestHelper(t, input, "foo.txt", expected)
}

func TestNoPragma(t *testing.T) {
	input := []byte("### FILE: foo.txt")

	file, ending := parsePragma(input)

	assert.Equal(t, file, "foo.txt")
	assert.Equal(t, ending, lineEndingUnknown)

	assert.Equal(t, ending.String(), "UNKNOWN")
	assert.Equal(t, parseLineEndingStyle("UNKNOWN"), lineEndingUnknown)
}

func TestInvalidPragma(t *testing.T) {
	input := []byte("```sh\n### FILE-CRFL: invalid.txt\nshould not be rendered\n```")
	var inputs [][]byte
	inputs = append(inputs, input)

	output, err := render(inputs)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 0)
}

func TestLineEndingDetection(t *testing.T) {
	// Empty input defaults to LF
	assert.Assert(t, detectLineEnding(nil) == lineEndingLF)

	// Empty lines
	assert.Assert(t, detectLineEnding([]byte("\r")) == lineEndingCR)
	assert.Assert(t, detectLineEnding([]byte("\n")) == lineEndingLF)
	assert.Assert(t, detectLineEnding([]byte("\r\n")) == lineEndingCRLF)

	// Mixed style
	assert.Assert(t, detectLineEnding([]byte("\n\r")) == lineEndingCR)
	assert.Assert(t, detectLineEnding([]byte("a\rb\n")) == lineEndingLF)
	assert.Assert(t, detectLineEnding([]byte("c\n\r\n")) == lineEndingCRLF)
}

func TestPragmaParser(t *testing.T) {
	n, e := parsePragma([]byte("### FILE: foo.sh"))
	assert.Assert(t, n == "foo.sh")
	assert.Assert(t, e == lineEndingUnknown)

	n, e = parsePragma([]byte("### FILE-CRLF: foo/bar/baz/lol.txt"))
	assert.Assert(t, n == "foo/bar/baz/lol.txt")
	assert.Assert(t, e == lineEndingCRLF)

	n, e = parsePragma([]byte(`### FILE-CRLF: foo\bar\baz\lol.txt`))
	if isWindows {
		assert.Assert(t, n == "foo/bar/baz/lol.txt")
	} else {
		assert.Assert(t, n == `foo\bar\baz\lol.txt`)
	}
	assert.Assert(t, e == lineEndingCRLF)

	n, e = parsePragma([]byte("### FILE-CRLF: ### FILE-LF: recursive.txt"))
	assert.Assert(t, n == "### FILE-LF: recursive.txt")
	assert.Assert(t, e == lineEndingCRLF)

	// systemd units can use \ in file names. But this means a directory with a file inside on windows
	// both is valid depending on the platform. *sigh*
	n, e = parsePragma([]byte("### FILE-LF: foo\\bar"))

	if isWindows {
		assert.Assert(t, n == "foo/bar")
	} else {
		assert.Assert(t, n == "foo\\bar")
	}

	assert.Assert(t, e == lineEndingLF)
}

func TestHasDirUp(t *testing.T) {
	assert.Assert(t, hasDirUp("..") == true)
	assert.Assert(t, hasDirUp("../foo") == true)
	assert.Assert(t, hasDirUp(`.\.foo`) == false)
	assert.Assert(t, hasDirUp("..foo") == false)
	assert.Assert(t, hasDirUp("f..oo") == false)
	assert.Assert(t, hasDirUp(`foo/../bar`) == true)
	assert.Assert(t, hasDirUp(`foo/..`) == true)
}

func TestNewScriptBlocksPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("newScriptBlocks(nil) did not panic")
		}
	}()

	newScriptBlocks(nil)
}

func TestScriptBlocksDoubleUsagePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Double usage of script blocks did not panic")
		}
	}()

	output := make(map[string]script)
	blocks := newScriptBlocks(output)

	goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&blocks,
			&blocks, // should panic
		),
	)
}
