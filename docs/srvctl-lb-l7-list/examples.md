A command to list L7 load balancers:

```
srvctl lb l7 list
```

A command to list all L7 load balancers (including all pages):

```
srvctl lb l7 list --all
```

A command to filter L7 load balancers by location:

```
srvctl lb l7 list --location-id "LON1"
```

A command to filter L7 load balancers by cluster:

```
srvctl lb l7 list --cluster-id ex4mp1eClusterID
```

A command to filter L7 load balancers by label:

```
srvctl lb l7 list --label-selector "env=prod"
```

A command to filter L7 load balancers by search pattern:

```
srvctl lb l7 list --search-pattern "web"
```
