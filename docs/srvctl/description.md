srvctl is a Command Line Interface to manage Servers.com services.

## Getting started

1) Get a Public API token from the Servers.com [Customer Portal](https://portal.servers.com/iam/api-tokens).

2) Open terminal and perform commands that are described below.

3) Log in to the context:

```
srvctl login <context-name>
```

There is no context yet; so, this operation will create a context and give it a name that you entered within this command.

4) Enter your Public API token and confirm.

The newly created context will be setup as default.

5) When a context is successfully created, it will be shown by the listing contexts command:

```
srvctl context list
```

6) Use the help command to explore other possibilities of the agent:

```
srvctl help
```
