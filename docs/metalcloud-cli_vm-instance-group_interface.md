## metalcloud-cli vm-instance-group interface

Get one network interface of a VM instance group

### Synopsis

Get detailed information about one network interface of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  interface_id                The numeric ID of the interface

Examples:
  # Get interface 2 of group 67890
  metalcloud-cli vm-instance-group interface my-infra 67890 2

  # Using the alias
  metalcloud-cli vmg get-interface 12345 67890 2

```
metalcloud-cli vm-instance-group interface infrastructure_id_or_label vm_instance_group_id interface_id [flags]
```

### Options

```
  -h, --help   help for interface
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

