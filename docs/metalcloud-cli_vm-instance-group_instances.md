## metalcloud-cli vm-instance-group instances

List VM instances within a VM instance group

### Synopsis

List all VM instances within a specific VM instance group.

This command displays all individual VM instances that belong to the specified
VM instance group. The output includes instance details such as ID, status,
type and disk size for each instance in the group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # List all instances in VM instance group 67890 from infrastructure 12345
  metalcloud-cli vm-instance-group instances 12345 67890

  # List instances using alias
  metalcloud-cli vmg instances-ls my-infra 67890

```
metalcloud-cli vm-instance-group instances infrastructure_id_or_label vm_instance_group_id [flags]
```

### Options

```
  -h, --help   help for instances
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

