## metalcloud-cli vm-instance variables

Show the variables of a VM instance

### Synopsis

Show the variables a VM instance exposes to extensions and OS templates.

The response is a deeply nested document, so every non native output format is
rendered as YAML.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Optional Flags:
  --usage string  Restrict the variables to one usage type.

Examples:
  # Show the variables of VM instance 67890
  metalcloud-cli vm-instance variables my-infra 67890

  # Show only the variables used by Ansible bundles
  metalcloud-cli vmi vars 12345 67890 --usage AnsibleBundle

```
metalcloud-cli vm-instance variables infrastructure_id_or_label vm_instance_id [flags]
```

### Options

```
  -h, --help           help for variables
      --usage string   Restrict the variables to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).
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

