A command to list SBM flavor models for location ID `1`:

```
srvctl sbm-models list --location-id 1
```

A command to filter SBM flavor models by search pattern:

```
srvctl sbm-models list --location-id 1 --search-pattern "standard"
```

A command to list all SBM flavors including unavailable ones:

```
srvctl sbm-models list --location-id 1 --show-all
```
