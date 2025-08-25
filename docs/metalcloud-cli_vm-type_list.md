## metalcloud-cli vm-type list

List all VM types with optional pagination

### Synopsis

List all VM types available in the MetalCloud platform.

This command displays all VM types with their specifications including CPU cores, RAM, 
experimental status, and whether they are restricted to unmanaged VMs only.

FLAGS:
  --limit string    Number of records per page (optional)
  --page string     Page number to retrieve (optional, requires --limit)

EXAMPLES:
  # List all VM types
  metalcloud vm-type list
  
  # List VM types with pagination (10 per page, page 1)
  metalcloud vm-type list --limit 10 --page 1
  
  # List first 5 VM types
  metalcloud vm-type list --limit 5

```
metalcloud-cli vm-type list [flags]
```

### Options

```
  -h, --help           help for list
      --limit string   Number of records per page
      --page string    Page number
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

* [metalcloud-cli vm-type](metalcloud-cli_vm-type.md)	 - Manage VM types and configurations

