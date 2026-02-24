A command to get the JSON structure required to update an L7 load balancer:

```
srvctl lb l7 update --skeleton
```

A command to update the L7 load balancer with the "ex4mp1eID" ID using a JSON file:

```
srvctl lb l7 update ex4mp1eID --input /path/to/input.json
```

A command to update the L7 load balancer using stdin:

```
cat input.json | srvctl lb l7 update ex4mp1eID --input -
```

JSON structure example:

```json
{
    "name": "",
    "cluster_id": "",
    "shared_cluster": "",
    "store_logs": true,
    "store_logs_region_id": "",
    "geoip": "true",
    "new_external_ips_count": "",
    "delete_external_ips": [],
    "vhost_zones": [
        {
            "id": "",
            "ports": [],
            "ssl": true,
            "http2": true,
            "ssl_certificate_id": "",
            "domains": [],
            "location_zones": [
                {
                    "location": "",
                    "upstream_path": "",
                    "upstream_id": ""
                }
            ],
            "real_ip_header": {
                "name": "",
                "networks": []
            }
        }
    ],
    "upstream_zones": [
        {
            "id": "",
            "ssl": true,
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
