## metalcloud-cli logical-network-interconnect link add

Link a logical network to a logical network interconnect

### Synopsis

Link a logical network to a logical network interconnect so that it is stretched
across the interconnected fabrics.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect
  logical_network_id         The numeric ID of the logical network

Examples:
  metalcloud logical-network-interconnect link add 4 44

```
metalcloud-cli logical-network-interconnect link add interconnect_id_or_label logical_network_id [flags]
```

### Options

```
  -h, --help   help for add
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

