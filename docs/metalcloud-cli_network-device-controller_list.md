## metalcloud-cli network-device-controller list

List network device controllers

### Synopsis

List all network device controllers.

Optional Flags:
  --filter-id strings                  Filter by controller ID. Repeatable or comma-separated.
  --filter-site-id strings             Filter by site ID. Repeatable or comma-separated.
  --filter-datacenter-name strings     Filter by datacenter name. Repeatable or comma-separated.
  --filter-management-address strings  Filter by management address. Repeatable or comma-separated.
  --filter-identifier-string strings   Filter by identifier (hostname). Repeatable or comma-separated.

Examples:
  metalcloud network-device-controller list
  metalcloud ndc list --filter-site-id 1

```
metalcloud-cli network-device-controller list [flags]
```

### Options

```
      --filter-datacenter-name strings      Filter by datacenter name.
      --filter-id strings                   Filter by controller ID.
      --filter-identifier-string strings    Filter by identifier (hostname).
      --filter-management-address strings   Filter by management address.
      --filter-site-id strings              Filter by site ID.
  -h, --help                                help for list
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

* [metalcloud-cli network-device-controller](metalcloud-cli_network-device-controller.md)	 - Network device controller management

