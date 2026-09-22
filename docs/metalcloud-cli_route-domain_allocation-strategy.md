## metalcloud-cli route-domain allocation-strategy

Manage route domain allocation strategies

### Synopsis

Manage the allocation strategies of a route domain.

Each strategy family is a separate collection. Supported families for this resource:
  l3-vlan, l3-vni, vrf

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

* [metalcloud-cli route-domain](metalcloud-cli_route-domain.md)	 - Manage route domains (tenant VRFs)
* [metalcloud-cli route-domain allocation-strategy add](metalcloud-cli_route-domain_allocation-strategy_add.md)	 - Add an allocation strategy
* [metalcloud-cli route-domain allocation-strategy config-example](metalcloud-cli_route-domain_allocation-strategy_config-example.md)	 - Print example allocation strategy configurations
* [metalcloud-cli route-domain allocation-strategy get](metalcloud-cli_route-domain_allocation-strategy_get.md)	 - Get one allocation strategy
* [metalcloud-cli route-domain allocation-strategy list](metalcloud-cli_route-domain_allocation-strategy_list.md)	 - List allocation strategies of one family
* [metalcloud-cli route-domain allocation-strategy remove](metalcloud-cli_route-domain_allocation-strategy_remove.md)	 - Remove an allocation strategy
* [metalcloud-cli route-domain allocation-strategy replace](metalcloud-cli_route-domain_allocation-strategy_replace.md)	 - Replace an allocation strategy

