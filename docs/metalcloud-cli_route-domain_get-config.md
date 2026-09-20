## metalcloud-cli route-domain get-config

Get the config of a route domain

### Synopsis

Display the config object of a route domain.

The config is a separate sub-resource holding the desired state of the route
domain: its kind, the auto route distinguisher / route target settings, and the
allocation strategies that the 'route-domain allocation-strategy' commands
manage. It carries its own revision, distinct from the route domain's.

Required Arguments:
  route_domain_id  The ID of the route domain

Examples:
  metalcloud-cli route-domain get-config 2
  metalcloud-cli route-domain get-config 2 -f json

```
metalcloud-cli route-domain get-config route_domain_id [flags]
```

### Options

```
  -h, --help   help for get-config
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

