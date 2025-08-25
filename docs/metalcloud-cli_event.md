## metalcloud-cli event

Manage and monitor system events

### Synopsis

Manage and monitor system events in MetalCloud.

Events represent important system activities such as infrastructure changes,
server deployments, job executions, and user actions. This command provides
tools to list, filter, search, and retrieve detailed information about events.

Available Commands:
  list    List events with filtering and search capabilities
  get     Retrieve detailed information about a specific event

Examples:
  # List all events
  metalcloud event list

  # Get details of a specific event
  metalcloud event get 12345

  # List events with filters
  metalcloud event list --filter-type deployment --filter-severity error

### Options

```
  -h, --help   help for event
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli event get](metalcloud-cli_event_get.md)	 - Get detailed information about a specific event
* [metalcloud-cli event list](metalcloud-cli_event_list.md)	 - List and filter system events

