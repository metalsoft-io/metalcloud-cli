## metalcloud-cli storage drives

List drives available in a storage pool

### Synopsis

List drives available in a specific storage pool.

This command retrieves all drives associated with a storage pool, showing their
configuration, status, and specifications. Results can be paginated using the
limit and page flags.

Arguments:
  storage_id    The numeric ID of the storage pool

Optional flags:
  --limit       Number of records per page (default: all records)
  --page        Page number for pagination (requires --limit)

Examples:
  # List all drives for storage pool 123
  metalcloud storage drives 123

  # List first 10 drives
  metalcloud storage drives 123 --limit 10

  # List second page with 10 drives per page
  metalcloud storage drives 123 --limit 10 --page 2

```
metalcloud-cli storage drives storage_id [flags]
```

### Options

```
  -h, --help           help for drives
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

* [metalcloud-cli storage](metalcloud-cli_storage.md)	 - Manage storage pools and related resources

