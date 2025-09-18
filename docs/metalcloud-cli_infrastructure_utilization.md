## metalcloud-cli infrastructure utilization

Get resource utilization report for infrastructures

### Synopsis

Get detailed utilization report for infrastructure resources within a specified time range. 
The report provides insights into resource usage patterns and capacity planning for infrastructures.

Required flags:
  --user-id       ID of the user to include in the report
  --start-time    Start time for the report (RFC3339 or date format)
  --end-time      End time for the report (RFC3339 or date format)

Optional flags:
  --site-id            Site IDs to include in the report (can be specified multiple times)
  --infrastructure-id  Infrastructure IDs to include in the report (can be specified multiple times)
  --show-all           Show all utilizations
  --show-instances     Show instance utilizations
  --show-drives        Show drive utilizations
  --show-subnets       Show subnet utilizations

Examples:
  # Get utilization report for user 123 for the last 7 days
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08

  # Get utilization for specific sites and infrastructures
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01T00:00:00Z --end-time 2025-08-08T23:59:59Z --site-id 1 --site-id 2 --infrastructure-id 100 --infrastructure-id 101

  # Get utilization for all infrastructures of a user in specific sites
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08 --site-id 1 --site-id 3

```
metalcloud-cli infrastructure utilization [flags]
```

### Options

```
      --end-time time            End time for the report. (default 2025-09-18T13:16:40.971700756+03:00)
  -h, --help                     help for utilization
      --infrastructure-id ints   Infrastructure IDs to include in the report.
      --show-all                 If set, will display all utilizations.
      --show-drives              If set, will display drive utilizations.
      --show-instances           If set, will display instance utilizations.
      --show-subnets             If set, will display subnet utilizations.
      --site-id ints             Site IDs to include in the report.
      --start-time time          Start time for the report. (default 2025-09-18T13:16:40.971660385+03:00)
      --user-id int              ID of the user to include in the report.
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API key
  -c, --config string          Config file path
  -d, --debug                  Set to enable debug logging
  -e, --endpoint string        MetalCloud API endpoint
  -f, --format string          Output format. Supported values are 'text','csv','md','json','yaml'. (default "text")
  -i, --insecure_skip_verify   Set to allow insecure transport
  -l, --log_file string        Log file path
  -v, --verbosity string       Log level verbosity (default "INFO")
```

### SEE ALSO

* [metalcloud-cli infrastructure](metalcloud-cli_infrastructure.md)	 - Manage infrastructure resources and configurations

