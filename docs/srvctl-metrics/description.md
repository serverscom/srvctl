You can get metrics for your hosts and private racks by performing commands listed in `srvctl metrics --help`.

Metrics support two output formats:

- `--output text` (default) folds the metrics into a table with one row per host or rack, with traffic humanized.
- `--output raw` prints the metrics in the Prometheus text exposition format exactly as returned by the API, which is handy to feed a Prometheus textfile collector.
