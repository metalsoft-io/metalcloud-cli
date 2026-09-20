## metalcloud-cli vm-pool update-cluster-host

Update one cluster host of a VM pool

### Synopsis

Update a cluster host of a VM pool from a JSON or YAML configuration.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

Required Flags:
  --config-source  Source of the cluster host update configuration.
                   Values: 'pipe' for stdin input, or path to a JSON/YAML file.

Configuration Fields (all optional):
  allowVMsToBeCreated         Allow new VMs to be created on this host
  allowContainersToBeCreated  Allow new containers to be created on this host

Examples:
  # Stop scheduling new VMs on a host
  echo '{"allowVMsToBeCreated": false}' | metalcloud-cli vm-pool update-cluster-host 123 456 --config-source pipe

  # Update from a file
  metalcloud-cli vm-pool update-cluster-host 123 456 --config-source host.json

```
metalcloud-cli vm-pool update-cluster-host vm_pool_id host_id [flags]
```

### Options

```
      --config-source string   Source of the cluster host update configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-cluster-host
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

