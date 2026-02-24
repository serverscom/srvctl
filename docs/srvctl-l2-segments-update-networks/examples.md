A command to show the skeleton JSON structure required to update L2 segment networks:

```
srvctl l2-segments update-networks --skeleton
```

A command to update networks for the L2 segment with the "ex4mp1eID" ID using an input file:

```
srvctl l2-segments update-networks ex4mp1eID --input <file name>
```

An example of the file's content:

```
{
  "create": [
    {
      "mask": "",
      "distribution_method": ""
    }
  ],
  "delete": [
    "id"
  ]
}
```
