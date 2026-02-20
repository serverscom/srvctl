A command to create a subnetwork in the network pool with the "ex4mp1ePoolID" ID using a CIDR:

```
srvctl network-pools add-subnet ex4mp1ePoolID --cidr 192.168.1.0/24 --title "My subnet"
```

A command to create a subnetwork using a mask:

```
srvctl network-pools add-subnet ex4mp1ePoolID --mask 24
```
