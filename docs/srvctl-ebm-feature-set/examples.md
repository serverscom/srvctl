A command to activate disaggregated public ports on the server with the "ex4mp1eID" ID:

```
srvctl ebm feature-set ex4mp1eID --feature disaggregated_public_ports --command activate
```

A command to deactivate disaggregated public ports on the server with the "ex4mp1eID" ID:

```
srvctl ebm feature-set ex4mp1eID --feature disaggregated_public_ports --command deactivate
```

A command to activate rescue mode with password and SSH key authentication:

```
srvctl ebm feature-set ex4mp1eID \
	--feature host_rescue_mode \
	--command activate \
	--auth-method password \
	--auth-method ssh_key \
	--ssh-key-fingerprint aa:bb:cc:dd:ee:ff
```

A command to activate private iPXE boot with a custom script:

```
srvctl ebm feature-set ex4mp1eID \
	--feature private_ipxe_boot \
	--command activate \
	--ipxe-config "#!ipxe\nchain http://boot.example.com/script.ipxe"
```

A command to deactivate rescue mode:

```
srvctl ebm feature-set ex4mp1eID --feature host_rescue_mode --command deactivate
```