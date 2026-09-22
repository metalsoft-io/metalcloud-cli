## metalcloud-cli server-type statistics

Get the server availability statistics of a site

### Synopsis

Get the server availability statistics of a site.

The report contains the number of available servers of each server type in the
site, the details of those servers and the utilization report grouped by RAM,
server type name, product name and owner. Use the json or yaml output format to
see the nested per-server information and the utilization report.

Optional Arguments:
  server_type_id...              Restrict the report to these server type IDs

Required Flags:
  --site-id                      The ID of the site to report on

Optional Flags:
  --user-id                      Only count the resources owned by this user ID
  --max-results-per-server-type  Maximum number of servers returned per server type
  --instance-array-id            Treat only the active instances of this instance array as available

Examples:
  # Statistics for every server type in site 1
  metalcloud-cli server-type statistics --site-id 1

  # Statistics for two server types only
  metalcloud-cli server-type statistics 12 13 --site-id 1

  # Statistics for the servers owned by user 5, at most 10 servers per type
  metalcloud-cli server-type statistics --site-id 1 --user-id 5 --max-results-per-server-type 10


```
metalcloud-cli server-type statistics [server_type_id...] [flags]
```

### Options

```
  -h, --help                              help for statistics
      --instance-array-id int             Treat only the active instances of this instance array as available.
      --max-results-per-server-type int   Maximum number of servers returned per server type.
      --site-id int                       The ID of the site to report on.
      --user-id int                       Only count the resources owned by this user ID.
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

* [metalcloud-cli server-type](metalcloud-cli_server-type.md)	 - Manage server types and hardware configurations

