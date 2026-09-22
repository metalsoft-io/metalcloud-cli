## metalcloud-cli logical-network-interconnect link remove

Remove a link from a logical network interconnect

### Synopsis

Remove a logical network link from a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect
  link_id                    The ID of the link to remove

Examples:
  metalcloud logical-network-interconnect link remove 4 9

```
metalcloud-cli logical-network-interconnect link remove interconnect_id_or_label link_id [flags]
```

### Options

```
  -h, --help   help for remove
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

