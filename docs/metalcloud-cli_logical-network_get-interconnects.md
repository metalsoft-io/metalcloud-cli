## metalcloud-cli logical-network get-interconnects

List the logical network interconnects of a logical network

### Synopsis

List the logical network interconnects attached to a logical network.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-interconnects 12

```
metalcloud-cli logical-network get-interconnects logical_network_id [flags]
```

### Options

```
  -h, --help   help for get-interconnects
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

