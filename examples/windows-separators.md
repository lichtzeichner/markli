# Handling of directory separators on Windows

Windows itself accepts both type of slashes, \ and /. markli internally normalizes and validates the path of the output file. So both of the codeblocks end up in the same output file `hello.bat`.

```bat
### FILE-CRLF: example/hello.bat
@echo off
echo Hello,
```

```bat
### FILE: example\hello.bat
echo Same File
```