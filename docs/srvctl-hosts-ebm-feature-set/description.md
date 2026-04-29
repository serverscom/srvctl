This command activates or deactivates a feature on the selected enterprise bare metal server.

The `--feature` and `--state` flags are required. The `--state` flag accepts `activate` or `deactivate`.

Supported feature names:

- `disaggregated_public_ports` - disaggregated public network ports.
- `disaggregated_private_ports` - disaggregated private network ports.
- `no_public_ip_address` - disable public IP address assignment.
- `no_private_ip` - disable private IP address assignment.
- `oob_public_access` - enable public access to out-of-band management.
- `no_public_network` - disable public network connectivity.
- `host_rescue_mode` - boot into rescue mode. When activating, use `--auth-method` (repeatable: `password`, `ssh_key`) and `--ssh-key-fingerprint` (repeatable) to configure access.
- `private_ipxe_boot` - boot from a private iPXE script. When activating, use `--ipxe-config` to supply the iPXE script.