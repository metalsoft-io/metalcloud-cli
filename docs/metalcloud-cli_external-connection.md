## metalcloud-cli external-connection

External connection management

### Synopsis

Manage external connections: the network device interfaces of a fabric that
carry traffic towards networks outside the MetalSoft managed infrastructure.

Command categories:
  Lifecycle:        list, get, create, update, delete, config-example
  Interfaces:       interface list|get|add|update|remove
  Logical networks: logical-network list|get|add|remove
  Discovery:        get-network-device-interfaces

### Options

```
  -h, --help   help for external-connection
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
* [metalcloud-cli external-connection config-example](metalcloud-cli_external-connection_config-example.md)	 - Print an example external connection create configuration
* [metalcloud-cli external-connection create](metalcloud-cli_external-connection_create.md)	 - Create an external connection
* [metalcloud-cli external-connection delete](metalcloud-cli_external-connection_delete.md)	 - Delete an external connection
* [metalcloud-cli external-connection get](metalcloud-cli_external-connection_get.md)	 - Get external connection details
* [metalcloud-cli external-connection get-network-device-interfaces](metalcloud-cli_external-connection_get-network-device-interfaces.md)	 - List a network device's interfaces and their external connections
* [metalcloud-cli external-connection interface](metalcloud-cli_external-connection_interface.md)	 - External connection interface management
* [metalcloud-cli external-connection list](metalcloud-cli_external-connection_list.md)	 - List external connections
* [metalcloud-cli external-connection logical-network](metalcloud-cli_external-connection_logical-network.md)	 - External connection logical network management
* [metalcloud-cli external-connection update](metalcloud-cli_external-connection_update.md)	 - Update an external connection

