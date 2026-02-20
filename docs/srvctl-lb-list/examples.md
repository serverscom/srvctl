A command to list all load balancers:

```
srvctl lb list
```

A command to list all load balancers (including all pages):

```
srvctl lb list --all
```

A command to filter load balancers by location:

```
srvctl lb list --location-id "LON1"
```

A command to filter load balancers by cluster:

```
srvctl lb list --cluster-id ex4mp1eClusterID
```

A command to filter load balancers by label:

```
srvctl lb list --label-selector "env=prod"
```

A command to filter load balancers by search pattern:

```
srvctl lb list --search-pattern "web"
```
