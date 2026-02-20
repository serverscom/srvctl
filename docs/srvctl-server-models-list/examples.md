A command to list server models for location ID `1`:

```
srvctl server-models list --location-id 1
```

A command to filter server models by search pattern:

```
srvctl server-models list --location-id 1 --search-pattern "E5"
```

A command to list only server models with a RAID controller:

```
srvctl server-models list --location-id 1 --has-raid-controller
```
