## metalcloud-cli point-to-point-link add-ipv4-strategy

Attach a manual IPv4 subnet allocation strategy to a link

### Synopsis

Attach a manual IPv4 subnet allocation strategy to an existing point-to-point
link. This is the repair path for a link that was created without its strategy;
new links should stage the strategy on create instead.

Arguments:
  link_id            The ID of the point-to-point link

Required Flags:
  --subnet-id        ID of the IPAM subnet to allocate from
  --binding          Which interface gets the first address: a_first, b_first, or auto

Optional Flags:
  --scope-kind       Allocation scope kind (default: global)
  --scope-resource-id   Resource id for non-global scopes (default: 0)

Examples:
  metalcloud-cli p2p add-ipv4-strategy 42 --subnet-id 12345 --binding a_first
  metalcloud-cli p2p add-ipv4-strategy 42 --subnet-id 12345 --binding b_first

```
metalcloud-cli point-to-point-link add-ipv4-strategy link_id [flags]
```

### Options

```
      --binding string          Interface A binding: a_first, b_first, or auto. (default "a_first")
  -h, --help                    help for add-ipv4-strategy
      --scope-kind string       Allocation scope kind. (default "global")
      --scope-resource-id int   Resource id for non-global scopes.
      --subnet-id int           ID of the IPAM subnet to allocate from.
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

* [metalcloud-cli point-to-point-link](metalcloud-cli_point-to-point-link.md)	 - Manage point-to-point links between network interfaces

