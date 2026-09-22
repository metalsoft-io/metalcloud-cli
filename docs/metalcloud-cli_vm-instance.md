## metalcloud-cli vm-instance

Manage individual VM instances within infrastructures

### Synopsis

Manage individual VM instances within infrastructures.

VM instances are individual virtual machines that can be created, managed,
and controlled independently. This includes operations like getting instance
details, listing instances, managing power states, and accessing configuration.

Available commands:
  list                  List all VM instances in an infrastructure
  get                   Get details of a specific VM instance
  config-example        Print an example VM instance configuration
  create                Create a new VM instance
  delete                Delete a VM instance
  config                Show the pending configuration of a VM instance
  update-config         Update the pending configuration of a VM instance
  update-meta           Update the metadata (tags) of a VM instance
  credentials           Show the credentials of a VM instance
  variables             Show the variables of a VM instance
  os-installation-data  Show the OS installation data of a VM instance
  power-status          Get VM instance power status
  start                 Start a VM instance
  shutdown              Shutdown a VM instance
  reboot                Reboot a VM instance
  apply-type            Apply a VM type on a VM instance

Examples:
  metalcloud-cli vm-instance list my-infra
  metalcloud-cli vmi get 12345 67890
  metalcloud-cli vm start 12345 67890

### Options

```
  -h, --help   help for vm-instance
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
* [metalcloud-cli vm-instance apply-type](metalcloud-cli_vm-instance_apply-type.md)	 - Apply a VM type on a VM instance
* [metalcloud-cli vm-instance config](metalcloud-cli_vm-instance_config.md)	 - Show the pending configuration of a VM instance
* [metalcloud-cli vm-instance config-example](metalcloud-cli_vm-instance_config-example.md)	 - Print a VM instance configuration example
* [metalcloud-cli vm-instance create](metalcloud-cli_vm-instance_create.md)	 - Create a new VM instance
* [metalcloud-cli vm-instance credentials](metalcloud-cli_vm-instance_credentials.md)	 - Get login credentials for a VM instance
* [metalcloud-cli vm-instance delete](metalcloud-cli_vm-instance_delete.md)	 - Delete a VM instance
* [metalcloud-cli vm-instance get](metalcloud-cli_vm-instance_get.md)	 - Get detailed information about a specific VM instance
* [metalcloud-cli vm-instance list](metalcloud-cli_vm-instance_list.md)	 - List all VM instances in an infrastructure
* [metalcloud-cli vm-instance os-installation-data](metalcloud-cli_vm-instance_os-installation-data.md)	 - Show the OS installation data of a VM instance
* [metalcloud-cli vm-instance power-status](metalcloud-cli_vm-instance_power-status.md)	 - Get VM instance power status
* [metalcloud-cli vm-instance reboot](metalcloud-cli_vm-instance_reboot.md)	 - Reboot a VM instance
* [metalcloud-cli vm-instance shutdown](metalcloud-cli_vm-instance_shutdown.md)	 - Shutdown a VM instance
* [metalcloud-cli vm-instance start](metalcloud-cli_vm-instance_start.md)	 - Start a VM instance
* [metalcloud-cli vm-instance update-config](metalcloud-cli_vm-instance_update-config.md)	 - Update the pending configuration of a VM instance
* [metalcloud-cli vm-instance update-meta](metalcloud-cli_vm-instance_update-meta.md)	 - Update the metadata of a VM instance
* [metalcloud-cli vm-instance variables](metalcloud-cli_vm-instance_variables.md)	 - Show the variables of a VM instance

