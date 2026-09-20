## metalcloud-cli logical-network-profile allocation-strategy

Manage logical network profile allocation strategies

### Synopsis

Manage the allocation strategies of a logical network profile.

Each strategy family is a separate collection. Supported families for this resource:
  ipv4, ipv6, vlan, vni, pkey, zone

Commands:
  list, get, config-example, add, replace, remove

### Options

```
  -h, --help   help for allocation-strategy
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

* [metalcloud-cli logical-network-profile](metalcloud-cli_logical-network-profile.md)	 - Manage logical network profiles for network configuration templates
* [metalcloud-cli logical-network-profile allocation-strategy add](metalcloud-cli_logical-network-profile_allocation-strategy_add.md)	 - Add an allocation strategy
* [metalcloud-cli logical-network-profile allocation-strategy config-example](metalcloud-cli_logical-network-profile_allocation-strategy_config-example.md)	 - Print example allocation strategy configurations
* [metalcloud-cli logical-network-profile allocation-strategy get](metalcloud-cli_logical-network-profile_allocation-strategy_get.md)	 - Get one allocation strategy
* [metalcloud-cli logical-network-profile allocation-strategy list](metalcloud-cli_logical-network-profile_allocation-strategy_list.md)	 - List allocation strategies of one family
* [metalcloud-cli logical-network-profile allocation-strategy remove](metalcloud-cli_logical-network-profile_allocation-strategy_remove.md)	 - Remove an allocation strategy
* [metalcloud-cli logical-network-profile allocation-strategy replace](metalcloud-cli_logical-network-profile_allocation-strategy_replace.md)	 - Replace an allocation strategy

