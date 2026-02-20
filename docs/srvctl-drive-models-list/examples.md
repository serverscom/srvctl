A command to list drive models for location ID `1` and server model ID `42`:

```
srvctl drive-models list --location-id 1 --server-model-id 42
```

A command to list all drive models (including all pages):

```
srvctl drive-models list --location-id 1 --server-model-id 42 --all
```

A command to filter drive models by search pattern:

```
srvctl drive-models list --location-id 1 --server-model-id 42 --search-pattern "SSD"
```

A command to filter drive models by media type:

```
srvctl drive-models list --location-id 1 --server-model-id 42 --media-type SSD
```

A command to filter drive models by interface:

```
srvctl drive-models list --location-id 1 --server-model-id 42 --interface SATA3
```
