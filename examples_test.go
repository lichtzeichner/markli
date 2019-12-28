// This file tests the "render" function using the examples, to verify that it
// - correctly handles all ways to specify paths
// - only handles valid paths
// - correctly combines files specified in various places
package main

import (
	"testing"

	"gotest.tools/assert"
)

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

	t.Log(output)

	if isWindows {
		assert.Assert(t, err == nil)
		assert.Assert(t, len(output) == 0)
	} else {
		// Unix filenames are crazy...
		assert.Assert(t, err == nil)
		assert.Assert(t, len(output) == 3)

		assert.Assert(t, output[`..\..\something.txt`] != nil)
		assert.Assert(t, output[`C:\temp\evil.bat`] != nil)
		assert.Assert(t, output["C:/temp/evil.bat"] != nil)
	}
}

func TestRenderWindowsSeparator(t *testing.T) {
	input := readExampleFile("windows-separators.md")

	output, err := render(input)

	assert.Assert(t, err == nil)
	if isWindows {
		assert.Assert(t, len(output) == 1)
		helloBat := "@echo off\r\necho Hello,\r\necho Same File\r\n"
		assertOutput(t, output["example/hello.bat"], helloBat)
	} else {
		assert.Assert(t, len(output) == 2)
		helloBat := "@echo off\r\necho Hello,\r\n"
		assertOutput(t, output["example/hello.bat"], helloBat)
		helloBat2 := "echo Same File\r\n"
		assertOutput(t, output[`example\hello.bat`], helloBat2)
	}
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
