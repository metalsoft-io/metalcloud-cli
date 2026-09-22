## metalcloud-cli point-to-point-link allocation-strategy config-example

Print example allocation strategy configurations

### Synopsis

Print one example create body per strategy kind (auto, manual, ...) for a family.
Edit one of them and pass it to 'add' or 'replace' via --config-source.

Required Arguments:
  family       One of: ipv4, ipv6

Examples:
  metalcloud point-to-point-link allocation-strategy config-example ipv4 -f yaml

```
metalcloud-cli point-to-point-link allocation-strategy config-example family [flags]
```

### Options

```
  -h, --help   help for config-example
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

