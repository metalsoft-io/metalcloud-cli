## metalcloud-cli container-instance

Manage container instances within infrastructures

### Synopsis

Manage container instances within infrastructures.

A container instance is the desired-state description of a single container in
an infrastructure. It belongs to a container instance group and is materialised
into a container when the infrastructure is deployed.

Available commands:
  list                  List the container instances of an infrastructure
  get                   Get details of a container instance
  config-example        Print an example container instance configuration
  create                Create a new container instance
  delete                Delete a container instance
  config                Show the pending configuration of a container instance
  update-config         Update the pending configuration of a container instance
  update-meta           Update the metadata (tags) of a container instance
  credentials           Show the credentials of a container instance
  variables             Show the variables of a container instance
  os-installation-data  Show the OS installation data of a container instance
  power-status          Show the power state of a container instance
  start                 Power on a container instance
  shutdown              Power off a container instance
  reboot                Reboot a container instance
  apply-type            Apply a container type on a container instance

Examples:
  metalcloud-cli container-instance list my-infra
  metalcloud-cli ci get my-infra 5678
  metalcloud-cli ci power-status my-infra 5678

### Options

```
  -h, --help   help for container-instance
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
* [metalcloud-cli container-instance apply-type](metalcloud-cli_container-instance_apply-type.md)	 - Apply a container type on a container instance
* [metalcloud-cli container-instance config](metalcloud-cli_container-instance_config.md)	 - Show the pending configuration of a container instance
* [metalcloud-cli container-instance config-example](metalcloud-cli_container-instance_config-example.md)	 - Print a container instance configuration example
* [metalcloud-cli container-instance create](metalcloud-cli_container-instance_create.md)	 - Create a new container instance
* [metalcloud-cli container-instance credentials](metalcloud-cli_container-instance_credentials.md)	 - Show the credentials of a container instance
* [metalcloud-cli container-instance delete](metalcloud-cli_container-instance_delete.md)	 - Delete a container instance
* [metalcloud-cli container-instance get](metalcloud-cli_container-instance_get.md)	 - Get container instance details
* [metalcloud-cli container-instance list](metalcloud-cli_container-instance_list.md)	 - List the container instances of an infrastructure
* [metalcloud-cli container-instance os-installation-data](metalcloud-cli_container-instance_os-installation-data.md)	 - Show the OS installation data of a container instance
* [metalcloud-cli container-instance power-status](metalcloud-cli_container-instance_power-status.md)	 - Get the power status of a container instance
* [metalcloud-cli container-instance reboot](metalcloud-cli_container-instance_reboot.md)	 - Reboot a container instance
* [metalcloud-cli container-instance shutdown](metalcloud-cli_container-instance_shutdown.md)	 - Power off a container instance
* [metalcloud-cli container-instance start](metalcloud-cli_container-instance_start.md)	 - Power on a container instance
* [metalcloud-cli container-instance update-config](metalcloud-cli_container-instance_update-config.md)	 - Update the pending configuration of a container instance
* [metalcloud-cli container-instance update-meta](metalcloud-cli_container-instance_update-meta.md)	 - Update the metadata of a container instance
* [metalcloud-cli container-instance variables](metalcloud-cli_container-instance_variables.md)	 - Show the variables of a container instance

