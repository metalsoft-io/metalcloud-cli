## metalcloud-cli vm-instance reboot

Reboot a VM instance

### Synopsis

Reboot a VM instance.

This command initiates a restart process for a running VM instance. The instance
will be gracefully shutdown and then automatically restarted. This is useful
for applying configuration changes or recovering from software issues.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to reboot

EXAMPLES:
  # Reboot VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance reboot 12345 67890
  
  # Reboot instance using alias
  metalcloud-cli vm reboot 12345 67890

```
metalcloud-cli vm-instance reboot infrastructure_id vm_instance_id [flags]
```

### Options

```
  -h, --help   help for reboot
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

