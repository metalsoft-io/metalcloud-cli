## metalcloud-cli vm-pool sync

Sync a VM pool with its hypervisor

### Synopsis

Start a synchronization job for a VM pool.

The sync discovers objects that were created directly on the hypervisor. On
VMware VCF, for example, it discovers new Virtual Distributed Switches. The
command returns the job information of the job that was started; use
'metalcloud-cli job get <job_id>' to follow it.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool to sync

Examples:
  # Sync VM pool 123
  metalcloud-cli vm-pool sync 123

```
metalcloud-cli vm-pool sync vm_pool_id [flags]
```

### Options

```
  -h, --help   help for sync
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

