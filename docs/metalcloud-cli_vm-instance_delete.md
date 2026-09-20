## metalcloud-cli vm-instance delete

Delete a VM instance

### Synopsis

Delete a VM instance from an infrastructure.

The instance is removed when the infrastructure is deployed. The current
instance revision is fetched automatically and sent as the If-Match header.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance to delete

Examples:
  # Delete VM instance 67890 from infrastructure 12345
  metalcloud-cli vm-instance delete 12345 67890

  # Using the alias
  metalcloud-cli vmi rm my-infra 67890

```
metalcloud-cli vm-instance delete infrastructure_id_or_label vm_instance_id [flags]
```

### Options

```
  -h, --help   help for delete
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

