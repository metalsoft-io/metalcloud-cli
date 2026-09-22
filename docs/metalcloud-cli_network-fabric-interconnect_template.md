## metalcloud-cli network-fabric-interconnect template

Show the BGP templates used for an interconnect type

### Synopsis

Show the global and neighbor BGP templates that are rendered when an
interconnect of the given type is activated or deactivated.

Required Arguments:
  interconnect_type   The interconnect type (e.g. dci-evpn)

Examples:
  metalcloud network-fabric-interconnect template dci-evpn

```
metalcloud-cli network-fabric-interconnect template interconnect_type [flags]
```

### Options

```
  -h, --help   help for template
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

