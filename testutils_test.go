package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"gotest.tools/assert"
)

var isWindows = runtime.GOOS == "windows"

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
