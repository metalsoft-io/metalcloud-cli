## metalcloud-cli vm-type

Manage VM types and configurations

### Synopsis

Manage VM types and their configurations in the MetalCloud platform.

VM types define the resource specifications (CPU cores, RAM) for virtual machines. 
They can be experimental or production-ready, and can be restricted to unmanaged VMs only.

Available commands:
  list          List all VM types with pagination support
  get           Get detailed information about a specific VM type
  create        Create a new VM type from configuration
  update        Update an existing VM type configuration
  delete        Delete a VM type
  vms           List all VMs using a specific VM type
  config-example Show an example configuration for creating VM types

Use "metalcloud vm-type [command] --help" for more information about a command.

### Options

```
  -h, --help   help for vm-type
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
* [metalcloud-cli vm-type config-example](metalcloud-cli_vm-type_config-example.md)	 - Show an example configuration for creating VM types
* [metalcloud-cli vm-type create](metalcloud-cli_vm-type_create.md)	 - Create a new VM type from configuration
* [metalcloud-cli vm-type delete](metalcloud-cli_vm-type_delete.md)	 - Delete a VM type
* [metalcloud-cli vm-type get](metalcloud-cli_vm-type_get.md)	 - Get detailed information about a specific VM type
* [metalcloud-cli vm-type list](metalcloud-cli_vm-type_list.md)	 - List all VM types with optional pagination
* [metalcloud-cli vm-type update](metalcloud-cli_vm-type_update.md)	 - Update an existing VM type configuration
* [metalcloud-cli vm-type vms](metalcloud-cli_vm-type_vms.md)	 - List all VMs using a specific VM type

