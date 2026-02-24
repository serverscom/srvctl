A command to list all SSL certificates for the account:

```
srvctl ssl list
```

A command to filter SSL certificates by label:

```
srvctl ssl list --label-selector environment=production
```

A command to search SSL certificates by name:

```
srvctl ssl list --search-pattern my-cert
```
