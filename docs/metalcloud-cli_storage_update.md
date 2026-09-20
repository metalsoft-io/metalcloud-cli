## metalcloud-cli storage update

Update an existing storage pool

### Synopsis

Update an existing storage pool from a JSON or YAML configuration.

Only the mutable properties of a storage pool can be changed (maintenance and
experimental flags, drive priorities, tags, network fabric, options and
credentials). The current revision of the storage pool is read first and sent
back as the If-Match entity tag, so a concurrent change is rejected by the API
instead of being silently overwritten.

Required Arguments:
  storage_id    The numeric ID of the storage pool to update

Required Flags:
  --config-source   Source of the storage update configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update a storage pool from a file
  metalcloud-cli storage update 123 --config-source ./storage-update.json

  # Put a storage pool in maintenance from piped input
  echo '{"inMaintenance": 1}' | metalcloud-cli storage update 123 --config-source pipe

```
metalcloud-cli storage update storage_id [flags]
```

### Options

```
      --config-source string   Source of the storage update configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

