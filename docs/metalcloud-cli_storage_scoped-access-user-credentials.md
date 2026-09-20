## metalcloud-cli storage scoped-access-user-credentials

Get the credentials of a scoped access user

### Synopsis

Get the credentials (username, password and/or API token) of a scoped access
user of a storage pool.

Required Arguments:
  storage_id    The numeric ID of the storage pool
  user_id       The numeric ID of the scoped access user

Examples:
  # Get the credentials of scoped access user 42 of storage pool 123
  metalcloud-cli storage scoped-access-user-credentials 123 42

  # Save the credentials as JSON
  metalcloud-cli storage scoped-access-user-credentials 123 42 -f json > creds.json

```
metalcloud-cli storage scoped-access-user-credentials storage_id user_id [flags]
```

### Options

```
  -h, --help   help for scoped-access-user-credentials
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

