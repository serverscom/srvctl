#### Create custom SSL certificate via input

A command to show the skeleton JSON structure required to create a custom SSL certificate:

```
srvctl ssl custom add --skeleton
```

An example of a command for a file in the same directory with srvctl:

```
srvctl ssl custom add --input <file name>
```

An example of the file's content:

```
{
    "name": "",
    "public_key": "",
    "private_key": "",
    "labels": {
        "key": "value"
    }
}
```

#### Create custom SSL certificate via flags

An example of a command to create a custom SSL certificate via flags:

```
srvctl ssl custom add \
	--name my-certificate \
	--public-key "$(cat cert.pem)" \
	--private-key "$(cat key.pem)" \
	--chain-key "$(cat chain.pem)" \
	--label environment=production
```
