![Actions Status](https://github.com/lichtzeichner/markli/workflows/tests/badge.svg)
![Actions Status](https://github.com/lichtzeichner/markli/workflows/lint/badge.svg)
![Actions Status](https://github.com/lichtzeichner/markli/workflows/build-binaries/badge.svg)

# markli

markli is a simple commandline tool to support [literate programming](https://en.wikipedia.org/wiki/Literate_programming). It's main focus is better support for setup scripts that configure build machines or there alike. Basically you have much more documentation than actual code.

## Basic Usage

This utility has a simple commandline syntax:

    markli -i your-markdown.md -o output-folder

When called like this, all code-blocks containing `###FILE: ` within the first line will be converted into standalone files contained within `output-folder`.

See the examples subfolder for some use cases

## Examples

See the examples folder for basic use cases and features of markli. 

**Note**: These example files are also used as tests, see `main_test.go`

## Acknowledgements

Thanks to [simonfxr](https://github.com/simonfxr) for sharing the idea!