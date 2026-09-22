## metalcloud-cli storage scoped-access-users

List the scoped access users of a storage pool

### Synopsis

List the scoped access users provisioned on a storage pool.

Scoped access users are per-infrastructure accounts created on the storage
system so that an infrastructure can only access its own storage resources.

Required Arguments:
  storage_id    The numeric ID of the storage pool

Examples:
  # List the scoped access users of storage pool 123
  metalcloud-cli storage scoped-access-users 123

```
metalcloud-cli storage scoped-access-users storage_id [flags]
```

### Options

```
  -h, --help   help for scoped-access-users
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

