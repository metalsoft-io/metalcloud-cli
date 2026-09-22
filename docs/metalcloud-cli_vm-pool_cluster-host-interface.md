## metalcloud-cli vm-pool cluster-host-interface

Get one network interface of a cluster host

### Synopsis

Get the details of a single network interface of a cluster host, including its
MAC address, management status, fabric and network device assignments.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Examples:
  # Get interface 789 of cluster host 456 in VM pool 123
  metalcloud-cli vm-pool cluster-host-interface 123 456 789

```
metalcloud-cli vm-pool cluster-host-interface vm_pool_id host_id interface_id [flags]
```

### Options

```
  -h, --help   help for cluster-host-interface
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

