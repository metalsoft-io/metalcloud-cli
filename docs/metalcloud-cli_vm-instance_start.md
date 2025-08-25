## metalcloud-cli vm-instance start

Start a VM instance

### Synopsis

Start a VM instance.

This command initiates the startup process for a VM instance that is currently
powered off or stopped. The instance will be powered on and begin booting
according to its configured operating system and startup settings.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to start

EXAMPLES:
  # Start VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance start 12345 67890
  
  # Start instance using alias
  metalcloud-cli vm start 12345 67890

```
metalcloud-cli vm-instance start infrastructure_id vm_instance_id [flags]
```

### Options

```
  -h, --help   help for start
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

