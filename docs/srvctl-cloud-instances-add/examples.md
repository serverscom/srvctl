A command to show the skeleton JSON structure required to create a cloud instance:

```
srvctl cloud-instances add --skeleton
```

An example of a command for a file in the same directory with srvctl:

```
srvctl cloud-instances add --input <file name>
```

An example of the file's content:

```
{
    "name": "",
    "flavor_id": "",
    "image_id": "",
    "gpn_enabled": false,
    "ipv6_enabled": true,
    "region_id": "",
    "user_data": "",
    "ssh_key_fingerprint": "",
    "labels": {
        "key": "value"
    }
}
```
