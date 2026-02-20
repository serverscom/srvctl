A command to list all cloud volumes for the account:

```
srvctl cloud-volumes list
```

A command to filter cloud volumes by region:

```
srvctl cloud-volumes list --region-id 1
```

A command to filter cloud volumes attached to a specific instance:

```
srvctl cloud-volumes list --instance-id ex4mp1eID
```

A command to filter cloud volumes by label:

```
srvctl cloud-volumes list --label-selector environment=production
```
