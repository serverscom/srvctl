A command to list all Let's Encrypt SSL certificates for the account:

```
srvctl ssl le list
```

A command to filter Let's Encrypt SSL certificates by label:

```
srvctl ssl le list --label-selector environment=production
```

A command to search Let's Encrypt SSL certificates by name:

```
srvctl ssl le list --search-pattern my-cert
```
