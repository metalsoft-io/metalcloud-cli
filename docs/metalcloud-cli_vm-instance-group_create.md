## metalcloud-cli vm-instance-group create

Create a new VM instance group in an infrastructure

### Synopsis

Create a new VM instance group in an infrastructure.

This command creates a new VM instance group with the specified configuration.
The group will contain multiple VM instances of the same type and configuration,
making it easier to manage and scale similar workloads.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_type_id                  The VM type ID defining CPU, memory and other specs
  disk_size_gb                The disk size in GB of each VM instance of the group
  instance_count              The number of VM instances to create in the group
  os_template_id              Optional. The OS template ID for the VM instances

Examples:
  # Create a VM instance group with 3 instances using a specific OS template
  metalcloud-cli vm-instance-group create 12345 5 100 3 7

  # Create a VM instance group without specifying an OS template
  metalcloud-cli vmg new my-infra 5 50 2

```
metalcloud-cli vm-instance-group create infrastructure_id_or_label vm_type_id disk_size_gb instance_count [os_template_id] [flags]
```

### Options

```
  -h, --help   help for create
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

