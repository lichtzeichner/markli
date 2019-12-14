# Multiple files

This is a more complexe example containing two files insid.

### Data

This can also be used to store example configuration data. This time an alternative syntax using four indents is used:

    ### FILE: data.json
    {
       "foo": "bar",
       "hello": "world"
    }


## Do it B

This is the second script.  
*Note: It needs the tool jq to run*

```sh
    ### FILE: show.sh
    #!/bin/bash
    cat data.json | jq .
```