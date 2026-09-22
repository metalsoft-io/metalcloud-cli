## metalcloud-cli logical-network-profile allocation-strategy replace

Replace an allocation strategy

### Synopsis

Replace the whole configuration of one allocation strategy of a logical network profile.

Required Arguments:
  logical_network_profile_id The ID of the logical network profile
  family       One of: ipv4, ipv6, vlan, vni, pkey, zone
  strategy_id  The ID of the allocation strategy

Required Flags:
  --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud logical-network-profile allocation-strategy replace 12 ipv4 3 --config-source strategy.json

```
metalcloud-cli logical-network-profile allocation-strategy replace logical_network_profile_id family strategy_id [flags]
```

### Options

```
      --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for replace
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

* [metalcloud-cli logical-network-profile allocation-strategy](metalcloud-cli_logical-network-profile_allocation-strategy.md)	 - Manage logical network profile allocation strategies

