
# Test filepath validation logic

## Empty paths

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

```
### FILE: ..\..\something.txt
Lorem ipsum
```

## Windows - Absolute paths (with backward \)

```bat
### FILE: C:\temp\evil.bat
@echo off
echo "EVIL"
```

## Windows - Absolute windows paths (With forward /)

```bat
### FILE: C:/temp/evil.bat
@echo off
echo "EVIL"
```
