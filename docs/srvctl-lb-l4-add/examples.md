A command to get the JSON structure required to create a new L4 load balancer:

```
srvctl lb l4 add --skeleton
```

A command to create a new L4 load balancer using a JSON file:

```
srvctl lb l4 add --input /path/to/input.json
```

A command to create a new L4 load balancer using stdin:

```
cat input.json | srvctl lb l4 add --input -
```

JSON structure example:

```json
{
    "name": "",
    "location_id": "",
    "cluster_id": "",
    "store_logs": true,
    "store_logs_region_id": "",
    "vhost_zones": [
        {
            "id": "",
            "ports": [],
            "udp": true,
            "proxy_protocol": true,
            "upstream_id": ""
        }
    ],
    "upstream_zones": [
        {
            "id": "",
            "udp": true,
            "upstreams": [
                {
                    "ip": "",
                    "port": "",
                    "weight": ""
                }
            ]
        }
    ],
    "labels": {
        "key": "value"
    }
}
```
