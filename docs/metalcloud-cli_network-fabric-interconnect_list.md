## metalcloud-cli network-fabric-interconnect list

List network fabric interconnects

### Synopsis

List all network fabric interconnects.

Optional Flags:
  --filter-status strings   Filter by status (e.g. draft, active). Repeatable or comma-separated.

Examples:
  metalcloud network-fabric-interconnect list
  metalcloud nfi list --filter-status active

```
metalcloud-cli network-fabric-interconnect list [flags]
```

### Options

```
      --filter-status strings   Filter by status (e.g. draft, active).
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

