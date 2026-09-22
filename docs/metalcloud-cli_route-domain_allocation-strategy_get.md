## metalcloud-cli route-domain allocation-strategy get

Get one allocation strategy

### Synopsis

Get the details of one allocation strategy of a route domain.

Required Arguments:
  route_domain_id The ID of the route domain
  family       One of: l3-vlan, l3-vni, vrf
  strategy_id  The ID of the allocation strategy

Examples:
  metalcloud route-domain allocation-strategy get 12 l3-vlan 3

```
metalcloud-cli route-domain allocation-strategy get route_domain_id family strategy_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli route-domain allocation-strategy](metalcloud-cli_route-domain_allocation-strategy.md)	 - Manage route domain allocation strategies

