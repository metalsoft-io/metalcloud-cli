## metalcloud-cli drive list

List all drives within an infrastructure

### Synopsis

List all drives within an infrastructure with optional status filtering.

This command displays a comprehensive list of all drives associated with the specified infrastructure,
including their IDs, configurations, current status, and metadata.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure

Flags:
  --filter-status strings       Filter drives by status (optional)
                               Multiple statuses can be provided as comma-separated values

Examples:
  # List all drives in an infrastructure
  metalcloud-cli drive list my-infrastructure

  # List drives with specific status
  metalcloud-cli drive list my-infrastructure --filter-status active

  # List drives with multiple statuses
  metalcloud-cli drive list my-infrastructure --filter-status active,pending

```
metalcloud-cli drive list infrastructure_id_or_label [flags]
```

### Options

```
      --filter-status strings   Filter the result by drive status.
  -h, --help                    help for list
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

* [metalcloud-cli drive](metalcloud-cli_drive.md)	 - Manage drives within infrastructures

