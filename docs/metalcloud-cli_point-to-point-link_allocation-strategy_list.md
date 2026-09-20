## metalcloud-cli point-to-point-link allocation-strategy list

List allocation strategies of one family

### Synopsis

List the allocation strategies of one family on a point-to-point link.

Required Arguments:
  link_id      The ID of the point-to-point link
  family       One of: ipv4, ipv6

Examples:
  metalcloud point-to-point-link allocation-strategy list 12 ipv4

```
metalcloud-cli point-to-point-link allocation-strategy list link_id family [flags]
```

### Options

```
  -h, --help   help for list
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

