#### Create L2 segment via input

A command to show the skeleton JSON structure required to create an L2 segment:

```
srvctl l2-segments add --skeleton
```

An example of a command for a file in the same directory with srvctl:

```
srvctl l2-segments add --input <file name>
```

An example of the file's content:

```
{
    "name": "",
    "type": "",
    "location_group_id": "",
    "members": [
        {
            "id": "",
            "mode": ""
        }
    ],
    "labels": {
        "key": "value"
    }
}
```

#### Create L2 segment via flags

An example of a command to create an L2 segment via flags:

```
srvctl l2-segments add \
	--type public \
	--member id=ex4mp1eID,mode=native \
	--member id=ex4mp1eID2,mode=trunk \
	--name my-segment \
	--location-group-id 42 \
	--label environment=production
```
