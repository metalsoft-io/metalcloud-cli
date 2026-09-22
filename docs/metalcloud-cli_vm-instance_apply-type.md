## metalcloud-cli vm-instance apply-type

Apply a VM type on a VM instance

### Synopsis

Apply a VM type on a VM instance.

The instance is reconfigured with the CPU, memory and GPU specification of the
given VM type. The current instance revision is fetched automatically and sent
as the If-Match header.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance
  vm_type_id                  The numeric ID of the VM type to apply

Examples:
  # Apply VM type 5 on VM instance 67890
  metalcloud-cli vm-instance apply-type my-infra 67890 5

  # Using the alias
  metalcloud-cli vmi set-type 12345 67890 5

```
metalcloud-cli vm-instance apply-type infrastructure_id_or_label vm_instance_id vm_type_id [flags]
```

### Options

```
  -h, --help   help for apply-type
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

