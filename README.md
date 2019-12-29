![Actions Status](https://github.com/lichtzeichner/markli/workflows/tests/badge.svg)
![Actions Status](https://github.com/lichtzeichner/markli/workflows/lint/badge.svg)
![Actions Status](https://github.com/lichtzeichner/markli/workflows/build-binaries/badge.svg)

# markli

markli (from **mark**down **li**terator) is a simple commandline tool to support [literate programming](https://en.wikipedia.org/wiki/Literate_programming). It's main focus is better support for setup scripts that configure build machines or there alike. Basically you have much more documentation than actual code.

## Basic Usage

This utility has a simple commandline syntax:

    markli -i your-markdown.md -o output-folder

When called like this, all code-blocks containing `###FILE: ` within the first line will be converted into standalone files contained within `output-folder`.

## Line Endings

For certain things, e. g. Bash Scripts, you want to be able to explicitely control the line ending of the output file. You can use the following pragma extensions to achieve this:

* `FILE-LF`: Produces unix style `\n`
* `FILE-CRLF`: Produces Windows style `\r\n`

For more details and usage examples, have a look at [examples/lineendings.md](examples/lineendings.md)

## Examples

See the examples folder for basic use cases and features of markli. 

**Note**: These example files are also used as tests, see [examples_test.go](examples_test.go)

## Acknowledgements

Thanks to [simonfxr](https://github.com/simonfxr) for sharing the idea!