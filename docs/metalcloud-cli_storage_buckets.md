## metalcloud-cli storage buckets

List object storage buckets in a storage pool

### Synopsis

List object storage buckets available in a specific storage pool.

This command retrieves all object storage buckets associated with a storage pool,
showing their configuration, status, and access information. Buckets are typically
used for S3-compatible object storage. Results can be paginated using the limit
and page flags.

Arguments:
  storage_id    The numeric ID of the storage pool

Optional flags:
  --limit       Number of records per page (default: all records)
  --page        Page number for pagination (requires --limit)

Examples:
  # List all buckets for storage pool 123
  metalcloud storage buckets 123

  # List first 10 buckets
  metalcloud storage buckets 123 --limit 10

  # List second page with 10 buckets per page
  metalcloud storage buckets 123 --limit 10 --page 2

```
metalcloud-cli storage buckets storage_id [flags]
```

### Options

```
  -h, --help           help for buckets
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

