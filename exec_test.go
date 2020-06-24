package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestExecEmpty(t *testing.T) {
	args := []string{"markli"}
	stdin := strings.NewReader("")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	err := exec(args, stdin, stdout, stderr)

	assert.Assert(t, errors.Is(err, NoInputError))
	assert.Assert(t, strings.Contains(stderr.String(), "Usage of"))
	assert.Assert(t, strings.Contains(stderr.String(), "-i, --input"))
}

func TestExecNonExisting(t *testing.T) {
	args := []string{"markli", "-i", "/nonexisting"}
	stdin := strings.NewReader("")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	err := exec(args, stdin, stdout, stderr)

	assert.Assert(t, os.IsNotExist(err))
}

func TestExecInvalidCmd(t *testing.T) {
	args := []string{"markli", "--foo"}
	stdin := strings.NewReader("")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	err := exec(args, stdin, stdout, stderr)

	assert.ErrorContains(t, err, "unknown flag")
}

func TestExecStdin(t *testing.T) {
	dir := getTempDir(t)
	args := []string{"markli", "-v", "-i -", "-o " + dir}
	stdin := strings.NewReader("Foo\n```sh\n### FILE-CRLF: foo.txt\nHello, World\n```\n")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	err := exec(args, stdin, stdout, stderr)

	assert.NilError(t, err)
	// Make sure -v works
	assert.Assert(t, stdout.String() != "")

	files := []string{
		"foo.txt",
	}

	validateFile(t, "foo.txt", []byte("Hello, World\r\n"))
	validateDirStruct(t, dir, files)
}

func TestExecEmptyPathPragmaStdin(t *testing.T) {
	dir := getTempDir(t)
	args := []string{"markli", "-v", "-i -", "-o " + dir}
	stdin := strings.NewReader("Foo\n```sh\n### FILE-CRLF: \nHello, World\n```\n")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	err := exec(args, stdin, stdout, stderr)

	assert.NilError(t, err)
	// Make sure -v works
	outstr := stdout.String()
	t.Logf("stdout: [%s]\n", outstr)
	assert.Assert(t, strings.Contains(outstr, "ignoring empty path"))
}

func TestExecFile(t *testing.T) {
	dir := getTempDir(t)

	hw := []byte("Foo\n```sh\n### FILE-CRLF: foo.txt\nHello, World\n```\n")
	err := ioutil.WriteFile(dir+"/input.md", hw, 0644)
	assert.NilError(t, err)

	args := []string{"markli", "-i " + dir + "/input.md", "-o " + dir}
	stdin := bytes.NewBufferString("")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	err = exec(args, stdin, stdout, stderr)

	assert.NilError(t, err)

	files := []string{
		"input.md",
		"foo.txt",
	}

	validateFile(t, "foo.txt", []byte("Hello, World\r\n"))
	validateDirStruct(t, dir, files)
}

func TestExecFileAndStdIn(t *testing.T) {
	dir := getTempDir(t)

	args := []string{"markli", "-i " + dir + "/input.md", "-i -", "-o " + dir}

	stdin := strings.NewReader("Foo\n```sh\n### FILE-CRLF: foo.txt\nSTDIN\n```\n")
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")

	hw := []byte("Foo\r\n```sh\r\n### FILE-CRLF: foo.txt\r\nFILE\r\n```\r\n")
	err := ioutil.WriteFile(dir+"/input.md", hw, 0644)
	assert.NilError(t, err)

	err = exec(args, stdin, stdout, stderr)

	assert.NilError(t, err)

	files := []string{
		"input.md",
		"foo.txt",
	}

	validateFile(t, "foo.txt", []byte("FILE\r\nSTDIN\r\n"))
	validateDirStruct(t, dir, files)
}
