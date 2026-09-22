## metalcloud-cli vm-instance shutdown

Shutdown a VM instance

### Synopsis

Shutdown a VM instance gracefully.

This command initiates a graceful shutdown process for a running VM instance.
The instance will receive a shutdown signal and will attempt to properly
terminate all running processes before powering off. This is the recommended
way to stop a VM instance to prevent data loss.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance to shutdown

Examples:
  # Shutdown VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance shutdown 12345 67890

  # Shutdown instance using alias
  metalcloud-cli vm shutdown my-infra 67890

```
metalcloud-cli vm-instance shutdown infrastructure_id_or_label vm_instance_id [flags]
```

### Options

```
  -h, --help   help for shutdown
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

