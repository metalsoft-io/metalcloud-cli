## metalcloud-cli vm-pool update-cluster-host-interface

Update one network interface of a cluster host

### Synopsis

Update a cluster host network interface from a JSON or YAML configuration.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Required Flags:
  --config-source  Source of the interface update configuration.
                   Values: 'pipe' for stdin input, or path to a JSON/YAML file.

Configuration Fields:
  status           Required. One of 'managed', 'unmanaged', 'inactive'.

Examples:
  # Take an interface under management
  echo '{"status":"managed"}' | metalcloud-cli vm-pool update-cluster-host-interface 123 456 789 --config-source pipe

  # Update from a file
  metalcloud-cli vm-pool update-cluster-host-interface 123 456 789 --config-source interface.json

```
metalcloud-cli vm-pool update-cluster-host-interface vm_pool_id host_id interface_id [flags]
```

### Options

```
      --config-source string   Source of the cluster host interface update configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-cluster-host-interface
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

