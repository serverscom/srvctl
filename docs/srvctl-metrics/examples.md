A command to get metrics of all hosts:

```
srvctl metrics hosts
```

A command to get metrics of all private racks:

```
srvctl metrics racks
```

A command to collect hosts metrics for a Prometheus textfile collector:

```
srvctl metrics hosts --output raw > /var/lib/node_exporter/textfile_collector/serverscom_hosts.prom
```
