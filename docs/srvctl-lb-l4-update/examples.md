A command to get the JSON structure required to update an L4 load balancer:

```
srvctl lb l4 update --skeleton
```

A command to update the L4 load balancer with the "ex4mp1eID" ID using a JSON file:

```
srvctl lb l4 update ex4mp1eID --input /path/to/input.json
```

A command to update the L4 load balancer using stdin:

```
cat input.json | srvctl lb l4 update ex4mp1eID --input -
```

JSON structure example:

```json
{
    "cluster_id": "",
    "shared_cluster": "",
    "name": "",
    "store_logs": true,
    "store_logs_region_id": "",
    "new_external_ips_count": "",
    "delete_external_ips": [],
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
