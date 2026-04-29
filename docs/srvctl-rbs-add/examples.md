An example of a command for a file in the same directory with srvctl:

```
srvctl rbs add --input <file name>
```

An example of the file's content:

```
{
    "name": "",
    "size": 0,
    "location_id": 0,
    "flavor_id": 0,
    "labels": {
        "key": "value"
    }
}
```

An example of a command to create a remote block storage volume via flags:

```
srvctl rbs add \
	--name my-volume \
	--size 100 \
	--flavor-id 1 \
	--location-id 1 \
	--label environment=production
```