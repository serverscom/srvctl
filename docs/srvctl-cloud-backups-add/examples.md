A command to create a cloud backup for a volume:

```
srvctl cloud-backups add --volume-id ex4mp1eID --name my-backup
```

A command to create an incremental cloud backup:

```
srvctl cloud-backups add --volume-id ex4mp1eID --name my-backup --incremental
```

A command to create a cloud backup with labels:

```
srvctl cloud-backups add --volume-id ex4mp1eID --name my-backup --label environment=production
```
