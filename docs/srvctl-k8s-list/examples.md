A command to list Kubernetes clusters:

```
srvctl k8s list
```

A command to filter Kubernetes clusters by location:

```
srvctl k8s list --location-id "LON1"
```

A command to filter Kubernetes clusters by label:

```
srvctl k8s list --label-selector "env=prod"
```

A command to filter Kubernetes clusters by search pattern:

```
srvctl k8s list --search-pattern "production"
```
