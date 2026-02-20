A command to get the JSON structure required to create a new L7 load balancer:

```
srvctl lb l7 add --skeleton
```

A command to create a new L7 load balancer using a JSON file:

```
srvctl lb l7 add --input /path/to/input.json
```

A command to create a new L7 load balancer using stdin:

```
cat input.json | srvctl lb l7 add --input -
```

JSON structure example:

```json
{
    "name": "",
    "location_id": "",
    "cluster_id": "",
    "store_logs": true,
    "store_logs_region_id": "",
    "geoip": "true",
    "vhost_zones": [
        {
            "id": "",
            "ports": [],
            "ssl": true,
            "http2": true,
            "ssl_certificate_id": "",
            "proxy_protocol": true,
            "upstream_id": "",
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
