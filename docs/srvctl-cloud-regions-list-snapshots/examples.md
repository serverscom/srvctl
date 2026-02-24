A command to list all snapshots for the cloud region with ID "42":

```
srvctl cloud-regions list-snapshots 42
```

A command to filter snapshots by instance:

```
srvctl cloud-regions list-snapshots 42 --instance-id ex4mp1eID
```

A command to filter snapshots by backup status:

```
srvctl cloud-regions list-snapshots 42 --is-backup
```
