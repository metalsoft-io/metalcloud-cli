## metalcloud-cli logical-network-interconnect get

Get logical network interconnect details

### Synopsis

Get the details of a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Examples:
  metalcloud logical-network-interconnect get 4
  metalcloud lni get dc1-dc2-logical-interconnect

```
metalcloud-cli logical-network-interconnect get interconnect_id_or_label [flags]
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

* [metalcloud-cli logical-network-interconnect](metalcloud-cli_logical-network-interconnect.md)	 - Logical network interconnect management

