## metalcloud-cli container list

List all containers

### Synopsis

List all containers.

Retrieves every container visible to the current user. The output can be
narrowed with the filter flags below; each filter accepts multiple values.

Optional Flags:
  --filter-id strings                   Filter by container ID
  --filter-site-id strings              Filter by site ID
  --filter-name strings                 Filter by container name
  --filter-host strings                 Filter by host name
  --filter-type-id strings              Filter by container type ID
  --filter-pool-id strings              Filter by VM pool ID
  --filter-administration-state strings Filter by administration state
  --filter-infrastructure-id strings    Filter by infrastructure ID

Examples:
  # List all containers
  metalcloud-cli container list

  # List the containers of one infrastructure
  metalcloud-cli containers ls --filter-infrastructure-id 1234

  # List the containers of a given type
  metalcloud-cli container ls --filter-type-id 42

```
metalcloud-cli container list [flags]
```

### Options

```
      --filter-administration-state strings   Filter by administration state.
      --filter-host strings                   Filter by host name.
      --filter-id strings                     Filter by container ID.
      --filter-infrastructure-id strings      Filter by infrastructure ID.
      --filter-name strings                   Filter by container name.
      --filter-pool-id strings                Filter by VM pool ID.
      --filter-site-id strings                Filter by site ID.
      --filter-type-id strings                Filter by container type ID.
  -h, --help                                  help for list
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

* [metalcloud-cli container](metalcloud-cli_container.md)	 - Manage provisioned containers

