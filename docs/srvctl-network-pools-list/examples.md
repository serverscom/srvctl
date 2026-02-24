A command to list network pools:

```
srvctl network-pools list
```

A command to list all network pools (including all pages):

```
srvctl network-pools list --all
```

A command to filter network pools by location:

```
srvctl network-pools list --location-id "LON1"
```

A command to filter network pools by type:

```
srvctl network-pools list --type public
```

A command to filter network pools by label:

```
srvctl network-pools list --label-selector "env=prod"
```

A command to filter network pools by search pattern:

```
srvctl network-pools list --search-pattern "main"
```
