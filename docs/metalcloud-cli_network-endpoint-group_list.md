## metalcloud-cli network-endpoint-group list

List network endpoint groups

### Synopsis

List all network endpoint groups.

Optional Flags:
  --filter-id strings        Filter by network endpoint group ID. Repeatable or comma-separated.
  --filter-name strings      Filter by name. Repeatable or comma-separated.
  --filter-site-id strings   Filter by site ID. Repeatable or comma-separated.

Examples:
  metalcloud network-endpoint-group list
  metalcloud neg list --filter-site-id 1

```
metalcloud-cli network-endpoint-group list [flags]
```

### Options

```
      --filter-id strings        Filter by network endpoint group ID.
      --filter-name strings      Filter by name.
      --filter-site-id strings   Filter by site ID.
  -h, --help                     help for list
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

* [metalcloud-cli network-endpoint-group](metalcloud-cli_network-endpoint-group.md)	 - Network endpoint group management

