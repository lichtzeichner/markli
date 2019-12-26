# Multiple Input Files

markli can process multiple markdown files at once, just specify the `-i` parameter multiple times.
All output will be combined, so if the same file name is specified in multiple inputs, the resulting file will be the sum of all.

# Example

For an example of this behavior, you can try this command:

```
markli -i examples/simple.md -i examples/multiple-inputs.md -o example-out
```

```sh
### FILE: hello.sh 
echo "Hello from second file"
```

Again, the first occurence of the file will specify it's output line endings (See [lineendings.md](lineendings.md)) for more details.