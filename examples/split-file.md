# Split files

This document illustrates, that one code file can be splitted into multiple code blocks.

The first block illustrates that powershell will echo back any literal string in the source of a script:

```powershell
### FILE-CRLF: splitted.ps1
"Hello World"
```

The second part shows how to print all environment variables using powershell

```powershell
### FILE: splitted.ps1
gci env:* | sort-object name
```