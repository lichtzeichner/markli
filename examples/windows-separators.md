# Handling of directory separators on Windows

Windows itself accepts both type of slashes, \ and /. So on widnows both of the
codeblocks end up in the same output file `hello.bat`. On other platforms markli
will write two different files.

```bat
### FILE-CRLF: example/hello.bat
@echo off
echo Hello,
```

```bat
### FILE-CRLF: example\hello.bat
echo Same File
```
