## metalcloud-cli vm-instance credentials

Get login credentials for a VM instance

### Synopsis

Get the login credentials of a VM instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Examples:
  # Show the credentials of VM instance 67890
  metalcloud-cli vm-instance credentials 12345 67890

  # Using the alias
  metalcloud-cli vmi creds my-infra 67890

```
metalcloud-cli vm-instance credentials infrastructure_id_or_label vm_instance_id [flags]
```

### Options

```
  -h, --help   help for credentials
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

