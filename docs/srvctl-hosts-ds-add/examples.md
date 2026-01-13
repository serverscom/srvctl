#### Create server via input 

The `-i`, `--input` allows to provide parameters of a created server in a local file. Parameters should be described as a request body of the [Public API request](https://developers.servers.com/api-documentation/v1/#tag/Dedicated-Server/operation/CreateADedicatedServer).

An example of a command for a file in the same directory with srvctl:
```
srvctl hosts ds add --input <file name>
```

An example of the file's content:
```
{
  "server_model_id": 10515,
  "location_id": 2,
  "ram_size": 32,
  "uplink_models": {
    "public": {
      "id": 10198,
      "bandwidth_model_id": 13744
    },
    "private": {
      "id": 10201
    }
  },
  "drives": {
    "slots": [
      {
        "position": 0,
        "drive_model_id": 10306
      },
      {
        "position": 1,
        "drive_model_id": 10306
      }
    ],
    "layout": [
      {
        "slot_positions": [
          0,
          1
        ],
        "raid": 1,
        "partitions": [
          {
            "target": "/",
            "size": 10240,
            "fill": false,
            "fs": "ext4"
          },
          {
            "target": "/boot",
            "size": 1024,
            "fill": false,
            "fs": "ext4"
          },
          {
            "target": "/home",
            "size": 1,
            "fill": true,
            "fs": "ext4"
          }
        ]
      }
    ]
  },
  "ipv6": false,
  "hosts": [
    {
      "hostname": "<give a name>"
    }
  ],
  "operating_system_id": 62,
  "ssh_key_fingerprints": [
    "<fingerprint of an SSH key>"
  ]
}
```

There is also an option to use standard input (stdin) when specifying the flag this way: `--input -`

#### Create server via flags

It's possible to pass server parameters via flags that are described in the **Flags** section. This is an example of a command to create a dedicated server:
```
srvctl hosts ds add \
	--location-id 2 \
	--server-model-id 10515 \
	--ram-size 32 \
	--operating-system-id 62 \
	--public-uplink-id 10198 \
	--public-bandwidth-id 13744 \
	--private-uplink-id 10201 \
	--drive-slots 1=10306 \
	--drive-slots 2=10306 \
	--layout=slot=0,slot=1,raid=1 \
	--partition=slot=0,slot=1,target=/,fs=ext4,size=10240,fill=false \
	--partition=slot=0,slot=1,target=/boot,fs=ext4,size=1024,fill=false \
	--partition=slot=0,slot=1,target=/home,fs=ext4,size=1,fill=true \
	--feature no_private_ip \
	--ipv6 \
	<hostname>
```

The only available authentication method is password. An SSH key can be added only via the input process (see **Create server via input**).
