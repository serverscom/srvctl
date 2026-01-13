A command to create a dedicated server. It allows to pass parameters of a server in two ways:

- Input - server parameters are described in a file, a path to the file is specified via the `-i` or `–input` flag. The path can be absolute or relative to the srvctl file.

- Flags  - parameters are specified via flags inside the command and hostnames are listed as position arguments. As many arguments, as many servers of this configuration will be created.
