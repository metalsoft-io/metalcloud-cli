## metalcloud-cli endpoint-instance-group network

Manage the network configuration of an endpoint instance group

### Synopsis

Manage the network configuration of an endpoint instance group.

Command categories:
  Configuration:  list, replace, config-example
  Connections:    connections, get, connect, update, disconnect
  Security:       acl (security rules of a connection)

Use "metalcloud-cli endpoint-instance-group network [command] --help" for detailed information about each command.

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

* [metalcloud-cli endpoint-instance-group](metalcloud-cli_endpoint-instance-group.md)	 - Endpoint instance group management
* [metalcloud-cli endpoint-instance-group network acl](metalcloud-cli_endpoint-instance-group_network_acl.md)	 - Manage the security rules of a network connection
* [metalcloud-cli endpoint-instance-group network config-example](metalcloud-cli_endpoint-instance-group_network_config-example.md)	 - Print an example network connection configuration
* [metalcloud-cli endpoint-instance-group network connect](metalcloud-cli_endpoint-instance-group_network_connect.md)	 - Connect an endpoint instance group to a logical network
* [metalcloud-cli endpoint-instance-group network connections](metalcloud-cli_endpoint-instance-group_network_connections.md)	 - List the network connections of an endpoint instance group
* [metalcloud-cli endpoint-instance-group network disconnect](metalcloud-cli_endpoint-instance-group_network_disconnect.md)	 - Remove a network connection from an endpoint instance group
* [metalcloud-cli endpoint-instance-group network get](metalcloud-cli_endpoint-instance-group_network_get.md)	 - Get a network connection of an endpoint instance group
* [metalcloud-cli endpoint-instance-group network list](metalcloud-cli_endpoint-instance-group_network_list.md)	 - Get the network configuration of an endpoint instance group
* [metalcloud-cli endpoint-instance-group network replace](metalcloud-cli_endpoint-instance-group_network_replace.md)	 - Create or replace the network configuration of an endpoint instance group
* [metalcloud-cli endpoint-instance-group network update](metalcloud-cli_endpoint-instance-group_network_update.md)	 - Update a network connection of an endpoint instance group

