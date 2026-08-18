Get metrics of all hosts of the account.

In the default text format each row represents a host with its monthly traffic split by public and private traffic. Only hosts that have traffic data are listed, as the API labels a host with its id only when it reports traffic counters for it. Use `--field-list` to see all available fields, `--field` to pick the ones you need and `--page-view` to print a field per line.

All the metrics come in a single API response, so `--per-page`, `--page` and `--all` are applied locally. 20 rows are printed per page by default, use `--all` to print all of them.

With `--output raw` the metrics are printed in the Prometheus text exposition format as returned by the API, including the hosts count metric that has no dedicated column in the table.
