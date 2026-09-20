## metalcloud-cli external-connection list

List external connections

### Synopsis

List all external connections.

Optional Flags:
  --filter-id strings          Filter by external connection ID. Repeatable or comma-separated.
  --filter-fabric-id strings   Filter by fabric ID. Repeatable or comma-separated.
  --filter-label strings       Filter by label. Repeatable or comma-separated.
  --filter-name strings        Filter by name. Repeatable or comma-separated.

Examples:
  metalcloud external-connection list
  metalcloud ext-conn list --filter-fabric-id 7

```
metalcloud-cli external-connection list [flags]
```

### Options

```
      --filter-fabric-id strings   Filter by fabric ID.
      --filter-id strings          Filter by external connection ID.
      --filter-label strings       Filter by label.
      --filter-name strings        Filter by name.
  -h, --help                       help for list
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

* [metalcloud-cli external-connection](metalcloud-cli_external-connection.md)	 - External connection management

