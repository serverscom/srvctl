# CLI Documentation

## Table of Contents

| Command | Type | Description |
|---------|------|-------------|
| [srvctl hosts](srvctl-hosts/description.md) | Hosts | A host is a bare metal-based service that includes dedicated servers, Kubernetes bare metal nodes and scalable bare metal. (ds, kbm, sbm) |
| [srvctl hosts ds](srvctl-hosts-ds/description.md) | Hosts / Dedicated Servers | This command allows to manage dedicated servers (Enterprise bare metal). |
| [srvctl hosts ds abort-release](srvctl-hosts-ds-abort-release/description.md) | Hosts / Dedicated Servers | This command cancels the scheduled release for the selected dedicated server. |
| [srvctl hosts ds add](srvctl-hosts-ds-add/description.md) | Hosts / Dedicated Servers | A command to create a dedicated server. |
| [srvctl hosts ds add-network](srvctl-hosts-ds-add-network/description.md) | Hosts / Dedicated Servers | This command adds a network to the selected server. |
| [srvctl hosts ds add-ptr](srvctl-hosts-ds-add-ptr/description.md) | Hosts / Dedicated Servers | This command adds a PTR record to the selected server. |
| [srvctl hosts ds delete-network](srvctl-hosts-ds-delete-network/description.md) | Hosts / Dedicated Servers | This command deletes a specified network for the selected server. |
| [srvctl hosts ds delete-ptr](srvctl-hosts-ds-delete-ptr/description.md) | Hosts / Dedicated Servers | This command deletes a specified PTR record for the selected server. |
| [srvctl hosts ds get](srvctl-hosts-ds-get/description.md) | Hosts / Dedicated Servers | This command provides information for the selected server. |
| [srvctl hosts ds get-network](srvctl-hosts-ds-get-network/description.md) | Hosts / Dedicated Servers | This command provides information about a specified network of the selected server. |
| [srvctl hosts ds get-oob-credentials](srvctl-hosts-ds-get-oob-credentials/description.md) | Hosts / Dedicated Servers | This command provides OOB credentials for the selected server. |
| [srvctl hosts ds list-connections](srvctl-hosts-ds-list-connections/description.md) | Hosts / Dedicated Servers | This command lists connections for the selected dedicated server. |
| [srvctl hosts ds list-drive-slots](srvctl-hosts-ds-list-drive-slots/description.md) | Hosts / Dedicated Servers | This command lists drive slots for the selected dedicated server. |
| [srvctl hosts ds list-features](srvctl-hosts-ds-list-features/description.md) | Hosts / Dedicated Servers | This command lists features for the selected dedicated server. |
| [srvctl hosts ds list-networks](srvctl-hosts-ds-list-networks/description.md) | Hosts / Dedicated Servers | This command lists networks for the selected dedicated server. |
| [srvctl hosts ds list-power-feeds](srvctl-hosts-ds-list-power-feeds/description.md) | Hosts / Dedicated Servers | This command lists power feeds for the selected dedicated server. |
| [srvctl hosts ds list-ptr](srvctl-hosts-ds-list-ptr/description.md) | Hosts / Dedicated Servers | This command lists PTR records for the selected dedicated server. |
| [srvctl hosts ds list-services](srvctl-hosts-ds-list-services/description.md) | Hosts / Dedicated Servers | This command lists services for the selected dedicated server. |
| [srvctl hosts ds ls](srvctl-hosts-ds-ls/description.md) | Hosts / Dedicated Servers | This command lists dedicated servers of the account. |
| [srvctl hosts ds power](srvctl-hosts-ds-power/description.md) | Hosts / Dedicated Servers | This command manages power operations for the selected dedicated server. |
| [srvctl hosts ds reinstall](srvctl-hosts-ds-reinstall/description.md) | Hosts / Dedicated Servers | This command reinstalls an operating system for the selected dedicated server. |
| [srvctl hosts ds schedule-release](srvctl-hosts-ds-schedule-release/description.md) | Hosts / Dedicated Servers | This command schedules release on YYYY-MM-DDTHH:MM:SS+HH:MM (dateTtime+time zone from UTC) for the selected dedicated server. |
| [srvctl hosts ds update](srvctl-hosts-ds-update/description.md) | Hosts / Dedicated Servers | This command updates parameters and labels for the selected dedicated server. |
| [srvctl l2-segments](srvctl-l2-segments/description.md) | L2 Segments | This command allows to manage L2 segments. |
| [srvctl l2-segments list](srvctl-l2-segments-list/description.md) | L2 Segments | This command lists L2 segments of the account. |
| [srvctl l2-segments list-groups](srvctl-l2-segments-list-groups/description.md) | L2 Segments | This command lists L2 segment location groups. |
| [srvctl l2-segments list-members](srvctl-l2-segments-list-members/description.md) | L2 Segments | This command lists members of the selected L2 segment. |
| [srvctl l2-segments list-networks](srvctl-l2-segments-list-networks/description.md) | L2 Segments | This command lists networks of the selected L2 segment. |
| [srvctl l2-segments get](srvctl-l2-segments-get/description.md) | L2 Segments | This command provides information for the selected L2 segment. |
| [srvctl l2-segments add](srvctl-l2-segments-add/description.md) | L2 Segments | A command to create a new L2 segment. |
| [srvctl l2-segments update](srvctl-l2-segments-update/description.md) | L2 Segments | This command updates the selected L2 segment. |
| [srvctl l2-segments update-networks](srvctl-l2-segments-update-networks/description.md) | L2 Segments | This command updates networks of the selected L2 segment. |
| [srvctl l2-segments delete](srvctl-l2-segments-delete/description.md) | L2 Segments | This command deletes the selected L2 segment. |
| [srvctl invoices](srvctl-invoices/description.md) | Invoices | This command allows to manage invoices. |
| [srvctl invoices list](srvctl-invoices-list/description.md) | Invoices | This command lists invoices of the account. |
| [srvctl invoices get](srvctl-invoices-get/description.md) | Invoices | This command provides information for the selected invoice. |
| [srvctl ssl](srvctl-ssl/description.md) | SSL Certificates | This command allows to manage SSL certificates of different types (custom, Let's Encrypt). |
| [srvctl ssl list](srvctl-ssl-list/description.md) | SSL Certificates | This command lists all SSL certificates of the account. |
| [srvctl ssl custom](srvctl-ssl-custom/description.md) | SSL Certificates / Custom | This command allows to manage custom SSL certificates. |
| [srvctl ssl custom list](srvctl-ssl-custom-list/description.md) | SSL Certificates / Custom | This command lists custom SSL certificates of the account. |
| [srvctl ssl custom get](srvctl-ssl-custom-get/description.md) | SSL Certificates / Custom | This command provides information for the selected custom SSL certificate. |
| [srvctl ssl custom add](srvctl-ssl-custom-add/description.md) | SSL Certificates / Custom | A command to create a new custom SSL certificate. |
| [srvctl ssl custom update](srvctl-ssl-custom-update/description.md) | SSL Certificates / Custom | This command updates labels for the selected custom SSL certificate. |
| [srvctl ssl custom delete](srvctl-ssl-custom-delete/description.md) | SSL Certificates / Custom | This command deletes the selected custom SSL certificate. |
| [srvctl ssl le](srvctl-ssl-le/description.md) | SSL Certificates / Let's Encrypt | This command allows to manage Let's Encrypt SSL certificates. |
| [srvctl ssl le list](srvctl-ssl-le-list/description.md) | SSL Certificates / Let's Encrypt | This command lists Let's Encrypt SSL certificates of the account. |
| [srvctl ssl le get](srvctl-ssl-le-get/description.md) | SSL Certificates / Let's Encrypt | This command provides information for the selected Let's Encrypt SSL certificate. |
| [srvctl ssl le update](srvctl-ssl-le-update/description.md) | SSL Certificates / Let's Encrypt | This command updates labels for the selected Let's Encrypt SSL certificate. |
| [srvctl ssl le delete](srvctl-ssl-le-delete/description.md) | SSL Certificates / Let's Encrypt | This command deletes the selected Let's Encrypt SSL certificate. |

