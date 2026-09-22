## metalcloud-cli logical-network-interconnect link list

List the links of a logical network interconnect

### Synopsis

List the logical networks linked to a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Optional Flags:
  --filter-logical-network-id strings   Filter by logical network ID. Repeatable or comma-separated.
  --filter-status strings               Filter by status. Repeatable or comma-separated.

Examples:
  metalcloud logical-network-interconnect link list 4
  metalcloud lni link list dc1-dc2-ln --filter-status active

```
metalcloud-cli logical-network-interconnect link list interconnect_id_or_label [flags]
```

### Options

```
      --filter-logical-network-id strings   Filter by logical network ID.
      --filter-status strings               Filter by status.
  -h, --help                                help for list
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

* [metalcloud-cli logical-network-interconnect link](metalcloud-cli_logical-network-interconnect_link.md)	 - Logical network interconnect link management

