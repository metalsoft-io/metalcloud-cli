## metalcloud-cli server-instance

Manage individual server instances

### Synopsis

Server Instance management commands.

Server Instances are individual compute resources within server instance groups.
They represent physical or virtual servers with specific hardware configurations
and network connections. Each instance inherits properties from its parent
instance group but can have individual characteristics and status.

Command categories:
  Lifecycle:      list, get, create, delete, config-example
  Configuration:  config, update-config, update-meta, variables, os-installation-data
  Power:          power, power-status-batch, power-set-batch
  Operations:     reset, reinstall-os, credentials
  Hardware:       drives, interfaces, interface, update-interface-config
  Reporting:      statistics

Use "metalcloud-cli server-instance [command] --help" for detailed information about each command.

### Options

```
  -h, --help   help for server-instance
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
* [metalcloud-cli server-instance config](metalcloud-cli_server-instance_config.md)	 - Get configuration of a server instance
* [metalcloud-cli server-instance config-example](metalcloud-cli_server-instance_config-example.md)	 - Print an example server instance create configuration
* [metalcloud-cli server-instance create](metalcloud-cli_server-instance_create.md)	 - Create a server instance in an infrastructure
* [metalcloud-cli server-instance credentials](metalcloud-cli_server-instance_credentials.md)	 - Get login credentials for a server instance
* [metalcloud-cli server-instance delete](metalcloud-cli_server-instance_delete.md)	 - Delete a server instance
* [metalcloud-cli server-instance drives](metalcloud-cli_server-instance_drives.md)	 - List the drives of a server instance
* [metalcloud-cli server-instance get](metalcloud-cli_server-instance_get.md)	 - Get detailed information about a server instance
* [metalcloud-cli server-instance interface](metalcloud-cli_server-instance_interface.md)	 - Get one interface of a server instance
* [metalcloud-cli server-instance interfaces](metalcloud-cli_server-instance_interfaces.md)	 - List the interfaces of a server instance
* [metalcloud-cli server-instance list](metalcloud-cli_server-instance_list.md)	 - List server instances
* [metalcloud-cli server-instance os-installation-data](metalcloud-cli_server-instance_os-installation-data.md)	 - Get the OS installation data of a server instance
* [metalcloud-cli server-instance power](metalcloud-cli_server-instance_power.md)	 - Set power state of a server instance
* [metalcloud-cli server-instance power-set-batch](metalcloud-cli_server-instance_power-set-batch.md)	 - Set the power state of several server instances at once
* [metalcloud-cli server-instance power-status-batch](metalcloud-cli_server-instance_power-status-batch.md)	 - Get the power status of several server instances at once
* [metalcloud-cli server-instance reinstall-os](metalcloud-cli_server-instance_reinstall-os.md)	 - Schedule OS reinstall for a server instance
* [metalcloud-cli server-instance reset](metalcloud-cli_server-instance_reset.md)	 - Reset the deployed server of a server instance
* [metalcloud-cli server-instance statistics](metalcloud-cli_server-instance_statistics.md)	 - Get server instance statistics
* [metalcloud-cli server-instance update-config](metalcloud-cli_server-instance_update-config.md)	 - Update the configuration of a server instance
* [metalcloud-cli server-instance update-interface-config](metalcloud-cli_server-instance_update-interface-config.md)	 - Update the configuration of a server instance interface
* [metalcloud-cli server-instance update-meta](metalcloud-cli_server-instance_update-meta.md)	 - Update the metadata of a server instance
* [metalcloud-cli server-instance variables](metalcloud-cli_server-instance_variables.md)	 - Get the variables of a server instance

