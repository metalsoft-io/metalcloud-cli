## metalcloud-cli vm-pool

Manage virtual machine pools and their resources

### Synopsis

Manage virtual machine pools including VMware vSphere, Hyper-V, and other hypervisor environments.

VM pools provide centralized management of virtualization infrastructure, allowing you to:
- Create and configure connections to hypervisor management systems
- Monitor VM and cluster host resources
- Manage credentials and certificates for secure access
- Control maintenance and experimental modes

Available commands support full lifecycle management from initial configuration
to ongoing monitoring and resource inspection.

### Options

```
  -h, --help   help for vm-pool
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
* [metalcloud-cli vm-pool cluster-host-interfaces](metalcloud-cli_vm-pool_cluster-host-interfaces.md)	 - List network interfaces for a cluster host in a VM pool
* [metalcloud-cli vm-pool cluster-host-vms](metalcloud-cli_vm-pool_cluster-host-vms.md)	 - List virtual machines on a specific cluster host with pagination
* [metalcloud-cli vm-pool cluster-hosts](metalcloud-cli_vm-pool_cluster-hosts.md)	 - List cluster hosts in a VM pool with pagination
* [metalcloud-cli vm-pool config-example](metalcloud-cli_vm-pool_config-example.md)	 - Display a complete VM pool configuration example
* [metalcloud-cli vm-pool create](metalcloud-cli_vm-pool_create.md)	 - Create a new VM pool from configuration file or pipe
* [metalcloud-cli vm-pool credentials](metalcloud-cli_vm-pool_credentials.md)	 - Retrieve authentication credentials for a VM pool
* [metalcloud-cli vm-pool delete](metalcloud-cli_vm-pool_delete.md)	 - Delete a VM pool permanently
* [metalcloud-cli vm-pool get](metalcloud-cli_vm-pool_get.md)	 - Get detailed information about a specific VM pool
* [metalcloud-cli vm-pool list](metalcloud-cli_vm-pool_list.md)	 - List all VM pools with optional filtering
* [metalcloud-cli vm-pool vms](metalcloud-cli_vm-pool_vms.md)	 - List virtual machines in a VM pool with pagination

