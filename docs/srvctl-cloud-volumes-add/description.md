A command to create a new cloud volume. It allows to pass parameters in two ways:

- Input - volume parameters are described in a file, a path to the file is specified via the `-i` or `--input` flag. The path can be absolute or relative to the srvctl file. There is also an option to use standard input (stdin) when specifying the flag this way: `--input -`

- Flags - parameters are specified via flags inside the command. The `--name` and `--region-id` flags are required.
