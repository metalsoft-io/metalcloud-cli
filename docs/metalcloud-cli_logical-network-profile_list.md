## metalcloud-cli logical-network-profile list

List logical network profiles with optional filtering

### Synopsis

List all logical network profiles in the system with optional filtering capabilities.

This command displays a tabular view of logical network profiles including their
ID, name, kind, label, fabric ID, and other metadata. Use the filter flags to
narrow down results based on specific criteria.

Flags:
  --filter-id          Filter profiles by one or more profile IDs
  --filter-label       Filter profiles by label pattern matching
  --filter-kind        Filter profiles by profile kind (e.g., 'cisco', 'juniper')
  --filter-name        Filter profiles by name pattern matching
  --filter-fabric-id   Filter profiles by associated fabric ID
  --sort-by            Sort results by specified fields with direction

Examples:
  # List all logical network profiles
  metalcloud-cli logical-network-profile list

  # Filter by profile kind
  metalcloud-cli lnp list --filter-kind cisco

  # Filter by multiple criteria
  metalcloud-cli network-profile list --filter-kind cisco --filter-label "prod"

  # Sort by name in descending order
  metalcloud-cli lnp ls --sort-by name:DESC

  # Filter by fabric ID and sort by ID
  metalcloud-cli lnp list --filter-fabric-id 100 --sort-by id:ASC

```
metalcloud-cli logical-network-profile list [flags]
```

### Options

```
      --filter-fabric-id strings   Filter by fabric ID.
      --filter-id strings          Filter by profile ID.
      --filter-kind strings        Filter by profile kind.
      --filter-label strings       Filter by profile label.
      --filter-name strings        Filter by profile name.
  -h, --help                       help for list
      --sort-by strings            Sort by fields (e.g., id:ASC, name:DESC).
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

* [metalcloud-cli logical-network-profile](metalcloud-cli_logical-network-profile.md)	 - Manage logical network profiles for network configuration templates

