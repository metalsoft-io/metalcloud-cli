## metalcloud-cli infrastructure utilization-summary

Get a summarized resource utilization report for infrastructures

### Synopsis

Get a summarized utilization report for infrastructure resources within a specified time range.

Unlike 'infrastructure utilization', which returns every metered resource of every
infrastructure, this command returns the aggregated quantity per resource type plus the
internet upload/download totals. Use --format json or yaml to get the full response,
including the metered waypoints and the reservation and license installments.

Required flags:
  --user-id       ID of the user the report is generated for
  --start-time    Start time for the report (RFC3339 or date format)
  --end-time      End time for the report (RFC3339 or date format)

Optional flags:
  --infrastructure-id  Infrastructure IDs to include in the report (can be specified multiple times)

Examples:
  # Summarized utilization for user 123 over a week
  metalcloud-cli infrastructure utilization-summary --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08

  # Summarized utilization for specific infrastructures, as JSON
  metalcloud-cli infrastructure utilization-summary --user-id 123 --start-time 2025-08-01T00:00:00Z --end-time 2025-08-08T23:59:59Z --infrastructure-id 100 --infrastructure-id 101 -f json

```
metalcloud-cli infrastructure utilization-summary [flags]
```

### Options

```
      --end-time time            End time for the report. (default 2026-09-20T15:26:54.636786915+03:00)
  -h, --help                     help for utilization-summary
      --infrastructure-id ints   Infrastructure IDs to include in the report.
      --start-time time          Start time for the report. (default 2026-09-20T15:26:54.636786415+03:00)
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

