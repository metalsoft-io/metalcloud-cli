## metalcloud-cli vm-pool containers

List the containers running in a VM pool

### Synopsis

List all containers deployed on the hosts of a VM pool.

The listing walks every page transparently.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool

Examples:
  # List containers of VM pool 123
  metalcloud-cli vm-pool containers 123

```
metalcloud-cli vm-pool containers vm_pool_id [flags]
```

### Options

```
  -h, --help   help for containers
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

