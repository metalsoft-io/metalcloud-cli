## metalcloud-cli vm-instance-group network connections

List the network connections of a VM instance group

### Synopsis

List all network connections of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # List the network connections of group 67890
  metalcloud-cli vm-instance-group network connections my-infra 67890

  # Using the alias
  metalcloud-cli vmg net list-connections 12345 67890

```
metalcloud-cli vm-instance-group network connections infrastructure_id_or_label vm_instance_group_id [flags]
```

### Options

```
  -h, --help   help for connections
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

* [metalcloud-cli vm-instance-group network](metalcloud-cli_vm-instance-group_network.md)	 - Manage the network connections of a VM instance group

