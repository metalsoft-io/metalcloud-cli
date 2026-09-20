## metalcloud-cli point-to-point-link allocation-strategy add

Add an allocation strategy

### Synopsis

Add an allocation strategy to a point-to-point link.

Required Arguments:
  link_id      The ID of the point-to-point link
  family       One of: ipv4, ipv6

Required Flags:
  --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud point-to-point-link allocation-strategy add 12 ipv4 --config-source strategy.json
  echo '{"kind":"auto","scope":{"kind":"global"}}' | metalcloud point-to-point-link allocation-strategy add 12 ipv4 --config-source pipe

```
metalcloud-cli point-to-point-link allocation-strategy add link_id family [flags]
```

### Options

```
      --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for add
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

