A command to add a PTR record to the cloud instance with the "ex4mp1eID" ID:

```
srvctl cloud-instances add-ptr ex4mp1eID --data hostname.example.com --ip 192.0.2.1
```

A command to add a PTR record with a custom TTL:

```
srvctl cloud-instances add-ptr ex4mp1eID --data hostname.example.com --ip 192.0.2.1 --ttl 3600
```
