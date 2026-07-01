## metalcloud-cli route-domain

Manage route domains (tenant VRFs)

### Synopsis

Manage route domains in the MetalCloud infrastructure.

A route domain is a tenant VRF: an EVPN-L3VPN / VRF-Lite routing instance that L3
logical networks attach to (via a logical network profile's routeDomainId). Use
these commands to list, create, update, and delete route domains.

### Options

```
  -h, --help   help for route-domain
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli route-domain config-example](metalcloud-cli_route-domain_config-example.md)	 - Display a route domain configuration example
* [metalcloud-cli route-domain create](metalcloud-cli_route-domain_create.md)	 - Create a new route domain
* [metalcloud-cli route-domain delete](metalcloud-cli_route-domain_delete.md)	 - Delete a route domain
* [metalcloud-cli route-domain get](metalcloud-cli_route-domain_get.md)	 - Get details about a specific route domain
* [metalcloud-cli route-domain list](metalcloud-cli_route-domain_list.md)	 - List all route domains
* [metalcloud-cli route-domain update](metalcloud-cli_route-domain_update.md)	 - Update an existing route domain

