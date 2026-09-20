## metalcloud-cli vm-pool cluster-host

Get details of one cluster host of a VM pool

### Synopsis

Get the full details of a single cluster host (ESXi host, Hyper-V server, ...)
that belongs to a VM pool, including its health status, roles and the flags that
control whether VMs and containers may be created on it.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

Examples:
  # Get cluster host 456 of VM pool 123
  metalcloud-cli vm-pool cluster-host 123 456

```
metalcloud-cli vm-pool cluster-host vm_pool_id host_id [flags]
```

### Options

```
  -h, --help   help for cluster-host
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

