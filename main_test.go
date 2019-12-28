// This file tests the "render" function, to verify that
// - correctly handles all ways to specify paths
// - only handles valid paths
// - correctly combines files specified in various places

package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func readExampleFile(name string) [][]byte {
	path := filepath.Join("./examples", name)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ret [][]byte
	ret = append(ret, bytes)
	return ret
}

func assertOutput(t *testing.T, outBytes []byte, reference string) {
	// If a test is successful, t.Log is ignored
	if !bytes.Equal(outBytes, []byte(reference)) {
		t.Logf("actual:\n%s\n", hex.Dump(outBytes))
		t.Logf("expected:\n%s\n", hex.Dump([]byte(reference)))
		t.Fatal("output differs from reference")
	}
}

func lineEndingTestHelper(t *testing.T, input string, expectedFilename string, expected string) {
	output, err := render([][]byte{[]byte(input)})
	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 1)
	assertOutput(t, output[expectedFilename], expected)
}

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

func TestRenderSimple(t *testing.T) {
	input := readExampleFile("simple.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 1)

	sh := "#!/usr/bin/env bash\necho \"Hello, World\"\n"
	assertOutput(t, output["hello.sh"], sh)
}

func TestRenderMultipleFiles(t *testing.T) {
	input := readExampleFile("multiple-files.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 2)

	dataJSON := "{\r\n    \"foo\": \"bar\",\r\n    \"hello\": \"world\"\r\n}\r\n"
	showSh := "#!/bin/bash\ncat data.json | jq .\n"

	assertOutput(t, output["data.json"], dataJSON)
	assertOutput(t, output["show.sh"], showSh)
}

func TestRenderSplitFile(t *testing.T) {
	input := readExampleFile("split-file.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 1)

	ps1 := "\"Hello World\"\r\ngci env:* | sort-object name\r\n"

	assertOutput(t, output["splitted.ps1"], ps1)
}

func TestRenderInvalid(t *testing.T) {
	input := readExampleFile("invalid.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 0)
}

func TestRenderWindowsSeparator(t *testing.T) {
	input := readExampleFile("windows-separators.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 1)

	helloBat := "@echo off\r\necho Hello,\r\necho Same File\r\n"
	assertOutput(t, output["example/hello.bat"], helloBat)
}

func TestRenderLineEndings(t *testing.T) {
	input := readExampleFile("lineendings.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 4)

	unixSh := "#!/usr/bin/env bash\necho \"Using LF on linux\"\n"
	assertOutput(t, output["unix.sh"], unixSh)

	windowsBat := "@echo off\r\necho For windows\r\necho Use CRLF\r\n"
	assertOutput(t, output["windows.bat"], windowsBat)

	splittedSh := "#!/usr/bin/env bash\necho \"This file, will use LF.\"\necho \"Because LF was specified first.\"\necho \"It's not important to keep all FILE-pragmas in sync.\"\n"
	assertOutput(t, output["splitted.sh"], splittedSh)

	exampleTxt := "This file uses \\r\ras line ending.\r"
	assertOutput(t, output["example.txt"], exampleTxt)
}

func TestInvalidPragma(t *testing.T) {
	input := []byte("```sh\n### FILE-CRFL: invalid.txt\nshould not be rendered\n```")
	var inputs [][]byte
	inputs = append(inputs, input)

	output, err := render(inputs)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 0)
}

func TestRenderMultipleInputFilesSingleOutput(t *testing.T) {
	input := readExampleFile("simple.md")
	input = append(input, readExampleFile("multiple-inputs.md")...)

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 1)

	sh := "#!/usr/bin/env bash\necho \"Hello, World\"\necho \"Hello from second file\"\n"
	assertOutput(t, output["hello.sh"], sh)
}

func TestRenderMultipleInputOutput(t *testing.T) {
	input := readExampleFile("multiple-files.md")
	input = append(input, readExampleFile("lineendings.md")...)
	input = append(input, readExampleFile("simple.md")...)
	input = append(input, readExampleFile("multiple-inputs.md")...)

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 7)

	dataJSON := "{\r\n    \"foo\": \"bar\",\r\n    \"hello\": \"world\"\r\n}\r\n"
	showSh := "#!/bin/bash\ncat data.json | jq .\n"

	assertOutput(t, output["data.json"], dataJSON)
	assertOutput(t, output["show.sh"], showSh)

	unixSh := "#!/usr/bin/env bash\necho \"Using LF on linux\"\n"
	assertOutput(t, output["unix.sh"], unixSh)

	windowsBat := "@echo off\r\necho For windows\r\necho Use CRLF\r\n"
	assertOutput(t, output["windows.bat"], windowsBat)

	splittedSh := "#!/usr/bin/env bash\necho \"This file, will use LF.\"\necho \"Because LF was specified first.\"\necho \"It's not important to keep all FILE-pragmas in sync.\"\n"
	assertOutput(t, output["splitted.sh"], splittedSh)

	exampleTxt := "This file uses \\r\ras line ending.\r"
	assertOutput(t, output["example.txt"], exampleTxt)

	sh := "#!/usr/bin/env bash\necho \"Hello, World\"\necho \"Hello from second file\"\n"
	assertOutput(t, output["hello.sh"], sh)
}
