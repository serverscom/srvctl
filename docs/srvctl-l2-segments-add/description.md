A command to create a new L2 segment. It allows to pass parameters in two ways:

- Input - segment parameters are described in a file, a path to the file is specified via the `-i` or `--input` flag. The path can be absolute or relative to the srvctl file. Parameters should be described as a request body of the Public API request. There is also an option to use standard input (stdin) when specifying the flag this way: `--input -`

- Flags - parameters are specified via flags inside the command. The `--type` and `--member` flags are required. Members are specified in `id=<string>,mode=<native|trunk>` format and can be repeated for multiple members.
