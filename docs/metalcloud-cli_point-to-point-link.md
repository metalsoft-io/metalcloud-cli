## metalcloud-cli point-to-point-link

Manage point-to-point links between network interfaces

### Synopsis

Manage point-to-point links between network device (and server) interfaces.

A point-to-point link connects two interfaces (or a single interface, for a
half-connected link) and can carry IPv4/IPv6 subnet allocation strategies that
assign the link's addresses. Links can be created fully staged (interfaces plus
a manual /31 strategy) in one call via the create command's config source.

Available Commands:
  list                  List point-to-point links
  get                   Get link details
  create                Create a link
  update                Update a link's properties
  delete                Delete a link
  get-config            Get the link's staged configuration
  update-config         Update the link's staged configuration (MTU)
  config-example        Show a create configuration example
  add-ipv4-strategy     Attach a manual IPv4 allocation strategy
  allocation-strategy   Manage the link's allocation strategies
  static-route          Manage the link's staged static routes

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
* [metalcloud-cli point-to-point-link allocation-strategy](metalcloud-cli_point-to-point-link_allocation-strategy.md)	 - Manage point-to-point link allocation strategies
* [metalcloud-cli point-to-point-link config-example](metalcloud-cli_point-to-point-link_config-example.md)	 - Display a point-to-point link configuration example
* [metalcloud-cli point-to-point-link create](metalcloud-cli_point-to-point-link_create.md)	 - Create a point-to-point link
* [metalcloud-cli point-to-point-link delete](metalcloud-cli_point-to-point-link_delete.md)	 - Delete a point-to-point link
* [metalcloud-cli point-to-point-link get](metalcloud-cli_point-to-point-link_get.md)	 - Get details about a specific point-to-point link
* [metalcloud-cli point-to-point-link get-config](metalcloud-cli_point-to-point-link_get-config.md)	 - Get the staged configuration of a point-to-point link
* [metalcloud-cli point-to-point-link list](metalcloud-cli_point-to-point-link_list.md)	 - List point-to-point links
* [metalcloud-cli point-to-point-link static-route](metalcloud-cli_point-to-point-link_static-route.md)	 - Manage the staged static routes of a point-to-point link
* [metalcloud-cli point-to-point-link update](metalcloud-cli_point-to-point-link_update.md)	 - Update a point-to-point link
* [metalcloud-cli point-to-point-link update-config](metalcloud-cli_point-to-point-link_update-config.md)	 - Update the staged configuration of a point-to-point link
* [metalcloud-cli point-to-point-link update-config-example](metalcloud-cli_point-to-point-link_update-config-example.md)	 - Display a point-to-point link config update example
* [metalcloud-cli point-to-point-link update-example](metalcloud-cli_point-to-point-link_update-example.md)	 - Display a point-to-point link update example

