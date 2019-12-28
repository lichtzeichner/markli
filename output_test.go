// This file tests the "writeRendered" function, to verify that
// - it only writes the specified files
// - it writes them to the specified output directory

package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

var tempDirs = make(map[string]string)
var baseDir string = "markli-testdir"

func getTempDir(t *testing.T) string {
	if dir, err := ioutil.TempDir(baseDir, t.Name()); err != nil {
		t.Fatal(err)
	} else {
		os.Mkdir(dir, 0644)
		tempDirs[t.Name()] = dir
		return dir
	}
	return ""
}

func validateFile(t *testing.T, path string, expected []byte) {
	dir := tempDirs[t.Name()]
	if dir == "" {
		t.Fatal("Could not find base dir for test: " + t.Name())
	}
	file := filepath.Join(dir, path)
	t.Logf("File: %s\n", file)
	info, err := os.Stat(file)
	assert.Assert(t, os.IsNotExist(err) == false)
	assert.Assert(t, !info.IsDir())
	actual, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("actual:\n%s\n", hex.Dump(actual))
	t.Logf("expected:\n%s\n", hex.Dump(expected))
	assert.Assert(t, bytes.Equal(actual, expected))
}

func validateDirStruct(t *testing.T, toScan string, expectedFiles []string) {
	var files []string

	err := filepath.Walk(toScan, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	onlyExpected := true
OUTER:
	for _, file := range files {
		if file == toScan {
			continue
		}
		for _, a := range expectedFiles {
			joined := filepath.Join(toScan, a)
			if joined == file {
				t.Logf("file: %s\n", file)
				continue OUTER
			}
		}
		t.Logf("UNEXPECTED file: %s\n", file)
		// Do not end the loop here to output the complete structure
		onlyExpected = false
	}
	assert.Assert(t, onlyExpected)
}

func TestMain(m *testing.M) {
	if err := os.Mkdir(baseDir, 0755); err != nil {
		panic(err)
	}

	result := m.Run()

	os.RemoveAll(baseDir)
	os.Exit(result)
}

func TestOutputFile(t *testing.T) {
	dir := getTempDir(t)

	output := make(map[string][]byte)
	output["foo.txt"] = []byte("bar")

	writeRendered(dir, output)

	validateFile(t, "foo.txt", []byte("bar"))

	files := []string{
		"foo.txt",
	}
	validateDirStruct(t, dir, files)
}

func TestOutputDir(t *testing.T) {
	dir := getTempDir(t)

	output := make(map[string][]byte)
	output["abc/def/foo.txt"] = []byte("bar")

	writeRendered(dir, output)

	validateFile(t, "abc/def/foo.txt", []byte("bar"))
	files := []string{
		"abc",
		"abc/def",
		"abc/def/foo.txt",
	}
	validateDirStruct(t, dir, files)
}

func TestOutputMulti(t *testing.T) {
	dir := getTempDir(t)

	output := make(map[string][]byte)
	output["abc/def/foo.txt"] = []byte("bar")
	output["abc/foo.txt"] = []byte("baz")
	output["foo.txt"] = []byte("zab")

	writeRendered(dir, output)

	validateFile(t, "abc/def/foo.txt", []byte("bar"))
	validateFile(t, "abc/foo.txt", []byte("baz"))
	validateFile(t, "foo.txt", []byte("zab"))

	files := []string{
		"foo.txt",
		"abc",
		"abc/foo.txt",
		"abc/def/",
		"abc/def/foo.txt",
	}
	validateDirStruct(t, dir, files)
}
