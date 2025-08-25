## metalcloud-cli vm-instance get

Get detailed information about a specific VM instance

### Synopsis

Get detailed information about a specific VM instance.

This command retrieves comprehensive information about a VM instance including
its current status, configuration, network details, disk information, and
associated metadata. This is useful for debugging, monitoring, and understanding
the current state of a virtual machine.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to retrieve details for

EXAMPLES:
  # Get details of VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance get 12345 67890
  
  # Get instance details using alias
  metalcloud-cli vmi show 12345 67890
  
  # Get instance details using short alias
  metalcloud-cli vm get 12345 67890

```
metalcloud-cli vm-instance get infrastructure_id vm_instance_id [flags]
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

* [metalcloud-cli vm-instance](metalcloud-cli_vm-instance.md)	 - Manage individual VM instances within infrastructures

