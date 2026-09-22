## metalcloud-cli vm list

List all virtual machines with optional filtering

### Synopsis

List every virtual machine visible to the current user, across all VM pools and
infrastructures. Results are paginated transparently - the command walks every
page and prints the combined list.

Optional Flags:
  --filter-id                     Filter by VM ID.
  --filter-site-id                Filter by site ID.
  --filter-name                   Filter by VM name.
  --filter-address                Filter by VM address.
  --filter-host                   Filter by host name.
  --filter-hosts                  Filter by hosts list.
  --filter-type-id                Filter by VM type ID.
  --filter-pool-id                Filter by VM pool ID.
  --filter-administration-state   Filter by administration state.
  --filter-numa-nodes             Filter by NUMA node.
  --filter-infrastructure-id      Filter by infrastructure ID.

Every filter accepts the API filter DSL ('$eq:value', '$not:$eq:value',
'$contains:value', ...); a bare value is treated as '$eq:value'. Repeat a flag
or pass a comma-separated list to combine several values.

Examples:
  # List all VMs
  metalcloud-cli vm list

  # List the VMs of one pool
  metalcloud-cli vm list --filter-pool-id 70

  # List the running VMs of one infrastructure
  metalcloud-cli vm list --filter-infrastructure-id 123 --filter-administration-state active

  # List VMs whose name contains 'prod'
  metalcloud-cli vm list --filter-name '$contains:prod'

```
metalcloud-cli vm list [flags]
```

### Options

```
      --filter-address strings                Filter the result by VM address.
      --filter-administration-state strings   Filter the result by administration state.
      --filter-host strings                   Filter the result by host name.
      --filter-hosts strings                  Filter the result by hosts list.
      --filter-id strings                     Filter the result by VM ID.
      --filter-infrastructure-id strings      Filter the result by infrastructure ID.
      --filter-name strings                   Filter the result by VM name.
      --filter-numa-nodes strings             Filter the result by NUMA node.
      --filter-pool-id strings                Filter the result by VM pool ID.
      --filter-site-id strings                Filter the result by site ID.
      --filter-type-id strings                Filter the result by VM type ID.
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

* [metalcloud-cli vm](metalcloud-cli_vm.md)	 - Manage virtual machines lifecycle and configuration

