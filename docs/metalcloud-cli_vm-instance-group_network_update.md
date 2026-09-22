## metalcloud-cli vm-instance-group network update

Update one network connection of a VM instance group

### Synopsis

Update one network connection of a VM instance group.

The updates can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  connection_id               The numeric ID of the network connection

Required Flags:
  --config-source string  Source of the connection updates.
                          Can be 'pipe' or path to a JSON/YAML file.
  --access-mode string    The network access mode (e.g. 'l2').
  --tagged string         Whether VLAN tagging is enabled (true/false).
  --redundancy string     The redundancy mode ('active-backup', 'active-active').
  --mtu string            The MTU of the network connection.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. At least one
  of --config-source, --access-mode, --tagged, --redundancy or --mtu is required.

Examples:
  # Change the access mode of connection 535
  metalcloud-cli vm-instance-group network update my-infra 67890 535 --access-mode l2

  # Enable VLAN tagging
  metalcloud-cli vmg net edit 12345 67890 535 --tagged true

```
metalcloud-cli vm-instance-group network update infrastructure_id_or_label vm_instance_group_id connection_id [flags]
```

### Options

```
      --access-mode string     Network connection access mode.
      --config-source string   Source of the network connection updates. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
      --mtu string             Network connection MTU.
      --redundancy string      Network connection redundancy mode.
      --tagged string          Network connection VLAN tagging (true/false).
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

