The Servers.com CLI tool or srvctl is a Command Line Interface to manage Servers.com services.

## Key terms

A **resource** is an object to be managed from CLI. It can be a service, product or another entity that you can manage. For example, a dedicated server.

A **command** is an operation to be performed on a resource or another action initiated from the CLI agent. For example, a command to create a dedicated server.

An **option** (or **flag**) is an additional argument to a command that makes its action more specific. For example, a flag to perform a command forcibly `-f`. A flag is just a standalone argument `-f`, an option implies adding some additional data to a flag `--label='env:production'`.

A **configuration** is a file containing a set of parameters.

A **context** is an environment of a specific configuration. A default context is the one you log in when a name of a specific context is not specified. Access to a context is verified by a Public API token.

A context has the following limitations:

- a context name can contain lower and upper case Latin letters, numbers, underscores, periods, dots, minus symbols;
- there can be only one default context;
- if a default context is not specified, the first one from a list will be used.

## Installation

The below instructions will guide you through the steps on how to install a Servers.com CLI agent.

### MacOS

1) Tap into the Servers.com repository:
```
brew tap serverscom/serverscom
```

2) Install the client:
```
brew install srvctl
```

### Linux

1) Download an archive from https://github.com/serverscom/srvctl/releases. An operating system and its architecture are specified in the name of an archive. For example, srvctl_0.1.0_linux_amd64.zip is an archive for Linux with the AMD architecture.

Then, open terminal and perform the below commands.

2) Extract the archive content:
```
unzip <archive_name>
```

3) Make the **srvctl** file executable:
```
chmod +x srvctl
```

4) Run the **srvctl** file to start using the agent:
```
./srvctl
```

For now, using of the CLI agent is possible when specifying a path to the srvctl file.

5) To use the CLI agent without directory mentioning, it's necessary to add a location of the srvctl executable file to **.bashrc**. Open the file in a text editor and add this line in the very end of the file:
```
export PATH=<srvctl directory>:$PATH
```

If srvctl is located in the `/home/user/srvctl_folder/`, the command will look like: `export PATH=/home/user/srvctl_folder/srvctl:$PATH`

6) Save the file and close it. To apply changes in the current session, perform:
```
source ~/.bashrc
```

### Windows

1) Download an archive from https://github.com/serverscom/srvctl/releases. An operating system and its architecture are specified in the name of an archive. For example, `srvctl_0.1.0_windows_amd64.zip` is an archive for Windows with the AMD architecture.

2) Extract the archive.

3) Run the **srvctl** executable file.

## Getting started

Once you have a Servers.com Public API token and has installed the CLI agent, you need to create a context.

1) Install the Servers.com CLI agent as instructed in the Installation section.

2) Get a Public API token from the Servers.com [Customer Portal](https://portal.servers.com/iam/api-tokens).

3) Open terminal and perform commands that are described below.

4) Log in to the context:
```
srvctl login <context-name>
```

There are no contexts yet; so, this operation will create a context and give it a name that you entered within this command.

5) Enter your Public API token and confirm.

The newly created context will be setup as default.

6) When a context is successfully created, it will be shown by the listing contexts command:
```
srvctl context list
```

7) Use the help command to explore other possibilities of the agent:
```
srvctl help
```
