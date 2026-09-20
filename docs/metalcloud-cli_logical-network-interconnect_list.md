## metalcloud-cli logical-network-interconnect list

List logical network interconnects

### Synopsis

List all logical network interconnects.

Optional Flags:
  --filter-id strings                      Filter by ID. Repeatable or comma-separated.
  --filter-label strings                   Filter by label. Repeatable or comma-separated.
  --filter-name strings                    Filter by name. Repeatable or comma-separated.
  --filter-kind strings                    Filter by kind (e.g. dci-evpn). Repeatable or comma-separated.
  --filter-status strings                  Filter by status. Repeatable or comma-separated.
  --filter-fabric-interconnect-id strings  Filter by network fabric interconnect ID. Repeatable or comma-separated.

Examples:
  metalcloud logical-network-interconnect list
  metalcloud lni list --filter-status active

```
metalcloud-cli logical-network-interconnect list [flags]
```

### Options

```
      --filter-fabric-interconnect-id strings   Filter by network fabric interconnect ID.
      --filter-id strings                       Filter by logical network interconnect ID.
      --filter-kind strings                     Filter by kind (e.g. dci-evpn).
      --filter-label strings                    Filter by label.
      --filter-name strings                     Filter by name.
      --filter-status strings                   Filter by status.
  -h, --help                                    help for list
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

* [metalcloud-cli logical-network-interconnect](metalcloud-cli_logical-network-interconnect.md)	 - Logical network interconnect management

