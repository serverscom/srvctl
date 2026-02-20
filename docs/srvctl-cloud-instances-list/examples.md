A command to list all cloud instances for the account:

```
srvctl cloud-instances list
```

A command to filter cloud instances by region:

```
srvctl cloud-instances list --region-id 1
```

A command to filter cloud instances by location:

```
srvctl cloud-instances list --location-id ex4mp1eID
```

A command to filter cloud instances by label:

```
srvctl cloud-instances list --label-selector environment=production
```
