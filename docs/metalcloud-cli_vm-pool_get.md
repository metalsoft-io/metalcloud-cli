## metalcloud-cli vm-pool get

Get detailed information about a specific VM pool

### Synopsis

Get comprehensive details about a virtual machine pool including configuration, 
status, and connection information.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool to retrieve

EXAMPLES:
  # Get details for VM pool with ID 123
  metalcloud-cli vm-pool get 123

  # Using alias
  metalcloud-cli vm-pool show 123

```
metalcloud-cli vm-pool get vm_pool_id [flags]
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

* [metalcloud-cli vm-pool](metalcloud-cli_vm-pool.md)	 - Manage virtual machine pools and their resources

