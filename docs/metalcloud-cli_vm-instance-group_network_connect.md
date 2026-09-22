## metalcloud-cli vm-instance-group network connect

Connect a VM instance group to a logical network

### Synopsis

Connect a VM instance group to a logical network.

The connection can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Required Flags:
  --config-source string       Source of the new connection configuration.
                               Can be 'pipe' or path to a JSON/YAML file.
  --logical-network-id string  The ID of the logical network to connect to.
  --access-mode string         The network access mode (e.g. 'l2').

Optional Flags:
  --tagged string      Whether VLAN tagging is enabled (true/false).
  --redundancy string  The redundancy mode ('active-backup', 'active-active').
  --mtu string         The MTU of the network connection.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. One of
  --config-source or --logical-network-id is required; --access-mode is
  required when --logical-network-id is used.

Examples:
  # Connect group 67890 to logical network 5 from flags
  metalcloud-cli vm-instance-group network connect my-infra 67890 --logical-network-id 5 --access-mode l2 --tagged true

  # Connect from a configuration file
  metalcloud-cli vmg net add 12345 67890 --config-source connection.json

```
metalcloud-cli vm-instance-group network connect infrastructure_id_or_label vm_instance_group_id [flags]
```

### Options

```
      --access-mode string          Network connection access mode.
      --config-source string        Source of the new network connection configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                        help for connect
      --logical-network-id string   The ID of the logical network to connect to.
      --mtu string                  Network connection MTU.
      --redundancy string           Network connection redundancy mode.
      --tagged string               Network connection VLAN tagging (true/false).
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

