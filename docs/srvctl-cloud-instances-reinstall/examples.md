A command to reinstall the cloud instance with the "ex4mp1eID" ID:

```
srvctl cloud-instances reinstall ex4mp1eID --image-id ex4mp1eImageID
```

A command to reinstall a cloud instance with custom user data:

```
srvctl cloud-instances reinstall ex4mp1eID --image-id ex4mp1eImageID --user-data "#cloud-config"
```
