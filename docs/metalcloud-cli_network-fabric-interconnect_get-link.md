## metalcloud-cli network-fabric-interconnect get-link

Get one link of an interconnect

### Synopsis

Get the details of a single link of a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect
  link_id                    The ID of the link

Examples:
  metalcloud network-fabric-interconnect get-link 12 3

```
metalcloud-cli network-fabric-interconnect get-link interconnect_id_or_label link_id [flags]
```

### Options

```
  -h, --help   help for get-link
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

