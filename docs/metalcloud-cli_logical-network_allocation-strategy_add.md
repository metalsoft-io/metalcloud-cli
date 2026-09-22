## metalcloud-cli logical-network allocation-strategy add

Add an allocation strategy

### Synopsis

Add an allocation strategy to a logical network.

Required Arguments:
  logical_network_id The ID of the logical network
  family       One of: ipv4, ipv6, vlan, vni, pkey, zone

Required Flags:
  --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud logical-network allocation-strategy add 12 ipv4 --config-source strategy.json
  echo '{"kind":"auto","scope":{"kind":"global"}}' | metalcloud logical-network allocation-strategy add 12 ipv4 --config-source pipe

```
metalcloud-cli logical-network allocation-strategy add logical_network_id family [flags]
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

* [metalcloud-cli logical-network allocation-strategy](metalcloud-cli_logical-network_allocation-strategy.md)	 - Manage logical network allocation strategies

