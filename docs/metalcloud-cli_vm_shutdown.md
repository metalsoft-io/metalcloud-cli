## metalcloud-cli vm shutdown

Gracefully shutdown or force stop a VM

### Synopsis

Shutdown a virtual machine gracefully or force stop it. The command attempts
a graceful shutdown first, which allows the guest OS to properly close running
applications and services. If the VM doesn't respond, it can be force stopped.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine to shutdown.

Prerequisites:
  - VM must be in 'running' state
  - User must have write permissions for the VM

Examples:
  # Gracefully shutdown a VM
  metalcloud-cli vm shutdown 12345
  
  # Using the alias 'stop'
  metalcloud-cli vm stop 12345
  
  # Shutdown multiple VMs
  for vm in 12345 12346 12347; do metalcloud-cli vm shutdown $vm; done

```
metalcloud-cli vm shutdown vm_id [flags]
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

* [metalcloud-cli vm](metalcloud-cli_vm.md)	 - Manage virtual machines lifecycle and configuration

