## metalcloud-cli vm-pool cluster-hosts

List cluster hosts in a VM pool with pagination

### Synopsis

List all cluster hosts (ESXi hosts, Hyper-V servers, etc.) in a specific VM pool with optional pagination.

This command displays the hypervisor hosts that are part of the specified VM pool,
including their status, resource utilization, and connection details.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool

OPTIONAL FLAGS:
  --limit          Number of records to return per page (default: all)
  --page           Page number to retrieve (1-based, default: 1)
                   Only effective when --limit is specified

PAGINATION:
When using pagination, specify both --limit and --page for best results.
The --limit flag controls how many records are returned, while --page
specifies which page of results to retrieve.

EXAMPLES:
  # List all cluster hosts in VM pool 123
  metalcloud-cli vm-pool cluster-hosts 123

  # List first 5 hosts
  metalcloud-cli vm-pool cluster-hosts 123 --limit 5

  # List second page of 5 hosts each
  metalcloud-cli vm-pool cluster-hosts 123 --limit 5 --page 2

```
metalcloud-cli vm-pool cluster-hosts vm_pool_id [flags]
```

### Options

```
  -h, --help           help for cluster-hosts
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

