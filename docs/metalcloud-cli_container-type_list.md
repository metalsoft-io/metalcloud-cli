## metalcloud-cli container-type list

List all container types

### Synopsis

List all container types.

Retrieves every container type visible to the current user. The output can be
narrowed with the filter flags below; each filter accepts multiple values.

Optional Flags:
  --filter-id strings            Filter by container type ID
  --filter-label strings         Filter by container type label
  --filter-name strings          Filter by container type name
  --filter-display-name strings  Filter by container type display name

Examples:
  # List all container types
  metalcloud-cli container-type list

  # List only the container types with a given label
  metalcloud-cli ct ls --filter-label small-container

  # List several container types by ID
  metalcloud-cli ct ls --filter-id 10 --filter-id 11

```
metalcloud-cli container-type list [flags]
```

### Options

```
      --filter-display-name strings   Filter by container type display name.
      --filter-id strings             Filter by container type ID.
      --filter-label strings          Filter by container type label.
      --filter-name strings           Filter by container type name.
  -h, --help                          help for list
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

* [metalcloud-cli container-type](metalcloud-cli_container-type.md)	 - Manage container types

