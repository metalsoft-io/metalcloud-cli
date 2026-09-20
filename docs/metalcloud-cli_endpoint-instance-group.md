## metalcloud-cli endpoint-instance-group

Endpoint instance group management

### Synopsis

Manage endpoint instance groups.

An endpoint instance group collects the endpoint instances of an infrastructure that
share the same configuration and network connections. Groups are created inside an
infrastructure but are addressed by their own ID afterwards.

Command categories:
  Lifecycle:      list, get, create, delete
  Configuration:  config, update-config, update-meta, config-example
  Members:        instances
  Networking:     network (connections and their security rules)

Use "metalcloud-cli endpoint-instance-group [command] --help" for detailed information about each command.

### Options

```
  -h, --help   help for endpoint-instance-group
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
* [metalcloud-cli endpoint-instance-group config](metalcloud-cli_endpoint-instance-group_config.md)	 - Get the pending configuration of an endpoint instance group
* [metalcloud-cli endpoint-instance-group config-example](metalcloud-cli_endpoint-instance-group_config-example.md)	 - Print an example endpoint instance group create configuration
* [metalcloud-cli endpoint-instance-group create](metalcloud-cli_endpoint-instance-group_create.md)	 - Create an endpoint instance group in an infrastructure
* [metalcloud-cli endpoint-instance-group delete](metalcloud-cli_endpoint-instance-group_delete.md)	 - Delete an endpoint instance group
* [metalcloud-cli endpoint-instance-group get](metalcloud-cli_endpoint-instance-group_get.md)	 - Get endpoint instance group details
* [metalcloud-cli endpoint-instance-group instances](metalcloud-cli_endpoint-instance-group_instances.md)	 - List the endpoint instances of a group
* [metalcloud-cli endpoint-instance-group list](metalcloud-cli_endpoint-instance-group_list.md)	 - List the endpoint instance groups of an infrastructure
* [metalcloud-cli endpoint-instance-group network](metalcloud-cli_endpoint-instance-group_network.md)	 - Manage the network configuration of an endpoint instance group
* [metalcloud-cli endpoint-instance-group update-config](metalcloud-cli_endpoint-instance-group_update-config.md)	 - Update the configuration of an endpoint instance group
* [metalcloud-cli endpoint-instance-group update-meta](metalcloud-cli_endpoint-instance-group_update-meta.md)	 - Update the metadata of an endpoint instance group

