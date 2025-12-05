## hosts

A host is a bare metal-based service that includes dedicated servers, Kubernetes bare metal nodes and scalable bare metal. This section describes all host-related commands.
### Usage
```
srvctl hosts [command] [flags]
```
### Commands

`ds` - opens a list of commands to manage dedicated servers.
```
srvctl hosts ds
```

`kbm` - opens a list of commands to manage Kubernetes Bare Metal nodes.
```
srvctl hosts kbm
```

`list` - lists all hosts.
```
srvctl hosts list
```

`sbm` - opens a list of commands to manage Scalable Bare Metal.
```
srvctl hosts sbm
```
### Flags

- `-f`, `--field` -  a string that selects a field to show. To display several fields, use this flag multiple times.

- `--field-list` - a list of available fields.

- `--page-view` - enables the page view format.

- `-t`, `--template` - a Go template string for advanced customization. Example: `--template "{{range .}}CustomTitle: {{.Title}}\n{{ end }}"` 

- Global flags
### Examples

A command to list all account hosts showing their ID, Title and Status:
```
srvctl hosts list -f ID -f Title -f Status --all
```
## hosts ds

A set of commands to manage dedicated servers.
### Usage
```
srvctl hosts ds [command] [flags]
```
### Commands

`abort-release` - cancels a scheduled release.
```
srvctl hosts ds abort-release ex4mp1eID
```

`add` - creates a dedicated server (see the **hosts ds add** section).

`add-network` - adds a network to the selected server.
```
srvctl hosts ds add-network ex4mp1eID --distribution-method route --mask 32 --type private
```

`add-ptr` - adds a PTR record to the selected server.
```
srvctl hosts ds add-ptr ex4mp1eID --domain example.com --ip X.X.X.X --priority 1 --ttl 360
```

`delete-network` - deletes a specified network for the selected server.
```
srvctl hosts ds delete-network ex4mp1eID_server --network-id ex4mp1eID_network
```

`delete-ptr` - deletes a specified PTR record for the selected server.
```
srvctl hosts ds delete-ptr ex4mp1eID_server --ptr-id ex4mp1eID_ptr
```

`get` - provides information for the selected server.
```
srvctl hosts ds get ex4mp1eID
```

`get-network` - provides information about a specified network of the selected server.
```
srvctl hosts ds get-network ex4mp1eID_server --network-id ex4mp1eID_network
```

`get-oob-credentials` - provides OOB credentials for the selected server. A GPG key fingerprint is needed.
```
srvctl hosts ds get-network ex4mp1eID_server --fingerprint GPGKEYEX4MP1E
```

`ls`, `list` - lists dedicated servers of the account. Use `–help` to see available flags.
```
srvctl hosts ds list ex4mp1eID --all
```

`list-connections` - lists connections for the selected dedicated server. Use `–help` to see available flags.
```
srvctl hosts ds list-connections ex4mp1eID
```

`list-drive-slots` - lists drive slots for the selected dedicated server. Use `–help` to see available flags.
```
srvctl hosts ds list-drive-slots ex4mp1eID
```

`list-features` - lists features for the selected dedicated server. Use `–help` to see available flags.
```
srvctl hosts ds list-features ex4mp1eID
```

`list-networks` - lists networks for the selected dedicated server. Use `–help` to see available flags.
```
srvctl hosts ds list-networks ex4mp1eID
```

`list-power-feeds` - lists power feeds for the selected dedicated server.
```
srvctl hosts ds list-power-feeds ex4mp1eID
```

`list-ptr` - lists PTR records for the selected dedicated server. Use `–help` to see available flags.
```
srvctl hosts ds list-ptr ex4mp1eID
```

`list-services` - lists services for the selected dedicated server.
```
srvctl hosts ds list-ptr ex4mp1eID
```

`power` - sends a specified power command for the selected dedicated server (see the **hosts ds power** section).

`reinstall` - reinstalls an operating system for the selected dedicated server (see the **hosts ds reinstall** section).

`schedule-release` - schedules release on YYYY-MM-DDTHH:MM:SS+HH:MM (dateTtime+time zone from UTC) for the selected dedicated server.
```
srvctl hosts ds schedule-release ex4mp1eID --release-after 2022-05-24T12:48:00+03:00
```

`update` - updates parameters for the selected dedicated server.
```
srvctl hosts ds update ex4mp1eID --label environment=production --label team=frontend
```
### Flags

