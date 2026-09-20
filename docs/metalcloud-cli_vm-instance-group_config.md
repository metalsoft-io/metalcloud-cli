## metalcloud-cli vm-instance-group config

Show the pending configuration of a VM instance group

### Synopsis

Show the pending configuration of a VM instance group.

The configuration holds the changes that are applied when the infrastructure is
deployed, together with the revision that guards concurrent updates.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # Show the configuration of group 67890
  metalcloud-cli vm-instance-group config my-infra 67890

  # Using the alias
  metalcloud-cli vmg get-config 12345 67890

```
metalcloud-cli vm-instance-group config infrastructure_id_or_label vm_instance_group_id [flags]
```

### Options

```
  -h, --help   help for config
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

