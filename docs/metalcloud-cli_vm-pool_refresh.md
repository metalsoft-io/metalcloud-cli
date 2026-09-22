## metalcloud-cli vm-pool refresh

Refresh the information of a VM pool

### Synopsis

Refresh the cached information MetalSoft holds about a VM pool.

On VMware VCF this reports any new datastores. The refreshed VM pool is printed
when the operation completes.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool to refresh

Examples:
  # Refresh VM pool 123
  metalcloud-cli vm-pool refresh 123

```
metalcloud-cli vm-pool refresh vm_pool_id [flags]
```

### Options

```
  -h, --help   help for refresh
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

