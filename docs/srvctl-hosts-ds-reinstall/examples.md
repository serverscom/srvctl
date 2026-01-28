An example of a command for the server with the "ex4mp1eID" ID to reinstall OS using a file in the same directory with srvctl:

```
srvctl hosts ds reinstall ex4mp1eID --input <file name>
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
