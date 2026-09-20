## metalcloud-cli endpoint-instance

Endpoint instance management

### Synopsis

Manage endpoint instances.

An endpoint instance is the deployment of an endpoint inside an infrastructure. It is
created within an infrastructure but is addressed by its own ID afterwards.

Command categories:
  Lifecycle:      list, get, create, delete
  Configuration:  config, update-config, update-meta, config-example

Use "metalcloud-cli endpoint-instance [command] --help" for detailed information about each command.

### Options

```
  -h, --help   help for endpoint-instance
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
* [metalcloud-cli endpoint-instance config](metalcloud-cli_endpoint-instance_config.md)	 - Get the pending configuration of an endpoint instance
* [metalcloud-cli endpoint-instance config-example](metalcloud-cli_endpoint-instance_config-example.md)	 - Print an example endpoint instance create configuration
* [metalcloud-cli endpoint-instance create](metalcloud-cli_endpoint-instance_create.md)	 - Create an endpoint instance in an infrastructure
* [metalcloud-cli endpoint-instance delete](metalcloud-cli_endpoint-instance_delete.md)	 - Delete an endpoint instance
* [metalcloud-cli endpoint-instance get](metalcloud-cli_endpoint-instance_get.md)	 - Get endpoint instance details
* [metalcloud-cli endpoint-instance list](metalcloud-cli_endpoint-instance_list.md)	 - List endpoint instances
* [metalcloud-cli endpoint-instance update-config](metalcloud-cli_endpoint-instance_update-config.md)	 - Update the configuration of an endpoint instance
* [metalcloud-cli endpoint-instance update-meta](metalcloud-cli_endpoint-instance_update-meta.md)	 - Update the metadata of an endpoint instance

