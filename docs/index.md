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
| [srvctl account](srvctl-account/description.md) | Account | This command allows to manage account operations. |
| [srvctl account balance](srvctl-account-balance/description.md) | Account | This command provides account balance information. |
| [srvctl cloud-backups](srvctl-cloud-backups/description.md) | Cloud Backups | This command allows to manage cloud backups. |
| [srvctl cloud-backups list](srvctl-cloud-backups-list/description.md) | Cloud Backups | This command lists cloud backups of the account. |
| [srvctl cloud-backups get](srvctl-cloud-backups-get/description.md) | Cloud Backups | This command provides information for the selected cloud backup. |
| [srvctl cloud-backups add](srvctl-cloud-backups-add/description.md) | Cloud Backups | A command to create a new cloud backup. |
| [srvctl cloud-backups update](srvctl-cloud-backups-update/description.md) | Cloud Backups | This command updates labels for the selected cloud backup. |
| [srvctl cloud-backups restore](srvctl-cloud-backups-restore/description.md) | Cloud Backups | This command restores a cloud backup to the specified volume. |
| [srvctl cloud-backups delete](srvctl-cloud-backups-delete/description.md) | Cloud Backups | This command deletes the selected cloud backup. |
| [srvctl cloud-instances](srvctl-cloud-instances/description.md) | Cloud Instances | This command allows to manage cloud instances. |
| [srvctl cloud-instances list](srvctl-cloud-instances-list/description.md) | Cloud Instances | This command lists cloud instances of the account. |
| [srvctl cloud-instances get](srvctl-cloud-instances-get/description.md) | Cloud Instances | This command provides information for the selected cloud instance. |
| [srvctl cloud-instances add](srvctl-cloud-instances-add/description.md) | Cloud Instances | A command to create a new cloud instance. |
| [srvctl cloud-instances update](srvctl-cloud-instances-update/description.md) | Cloud Instances | This command updates parameters and labels for the selected cloud instance. |
| [srvctl cloud-instances delete](srvctl-cloud-instances-delete/description.md) | Cloud Instances | This command deletes the selected cloud instance. |
| [srvctl cloud-instances reinstall](srvctl-cloud-instances-reinstall/description.md) | Cloud Instances | This command reinstalls the operating system on the selected cloud instance. |
| [srvctl cloud-instances upgrade](srvctl-cloud-instances-upgrade/description.md) | Cloud Instances | This command upgrades the selected cloud instance to a new flavor. |
| [srvctl cloud-instances upgrade-approve](srvctl-cloud-instances-upgrade-approve/description.md) | Cloud Instances | This command approves a pending upgrade for the selected cloud instance. |
| [srvctl cloud-instances upgrade-revert](srvctl-cloud-instances-upgrade-revert/description.md) | Cloud Instances | This command reverts a pending upgrade for the selected cloud instance. |
| [srvctl cloud-instances reboot](srvctl-cloud-instances-reboot/description.md) | Cloud Instances | This command reboots the selected cloud instance. |
| [srvctl cloud-instances rescue-mode](srvctl-cloud-instances-rescue-mode/description.md) | Cloud Instances | This command activates or deactivates rescue mode for the selected cloud instance. |
| [srvctl cloud-instances power](srvctl-cloud-instances-power/description.md) | Cloud Instances | This command powers on or off the selected cloud instance. |
| [srvctl cloud-instances list-ptr](srvctl-cloud-instances-list-ptr/description.md) | Cloud Instances | This command lists PTR records for the selected cloud instance. |
| [srvctl cloud-instances add-ptr](srvctl-cloud-instances-add-ptr/description.md) | Cloud Instances | This command adds a PTR record to the selected cloud instance. |
| [srvctl cloud-instances delete-ptr](srvctl-cloud-instances-delete-ptr/description.md) | Cloud Instances | This command deletes a PTR record from the selected cloud instance. |
| [srvctl cloud-regions](srvctl-cloud-regions/description.md) | Cloud Regions | This command allows to manage cloud regions. |
| [srvctl cloud-regions list](srvctl-cloud-regions-list/description.md) | Cloud Regions | This command lists available cloud regions. |
| [srvctl cloud-regions get-credentials](srvctl-cloud-regions-get-credentials/description.md) | Cloud Regions | This command provides credentials for the selected cloud region. |
| [srvctl cloud-regions list-flavors](srvctl-cloud-regions-list-flavors/description.md) | Cloud Regions | This command lists available flavors for the selected cloud region. |
| [srvctl cloud-regions list-images](srvctl-cloud-regions-list-images/description.md) | Cloud Regions | This command lists available images for the selected cloud region. |
| [srvctl cloud-regions list-snapshots](srvctl-cloud-regions-list-snapshots/description.md) | Cloud Regions | This command lists snapshots for the selected cloud region. |
| [srvctl cloud-regions add-snapshot](srvctl-cloud-regions-add-snapshot/description.md) | Cloud Regions | A command to create a new snapshot for the selected cloud region. |
| [srvctl cloud-regions delete-snapshot](srvctl-cloud-regions-delete-snapshot/description.md) | Cloud Regions | This command deletes a snapshot from the selected cloud region. |
| [srvctl cloud-volumes](srvctl-cloud-volumes/description.md) | Cloud Volumes | This command allows to manage cloud volumes. |
| [srvctl cloud-volumes list](srvctl-cloud-volumes-list/description.md) | Cloud Volumes | This command lists cloud volumes of the account. |
| [srvctl cloud-volumes get](srvctl-cloud-volumes-get/description.md) | Cloud Volumes | This command provides information for the selected cloud volume. |
| [srvctl cloud-volumes add](srvctl-cloud-volumes-add/description.md) | Cloud Volumes | A command to create a new cloud volume. |
| [srvctl cloud-volumes update](srvctl-cloud-volumes-update/description.md) | Cloud Volumes | This command updates parameters and labels for the selected cloud volume. |
| [srvctl cloud-volumes delete](srvctl-cloud-volumes-delete/description.md) | Cloud Volumes | This command deletes the selected cloud volume. |
| [srvctl cloud-volumes volume-attach](srvctl-cloud-volumes-volume-attach/description.md) | Cloud Volumes | This command attaches a cloud volume to a cloud instance. |
| [srvctl cloud-volumes volume-detach](srvctl-cloud-volumes-volume-detach/description.md) | Cloud Volumes | This command detaches a cloud volume from a cloud instance. |

