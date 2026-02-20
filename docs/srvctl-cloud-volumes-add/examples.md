A command to show the skeleton JSON structure required to create a cloud volume:

```
srvctl cloud-volumes add --skeleton
```

An example of a command for a file in the same directory with srvctl:

```
srvctl cloud-volumes add --input <file name>
```

An example of the file's content:

```
{
    "name": "",
    "description": "",
    "attach_instance_id": "",
    "region_id": "",
    "size": "",
    "labels": {
        "key": "value"
    }
}
```

An example of a command to create a cloud volume via flags:

```
srvctl cloud-volumes add \
	--name my-volume \
	--region-id 1 \
	--size 50 \
	--label environment=production
```
