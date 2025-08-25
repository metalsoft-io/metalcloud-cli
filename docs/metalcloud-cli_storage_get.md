## metalcloud-cli storage get

Get detailed information about a specific storage pool

### Synopsis

Get detailed information about a specific storage pool by its ID.

This command displays comprehensive information about a storage pool including its
configuration, status, driver details, technologies, and associated metadata.

Arguments:
  storage_id    The numeric ID of the storage pool to retrieve

Examples:
  # Get details for storage pool with ID 123
  metalcloud storage get 123

  # Using the show alias
  metalcloud storage show 456

```
metalcloud-cli storage get storage_id [flags]
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

* [metalcloud-cli storage](metalcloud-cli_storage.md)	 - Manage storage pools and related resources

