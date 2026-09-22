## metalcloud-cli vm-pool cluster-host-interface-network-device get

Get one network device assignment of a cluster host interface

### Synopsis

Get one link between a VM pool cluster host interface and a network device
(switch) interface.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface
  assignment_id    The numeric ID of the network device assignment

Examples:
  # Get assignment 5 of interface 789
  metalcloud-cli vm-pool cluster-host-interface-network-device get 123 456 789 5

```
metalcloud-cli vm-pool cluster-host-interface-network-device get vm_pool_id host_id interface_id assignment_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli vm-pool cluster-host-interface-network-device](metalcloud-cli_vm-pool_cluster-host-interface-network-device.md)	 - Manage the network device assignments of a cluster host interface

