A command to list subnets of the network pool with the "ex4mp1eID" ID:

```
srvctl network-pools list-subnets ex4mp1eID
```

A command to filter only private subnets:

```
srvctl network-pools list-subnets ex4mp1eID --type private
```

A command to filter subnets by search pattern:

```
srvctl network-pools list-subnets ex4mp1eID --search-pattern "10.0"
```
