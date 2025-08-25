## metalcloud-cli storage file-shares

List file shares in a storage pool

### Synopsis

List file shares available in a specific storage pool.

This command retrieves all file shares associated with a storage pool, showing their
configuration, status, and access information. File shares are typically used for
NFS or CIFS/SMB file storage. Results can be paginated using the limit and page flags.

Arguments:
  storage_id    The numeric ID of the storage pool

Optional flags:
  --limit       Number of records per page (default: all records)
  --page        Page number for pagination (requires --limit)

Examples:
  # List all file shares for storage pool 123
  metalcloud storage file-shares 123

  # List first 5 file shares
  metalcloud storage file-shares 123 --limit 5

  # List third page with 5 file shares per page
  metalcloud storage file-shares 123 --limit 5 --page 3

```
metalcloud-cli storage file-shares storage_id [flags]
```

### Options

```
  -h, --help           help for file-shares
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

