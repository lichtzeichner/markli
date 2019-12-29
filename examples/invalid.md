
# Test filepath validation logic

This document shows various invalid or misleading usages of the `### FILE:` pragma.

## Empty paths

An empty file pragma will be ignored.

```sh
### FILE:
Nope!
```

## Unix - Absolute paths

Using absolute paths is not allowed:

```sh
### FILE: /etc/passwd
pwned::0:0:root:/root:/bin/sh
```

## Unix - Paths with ..

Using .. in paths is also forbidden:

```sh
### FILE: ../resolve.conf
nameserver injecting.dns.org
```

## Windows - relative paths (with backward \)

**Note**: This will work on unix, as it's actually a valid filename.

```
### FILE: ..\..\something.txt
Lorem ipsum
```

## Windows - Absolute paths (with backward \)

**Note**: This will work on unix, as it's actually a valid filename.

```bat
### FILE: C:\temp\evil.bat
@echo off
echo "EVIL"
```

## Windows - Absolute windows paths (With forward /)

**Note**: This will work on unix, as it's actually valid relative path with filename.

```bat
### FILE: C:/temp/evil.bat
@echo off
echo "EVIL"
```
