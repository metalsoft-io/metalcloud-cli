## metalcloud-cli storage scoped-access-user

Get one scoped access user of a storage pool

### Synopsis

Get the details of a single scoped access user of a storage pool.

Required Arguments:
  storage_id    The numeric ID of the storage pool
  user_id       The numeric ID of the scoped access user

Examples:
  # Get scoped access user 42 of storage pool 123
  metalcloud-cli storage scoped-access-user 123 42

```
metalcloud-cli storage scoped-access-user storage_id user_id [flags]
```

### Options

```
  -h, --help   help for scoped-access-user
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

