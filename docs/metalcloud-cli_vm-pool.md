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

Pool level commands:
  list, get, create, update, delete, config-example, update-config-example,
  credentials, statistics, containers, vms, sync, refresh, import-vms

Cluster host commands (flat, one command per operation):
  cluster-hosts, cluster-host, update-cluster-host, cluster-host-statistics,
  cluster-host-containers, cluster-host-vms, cluster-host-interfaces,
  cluster-host-interface, update-cluster-host-interface

Cluster host interface network device assignments live in their own sub-group,
because they are a full CRUD set of their own:
  cluster-host-interface-network-device (alias: chi-network-device, chind)
    list | get | add | remove | config-example

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
* [metalcloud-cli vm-pool cluster-host](metalcloud-cli_vm-pool_cluster-host.md)	 - Get details of one cluster host of a VM pool
* [metalcloud-cli vm-pool cluster-host-containers](metalcloud-cli_vm-pool_cluster-host-containers.md)	 - List the containers running on one cluster host
* [metalcloud-cli vm-pool cluster-host-interface](metalcloud-cli_vm-pool_cluster-host-interface.md)	 - Get one network interface of a cluster host
* [metalcloud-cli vm-pool cluster-host-interface-network-device](metalcloud-cli_vm-pool_cluster-host-interface-network-device.md)	 - Manage the network device assignments of a cluster host interface
* [metalcloud-cli vm-pool cluster-host-interfaces](metalcloud-cli_vm-pool_cluster-host-interfaces.md)	 - List network interfaces for a cluster host in a VM pool
* [metalcloud-cli vm-pool cluster-host-statistics](metalcloud-cli_vm-pool_cluster-host-statistics.md)	 - Show the resource usage statistics of one cluster host
* [metalcloud-cli vm-pool cluster-host-vms](metalcloud-cli_vm-pool_cluster-host-vms.md)	 - List virtual machines on a specific cluster host with pagination
* [metalcloud-cli vm-pool cluster-hosts](metalcloud-cli_vm-pool_cluster-hosts.md)	 - List cluster hosts in a VM pool with pagination
* [metalcloud-cli vm-pool config-example](metalcloud-cli_vm-pool_config-example.md)	 - Display a complete VM pool configuration example
* [metalcloud-cli vm-pool containers](metalcloud-cli_vm-pool_containers.md)	 - List the containers running in a VM pool
* [metalcloud-cli vm-pool create](metalcloud-cli_vm-pool_create.md)	 - Create a new VM pool from configuration file or pipe
* [metalcloud-cli vm-pool credentials](metalcloud-cli_vm-pool_credentials.md)	 - Retrieve authentication credentials for a VM pool
* [metalcloud-cli vm-pool delete](metalcloud-cli_vm-pool_delete.md)	 - Delete a VM pool permanently
* [metalcloud-cli vm-pool get](metalcloud-cli_vm-pool_get.md)	 - Get detailed information about a specific VM pool
* [metalcloud-cli vm-pool import-vms](metalcloud-cli_vm-pool_import-vms.md)	 - Import VMs from the hypervisor into a VM pool
* [metalcloud-cli vm-pool list](metalcloud-cli_vm-pool_list.md)	 - List all VM pools with optional filtering
* [metalcloud-cli vm-pool refresh](metalcloud-cli_vm-pool_refresh.md)	 - Refresh the information of a VM pool
* [metalcloud-cli vm-pool statistics](metalcloud-cli_vm-pool_statistics.md)	 - Show the resource usage statistics of a VM pool
* [metalcloud-cli vm-pool sync](metalcloud-cli_vm-pool_sync.md)	 - Sync a VM pool with its hypervisor
* [metalcloud-cli vm-pool update](metalcloud-cli_vm-pool_update.md)	 - Update a VM pool from a configuration file or pipe
* [metalcloud-cli vm-pool update-cluster-host](metalcloud-cli_vm-pool_update-cluster-host.md)	 - Update one cluster host of a VM pool
* [metalcloud-cli vm-pool update-cluster-host-config-example](metalcloud-cli_vm-pool_update-cluster-host-config-example.md)	 - Display a cluster host update configuration example
* [metalcloud-cli vm-pool update-cluster-host-interface](metalcloud-cli_vm-pool_update-cluster-host-interface.md)	 - Update one network interface of a cluster host
* [metalcloud-cli vm-pool update-cluster-host-interface-config-example](metalcloud-cli_vm-pool_update-cluster-host-interface-config-example.md)	 - Display a cluster host interface update configuration example
* [metalcloud-cli vm-pool update-config-example](metalcloud-cli_vm-pool_update-config-example.md)	 - Display a VM pool update configuration example
* [metalcloud-cli vm-pool vms](metalcloud-cli_vm-pool_vms.md)	 - List virtual machines in a VM pool with pagination

