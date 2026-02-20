A command to list all cloud backups for the account:

```
srvctl cloud-backups list --all
```

A command to filter cloud backups by region:

```
srvctl cloud-backups list --region-id 1
```

A command to filter cloud backups by label:

```
srvctl cloud-backups list --label-selector environment=production
```
