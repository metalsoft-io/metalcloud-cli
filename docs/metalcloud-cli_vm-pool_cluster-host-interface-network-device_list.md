## metalcloud-cli vm-pool cluster-host-interface-network-device list

List the network device assignments of a cluster host interface

### Synopsis

List every network device (switch) interface a VM pool cluster host interface is
linked to.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Examples:
  # List the assignments of interface 789
  metalcloud-cli vm-pool cluster-host-interface-network-device list 123 456 789

  # Using the group alias
  metalcloud-cli vm-pool chind ls 123 456 789

```
metalcloud-cli vm-pool cluster-host-interface-network-device list vm_pool_id host_id interface_id [flags]
```

### Options

```
  -h, --help   help for list
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

