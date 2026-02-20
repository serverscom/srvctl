A command to list nodes of the Kubernetes cluster with the "ex4mp1eClusterID" ID:

```
srvctl k8s list-nodes ex4mp1eClusterID
```

A command to filter nodes by search pattern:

```
srvctl k8s list-nodes ex4mp1eClusterID --search-pattern "worker"
```

A command to filter nodes by label:

```
srvctl k8s list-nodes ex4mp1eClusterID --label-selector "role=worker"
```
