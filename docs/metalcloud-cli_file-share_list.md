## metalcloud-cli file-share list

List all file shares for an infrastructure

### Synopsis

List all file shares associated with the specified infrastructure.

This command displays file shares with their basic information including ID, name, 
status, and other key attributes. Results can be filtered by status.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label to list file shares for

Optional Flags:
  --filter-status               Filter results by file share status (can be used multiple times)
                               Common statuses: active, inactive, creating, deleting, error

Examples:
  # List all file shares for an infrastructure
  metalcloud-cli file-share list my-infrastructure

  # List file shares with ID
  metalcloud-cli file-share list 12345

  # Filter by status
  metalcloud-cli file-share list my-infrastructure --filter-status active
  
  # Filter by multiple statuses
  metalcloud-cli file-share list my-infrastructure --filter-status active --filter-status creating

```
metalcloud-cli file-share list infrastructure_id_or_label [flags]
```

### Options

```
      --filter-status strings   Filter the result by file share status.
  -h, --help                    help for list
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

* [metalcloud-cli file-share](metalcloud-cli_file-share.md)	 - Manage file shares for infrastructure resources

