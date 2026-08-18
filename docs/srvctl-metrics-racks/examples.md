A command to get private racks metrics as a table:

```
srvctl metrics racks
```

A command to get power metrics of private racks:

```
srvctl metrics racks --field RackID --field PduWatts --field PduAmperes --field AtsWatts --field AtsAmperes
```

A command to get private racks metrics in the Prometheus text exposition format:

```
srvctl metrics racks --output raw
```
