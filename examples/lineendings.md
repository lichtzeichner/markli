# Specify Line Endings

Sometimes, a specific type of line-ending is desired for the output. E. g. Shell scripts won't work if the files use `\r\n` (or CRLF), whereas Batch files have strange behavior if written using `\n` (LF).

Markdown files however support either, and usually it's not important. With markli, there could even be a mixture of files needing CRLF and LF in the *same* document - thus it's necessary to specify which type of line-endings to support.

Just specify the line ending to use after the `### FILE` pragma, it can either be `LF` or `CRLF` or `CR`:

```sh
### FILE-LF: unix.sh
#!/usr/bin/env bash
echo "Using LF on linux"
```

```bat
### FILE-CRLF: windows.bat
@echo off
echo For windows
echo Use CRLF
```


# Split files


```sh
### FILE-LF: splitted.sh
#!/usr/bin/env bash
echo "This file, will use LF."
```

For split-files, the first encountered line ending will win.

```sh
### FILE-CRLF: splitted.sh
echo "Because LF was specified first."
```

This means, that you can omit the ending type on later occurrences altogether:

```sh
### FILE: splitted.sh
echo "It's not important to keep all FILE-pragmas in sync."
```

# CR

Even `CR` is supported:

```
### FILE-CR: example.txt
This file uses \r
as line ending.
```