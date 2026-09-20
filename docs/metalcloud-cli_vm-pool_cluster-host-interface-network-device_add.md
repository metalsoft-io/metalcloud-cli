## metalcloud-cli vm-pool cluster-host-interface-network-device add

Link a cluster host interface to a network device interface

### Synopsis

Link a VM pool cluster host interface to an interface of a network device
(switch). The network device must be active and in a leaf position.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Required Flags (either the pair of flags, or --config-source):
  --network-device            ID or label of the network device (switch)
  --network-device-interface  Name of the interface on the network device
  --config-source             Source of the assignment configuration.
                              Values: 'pipe' for stdin input, or path to a
                              JSON/YAML file. Mutually exclusive with the flags
                              above; the file carries a numeric networkDeviceId.

Examples:
  # Link by network device label
  metalcloud-cli vm-pool chind add 123 456 789 --network-device leaf-su00-r0 --network-device-interface Ethernet1/1

  # Link by network device ID
  metalcloud-cli vm-pool chind add 123 456 789 --network-device 42 --network-device-interface Ethernet1/1

  # Link from a configuration file
  metalcloud-cli vm-pool chind add 123 456 789 --config-source assignment.json

```
metalcloud-cli vm-pool cluster-host-interface-network-device add vm_pool_id host_id interface_id [flags]
```

### Options

```
      --config-source string              Source of the network device assignment configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                              help for add
      --network-device string             ID or label of the network device (switch) to link the interface to.
      --network-device-interface string   Name of the interface on the network device.
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

