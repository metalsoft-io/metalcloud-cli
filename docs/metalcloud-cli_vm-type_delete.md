## metalcloud-cli vm-type delete

Delete a VM type

### Synopsis

Delete a VM type from the MetalCloud platform.

WARNING: This action is irreversible. Ensure that no VMs are currently using this VM type
before deletion, as this could cause issues with existing deployments.

ARGUMENTS:
  vm_type_id    The numeric ID of the VM type to delete

EXAMPLES:
  # Delete VM type with ID 123
  metalcloud vm-type delete 123
  
  # Check VMs using a VM type before deletion
  metalcloud vm-type vms 123
  metalcloud vm-type delete 123

```
metalcloud-cli vm-type delete vm_type_id [flags]
```

### Options

```
  -h, --help   help for delete
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

* [metalcloud-cli vm-type](metalcloud-cli_vm-type.md)	 - Manage VM types and configurations

