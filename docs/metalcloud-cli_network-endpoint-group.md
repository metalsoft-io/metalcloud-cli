## metalcloud-cli network-endpoint-group

Network endpoint group management

### Synopsis

Manage network endpoint groups: named sets of logical network connections that
server instance groups and VM instance groups attach to.

Command categories:
  Lifecycle:        list, get, create, update, delete, config-example
  Logical networks: logical-network list|get|add|update|remove

### Options

```
  -h, --help   help for network-endpoint-group
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
* [metalcloud-cli network-endpoint-group config-example](metalcloud-cli_network-endpoint-group_config-example.md)	 - Print an example network endpoint group create configuration
* [metalcloud-cli network-endpoint-group create](metalcloud-cli_network-endpoint-group_create.md)	 - Create a network endpoint group
* [metalcloud-cli network-endpoint-group delete](metalcloud-cli_network-endpoint-group_delete.md)	 - Delete a network endpoint group
* [metalcloud-cli network-endpoint-group get](metalcloud-cli_network-endpoint-group_get.md)	 - Get network endpoint group details
* [metalcloud-cli network-endpoint-group list](metalcloud-cli_network-endpoint-group_list.md)	 - List network endpoint groups
* [metalcloud-cli network-endpoint-group logical-network](metalcloud-cli_network-endpoint-group_logical-network.md)	 - Network endpoint group logical network management
* [metalcloud-cli network-endpoint-group update](metalcloud-cli_network-endpoint-group_update.md)	 - Update a network endpoint group

