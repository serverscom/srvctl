A command to list all custom SSL certificates for the account:

```
srvctl ssl custom list
```

A command to filter custom SSL certificates by label:

```
srvctl ssl custom list --label-selector environment=production
```

A command to search custom SSL certificates by name:

```
srvctl ssl custom list --search-pattern my-cert
```
