## metalcloud-cli storage statistics

Show capacity statistics for one or for all storage pools

### Synopsis

Show storage capacity statistics.

Called without an argument the command returns the global statistics of all
storage pools (counts per state and per type, total used and free space).
Called with a storage pool ID it returns the total, used and free capacity of
that single storage pool.

Optional Arguments:
  storage_id    The numeric ID of a storage pool. When omitted, global
                statistics for all storage pools are returned.

Examples:
  # Global statistics for all storage pools
  metalcloud-cli storage statistics

  # Statistics for storage pool 123
  metalcloud-cli storage statistics 123

```
metalcloud-cli storage statistics [storage_id] [flags]
```

### Options

```
  -h, --help   help for statistics
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

