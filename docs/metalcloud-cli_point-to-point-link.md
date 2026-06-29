## metalcloud-cli point-to-point-link

Manage point-to-point links between network interfaces

### Synopsis

Manage point-to-point links between network device (and server) interfaces.

A point-to-point link connects two interfaces (or a single interface, for a
half-connected link) and can carry IPv4/IPv6 subnet allocation strategies that
assign the link's addresses. Links can be created fully staged (interfaces plus
a manual /31 strategy) in one call via the create command's config source.

### Options

```
  -h, --help   help for point-to-point-link
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
* [metalcloud-cli point-to-point-link add-ipv4-strategy](metalcloud-cli_point-to-point-link_add-ipv4-strategy.md)	 - Attach a manual IPv4 subnet allocation strategy to a link
* [metalcloud-cli point-to-point-link config-example](metalcloud-cli_point-to-point-link_config-example.md)	 - Display a point-to-point link configuration example
* [metalcloud-cli point-to-point-link create](metalcloud-cli_point-to-point-link_create.md)	 - Create a point-to-point link
* [metalcloud-cli point-to-point-link delete](metalcloud-cli_point-to-point-link_delete.md)	 - Delete a point-to-point link
* [metalcloud-cli point-to-point-link get](metalcloud-cli_point-to-point-link_get.md)	 - Get details about a specific point-to-point link
* [metalcloud-cli point-to-point-link list](metalcloud-cli_point-to-point-link_list.md)	 - List point-to-point links

