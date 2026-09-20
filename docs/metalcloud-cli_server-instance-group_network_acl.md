## metalcloud-cli server-instance-group network acl

Manage the security rules of a network connection

### Synopsis

Manage the security rules (ACLs) of one network connection of a server instance
group.

Available commands:
- list: List the security rules of a network connection
- get: Get one security rule
- add: Add a security rule
- update: Update a security rule
- remove: Remove a security rule
- config-example: Print an example security rule configuration

Use "metalcloud-cli server-instance-group network acl [command] --help" for detailed information about each command.

### Options

```
  -h, --help   help for acl
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

* [metalcloud-cli server-instance-group network](metalcloud-cli_server-instance-group_network.md)	 - Manage network connections for server instance groups
* [metalcloud-cli server-instance-group network acl add](metalcloud-cli_server-instance-group_network_acl_add.md)	 - Add a security rule to a network connection
* [metalcloud-cli server-instance-group network acl config-example](metalcloud-cli_server-instance-group_network_acl_config-example.md)	 - Print an example security rule configuration
* [metalcloud-cli server-instance-group network acl get](metalcloud-cli_server-instance-group_network_acl_get.md)	 - Get a security rule of a network connection
* [metalcloud-cli server-instance-group network acl list](metalcloud-cli_server-instance-group_network_acl_list.md)	 - List the security rules of a network connection
* [metalcloud-cli server-instance-group network acl remove](metalcloud-cli_server-instance-group_network_acl_remove.md)	 - Remove a security rule from a network connection
* [metalcloud-cli server-instance-group network acl update](metalcloud-cli_server-instance-group_network_acl_update.md)	 - Update a security rule of a network connection

