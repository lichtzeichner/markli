
# Test filepath validation logic

## Case 1 - Absolute paths

Using absolute paths is not allowed:

```sh
### FILE: /etc/passwd
pwned::0:0:root:/root:/bin/sh
```

## Case 2 - Paths with ..

Using .. in paths is also forbidden:

```sh
### FILE: ../resolve.conf
nameserver injecting.dns.org
```
