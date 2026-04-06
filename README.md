# srvctl - CLI for Servers.com infrastructure

[![Go Reference](https://pkg.go.dev/badge/github.com/serverscom/srvctl.svg)](https://pkg.go.dev/github.com/serverscom/srvctl) [![Go Report Card](https://goreportcard.com/badge/github.com/serverscom/srvctl)](https://goreportcard.com/report/github.com/serverscom/srvctl)

## Description

Manage your Servers.com bare metal servers, cloud servers, and infrastructure directly from the terminal.

`srvctl` wraps the Servers.com Public API into a fast, scriptable command-line interface built in Go.

## Quick Start

### Homebrew (macOS & Linux):

```sh
brew tap serverscom/serverscom
brew install srvctl
```

### Docker:

```sh
docker pull ghcr.io/serverscom/srvctl:latest
docker run --rm -it ghcr.io/serverscom/srvctl:latest --help
```

### Binary releases:
Download the latest release for your OS and architecture from the [Releases page](https://github.com/serverscom/srvctl/releases).

## Usage

### Authentication

Create a context with your [Servers.com API token](https://portal.servers.com/login):

```sh
$ srvctl login default 
Enter API token: *****

Successfully logged in with context "default"
Context "default" set as default
```

### Configuration

The config file is stored at `$XDG_CONFIG_HOME/srvctl/config.yaml`, if XDG_CONFIG_HOME exists.
Otherwise it will rely on `$HOME/.config/srvctl/config.yaml`. 
You can override this with the `SRVCTL_CONFIG_PATH` environment variable.

`srvctl` supports multiple contexts, allowing you to manage several Servers.com accounts or API endpoints from a single installation:

```yaml
globalConfig: {}
defaultContext: default
contexts:
    - name: default
      endpoint: https://api.servers.com/v1
      token: <YOUR_API_TOKEN>
      config: {}
    - name: different-context
      endpoint: https://api.servers.com/v2
      token: <2ND_API_TOKEN>
      config: {
        proxy: "",
        http-timeout: 30,
        verbose: true, /* (true|false) */
        output: "json" /* (text|json|yaml) */
      }
```

You can adjust the context later on:

```bash
# changing the context name
srvctl context update <context-name> --name=<new-name>

# setting context to act as default
srvctl context update <context-name> --default

# delete specific context
srvctl context delete <context-name>
```

## Documentation

Documentation is accessible via `man` or via `--help` flag, for example:

```bash
# man option
$ man srvctl-hosts-ds-list

# help option, short command help
$ srvctl hosts ds list --help
```

Man pages are based on the documentation info located in `/docs` directory.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

`srvctl` is released under the Apache 2.0 License.
