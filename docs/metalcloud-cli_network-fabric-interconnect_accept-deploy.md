## metalcloud-cli network-fabric-interconnect accept-deploy

Accept a pending interconnect deploy

### Synopsis

Accept a deploy that was started with --require-confirmation so that it proceeds.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect accept-deploy 12

```
metalcloud-cli network-fabric-interconnect accept-deploy interconnect_id_or_label [flags]
```

### Options

```
  -h, --help   help for accept-deploy
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

