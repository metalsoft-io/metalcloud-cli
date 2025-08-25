## metalcloud-cli resource-pool get

Get detailed information about a specific resource pool

### Synopsis

Get detailed information about a specific resource pool by its ID.

This command retrieves and displays comprehensive information about a resource pool
including its ID, label, description, and any associated metadata.

Arguments:
  pool_id    The numeric ID of the resource pool to retrieve

Examples:
  # Get information about resource pool with ID 123
  metalcloud-cli resource-pool get 123

  # Get resource pool details using alias
  metalcloud-cli rp get 456

```
metalcloud-cli resource-pool get <pool_id> [flags]
```

### Options

```
  -h, --help   help for get
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

