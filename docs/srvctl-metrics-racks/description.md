Get metrics of all private racks of the account.

In the default text format each row represents a rack with the number of hosts in it, its monthly traffic and the power draw of its PDU and ATS devices. PDU and ATS values are summed per rack and kept in separate columns, as an ATS feeds the PDUs and summing them would count the same draw twice. Use `--field-list` to see all available fields, `--field` to pick the ones you need and `--page-view` to print a field per line.

All the metrics come in a single API response, so `--per-page`, `--page` and `--all` are applied locally. 20 rows are printed per page by default, use `--all` to print all of them.

With `--output raw` the metrics are printed in the Prometheus text exposition format as returned by the API, with power and current reported per device.
