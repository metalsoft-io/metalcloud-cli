## metalcloud-cli vm-pool statistics

Show the resource usage statistics of a VM pool

### Synopsis

Show the aggregated RAM, disk and GPU usage of a VM pool.

In the tabular formats (text, csv, md) the GPU list is flattened into a single
column; json and yaml return the full statistics object.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool

Examples:
  # Show statistics for VM pool 123
  metalcloud-cli vm-pool statistics 123

  # Get the raw statistics object
  metalcloud-cli vm-pool stats 123 -f json

```
metalcloud-cli vm-pool statistics vm_pool_id [flags]
```

### Options

```
  -h, --help   help for statistics
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

