## metalcloud-cli storage delete

Delete an existing storage pool

### Synopsis

Delete an existing storage pool by its ID.

This command permanently removes a storage pool from the MetalCloud infrastructure.
Warning: This action is irreversible and will remove all associated data.

Arguments:
  storage_id    The numeric ID of the storage pool to delete

Examples:
  # Delete storage pool with ID 123
  metalcloud storage delete 123

  # Delete storage pool with confirmation
  metalcloud storage delete 456

```
metalcloud-cli storage delete storage_id [flags]
```

### Options

```
  -h, --help   help for delete
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

