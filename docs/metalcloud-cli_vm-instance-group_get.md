## metalcloud-cli vm-instance-group get

Get details of a specific VM instance group

### Synopsis

Get detailed information about a specific VM instance group.

This command retrieves comprehensive information about a VM instance group
including its configuration, current status, instances, and associated metadata.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # Get details of VM instance group 67890 in infrastructure 12345
  metalcloud-cli vm-instance-group get 12345 67890

  # Get group details using alias
  metalcloud-cli vmg show my-infra 67890

```
metalcloud-cli vm-instance-group get infrastructure_id_or_label vm_instance_group_id [flags]
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

* [metalcloud-cli vm-instance-group](metalcloud-cli_vm-instance-group.md)	 - Manage VM instance groups within infrastructures

