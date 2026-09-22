## metalcloud-cli vm-instance-group network disconnect

Remove one network connection of a VM instance group

### Synopsis

Remove one network connection of a VM instance group.

This disconnects every instance of the group from the logical network.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  connection_id               The numeric ID of the network connection to remove

Examples:
  # Remove network connection 535 of group 67890
  metalcloud-cli vm-instance-group network disconnect my-infra 67890 535

  # Using an alias
  metalcloud-cli vmg net rm 12345 67890 535

```
metalcloud-cli vm-instance-group network disconnect infrastructure_id_or_label vm_instance_group_id connection_id [flags]
```

### Options

```
  -h, --help   help for disconnect
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

