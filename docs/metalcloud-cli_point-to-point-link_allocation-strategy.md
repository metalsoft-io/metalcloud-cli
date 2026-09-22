## metalcloud-cli point-to-point-link allocation-strategy

Manage point-to-point link allocation strategies

### Synopsis

Manage the allocation strategies of a point-to-point link.

Each strategy family is a separate collection. Supported families for this resource:
  ipv4, ipv6

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

* [metalcloud-cli point-to-point-link](metalcloud-cli_point-to-point-link.md)	 - Manage point-to-point links between network interfaces
* [metalcloud-cli point-to-point-link allocation-strategy add](metalcloud-cli_point-to-point-link_allocation-strategy_add.md)	 - Add an allocation strategy
* [metalcloud-cli point-to-point-link allocation-strategy config-example](metalcloud-cli_point-to-point-link_allocation-strategy_config-example.md)	 - Print example allocation strategy configurations
* [metalcloud-cli point-to-point-link allocation-strategy get](metalcloud-cli_point-to-point-link_allocation-strategy_get.md)	 - Get one allocation strategy
* [metalcloud-cli point-to-point-link allocation-strategy list](metalcloud-cli_point-to-point-link_allocation-strategy_list.md)	 - List allocation strategies of one family
* [metalcloud-cli point-to-point-link allocation-strategy remove](metalcloud-cli_point-to-point-link_allocation-strategy_remove.md)	 - Remove an allocation strategy
* [metalcloud-cli point-to-point-link allocation-strategy replace](metalcloud-cli_point-to-point-link_allocation-strategy_replace.md)	 - Replace an allocation strategy

