## metalcloud-cli container-instance-group

Manage container instance groups within infrastructures

### Synopsis

Manage container instance groups within infrastructures.

A container instance group is a collection of identically configured container
instances. The group owns the shared sizing, the OS template and the network
connections of its instances.

Available commands:
  list            List the container instance groups of an infrastructure
  get             Get details of a container instance group
  config-example  Print an example container instance group configuration
  create          Create a new container instance group
  delete          Delete a container instance group
  config          Show the pending configuration of a group
  update-config   Update the pending configuration of a group
  update-meta     Update the metadata (tags) of a group
  instances       List the container instances of a group
  interfaces      List the network interfaces of a group
  interface       Get one network interface of a group
  apply-type      Apply a container type on a group
  network         Manage the network connections of a group

Examples:
  metalcloud-cli container-instance-group list my-infra
  metalcloud-cli cig get my-infra 77
  metalcloud-cli cig network list my-infra 77

### Options

```
  -h, --help   help for container-instance-group
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
* [metalcloud-cli container-instance-group apply-type](metalcloud-cli_container-instance-group_apply-type.md)	 - Apply a container type on a container instance group
* [metalcloud-cli container-instance-group config](metalcloud-cli_container-instance-group_config.md)	 - Show the pending configuration of a container instance group
* [metalcloud-cli container-instance-group config-example](metalcloud-cli_container-instance-group_config-example.md)	 - Print a container instance group configuration example
* [metalcloud-cli container-instance-group create](metalcloud-cli_container-instance-group_create.md)	 - Create a new container instance group
* [metalcloud-cli container-instance-group delete](metalcloud-cli_container-instance-group_delete.md)	 - Delete a container instance group
* [metalcloud-cli container-instance-group get](metalcloud-cli_container-instance-group_get.md)	 - Get container instance group details
* [metalcloud-cli container-instance-group instances](metalcloud-cli_container-instance-group_instances.md)	 - List the container instances of a group
* [metalcloud-cli container-instance-group interface](metalcloud-cli_container-instance-group_interface.md)	 - Get one network interface of a group
* [metalcloud-cli container-instance-group interfaces](metalcloud-cli_container-instance-group_interfaces.md)	 - List the network interfaces of a group
* [metalcloud-cli container-instance-group list](metalcloud-cli_container-instance-group_list.md)	 - List the container instance groups of an infrastructure
* [metalcloud-cli container-instance-group network](metalcloud-cli_container-instance-group_network.md)	 - Manage the network connections of a container instance group
* [metalcloud-cli container-instance-group update-config](metalcloud-cli_container-instance-group_update-config.md)	 - Update the pending configuration of a container instance group
* [metalcloud-cli container-instance-group update-meta](metalcloud-cli_container-instance-group_update-meta.md)	 - Update the metadata of a container instance group

