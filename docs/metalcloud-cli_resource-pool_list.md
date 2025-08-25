## metalcloud-cli resource-pool list

List all resource pools with optional filtering and pagination

### Synopsis

List all resource pools in the system with optional filtering and pagination support.

This command displays a table of resource pools showing their ID, label, and description.
You can filter results using search terms and control the output with pagination parameters.

Flags:
  --page int      Page number for pagination (default: 0, shows all results)
  --limit int     Maximum number of records per page (default: 0, max: 100)
  --search string Search term to filter results by label or description

The search parameter performs a case-insensitive substring match against both
the resource pool label and description fields.

Examples:
  # List all resource pools
  metalcloud-cli resource-pool list

  # List resource pools with pagination
  metalcloud-cli resource-pool list --page 1 --limit 10

  # Search for resource pools containing "production"
  metalcloud-cli resource-pool list --search "production"

  # Combine search with pagination
  metalcloud-cli resource-pool list --search "dev" --page 1 --limit 5

```
metalcloud-cli resource-pool list [flags]
```

### Options

```
  -h, --help            help for list
      --limit int       Number of records per page (max 100)
      --page int        Page number
      --search string   Search term to filter results
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

* [metalcloud-cli resource-pool](metalcloud-cli_resource-pool.md)	 - Manage resource pools and their associated resources

