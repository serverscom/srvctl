# srvctl

## Description

`srvctl` is a command line app to manage servers.com resources.

## Installation

Homebrew:

```bash
brew tap serverscom/serverscom
brew install srvctl
```

## Usage

We need to define the context for your credentials, and global settings before we start using `srvctl`. 
In order to initiate CLI usage, we have to define the context, that will be used by the CLI.

It's done by using `srvctl login <context-name>` command. For the example below, we're using `default` context.

```bash
$ srvctl login default #prompts to enter your API token. 
Enter API token: ....

Successfully logged in with context "default"
Context "default" set as default
```

Config file will be located in `$XDG_CONFIG_HOME/srvctl/config.yaml`, if XDG_CONFIG_HOME exists.
Otherwise it will rely on `$HOME/.config/srvctl/config.yaml`. 
Additionally, you can define a custom path `SRVCTL_CONFIG_PATH`.

Config file supports multiple context options, allowing you to use various configs for each of them:

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
