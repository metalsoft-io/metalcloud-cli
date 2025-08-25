## metalcloud-cli vm-instance-group list

List all VM instance groups in an infrastructure

### Synopsis

List all VM instance groups in an infrastructure.

This command retrieves and displays all VM instance groups that exist within
the specified infrastructure. The output includes group details such as ID,
label, instance count, VM type, and current status.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure to list groups from

EXAMPLES:
  # List all VM instance groups in infrastructure 12345
  metalcloud-cli vm-instance-group list 12345
  
  # List groups using alias
  metalcloud-cli vmg ls 12345

```
metalcloud-cli vm-instance-group list infrastructure_id [flags]
```

### Options

```
  -h, --help   help for list
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

