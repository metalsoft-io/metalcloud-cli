## metalcloud-cli logical-network-interconnect link get

Get one link of a logical network interconnect

### Synopsis

Get the details of a single logical network interconnect link.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect
  link_id                    The ID of the link

Examples:
  metalcloud logical-network-interconnect link get 4 9

```
metalcloud-cli logical-network-interconnect link get interconnect_id_or_label link_id [flags]
```

### Options

```
  -h, --help   help for get
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

