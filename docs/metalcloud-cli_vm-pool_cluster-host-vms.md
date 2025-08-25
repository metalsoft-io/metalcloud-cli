## metalcloud-cli vm-pool cluster-host-vms

List virtual machines on a specific cluster host with pagination

### Synopsis

List all virtual machines running on a specific cluster host within a VM pool with optional pagination.

This command displays VMs that are currently running on the specified cluster host,
including their status, resource allocation, and configuration details.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

OPTIONAL FLAGS:
  --limit          Number of records to return per page (default: all)
  --page           Page number to retrieve (1-based, default: 1)
                   Only effective when --limit is specified

PAGINATION:
When using pagination, specify both --limit and --page for best results.
The --limit flag controls how many records are returned, while --page
specifies which page of results to retrieve.

EXAMPLES:
  # List all VMs on cluster host 456 in VM pool 123
  metalcloud-cli vm-pool cluster-host-vms 123 456

  # List first 5 VMs on the host
  metalcloud-cli vm-pool cluster-host-vms 123 456 --limit 5

  # List second page of 5 VMs each
  metalcloud-cli vm-pool cluster-host-vms 123 456 --limit 5 --page 2

```
metalcloud-cli vm-pool cluster-host-vms vm_pool_id host_id [flags]
```

### Options

```
  -h, --help           help for cluster-host-vms
      --limit string   Number of records per page
      --page string    Page number
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

