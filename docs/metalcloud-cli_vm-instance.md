## metalcloud-cli vm-instance

Manage individual VM instances within infrastructures

### Synopsis

Manage individual VM instances within infrastructures.

VM instances are individual virtual machines that can be created, managed,
and controlled independently. This includes operations like getting instance
details, listing instances, managing power states, and accessing configuration.

Available commands:
  get          Get details of a specific VM instance
  list         List all VM instances in an infrastructure
  config       Get VM instance configuration
  start        Start a VM instance
  shutdown     Shutdown a VM instance
  reboot       Reboot a VM instance
  power-status Get VM instance power status

Examples:
  metalcloud-cli vm-instance list 12345
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
* [metalcloud-cli vm-instance config](metalcloud-cli_vm-instance_config.md)	 - Get VM instance configuration
* [metalcloud-cli vm-instance get](metalcloud-cli_vm-instance_get.md)	 - Get detailed information about a specific VM instance
* [metalcloud-cli vm-instance list](metalcloud-cli_vm-instance_list.md)	 - List all VM instances in an infrastructure
* [metalcloud-cli vm-instance power-status](metalcloud-cli_vm-instance_power-status.md)	 - Get VM instance power status
* [metalcloud-cli vm-instance reboot](metalcloud-cli_vm-instance_reboot.md)	 - Reboot a VM instance
* [metalcloud-cli vm-instance shutdown](metalcloud-cli_vm-instance_shutdown.md)	 - Shutdown a VM instance
* [metalcloud-cli vm-instance start](metalcloud-cli_vm-instance_start.md)	 - Start a VM instance

