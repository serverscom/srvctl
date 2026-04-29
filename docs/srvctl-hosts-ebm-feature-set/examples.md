A command to activate disaggregated public ports on the server with the "ex4mp1eID" ID:

```
srvctl hosts ebm feature-set ex4mp1eID --feature disaggregated_public_ports --state activate
```

A command to deactivate disaggregated public ports on the server with the "ex4mp1eID" ID:

```
srvctl hosts ebm feature-set ex4mp1eID --feature disaggregated_public_ports --state deactivate
```

A command to activate rescue mode with password and SSH key authentication:

```
srvctl hosts ebm feature-set ex4mp1eID \
	--feature host_rescue_mode \
	--state activate \
	--auth-method password \
	--auth-method ssh_key \
	--ssh-key-fingerprint aa:bb:cc:dd:ee:ff
```

A command to activate private iPXE boot with a custom script:

```
srvctl hosts ebm feature-set ex4mp1eID \
	--feature private_ipxe_boot \
	--state activate \
	--ipxe-config "#!ipxe\nchain http://boot.example.com/script.ipxe"
```

A command to deactivate rescue mode:

```
srvctl hosts ebm feature-set ex4mp1eID --feature host_rescue_mode --state deactivate
```