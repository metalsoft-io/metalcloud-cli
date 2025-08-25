## metalcloud-cli event get

Get detailed information about a specific event

### Synopsis

Retrieve detailed information about a specific event by its ID.

This command displays comprehensive information about a single event, including
its metadata, timestamp, severity, type, associated resources, and full description.

REQUIRED ARGUMENTS:
  event_id                    The unique identifier of the event to retrieve

Examples:
  # Get details of event with ID 12345
  metalcloud event get 12345

  # Get event details using the 'show' alias
  metalcloud event show 67890

The output includes:
- Event ID and timestamp
- Event type and severity level
- Associated infrastructure, server, job, or site information
- Full event description and metadata
- User who triggered the event (if applicable)

```
metalcloud-cli event get event_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli event](metalcloud-cli_event.md)	 - Manage and monitor system events

