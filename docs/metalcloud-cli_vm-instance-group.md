## metalcloud-cli vm-instance-group

Manage VM instance groups within infrastructures

### Synopsis

Manage VM instance groups within infrastructures.

VM instance groups allow you to create and manage collections of virtual machines
with similar configurations. This provides easier scaling and management of
multiple VM instances that serve the same purpose.

Available commands:
  list         List all VM instance groups in an infrastructure
  get          Get details of a specific VM instance group
  create       Create a new VM instance group
  update       Update VM instance group configuration
  delete       Delete a VM instance group
  instances    List instances within a VM instance group
  config       Show the pending configuration of a VM instance group
  update-meta  Update the metadata (tags) of a VM instance group
  apply-type   Apply a VM type on a VM instance group
  interfaces   List the network interfaces of a VM instance group
  interface    Get one network interface of a VM instance group
  network      Manage the network connections of a VM instance group

Examples:
  metalcloud-cli vm-instance-group list my-infra
  metalcloud-cli vmg get 12345 67890
  metalcloud-cli vm-group create 12345 5 100 3 7

### Options

```
  -h, --help   help for vm-instance-group
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
* [metalcloud-cli vm-instance-group apply-type](metalcloud-cli_vm-instance-group_apply-type.md)	 - Apply a VM type on a VM instance group
* [metalcloud-cli vm-instance-group config](metalcloud-cli_vm-instance-group_config.md)	 - Show the pending configuration of a VM instance group
* [metalcloud-cli vm-instance-group create](metalcloud-cli_vm-instance-group_create.md)	 - Create a new VM instance group in an infrastructure
* [metalcloud-cli vm-instance-group delete](metalcloud-cli_vm-instance-group_delete.md)	 - Delete a VM instance group
* [metalcloud-cli vm-instance-group get](metalcloud-cli_vm-instance-group_get.md)	 - Get details of a specific VM instance group
* [metalcloud-cli vm-instance-group instances](metalcloud-cli_vm-instance-group_instances.md)	 - List VM instances within a VM instance group
* [metalcloud-cli vm-instance-group interface](metalcloud-cli_vm-instance-group_interface.md)	 - Get one network interface of a VM instance group
* [metalcloud-cli vm-instance-group interfaces](metalcloud-cli_vm-instance-group_interfaces.md)	 - List the network interfaces of a VM instance group
* [metalcloud-cli vm-instance-group list](metalcloud-cli_vm-instance-group_list.md)	 - List all VM instance groups in an infrastructure
* [metalcloud-cli vm-instance-group network](metalcloud-cli_vm-instance-group_network.md)	 - Manage the network connections of a VM instance group
* [metalcloud-cli vm-instance-group update](metalcloud-cli_vm-instance-group_update.md)	 - Update VM instance group configuration
* [metalcloud-cli vm-instance-group update-meta](metalcloud-cli_vm-instance-group_update-meta.md)	 - Update the metadata of a VM instance group

