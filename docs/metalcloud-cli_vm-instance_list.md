## metalcloud-cli vm-instance list

List all VM instances in an infrastructure

### Synopsis

List all VM instances in an infrastructure.

This command retrieves and displays all VM instances that exist within the
specified infrastructure. The output includes instance details such as ID,
status, VM type, disk size and other relevant information for each instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List all VM instances in infrastructure 12345
  metalcloud-cli vm-instance list 12345

  # List instances by infrastructure label
  metalcloud-cli vmi ls my-infra

```
metalcloud-cli vm-instance list infrastructure_id_or_label [flags]
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

* [metalcloud-cli vm-instance](metalcloud-cli_vm-instance.md)	 - Manage individual VM instances within infrastructures

