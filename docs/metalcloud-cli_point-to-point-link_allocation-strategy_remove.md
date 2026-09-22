## metalcloud-cli point-to-point-link allocation-strategy remove

Remove an allocation strategy

### Synopsis

Remove one allocation strategy from a point-to-point link.

Required Arguments:
  link_id      The ID of the point-to-point link
  family       One of: ipv4, ipv6
  strategy_id  The ID of the allocation strategy

Examples:
  metalcloud point-to-point-link allocation-strategy remove 12 ipv4 3

```
metalcloud-cli point-to-point-link allocation-strategy remove link_id family strategy_id [flags]
```

### Options

```
  -h, --help   help for remove
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

* [metalcloud-cli point-to-point-link allocation-strategy](metalcloud-cli_point-to-point-link_allocation-strategy.md)	 - Manage point-to-point link allocation strategies

