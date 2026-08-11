A command to get hosts metrics as a table:

```
srvctl metrics hosts
```

A command to get hosts metrics with specific fields:

```
srvctl metrics hosts --field HostID --field ChassisName --field TotalSent
```

A command to get the second page of hosts metrics:

```
srvctl metrics hosts --page 2
```

A command to get metrics of all the hosts at once:

```
srvctl metrics hosts -A
```

A command to get hosts metrics in the Prometheus text exposition format:

```
srvctl metrics hosts --output raw
```
