## metalcloud-cli vm-instance os-installation-data

Show the OS installation data of a VM instance

### Synopsis

Show the values the OS installer of a VM instance is rendered with.

The response is a deeply nested document, so every non native output format is
rendered as YAML.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Optional Flags:
  --usage string   Restrict the OS installation data to one usage type.
  --remove-empty   Omit the empty entries from the response.

Examples:
  # Show the OS installation data of VM instance 67890
  metalcloud-cli vm-instance os-installation-data my-infra 67890

  # Skip the empty entries
  metalcloud-cli vmi os-data 12345 67890 --remove-empty

```
metalcloud-cli vm-instance os-installation-data infrastructure_id_or_label vm_instance_id [flags]
```

### Options

```
  -h, --help           help for os-installation-data
      --remove-empty   Omit the empty entries from the response.
      --usage string   Restrict the OS installation data to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).
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

