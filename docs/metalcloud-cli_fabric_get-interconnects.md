## metalcloud-cli fabric get-interconnects

List the network fabric interconnects of a fabric

### Synopsis

List the network fabric interconnects this fabric takes part in.

Required Arguments:
  fabric_id    The ID or name of the fabric

Examples:
  metalcloud-cli fabric get-interconnects 12345
  metalcloud-cli fabric interconnects my-fabric

```
metalcloud-cli fabric get-interconnects fabric_id [flags]
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