Global flags
## hosts ds add

A command to create a dedicated server. It allows to pass parameters of the server in two ways:

- Input - server parameters are described in a file, a path to the file is specified via the `-i` or `–input` flag. The path can be absolute or relative to the srvctl file.

- Flags  - parameters are specified via flags inside the command and hostnames are listed as position arguments. As many arguments, as many servers of this configuration will be created.
### Usage
```
srvctl hosts ds add -i [path]
srvctl hosts ds add [flags] [hostname1] [hostnameN]
```
### Flags

- `--location-id` - a unique identifier of a location. Use `srvctl locations list` to list all location IDs.

- `--server-model-id` - a unique identifier of a server model. Use `srvctl server-models list --location-id <id>` to list all server models.

- `--ram-size` - an integer value of RAM in GB.

- `--operating-system-id` - a unique identifier of an operating system. Use `srvctl server-os-options list --location-id <id> --server-model-id <id>` to list all OS IDs.

- `--public-uplink-id` - a unique identifier of a public uplink. Use `srvctl uplink-models list --location-id <id> --server-model-id <id>` to list all public uplinks.

- `--public-bandwidth-id` - a unique identifier of a public uplink bandwidth. Use `srvctl uplink-bandwidths list --location-id <id> --server-model-id <id> --uplink-model-id <id>` to list all bandwidth IDs.

- `--private-uplink-id` - a unique identifier of a private uplink. Use `srvctl uplink-models list --location-id <id> --server-model-id <id>` to list all private uplinks.

- `--drive-slots` - a parameter to specify a disk model ID for a slot. Slots start from 0. This is the flag pattern for two disk slots: `--drive-slots <slot position>=<disk model id> --drive-slots <slot position>=<disk model id>`. Use `srvctl drive-models list --location-id <id> --server-model-id <id>` to list disk model IDs.

- `--layout` - a RAID level for disks in slots. An example of pattern for two disks: `--layout=slot=<position>,slot=<position>,raid=<level>`. An example for RAID 1: `--layout=slot=0,slot=1,raid=1`.

- `--partition` - partition parameters (target, fs, fill, size) for a disk in the specified slot. If two disks are within one RAID, both of them should be listed. This is an example for a boot partition and RAID1 of two disks: `--partition=slot=0,slot=1,target=/boot,fs=ext4,size=1024,fill=false`

- `--feature` - a string with features that has to be named as per [Public API](https://developers.servers.com/api-documentation/v1/#tag/Location/operation/ListLocations). For example, `--feature no_private_ip`

- `--ipv6` - a flag to add a public IPv6 address. It's necessary to specify this option `--ipv6` without additional values to activate the feature: the presence of the flag itself implies IPv6 enabling.

- `--labels` - a flag to assign labels to a server. For example: `--labels environment=production`

- `--user-data` - a flag that processes user data from a string: `--user-data <user data>`

- `--user-data-file` - a flag that collects user data from a specified directory. An example for a file located in the same directory with srvctl: `--user-data-file example.txt`

- `-i`, `--input` - a flag to specify a file with order parameters (see **Create server via input**).
### Examples
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

 The only available authentication method is password. An SSH key can be added only via the input process (see **Create server via input**)
## hosts ds power

A command to manage power operations on a server.
### Usage
```
srvctl hosts ds power <id> [flags]
```
### Flags

- `--command off` - a flag to power off a server.

- `--command on` - a flag to power on a server.

- `--command cycle` - a flag for the power cycle command.
### Examples

An example of a command to switch a server off:
```
srvctl hosts ds power <id> --command off
```
## hosts ds reinstall

A command to reinstall an operating system on a dedicated server.
### Usage
```
srvctl hosts ds reinstall <id> [flags]
```
### Flags

- `-i`, `--input` - a flag to specify a path with a file that contains reinstall parameters.
### Examples

The `-i`, `--input` allows to provide parameters of a created server in a local file. Parameters should be described as a request body of the [Public API request](https://developers.servers.com/api-documentation/v1/#tag/Dedicated-Server/operation/CreateADedicatedServer).

An example of a command for a file in the same directory with srvctl:
```
srvctl hosts ds reinstall <id> --input <file name>
```

An example of the file's content:
```
{
  "hostname": "<give a name>",
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
  "operating_system_id": 49
}
```

There is also an option to use standard input (stdin) when specifying the flag this way: `--input -`
