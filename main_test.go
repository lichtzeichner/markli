package main

import (
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

func TestRenderSimple(t *testing.T) {
	input := readExampleFile("simple.md")

	output, err := render(input)
	assert.Assert(t, err == nil)
	assert.Assert(t, len(output) == 1)
	assert.Assert(t, output["hello.sh"] != nil)
}
