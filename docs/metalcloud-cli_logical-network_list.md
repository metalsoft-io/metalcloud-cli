## metalcloud-cli logical-network list

List logical networks with optional filtering and sorting

### Synopsis

List all logical networks, optionally filtered by fabric and other criteria.

This command displays logical networks in a tabular format. You can optionally provide
a fabric ID or label to filter results to networks within that specific fabric.

Arguments:
  fabric_id_or_label  Optional fabric identifier to filter networks (can be ID or label)

Flags:
  --filter-id                Filter results by logical network ID(s) (can be used multiple times)
  --filter-label             Filter results by logical network label(s) (can be used multiple times) 
  --filter-fabric-id         Filter results by fabric ID(s) (can be used multiple times)
  --filter-infrastructure-id Filter results by infrastructure ID(s) (can be used multiple times). Use 'null' to filter public logical networks
  --filter-kind              Filter results by network kind(s) like 'vlan', 'vxlan' (can be used multiple times)
  --sort-by                  Sort results by field(s) with direction (e.g., id:ASC, name:DESC)
  --page                     Page number to retrieve (default: 1)
  --limit                    Number of records per page (default: 20, max: 100)

Examples:
  # List all logical networks
  metalcloud-cli logical-network list

  # List networks in a specific fabric
  metalcloud-cli logical-network list fabric-production

  # Filter by network kind
  metalcloud-cli logical-network list --filter-kind vlan

  # Filter by multiple criteria
  metalcloud-cli logical-network list --filter-kind vlan --filter-label test

  # Sort by name descending
  metalcloud-cli logical-network list --sort-by name:DESC

  # Paginate results (get page 2 with 50 records per page)
  metalcloud-cli logical-network list --page 2 --limit 50

  # Combine fabric filter with additional filters and pagination
  metalcloud-cli logical-network list fabric-1 --filter-kind vxlan --sort-by id:ASC --page 1 --limit 10

```
metalcloud-cli logical-network list [fabric_id_or_label] [flags]
```

### Options

```
      --filter-fabric-id strings           Filter by fabric ID.
      --filter-id strings                  Filter by logical network ID.
      --filter-infrastructure-id strings   Filter by infrastructure ID. Use 'null' to filter public logical networks.
      --filter-kind strings                Filter by logical network kind.
      --filter-label strings               Filter by logical network label.
  -h, --help                               help for list
      --limit int                          Number of records per page (default: 20, max: 100). (default 20)
      --page int                           Page number to retrieve (default: 1). (default 1)
      --sort-by strings                    Sort by fields (e.g., id:ASC, name:DESC).
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

* [metalcloud-cli logical-network](metalcloud-cli_logical-network.md)	 - Manage logical networks within fabrics

