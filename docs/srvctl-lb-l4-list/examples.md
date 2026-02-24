A command to list L4 load balancers:

```
srvctl lb l4 list
```

A command to filter L4 load balancers by location:

```
srvctl lb l4 list --location-id "LON1"
```

A command to filter L4 load balancers by cluster:

```
srvctl lb l4 list --cluster-id ex4mp1eClusterID
```

A command to filter L4 load balancers by label:

```
srvctl lb l4 list --label-selector "env=prod"
```

A command to filter L4 load balancers by search pattern:

```
srvctl lb l4 list --search-pattern "web"
```
