## metalcloud-cli server-instance-group network

Manage network connections for server instance groups

### Synopsis

Manage network connections for server instance groups.

This command group provides operations for managing network connections between
server instance groups and networks. You can list, view, create, update, and
delete network connections.

Available commands:
- list: List all network connections for a server instance group
- get: Get details of a specific network connection
- connect: Connect a server instance group to a network
- update: Update an existing network connection
- disconnect: Remove a network connection

Use "metalcloud-cli server-instance-group network [command] --help" for detailed information about each command.

### Options

```
  -h, --help   help for network
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

* [metalcloud-cli server-instance-group](metalcloud-cli_server-instance-group.md)	 - Manage server instance groups within infrastructures
* [metalcloud-cli server-instance-group network connect](metalcloud-cli_server-instance-group_network_connect.md)	 - Connect a server instance group to a network
* [metalcloud-cli server-instance-group network disconnect](metalcloud-cli_server-instance-group_network_disconnect.md)	 - Remove a network connection from a server instance group
* [metalcloud-cli server-instance-group network get](metalcloud-cli_server-instance-group_network_get.md)	 - Get network connection details for a server instance group
* [metalcloud-cli server-instance-group network list](metalcloud-cli_server-instance-group_network_list.md)	 - List all network connections for a server instance group
* [metalcloud-cli server-instance-group network update](metalcloud-cli_server-instance-group_network_update.md)	 - Update network connection for a server instance group

