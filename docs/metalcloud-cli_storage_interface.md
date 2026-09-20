## metalcloud-cli storage interface

Get one interface of a storage pool

### Synopsis

Get the details of a single storage pool interface.

Required Arguments:
  storage_id      The numeric ID of the storage pool
  interface_id    The numeric ID of the storage interface

Examples:
  # Get interface 7 of storage pool 123
  metalcloud-cli storage interface 123 7

```
metalcloud-cli storage interface storage_id interface_id [flags]
```

### Options

```
  -h, --help   help for interface
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

