A command to rename the cloud instance with the "ex4mp1eID" ID:

```
srvctl cloud-instances update ex4mp1eID --name new-name
```

A command to assign labels for the cloud instance with the "ex4mp1eID" ID:

```
srvctl cloud-instances update ex4mp1eID --label environment=production --label team=frontend
```

A command to enable IPv6 and set backup copies for the cloud instance:

```
srvctl cloud-instances update ex4mp1eID --ipv6-enabled --backup-copies 3
```
