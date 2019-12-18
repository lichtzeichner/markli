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

func readExampleFile(name string) []byte {

	path := filepath.Join("./examples", name)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func assertOutput(t *testing.T, outBytes []byte, reference string) {
	// If a test is successful, t.Log is ignored
	t.Logf("actual:\n%s\n", hex.Dump(outBytes))
	t.Logf("expected:\n%s\n", hex.Dump([]byte(reference)))
	assert.Assert(t, bytes.Compare(outBytes, []byte(reference)) == 0)
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
	showSh := "#!/bin/bash\r\ncat data.json | jq .\r\n"

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

	output, err := render(input)

	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 0)
}
