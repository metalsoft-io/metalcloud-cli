## metalcloud-cli vm-pool cluster-host-interface-network-device

Manage the network device assignments of a cluster host interface

### Synopsis

Manage the links between a VM pool cluster host interface and the interfaces of
the network devices (switches) it is cabled to.

Every command takes the VM pool, cluster host and cluster host interface as
positional arguments; 'get' and 'remove' additionally take the numeric ID of the
assignment itself (not the network device ID).

Available Commands:
  list            List the network device assignments of an interface
  get             Get one network device assignment
  add             Link an interface to a network device interface
  remove          Delete one network device assignment
  config-example  Display an example 'add' payload

### Options

```
  -h, --help   help for cluster-host-interface-network-device
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

* [metalcloud-cli vm-pool](metalcloud-cli_vm-pool.md)	 - Manage virtual machine pools and their resources
* [metalcloud-cli vm-pool cluster-host-interface-network-device add](metalcloud-cli_vm-pool_cluster-host-interface-network-device_add.md)	 - Link a cluster host interface to a network device interface
* [metalcloud-cli vm-pool cluster-host-interface-network-device config-example](metalcloud-cli_vm-pool_cluster-host-interface-network-device_config-example.md)	 - Display a network device assignment configuration example
* [metalcloud-cli vm-pool cluster-host-interface-network-device get](metalcloud-cli_vm-pool_cluster-host-interface-network-device_get.md)	 - Get one network device assignment of a cluster host interface
* [metalcloud-cli vm-pool cluster-host-interface-network-device list](metalcloud-cli_vm-pool_cluster-host-interface-network-device_list.md)	 - List the network device assignments of a cluster host interface
* [metalcloud-cli vm-pool cluster-host-interface-network-device remove](metalcloud-cli_vm-pool_cluster-host-interface-network-device_remove.md)	 - Delete one network device assignment of a cluster host interface

